package repository

import (
	"database/sql"

	"github.com/google/uuid"
	
	"github.com/NhatPixel/cinema-notification-service/internal/model"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepo(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(n *model.Notification) error {
	n.ID = uuid.NewString()

	query := `
		INSERT INTO notifications (id, user_id, title, content)
		VALUES (?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, n.ID, n.UserID, n.Title, n.Content)
	return err
}