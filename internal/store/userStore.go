package store

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SqliteUserStore struct {
	db *sql.DB
}

func NewSqliteUserStore(db *sql.DB) *SqliteUserStore {
	return &SqliteUserStore{db: db}
}

type UserStore interface {
	CreateUser(*User) (*UserResponse, error)
	GetUser(userId uint64) (*User, error)
	Login(loginData UserData) (string, string, error)
}

type UserResponse struct {
	Username string
	Email    string
}

func (uh *SqliteUserStore) CreateUser(user *User) (*UserResponse, error) {
	trans, err := uh.db.Begin()
	if err != nil {
		return nil, err
	}

	defer trans.Rollback()

	checkUserExistsQuery :=
		`
		SELECT id FROM users
		WHERE email = $2 OR username = $2
		`

	var userID *int
	err = trans.QueryRow(checkUserExistsQuery, user.Email, user.Username).Scan(&userID)
	if err != sql.ErrNoRows && err != nil {
		fmt.Println("Error, couldn't verify user exists:", err)
		return nil, err
	}

	if userID != nil {
		fmt.Println("The username or email is already taken.")
		return nil, errors.New("The username or email is already taken.")
	}

	// Hash the user password...
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
	if err != nil {
		fmt.Println("Error, couldn't hash the password:", err)
		return nil, err
	}

	// Create user query
	query :=
		`
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id
		`

	err = trans.QueryRow(query, user.Username, user.Email, hashedPassword).Scan(&user.ID)
	println(err)
	if err != nil {
		return nil, err
	}

	err = trans.Commit()
	if err != nil {
		return nil, err
	}

	return &UserResponse{user.Username, user.Email}, nil
}

func (uh *SqliteUserStore) GetUser(userId uint64) (*User, error) {
	trans, err := uh.db.Begin()
	if err != nil {
		return nil, err
	}

	user := &User{}

	defer trans.Rollback()

	query :=
		`
		SELECT id, email, username FROM USERS
		WHERE id = $1
		`

	err = trans.QueryRow(query, userId).Scan(&user.ID, &user.Email, &user.Username)

	fmt.Println(err)
	if err != nil {
		return nil, err
	}

	return user, nil
}

type UserData struct {
	ID       int
	Username string `json:"username"`
	Password string `json:"password"`
}

func (uh *SqliteUserStore) Login(loginData UserData) (string, string, error) {
	trans, err := uh.db.Begin()
	if err != nil {
		return "", "", err
	}

	defer trans.Rollback()

	query :=
		`
		SELECT id, username, password_hash FROM users
		WHERE username = $1
		`

	// Check if the username is legitimate or not
	user := &UserData{}
	fmt.Printf("username: %v passowrd: %v \n", loginData.Username, loginData.Password)
	err = trans.QueryRow(query, loginData.Username).Scan(&user.ID, &user.Username, &user.Password)

	if errors.Is(err, sql.ErrNoRows) || !checkPasswordHash([]byte(loginData.Password), []byte(user.Password)) {
		fmt.Println(err)
		invalidUsername := errors.New("Invalid username or password")
		return "", "", invalidUsername
	}

	sessionToken := generateToken(32)
	csrfToken := generateToken(32)

	sessionQuery :=
		`
		UPDATE users SET
		session_token = $1,
		csrf_token = $2,
		updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		`

	fmt.Printf("sessionToken: %v , csrfToken: %v %v \n", sessionToken, csrfToken, user.ID)
	_, err = trans.Exec(sessionQuery, sessionToken, csrfToken, user.ID)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	trans.Commit() // applies the update

	return sessionToken, csrfToken, nil
}

func checkPasswordHash(password, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	return err == nil
}

func generateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}

	return base64.URLEncoding.EncodeToString(bytes)
}
