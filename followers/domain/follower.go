package domain

type Follower struct {
	FollowerID int
	FolloweeID int
}

func NewFollower(followerID, followeeID int) *Follower {
	return &Follower{
		FollowerID: followerID,
		FolloweeID: followeeID,
	}
}
