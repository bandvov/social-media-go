package domain

type NotificationRepository interface {
	Save(notification Notification) error
	GetNotifications(userID string, limit, offset int) ([]Notification, error)
	MarkAsRead(notificationIDs []int) error
	CountByUserID(userID string) (int, error)
}
