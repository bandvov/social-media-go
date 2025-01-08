package domain

type FollowerRepository interface {
	AddFollower(follower *Follower) error
	RemoveFollower(follower *Follower) error
	GetFollowers(userID, otherUser, limit, offset int, sort, orderBy, search string) ([]User, error)
	GetFollowees(userID, otherUser, limit, offset int, sort, orderBy, search string) ([]User, error)
}
