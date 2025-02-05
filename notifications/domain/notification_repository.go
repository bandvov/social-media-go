package domain

type NotificationRepository interface {
	Save(notification Notification) error
	GetUnreadMessages(userID string) ([]Notification, error)
	MarkAsRead(notificationIDs []int) error
}
