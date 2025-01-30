package application

import (
	"activity-service/domain"
	"encoding/json"
	"time"
)

// ActivityService provides methods to manage the activity feed.
type ActivityService struct {
	repo domain.ActivityRepository
}

// NewActivityService creates a new instance of ActivityService.
func NewActivityService(repo domain.ActivityRepository) *ActivityService {
	return &ActivityService{repo: repo}
}

// AddActivity records a new user activity with event data.
func (s *ActivityService) AddActivity(userID, action, targetID string, eventData map[string]interface{}) error {
	eventDataJSON, err := json.Marshal(eventData)
	if err != nil {
		return err
	}

	activity := &domain.Activity{
		UserID:    userID,
		Action:    action,
		TargetID:  targetID,
		EventData: eventDataJSON,
		CreatedAt: time.Now(),
	}
	return s.repo.Save(activity)
}

// GetRecentActivities retrieves the latest activities for a user.
func (s *ActivityService) GetRecentActivities(userID string, limit int) ([]domain.Activity, error) {
	return s.repo.GetRecentActivities(userID, limit)
}
