package domain

import "context"

type UserRepository interface {
	CreateUser(user *User) error
	GetUserByUsername(username string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	GetPublicProfiles(offset, limit int) ([]User, error)
	GetAdminProfiles(limit, offset int) ([]User, error)
	GetUserProfileInfo(id, otherUser int) (*User, error)
	UpdateUser(user *User) error
	GetUsersByID(ctx context.Context, userIDs []int) ([]User, error)
}
