package application

import "auth-service/domain"

type AuthApplication struct {
	authService domain.AuthService
}

func NewAuthApplication(authService domain.AuthService) *AuthApplication {
	return &AuthApplication{authService: authService}
}

func (a *AuthApplication) Authenticate(token string) (string, error) {
	return a.authService.VerifyToken(token)
}
