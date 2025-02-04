package domain

type NotificationRepository interface {
	Save(notification Notification) error
	GetUnsentMessages(userID string) ([]Notification, error)
	MarkAsSent(notificationIDs []int) error
	MarkAsRead(notificationIDs []int) error
}
