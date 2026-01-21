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
		INSERT INTO notifications (id, user_id, title, content, is_read, expires_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, n.ID, n.UserID, n.Title, n.Content, n.IsRead, n.ExpiresAt)
	return err
}

func (r *NotificationRepository) FindByUserID(uid string) ([]model.Notification, error) {
	query := `
		SELECT id, user_id, title, content, is_read, expires_at 
		FROM notifications 
		WHERE user_id = ? AND (expires_at IS NULL OR expires_at > NOW())
		ORDER BY id DESC
	`
	rows, err := r.db.Query(query, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []model.Notification

	for rows.Next() {
		var n model.Notification
		var expiresAt sql.NullTime
		err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Content, &n.IsRead, &expiresAt)
		if err != nil {
			return nil, err
		}
		if expiresAt.Valid {
			n.ExpiresAt = &expiresAt.Time
		} else {
			n.ExpiresAt = nil
		}
		notifications = append(notifications, n)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}

func (r *NotificationRepository) UpdateReadStatus(id string) error {
	query := `
		UPDATE notifications
		SET is_read = 1
		WHERE id = ?
	`
	_, err := r.db.Exec(query, id)
	return err
}
