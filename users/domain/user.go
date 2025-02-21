package domain

import (
	"errors"
	"time"
)

type User struct {
	ID                 int       `json:"id"`
	Username           *string   `json:"username,omitempty"`
	Password           string    `json:"password,omitempty"`
	Email              string    `json:"email,omitempty"`
	Status             string    `json:"status,omitempty"` // "active", "inactive", "banned"
	Role               string    `json:"role,omitempty"`   // "user", "admin", "moderator"
	FirstName          *string   `json:"first_name,omitempty"`
	LastName           *string   `json:"last_name,omitempty"`
	ProfilePic         *string   `json:"profile_pic,omitempty"` // URL to profile picture
	Bio                *string   `json:"bio,omitempty"`         // Short biography
	CreatedAt          time.Time `json:"created_at,omitempty"`  // Account creation timestamp
	UpdatedAt          time.Time `json:"updated_at,omitempty"`  // Last update timestamp
	PostsCount         int       `json:"posts_count,omitempty"`
	FollowersCount     int       `json:"followers_count,omitempty"`
	FolloweesCount     int       `json:"followees_count,omitempty"`
	FollowsFollower    bool      `json:"follows_follower,omitempty"`
	FollowedByFollower bool      `json:"followed_by_follower,omitempty"`
	IsFollowee         bool      `json:"is_followee,omitempty"`
	IsFollower         bool      `json:"is_follower,omitempty"`
}

type CreateUserRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (u *User) UpdateEmail(newEmail string) {
	u.Email = newEmail
}

func (u *User) UpdatePassword(newPassword string) {
	u.Password = newPassword
}

func (u *User) ChangeStatus(newStatus string, isAdmin bool) error {
	if !isAdmin {
		return errors.New("only admin can change status")
	}
	u.Status = newStatus
	return nil
}
