package infrastructure

import (
	"activity-service/domain"
	"database/sql"
)

// PostgresActivityRepository implements ActivityRepository.
type PostgresActivityRepository struct {
	db *sql.DB
}

// NewPostgresActivityRepository initializes a new repository.
func NewPostgresActivityRepository(db *sql.DB) *PostgresActivityRepository {
	return &PostgresActivityRepository{db: db}
}

// Save inserts an activity with event data into the database.
func (r *PostgresActivityRepository) Save(activity *domain.Activity) error {
	_, err := r.db.Exec(
		"INSERT INTO activities ( user_id, action, target_id, event_data, created_at) VALUES ($1, $2, $3, $4, $5)",
		activity.UserID, activity.Action, activity.TargetID, activity.EventData, activity.CreatedAt,
	)
	return err
}

// GetRecentActivities retrieves recent activities for a user.
func (r *PostgresActivityRepository) GetRecentActivities(userID, limit int) ([]domain.Activity, error) {
	rows, err := r.db.Query(
		"SELECT id, user_id, action, target_id, event_data, created_at FROM activities WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2",
		userID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []domain.Activity
	for rows.Next() {
		var activity domain.Activity
		err := rows.Scan(&activity.ID, &activity.UserID, &activity.Action, &activity.TargetID, &activity.EventData, &activity.CreatedAt)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}
	return activities, nil
}
