package application

import (
	"errors"
	"testing"
	"time"

	"github.com/bandvov/social-media-go/domain"
	"github.com/bandvov/social-media-go/infrastructure"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterUser(t *testing.T) {
	tests := []struct {
		name           string
		input          domain.CreateUserRequest
		mockRepoFunc   func(user *domain.User) error
		expectedErr    error
		validateOutput func(t *testing.T, user *domain.User, input domain.CreateUserRequest)
	}{
		{
			name: "successful registration",
			input: domain.CreateUserRequest{
				Email:    "test@example.com",
				Password: "securepassword",
			},
			mockRepoFunc: func(user *domain.User) error {
				return nil
			},
			expectedErr: nil,
			validateOutput: func(t *testing.T, user *domain.User, input domain.CreateUserRequest) {
				if user.Email != input.Email {
					t.Errorf("expected email %s, got %s", input.Email, user.Email)
				}
				if user.Status != "active" {
					t.Errorf("expected status 'active', got %s", user.Status)
				}
				if user.Role != "user" {
					t.Errorf("expected role 'user', got %s", user.Role)
				}

				// Validate password hashing
				err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
				if err != nil {
					t.Errorf("expected password to be hashed correctly, got error: %v", err)
				}
			},
		},
		{
			name: "repository failure",
			input: domain.CreateUserRequest{
				Email:    "test@example.com",
				Password: "securepassword",
			},
			mockRepoFunc: func(user *domain.User) error {
				return errors.New("database error")
			},
			expectedErr:    errors.New("database error"),
			validateOutput: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &infrastructure.MockUserRepository{
				CreateUserFunc: tt.mockRepoFunc,
			}

			userService := NewUserService(mockRepo, &infrastructure.MockRedisCache{
				SetFunc: func(key string, value interface{}, ttl time.Duration) error {
					return nil
				},
				GetFunc: func(key string) (interface{}, error) {
					return nil, nil
				},
			})

			err := userService.RegisterUser(tt.input)

			// Assert errors
			if (err != nil && tt.expectedErr == nil) || (err == nil && tt.expectedErr != nil) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("expected error message '%s', got '%s'", tt.expectedErr.Error(), err.Error())
			}

			// Validate output if applicable
			if tt.validateOutput != nil {
				mockRepo.CreateUserFunc = func(user *domain.User) error {
					tt.validateOutput(t, user, tt.input)
					return nil
				}
			}
		})
	}
}
