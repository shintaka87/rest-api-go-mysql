package models

import "time"

type Post struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
