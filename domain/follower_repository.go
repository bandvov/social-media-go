package domain

type FollowerRepository interface {
	AddFollower(follower *Follower) error
	RemoveFollower(follower *Follower) error
	GetFollowers(userID int) ([]User, error)
}
