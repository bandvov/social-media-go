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

	message := ""

	switch n.Type {
	case domain.NewMention:
		message = fmt.Sprintf("%d mentioned you in a %s.", n.SenderID, n.EntityType)
	case domain.NewReaction:
		message = fmt.Sprintf("%d reacted to your %s.", n.SenderID, n.EntityType)
	case domain.NewPostComment:
		message = fmt.Sprintf("%d commented on your %s.", n.SenderID, n.EntityType)
	case domain.NewCommentReply:
		message = fmt.Sprintf("%d replied to your comment.", n.SenderID)
	case domain.NewFollower:
		message = fmt.Sprintf("%d started following you.", n.SenderID)
	case domain.NewDirectMessage:
		message = fmt.Sprintf("You received a message from %d.", n.SenderID)
	}
	notification := domain.Notification{
		UserID:     n.UserID,
		Type:       n.Type,
		Message:    message,
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
