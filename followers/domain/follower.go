package domain

type UserRelationship struct {
	IsFollower bool
	IsFollowed bool
}
type Follower struct {
	FollowerID int `json:"follower_id"`
	FolloweeID int `json:"followee_id"`
}

type Follow struct {
	ID                 int  `json:"id"`
	FollowedByFollower bool `json:"followed_by_follower,omitempty"`
	FollowsFollower    bool `json:"follows_follower,omitempty"`
}

func NewFollower(followerID, followeeID int) *Follower {
	return &Follower{
		FollowerID: followerID,
		FolloweeID: followeeID,
	}
}
