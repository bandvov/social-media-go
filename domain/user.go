package domain

import (
	"errors"
	"time"
)

type User struct {
	ID         int64
	Username   string
	Password   string
	Email      string
	Status     string // "active", "inactive", "banned"
	Role       string // "user", "admin", "moderator"
	FirstName  string
	LastName   string
	ProfilePic string    // URL to profile picture
	Bio        string    // Short biography
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
