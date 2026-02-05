package dto

import (
	"github.com/NhatPixel/cinema-notification-service/internal/model"
)

type SSEResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	IsRead    bool      `json:"is_read"`
}

func (s *SSEResponse) FromModel(n model.Notification) {
	s.ID = n.ID
	s.Title = n.Title
	s.Content = n.Content
	s.IsRead = n.IsRead
}
