package application

import (
	"encoding/json"
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

func (s *NotificationService) SendNotification(n domain.NotificationRequest) error {

	existing, err := s.repo.FindRecentNotification(n.EntityID, string(n.Type))
	if err != nil {
		return err
	}

	if existing != nil {
		existing.ActorIDs = append(existing.ActorIDs, n.SenderId)
		existing.Message = existing.GenerateMessage()
		s.repo.Update(existing)

		jsonNotification, err := json.Marshal(existing)
		if err != nil {
			return fmt.Errorf("Failed to marshal: %v", err)
		}

		strUserId := strconv.Itoa(n.UserID)

		// Publish event to Redis
		if err := s.events.Publish("notifications:"+strUserId, string(jsonNotification)); err != nil {
			return fmt.Errorf("Failed to publish event: %v", err)
		}

	} else {
		notification := domain.Notification{
			BaseNotification: domain.BaseNotification{
				UserID:     n.UserID,
				Type:       n.Type,
				EntityType: n.EntityType,
				EntityID:   n.EntityID,
			},
			ActorIDs:  []int{n.SenderId},
			CreatedAt: time.Now().Format(time.RFC3339),
		}
		notification.Message = notification.GenerateMessage()

		if err := s.repo.Save(notification); err != nil {
			return err
		}

		jsonNotification, err := json.Marshal(notification)
		if err != nil {
			return fmt.Errorf("Failed to marshal: %v", err)
		}

		strUserId := strconv.Itoa(n.UserID)
		// Publish event to Redis
		if err := s.events.Publish("notifications:"+strUserId, string(jsonNotification)); err != nil {
			return fmt.Errorf("Failed to publish event: %v", err)
		}
	}

	return nil
}

// Fetch unsent messages and mark them as sent
func (s *NotificationService) FetchNotifications(userID string, limit, offset int) ([]domain.Notification, error) {
	return s.repo.GetNotifications(userID, limit, offset)
}

// Subscribe to real-time notifications for a specific user
func (s *NotificationService) SubscribeToNotifications(userID string, handler func(string)) error {
	return s.events.Subscribe("notifications:"+userID, handler)

}

func (s *NotificationService) MarkAsRead(notificationIDs []int) error {
	return s.repo.MarkAsRead(notificationIDs)
}
func (s *NotificationService) CountByUserID(userId string) (int, error) {
	return s.repo.CountByUserID(userId)
}
