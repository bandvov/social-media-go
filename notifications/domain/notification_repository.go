package domain

type NotificationRepository interface {
	Save(notification Notification) error
	GetUnsentMessages(userID string) ([]Notification, error)
	MarkAsSent(notificationID string) error
}
