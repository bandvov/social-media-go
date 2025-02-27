package domain

type FollowerRepository interface {
	AddFollower(follower *Follower) error
	RemoveFollower(follower *Follower) error
	GetFollowers(userID, otherUser, limit, offset int, sort, orderBy, search string) ([]Follow, error)
	GetFollowees(userID, otherUser, limit, offset int, sort, orderBy, search string) ([]Follow, error)
	GetFollowerStats(userID int) (int, int, error)
}
