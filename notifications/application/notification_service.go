package application

import (
	"fmt"
	"n/domain"
	"strconv"
	"time"
)

type NotificationService struct {
	repo   domain.NotificationRepository
	events domain.EventListener
}

func NewNotificationService(repo domain.NotificationRepository, events domain.EventListener) *NotificationService {
	return &NotificationService{repo: repo, events: events}
}

func (s *NotificationService) SendNotification(n domain.Notification) error {
	notification := domain.Notification{
		UserID:     n.UserID,
		Type:       n.Type,
		Message:    n.Message,
		EntityType: n.EntityType,
		EntityID:   n.EntityID,
		CreatedAt:  time.Now().Format(time.RFC3339),
	}

	if err := s.repo.Save(notification); err != nil {
		return err
	}
	strUserId := strconv.Itoa(n.UserID)
	// Publish event to Redis
	if err := s.events.Publish("notifications:"+strUserId, fmt.Sprint(notification)); err != nil {
		return fmt.Errorf("Failed to publish event: %v", err)
	}

	return nil
}

// Fetch unsent messages and mark them as sent
func (s *NotificationService) FetchUnReadMessages(userID string) ([]domain.Notification, error) {
	return s.repo.GetUnreadMessages(userID)
}

// Subscribe to real-time notifications for a specific user
func (s *NotificationService) SubscribeToNotifications(userID string, handler func(string)) error {
	return s.events.Subscribe("notifications:"+userID, handler)

}

func (s *NotificationService) MarkAsRead(notificationIDs []int) error {
	return s.repo.MarkAsRead(notificationIDs)
}
