package domain

import "errors"

var ErrInvalidToken = errors.New("invalid token")

type AuthRepo interface {
	DecodeToken(token string) (string, error)
	CheckUser(userId string) (*int, error)
}
