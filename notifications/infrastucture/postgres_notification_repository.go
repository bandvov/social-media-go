package infrastructure

import (
	"database/sql"
	"n/domain"

	pg "github.com/lib/pq"
)

type PostgresNotificationRepository struct {
	db *sql.DB
}

func NewPostgresNotificationRepository(db *sql.DB) *PostgresNotificationRepository {
	return &PostgresNotificationRepository{db: db}
}

func (r *PostgresNotificationRepository) Save(notification domain.Notification) error {
	_, err := r.db.Exec(
		"INSERT INTO notifications (user_id, message, created_at) VALUES ($1, $2, $3)",
		notification.UserID, notification.Message, notification.CreatedAt,
	)
	return err
}

// Get all unsent messages for a user
func (r *PostgresNotificationRepository) GetUnsentMessages(userID string) ([]domain.Notification, error) {
	rows, err := r.db.Query("SELECT id, user_id, message, created_at FROM notifications WHERE user_id = $1 AND sent = false", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.Notification
	for rows.Next() {
		var msg domain.Notification
		if err := rows.Scan(&msg.ID, &msg.UserID, &msg.Message, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// Mark message as sent
func (r *PostgresNotificationRepository) MarkAsSent(notificationIDs []int) error {
	_, err := r.db.Exec("UPDATE notifications SET sent = true WHERE id = ANY($1)", pg.Array(notificationIDs))
	return err
}

// Mark message as read
func (r *PostgresNotificationRepository) MarkAsRead(notificationIDs []int) error {
	_, err := r.db.Exec("UPDATE notifications SET read = true WHERE id = ANY($1)", pg.Array(notificationIDs))
	return err
}
