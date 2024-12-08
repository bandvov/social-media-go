package domain

import (
	"errors"
	"time"
)

type User struct {
	ID         int
	Username   string `json:"username"`
	Password   string `json:"password"`
	Email      string
	Status     string    // "active", "inactive", "banned"
	Role       string    // "user", "admin", "moderator"
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	ProfilePic string    `json:"profile_pic"` // URL to profile picture
	Bio        string    `json:"bio"`         // Short biography
	CreatedAt  time.Time // Account creation timestamp
	UpdatedAt  time.Time // Last update timestamp
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

func (u *User) ChangeRole(newRole string, isAdmin bool) error {
	if !isAdmin {
		return errors.New("only admin can change roles")
	}
	u.Role = newRole
	return nil
}
