package models

import (
	"net/http"
	"time"
)


type Balance struct {
	UserId int `json:"user_id"`
	Balance float64 `json:"balance"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (*Balance) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}