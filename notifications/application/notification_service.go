package application

import (
	"fmt"
	"n/domain"
	"time"
)

type NotificationService struct {
	repo   domain.NotificationRepository
	events domain.EventListener
}

func NewNotificationService(repo domain.NotificationRepository, events domain.EventListener) *NotificationService {
	return &NotificationService{repo: repo, events: events}
}

func (s *NotificationService) SendNotification(userID, message string) error {
	notification := domain.Notification{
		UserID:    userID,
		Message:   message,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Save(notification); err != nil {
		return err
	}

	// Publish event to Redis
	if err := s.events.Publish("notifications", message); err != nil {
		return fmt.Errorf("Failed to publish event: %v", err)
	}

	return nil
}
