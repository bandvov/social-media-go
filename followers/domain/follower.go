package domain

import "time"

type Follower struct {
	FollowerID int
	FolloweeID int
}

type Follow struct {
	ID                 int       `json:"id"`
	FollowedByFollower time.Time `json:"followed_by_follower,omitempty"`
	FollowsFollower    time.Time `json:"follows_follower,omitempty"`
}

func NewFollower(followerID, followeeID int) *Follower {
	return &Follower{
		FollowerID: followerID,
		FolloweeID: followeeID,
	}
}
