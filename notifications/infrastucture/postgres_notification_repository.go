package infrastructure

import (
	"database/sql"
	"n/domain"
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

func (r *PostgresNotificationRepository) FindByUser(userID string) ([]domain.Notification, error) {
	rows, err := r.db.Query("SELECT id, user_id, message, created_at FROM notifications WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []domain.Notification
	for rows.Next() {
		var n domain.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Message, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}
