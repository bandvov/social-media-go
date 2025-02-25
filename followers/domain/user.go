package domain

import "time"

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
