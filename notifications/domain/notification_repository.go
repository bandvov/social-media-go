package domain

type NotificationRepository interface {
	Save(notification Notification) error
	FindByUser(userID string) ([]Notification, error)
}
