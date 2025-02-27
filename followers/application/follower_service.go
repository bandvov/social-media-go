package application

import (
	"context"
	"errors"
	"followers/domain"
)

// FollowServiceInterface defines methods for tags-related operations.
type FollowServiceInterface interface {
	AddFollower(ctx context.Context, followerID, followeeID int) error
	RemoveFollower(ctx context.Context, followerID, followeeID int) error
	GetFollowers(ctx context.Context, userID, otherUser, limit, offset int, sort, orderBy, search string) ([]domain.Follow, error)
	GetFollowees(ctx context.Context, userID, otherUser, limit, offset int, sort, orderBy, search string) ([]domain.Follow, error)
	GetFollowerStats(ctx context.Context, userID int) (int, int, error)
	CheckFollowStatus(ctx context.Context, userID, targetUserID int64) (*domain.UserRelationship, error)
}

type FollowService struct {
	repo domain.FollowRepository
}

func NewFollowService(repo domain.FollowRepository) *FollowService {
	return &FollowService{repo: repo}
}

// AddFollower adds a follower for a given user
func (s *FollowService) AddFollower(ctx context.Context, followerID, followeeID int) error {
	// Business logic to prevent self-following
	if followerID == followeeID {
		return errors.New("user cannot follow themselves")
	}

	follower := domain.NewFollower(followerID, followeeID)
	return s.repo.AddFollower(ctx, follower)
}

// RemoveFollower removes a follower from a given user
func (s *FollowService) RemoveFollower(ctx context.Context, followerID, followeeID int) error {
	follower := domain.NewFollower(followerID, followeeID)
	return s.repo.RemoveFollower(ctx, follower)
}

// GetFollowers retrieves all followers for a user
func (s *FollowService) GetFollowers(ctx context.Context, userID, otherUser, limit, offset int, sort, orderBy, search string) ([]domain.Follow, error) {
	return s.repo.GetFollowers(ctx, userID, otherUser, limit, offset, sort, orderBy, search)
}

// GetFollowers retrieves all followers for a user
func (s *FollowService) GetFollowees(ctx context.Context, userID, otherUser, limit, offset int, sort, orderBy, search string) ([]domain.Follow, error) {
	return s.repo.GetFollowees(ctx, userID, otherUser, limit, offset, sort, orderBy, search)
}

func (s *FollowService) GetFollowerStats(ctx context.Context, userID int) (int, int, error) {
	return s.repo.GetFollowerStats(ctx, userID)
}
func (s *FollowService) CheckFollowStatus(ctx context.Context, userID, targetUserID int64) (*domain.UserRelationship, error) {
	if userID == targetUserID {
		return nil, errors.New("user cannot check relationship with themselves")
	}
	return s.repo.CheckFollowStatus(ctx, userID, targetUserID)
}
