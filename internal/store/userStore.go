package store

import (
	"database/sql"
	"errors"
	"fmt"

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
	Login(loginData UserData, sessionToken, csrfToken string) error // optimize
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

func (uh *SqliteUserStore) Login(loginData UserData, sessionToken, csrfToken string) error {
	trans, err := uh.db.Begin()
	if err != nil {
		return err
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
		return invalidUsername
	}

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
		return err
	}

	trans.Commit() // applies the update

	return nil
}

func checkPasswordHash(password, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	return err == nil
}

type AuthorizeData struct {
	Username     string `json:"username"`
	SessionToken string `json:"session_token"`
	CSRFToken    string `json:"csrf_token"`
}

func (uh *SqliteUserStore) AuthorizeUser(authData AuthorizeData) error {
	trans, err := uh.db.Begin()
	if err != nil {
		return err
	}

	defer trans.Rollback()

	query :=
		`
		SELECT session_token, csrf_token FROM users
		WHERE username = $1
		`

	type InHouseTokens struct {
		SessionToken string
		CSRFToken    string
	}

	var tokens InHouseTokens
	err = trans.QueryRow(query, authData.Username).Scan(&tokens.SessionToken, &tokens.CSRFToken)
	if errors.Is(err, sql.ErrNoRows) {
		return err
	}

	sessionTokenValid := tokens.SessionToken != "" && tokens.SessionToken == authData.SessionToken
	csrfTokenValid := tokens.CSRFToken != "" && tokens.CSRFToken == authData.CSRFToken

	fmt.Printf("Sess: %v, csrf: %v\n", tokens.SessionToken, tokens.CSRFToken)
	fmt.Printf("Sess: %v, csrf: %v\n", authData.SessionToken, authData.CSRFToken)
	if !sessionTokenValid || !csrfTokenValid {
		return errors.New("Unathorized access of user")
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
