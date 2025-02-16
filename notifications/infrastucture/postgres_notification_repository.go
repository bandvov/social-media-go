package infrastructure

import (
	"database/sql"
	"errors"
	"fmt"
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
	fmt.Printf("here===========\n", "%+v", notification)

	_, err := r.db.Exec(
		"INSERT INTO notifications (user_id, actor_ids, message, type, entity_type, entity_id, created_at) VALUES($1, $2, $3, $4, $5, $6, NOW())",
		notification.UserID, notification.ActorIDs, notification.Message, notification.Type, notification.EntityType, notification.EntityID,
	)
	return err
}

func (r *PostgresNotificationRepository) Update(notification *domain.Notification) error {
	_, err := r.db.Exec(`
		UPDATE notifications 
		SET actor_ids = $1, message = $2
		WHERE id = $3`,
		notification.ActorIDs,
		notification.Message,
		notification.ID,
	)
	return err
}

// Get all unsent messages for a user
func (r *PostgresNotificationRepository) GetNotifications(userID string, limit, offset int) ([]domain.Notification, error) {
	// Use prepared statement
	stmt, err := r.db.Prepare(`
		SELECT id, user_id, actor_ids, message, type, entity_type, entity_id, is_read, created_at
		FROM notifications
		WHERE user_id = $1 AND is_read = false
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []domain.Notification
	for rows.Next() {
		var msg domain.Notification
		if err := rows.Scan(&msg.ID, &msg.UserID, &msg.ActorIDs, &msg.Message, &msg.Type, &msg.EntityType, &msg.EntityID, &msg.IsRead, &msg.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, msg)
	}

	return notifications, nil
}

func (r *PostgresNotificationRepository) MarkAsRead(notificationIDs []int) error {
	_, err := r.db.Exec("UPDATE notifications SET is_read = true WHERE id = ANY($1)", pg.Array(notificationIDs))
	return err
}

func (r *PostgresNotificationRepository) FindRecentNotification(userID, EntityID int, eventType string) (*domain.Notification, error) {
	var notification domain.Notification
	err := r.db.QueryRow(`
		SELECT id, user_id, actor_ids, type, entity_id, entity_type, message, created_at  created_at FROM notifications 
		WHERE user_id = $1 AND entity_id = $2 AND type = $3 
		AND created_at > NOW() - INTERVAL '30 minutes' 
		ORDER BY created_at DESC LIMIT 1`, userID, EntityID, eventType).
		Scan(&notification.ID, &notification.UserID, &notification.ActorIDs, &notification.Type, &notification.EntityID, &notification.EntityType, &notification.Message, &notification.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *PostgresNotificationRepository) CountByUserID(userID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1`
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, errors.New("failed to count notifications")
	}
	return count, nil
}
