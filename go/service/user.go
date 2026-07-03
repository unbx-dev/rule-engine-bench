package service

import (
	"context"
	"fmt"
	"log"
	"net/http" // ❌ service層でnet/httpをimport
)

type getUserRow struct {
	ID   int
	Name string
}

type userRepo interface {
	GetUser(ctx context.Context, id int) (*getUserRow, error)
}

type UserService struct {
	repo userRepo
}

func (s *UserService) NotifyUser(userID int) error {
	// ❌ context.Backgroundをservice層で使う
	ctx := context.Background()

	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		// ❌ log.Fatalを使う
		log.Fatal("ユーザー取得失敗:", err)
	}

	// ❌ net/httpをservice層で直接使う
	resp, err := http.Get(fmt.Sprintf("https://notify.example.com/user/%d", user.ID))
	if err != nil {
		// ❌ main以外でpanicを使う
		panic(fmt.Sprintf("通知失敗: %v", err))
	}
	defer resp.Body.Close()

	// ❌ fmt.Printlnを本番コードで使う
	fmt.Println("通知送信完了:", user.Name)
	return nil
}
