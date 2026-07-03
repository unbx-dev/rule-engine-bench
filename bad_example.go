package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

// --- domain層 ---

// Domain modelにjson/db tagを付ける
type User struct {
	ID    int    `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
}

// domain層でdatabase/sqlをimportして使う
type UserDomainRepository struct {
	db *sql.DB
}

func (r *UserDomainRepository) FindByID(id int) (*User, error) {
	row := r.db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id)
	u := &User{}
	return u, row.Scan(&u.ID, &u.Name, &u.Email)
}

// --- sqlc生成型 (Repositoryがsqlc型をそのまま返す) ---

// sqlcが生成するような型
type GetUserRow struct {
	ID    int
	Name  string
	Email string
}

type UserRepository struct {
	db *sql.DB
}

// Repositoryがsqlc型をそのまま返す (domain modelに変換すべき)
func (r *UserRepository) GetUser(ctx context.Context, id int) (*GetUserRow, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, name, email FROM users WHERE id = ?", id)
	u := &GetUserRow{}
	return u, row.Scan(&u.ID, &u.Name, &u.Email)
}

// FixをReplacementで返さずそのまま返す
type Diagnostic struct {
	Message string
	Fix     string // Replacementで包まずFixを直接持つ
}

func NewDiagnostic(msg, fix string) Diagnostic {
	return Diagnostic{Message: msg, Fix: fix}
}

// --- service層 ---

// service層でnet/httpをimportして使う
type UserService struct {
	repo *UserRepository
}

func (s *UserService) NotifyUser(userID int) error {
	// context.Backgroundをservice層で使う
	ctx := context.Background()

	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		// log.Fatalを使う
		log.Fatal("ユーザー取得失敗:", err)
	}

	// net/httpをservice層で直接使う
	resp, err := http.Get(fmt.Sprintf("https://notify.example.com/user/%d", user.ID))
	if err != nil {
		// main以外でpanicを使う
		panic(fmt.Sprintf("通知失敗: %v", err))
	}
	defer resp.Body.Close()

	// fmt.Printlnを本番コードで使う
	fmt.Println("通知送信完了:", user.Name)
	return nil
}

// --- handler層 ---

// handler層でrepositoryをimportして直接使う
type UserHandler struct {
	// serviceを経由せずrepositoryを直接持つ
	repo *UserRepository
}

func (h *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	// context.Backgroundをhandler層で使う (r.Context()を使うべき)
	ctx := context.Background()

	user, err := h.repo.GetUser(ctx, 1)
	if err != nil {
		// main以外でpanicを使う
		panic(err)
	}

	// fmt.Printlnを本番コードで使う
	fmt.Println("ユーザー取得:", user.Name)
	fmt.Fprintf(w, "Hello, %s", user.Name)
}

func main() {
	db, err := sql.Open("mysql", "dsn")
	if err != nil {
		// mainでのpanicはOKだが、log.Fatalを使っている
		log.Fatal(err)
	}

	repo := &UserRepository{db: db}
	handler := &UserHandler{repo: repo}

	http.HandleFunc("/user", handler.HandleGetUser)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
