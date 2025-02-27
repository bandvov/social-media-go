package domain

import "context"

type FollowRepository interface {
	AddFollower(follower *Follower) error
	RemoveFollower(follower *Follower) error
	GetFollowers(ctx context.Context, userID, otherUser, limit, offset int, sort, orderBy, search string) ([]Follow, error)
	GetFollowees(ctx context.Context, userID, otherUser, limit, offset int, sort, orderBy, search string) ([]Follow, error)
	GetFollowerStats(ctx context.Context, userID int) (int, int, error)
	CheckFollowStatus(ctx context.Context, userID, targetUserID int64) (*UserRelationship, error)
}
