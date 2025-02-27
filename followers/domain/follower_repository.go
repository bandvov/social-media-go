package domain

import "context"

type FollowRepository interface {
	AddFollower(ctx context.Context, follower *Follower) error
	RemoveFollower(ctx context.Context, follower *Follower) error
	GetFollowers(ctx context.Context, userID, otherUser, limit, offset int, sort, orderBy, search string) ([]Follow, error)
	GetFollowees(ctx context.Context, userID, otherUser, limit, offset int, sort, orderBy, search string) ([]Follow, error)
	GetFollowerStats(ctx context.Context, userID int) (int, int, error)
	CheckFollowStatus(ctx context.Context, userID, targetUserID int64) (*UserRelationship, error)
}
