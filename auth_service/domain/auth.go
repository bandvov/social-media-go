package domain

import "errors"

var ErrInvalidToken = errors.New("invalid token")

type AuthService interface {
	VerifyToken(token string) (string, error)
}
