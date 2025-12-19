package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

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
	Login(loginData UserData, sessionToken string, tokenExpiresAt time.Duration) (*UserLoginResponse, error) // optimize
	AuthorizeUser(AuthorizeData) error
	Logout(username, sessionCookie string) (string, error)
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

type UserLoginResponse struct {
	UserID   int
	Username string
}

func (uh *SqliteUserStore) Login(loginData UserData, sessionToken string, tokenExpiresAt time.Duration) (*UserLoginResponse, error) {
	trans, err := uh.db.Begin()
	if err != nil {
		return nil, err
	}

	defer trans.Rollback()

	query :=
		`
		SELECT id, username, password_hash FROM users
		WHERE username = $1
		`

	// Check if the username is legitimate or not...
	user := &UserData{}
	err = trans.QueryRow(query, loginData.Username).Scan(&user.ID, &user.Username, &user.Password)

	if errors.Is(err, sql.ErrNoRows) || !checkPasswordHash([]byte(loginData.Password), []byte(user.Password)) {
		invalidUsername := errors.New("Invalid username or password")
		return nil, invalidUsername
	}

	updateSessionQuery :=
		`
		UPDATE sessions SET expires_at = DATETIME('now', '+24 hours')
			WHERE user_id = $1
			AND expires_at > DATETIME('now')
			RETURNING user_id
		`

	var hasToken string
	_, err = trans.Exec(updateSessionQuery, user.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if hasToken != "" {
		return &UserLoginResponse{user.ID, user.Username}, nil
	}

	// Create the new session for the user that is logging in...
	sessionQuery :=
		`
		INSERT INTO sessions (user_id, session_token, expires_at)
			VALUES ($1, $2, $3)
			RETURNING session_token;
		`

	// Set the expiration date from "now" to be 24 hours ahead
	_, err = trans.Exec(sessionQuery, user.ID, sessionToken, time.Now().Add(tokenExpiresAt))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	trans.Commit() // applies the update
	return &UserLoginResponse{user.ID, user.Username}, nil
}

func checkPasswordHash(password, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	return err == nil
}

type AuthorizeData struct {
	UserID       int
	SessionToken string
}

func (uh *SqliteUserStore) AuthorizeUser(authData AuthorizeData) error {
	trans, err := uh.db.Begin()
	if err != nil {
		return err
	}

	defer trans.Rollback()

	query :=
		`
		SELECT session_token FROM sessions
			WHERE user_id = $1
			AND expires_at > DATETIME('now')
		`

	type InHouseTokens struct {
		SessionToken string
		CSRFToken    string
	}

	fmt.Println("pre-token")
	var token string
	err = trans.QueryRow(query, authData.UserID).Scan(&token)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		fmt.Println(err)
		return err
	}
	fmt.Println("Post-token")

	isSessionTokenValid := token == authData.SessionToken

	fmt.Printf("Sess: %v\n", token)
	fmt.Printf("Sess: %v\n", authData.SessionToken)
	if !isSessionTokenValid {
		return errors.New("Unathorized user")
	}

	return nil
}

func (uh *SqliteUserStore) Logout(username, sessionCookie string) (string, error) {
	trans, err := uh.db.Begin()
	if err != nil {
		return "", err
	}

	defer trans.Rollback()

	query :=
		`
		SELECT session_token FROM users
		WHERE username = $1
		`

	var sessionToken string
	err = trans.QueryRow(query, username).Scan(&sessionToken)
	if errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	return "", nil
}
