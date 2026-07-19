package domain

import "database/sql" // ❌ domain層でdatabase/sqlをimport

// ❌ Domain modelにjson/db tagを付ける
type User struct {
	ID    int    `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
}

// ❌ domain層でdatabase/sqlを使う
type UserRepository struct {
	db *sql.DB
}

func (r *UserRepository) FindByID(id int) (*User, error) {
	row := r.db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id)
	u := &User{}
	return u, row.Scan(&u.ID, &u.Name, &u.Email)
}
