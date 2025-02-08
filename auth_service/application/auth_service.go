package application

import "auth-service/domain"

type AuthApplication struct {
	AuthRepo domain.AuthRepo
}

func NewAuthApplication(authService domain.AuthRepo) *AuthApplication {
	return &AuthApplication{AuthRepo: authService}
}

func (a *AuthApplication) Authenticate(token string) (*int, error) {
	userId, err := a.AuthRepo.DecodeToken(token)
	if err != nil {
		return nil, err
	}
	exists, err := a.AuthRepo.CheckUser(userId)
	if err != nil {
		return nil, err
	}
	return exists, nil
}
