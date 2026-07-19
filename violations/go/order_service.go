package service

import (
	"context"
	"fmt"
	"net/http"
)

// VIOLATION: service layer imports net/http
func NotifyOrderStatus(orderID string, status string) error {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost,
		"https://internal.example.com/notify", nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println("notified:", resp.Status)
	return nil
}
