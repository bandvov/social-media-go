package interfaces

import (
	"errors"
	"regexp"
)

func ValidateEmail(email string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
}

func ValidateRole(role string) error {
	validRoles := map[string]bool{"user": true, "admin": true, "moderator": true}
	if !validRoles[role] {
		return errors.New("invalid role")
	}
	return nil
}
