package domain

import (
	"encoding/json"
	"errors"
	"time"
)

type User struct {
	ID         int             `json:"id"`
	Username   string          `json:"username,omitempty"`
	Password   string          `json:"password,omitempty"`
	Email      string          `json:"email"`
	Status     string          `json:"status"` // "active", "inactive", "banned"
	Role       string          `json:"role"`   // "user", "admin", "moderator"
	FirstName  string          `json:"first_name,omitempty"`
	LastName   string          `json:"last_name,omitempty"`
	ProfilePic string          `json:"profile_pic,omitempty"` // URL to profile picture
	Bio        string          `json:"bio,omitempty"`         // Short biography
	CreatedAt  time.Time       `json:"created_at"`            // Account creation timestamp
	UpdatedAt  time.Time       `json:"updated_at"`            // Last update timestamp
	PostsCount int             `json:"posts_count"`
	Followers  json.RawMessage `json:"followers"`
	Followeees json.RawMessage `json:"followees"`
}

type CreateUserRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserSearchOptions struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Sort   string `json:"sort"`
	Search string `json:"search"`
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

