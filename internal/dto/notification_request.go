package dto

import (
	"time"

	"github.com/NhatPixel/cinema-notification-service/internal/model"
)

type CreateRequest struct {
	UserID    string     `json:"user_id"`
	Title     string     `json:"title" binding:"required,min=1,max=100"`
	Content   string     `json:"content"`
	ExpiresAt *time.Time `json:"expires_at" binding:"omitempty,future"`
}

func (r *CreateRequest) ToModel() model.Notification {
	return model.Notification{
		UserID:    r.UserID,
		Title:     r.Title,
		Content:   r.Content,
		ExpiresAt: r.ExpiresAt,
	}
}
