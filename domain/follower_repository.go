package domain

type FollowerRepository interface {
	AddFollower(follower *Follower) error
	RemoveFollower(follower *Follower) error
	GetFollowers(userID, limit, offset int, sort, orderBy, search string) ([]User, error)
	GetFollowees(userID, limit, offset int, sort, orderBy, search string) ([]User, error)
}
