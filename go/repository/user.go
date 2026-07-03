package repository

import (
	"context"
	"database/sql"
)

// sqlcが生成するような型
type GetUserRow struct {
	ID    int
	Name  string
	Email string
}

type UserRepository struct {
	db *sql.DB
}

// ❌ Repositoryがsqlc型をそのまま返す (domain modelに変換すべき)
func (r *UserRepository) GetUser(ctx context.Context, id int) (*GetUserRow, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, name, email FROM users WHERE id = ?", id)
	u := &GetUserRow{}
	return u, row.Scan(&u.ID, &u.Name, &u.Email)
}

// ❌ FixをReplacementで返さずそのまま返す
type Diagnostic struct {
	Message string
	Fix     string
}

func NewDiagnostic(msg, fix string) Diagnostic {
	return Diagnostic{Message: msg, Fix: fix}
}
