package handler

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/example/app/repository"
	db "github.com/example/app/db/sqlc"
)

// VIOLATION: fmt.Println in production code
func GetUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetUser called")
	userID := r.PathValue("id")
	fmt.Println("userID:", userID)

	// VIOLATION: context.Background() in handler
	ctx := context.Background()

	// VIOLATION: handler imports and uses repository directly (bypassing service)
	repo := repository.NewUserRepository()
	user, err := repo.FindByID(ctx, userID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("user:", user)
	w.WriteHeader(http.StatusOK)
}

// VIOLATION: panic outside main
func mustParseConfig(raw string) Config {
	cfg, err := parseConfig(raw)
	if err != nil {
		panic("failed to parse config: " + err.Error())
	}
	return cfg
}

// VIOLATION: log.Fatal outside main
func initDB(dsn string) *sql.DB {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

// VIOLATION: domain model with json/db tags (shown as return value from repository)
type UserDomain struct {
	ID    int64  `json:"id"    db:"id"`
	Name  string `json:"name"  db:"name"`
	Email string `json:"email" db:"email"`
}

// VIOLATION: Repository returns sqlc type directly
func (r *UserRepository) FindByIDRaw(ctx context.Context, id int64) (db.User, error) {
	return r.queries.GetUser(ctx, id)
}
