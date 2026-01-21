package model

import "time"

type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	IsRead    bool   `json:"is_read"`
	ExpiresAt *time.Time  `json:"expires_at"`
}