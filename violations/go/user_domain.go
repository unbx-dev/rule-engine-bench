package domain

import (
	"database/sql"
)

// VIOLATION: domain layer imports database/sql
type UserRepository interface {
	FindByID(db *sql.DB, id int64) (*User, error)
}

type User struct {
	ID    int64
	Name  string
	Email string
}
