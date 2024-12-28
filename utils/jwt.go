package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	JWTSecretKey []byte
	once         sync.Once
)

// Claims defines the custom claims structure.
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// loadSecret initializes the JWT secret key from the environment only once.
func loadSecret() {
	once.Do(func() {
		if JWTSecretKey == nil {
			jwtSecret := os.Getenv("JWT_SECRET")
			if jwtSecret == "" {
				log.Fatal("JWT_SECRET is not set in the environment")
			}
			JWTSecretKey = []byte(jwtSecret)
		}
	})
}

// GenerateJWT generates a new JWT token for a user.
func GenerateJWT(userID int) (string, error) {
	loadSecret()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecretKey)
}

// ValidateJWT validates a JWT token and returns the user claims.
func ValidateJWT(tokenString string) (*Claims, error) {
	loadSecret()

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JWTSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
