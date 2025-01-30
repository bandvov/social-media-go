package domain

// ActivityRepository defines operations to persist activities.
type ActivityRepository interface {
	Save(activity *Activity) error
	GetRecentActivities(userID string, limit int) ([]Activity, error)
}
