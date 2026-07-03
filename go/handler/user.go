package handler

import (
	"context"
	"fmt"
	"net/http"

	"rule-engine-bench/repository" // ❌ handler層でrepositoryを直接import
)

// ❌ serviceを経由せずrepositoryを直接持つ
type UserHandler struct {
	repo *repository.UserRepository
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// ❌ context.Backgroundをhandler層で使う (r.Context()を使うべき)
	ctx := context.Background()

	user, err := h.repo.GetUser(ctx, 1)
	if err != nil {
		// ❌ main以外でpanicを使う
		panic(err)
	}

	// ❌ fmt.Printlnを本番コードで使う
	fmt.Println("ユーザー取得:", user.Name)
	fmt.Fprintf(w, "Hello, %s", user.Name)
}
