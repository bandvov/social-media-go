package application

import (
	"errors"

	"github.com/bandvov/social-media-go/domain"
)

type FollowerService struct {
	repo domain.FollowerRepository
}

func NewFollowerService(repo domain.FollowerRepository) *FollowerService {
	return &FollowerService{repo: repo}
}

// AddFollower adds a follower for a given user
func (s *FollowerService) AddFollower(followerID, followeeID int) error {
	// Business logic to prevent self-following
	if followerID == followeeID {
		return errors.New("user cannot follow themselves")
	}

	follower := domain.NewFollower(followerID, followeeID)
	return s.repo.AddFollower(follower)
}

// RemoveFollower removes a follower from a given user
func (s *FollowerService) RemoveFollower(followerID, followeeID int) error {
	follower := domain.NewFollower(followerID, followeeID)
	return s.repo.RemoveFollower(follower)
}

// GetFollowers retrieves all followers for a user
func (s *FollowerService) GetFollowers(userID int) ([]domain.User, error) {
	return s.repo.GetFollowers(userID)
}
