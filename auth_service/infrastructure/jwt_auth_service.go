package infrastructure

import (
	"auth-service/domain"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthService struct {
	secretKey           string
	userMicroserviceURL string
}

func NewJWTAuthService(secretKey, userMicroserviceURL string) *JWTAuthService {
	return &JWTAuthService{secretKey: secretKey, userMicroserviceURL: userMicroserviceURL}
}

// generateToken creates a JWT token with claims and signs it with the secret key
func (s *JWTAuthService) generateToken(data string, expirationTime time.Duration) (string, error) {
	// Create claims
	claims := jwt.MapClaims{
		"sub": data,                                  // The user ID (subject)
		"exp": time.Now().Add(expirationTime).Unix(), // Expiration time
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *JWTAuthService) DecodeToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secretKey), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok || int64(exp) < time.Now().Unix() {
		return "", errors.New("token expired")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid token subject")
	}

	return userID, nil
}

func (s *JWTAuthService) CheckUser(userId string) (*int, error) {
	token, err := s.generateToken("auth-service", time.Hour)
	if err != nil {
		return nil, err
	}
	// Create the request URL for the user microservice
	url := fmt.Sprintf("%s/users/check", s.userMicroserviceURL)

	req, _ := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Set("Authorization", "Bearer "+token)

	// Make the HTTP GET request to the user microservice
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("user not found")
	}

	// Decode the response body into a User object
	var user domain.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user.Id, nil

}
