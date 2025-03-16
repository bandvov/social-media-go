package domain

import (
	"errors"
	"time"
)

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type AggregatedResponse struct {
	UserID         string    `json:"id"`
	Email          string    `json:"email,omitempty"`
	Username       string    `json:"username"`
	PostCount      int       `json:"post_count"`
	FollowersCount int       `json:"followers_count"`
	FolloweesCount int       `json:"followees_count"`
	ProfilePic     *string   `json:"profile_pic,omitempty"` // URL to profile picture
	Bio            *string   `json:"bio,omitempty"`         // Short biography
	CreatedAt      time.Time `json:"created_at,omitempty"`  // Account creation timestamp
}
type User struct {
	ID         int       `json:"id"`
	Username   *string   `json:"username,omitempty"`
	Password   string    `json:"password,omitempty"`
	Email      string    `json:"email,omitempty"`
	Status     string    `json:"status,omitempty"` // "active", "inactive", "banned"
	Role       string    `json:"role,omitempty"`   // "user", "admin", "moderator"
	FirstName  *string   `json:"first_name,omitempty"`
	LastName   *string   `json:"last_name,omitempty"`
	ProfilePic *string   `json:"profile_pic,omitempty"` // URL to profile picture
	Bio        *string   `json:"bio,omitempty"`         // Short biography
	CreatedAt  time.Time `json:"created_at,omitempty"`  // Account creation timestamp
	UpdatedAt  time.Time `json:"updated_at,omitempty"`  // Last update timestamp
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

type ErrorMessage struct {
	Message string `json:"message"`
}
