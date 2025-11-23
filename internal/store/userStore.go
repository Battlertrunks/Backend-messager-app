package store

import "database/sql"

type User struct {
	ID 				int `json:"id"`
	Username 	string `json:"username"`
	Email 		string `json:"email"`
	Password 	string `json:"password"`
}

type SqliteUserStore struct {
	db *sql.DB
}

func NewSqliteUserStore(db *sql.DB) *SqliteUserStore {
	return &SqliteUserStore{db: db}
}

type UserStore interface {
	CreateUser(*User) (*User, error)
}

func (uh *SqliteUserStore) CreateUser(user *User) (*User, error) {
	
	return user, nil
}