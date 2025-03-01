package domain

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id int) (*User, error)
	GetPublicProfiles(ctx context.Context, offset, limit int) ([]User, error)
	GetAdminProfiles(ctx context.Context, limit, offset int) ([]User, error)
	GetUserProfileInfo(ctx context.Context, id int) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	GetUsersByIDs(ctx context.Context, userIDs []int) ([]User, error)
}
