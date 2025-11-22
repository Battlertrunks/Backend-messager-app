package database

import "database/sql"

type User struct {
	ID			 int 		`json:"id"`
	Username string `json:"username"`
	Email		 string `json:"email"`
	Password string `json:"password"`
}

type SqliteUserStore struct {
	db *sql.DB
}

func NewSqlUserStore(db *sql.DB) *SqliteUserStore {
	return &SqliteUserStore{db: db}
}

type UserStore interface {
	CreateUser(*User) (*User, error)
}

func (lite *SqliteUserStore) CreateUser(user *User) (*User, error) {
	trans, err := lite.db.Begin()

	if err != nil {
		return nil, err
	}

	// Rollback any bad changes when exiting
	defer trans.Rollback()

	query := 
	`INSERT INTO users (username, email, password)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	err = trans.QueryRow(query, user.Username, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		return nil, err
	}

	err = trans.Commit()
	if err != nil {
		return nil, err
	}

	return user, nil
}