package application

import (
	"context"
	"errors"

	"github.com/bandvov/social-media-go/domain"
	"golang.org/x/crypto/bcrypt"
)

// UserServiceInterface defines methods for user-related operations.
type UserServiceInterface interface {
	Authenticate(email, password string) (*domain.User, error)
	RegisterUser(user domain.CreateUserRequest) error
	UpdateUserData(*domain.User) error
	ChangeUserRole(userID int, newRole string, isAdmin bool) error
	GetUserByID(id int) (*domain.User, error)
	GetPublicProfiles(limit, offset int) ([]domain.User, error)
	GetAdminProfiles(limit, offset int) ([]domain.User, error)
	GetUserProfileInfo(id, otherUser int) (*domain.User, error)
	GetUsersByIDs(userIDs []int) (map[int]domain.User, error)
}
type UserService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) RegisterUser(u domain.CreateUserRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		Password: string(hashedPassword),
		Email:    u.Email,
		Status:   "active",
		Role:     "user",
	}

	return s.userRepo.CreateUser(user)
}

func (s *UserService) Authenticate(email, password string) (*domain.User, error) {
	// Retrieve user by email
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}
	user.Password = ""

	return user, nil
}

func (s *UserService) UpdateUserData(userData *domain.User) error {
	_, err := s.userRepo.GetUserByID(userData.ID)
	if err != nil {
		return err
	}

	if userData.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		userData.UpdatePassword(string(hashedPassword))
	}

	return s.userRepo.UpdateUser(userData)
}

func (s *UserService) ChangeUserRole(userID int, newRole string, isAdmin bool) error {
	return s.userRepo.UpdateUser(&domain.User{
		ID:   userID,
		Role: newRole,
	})
}

func (s *UserService) GetUserByID(id int) (*domain.User, error) {
	return s.userRepo.GetUserByID(id)
}

// GetPublicProfiles retrieves public profiles with pagination
func (s *UserService) GetPublicProfiles(limit, offset int) ([]domain.User, error) {
	return s.userRepo.GetPublicProfiles(limit, offset)
}

// GetAdminProfiles retrieves admin profiles with pagination
func (s *UserService) GetAdminProfiles(limit, offset int) ([]domain.User, error) {
	return s.userRepo.GetAdminProfiles(limit, offset)
}

func (s *UserService) GetUserProfileInfo(id, otherUser int) (*domain.User, error) {
	return s.userRepo.GetUserProfileInfo(id, otherUser)
}

func (s *UserService) GetUsersByIDs(userIDs []int) (map[int]domain.User, error) {
	userMap := make(map[int]domain.User)

	userDetails, err := s.userRepo.GetUsersByID(context.Background(), userIDs)
	if err != nil {
		return nil, err
	}

	for _, user := range userDetails {
		userMap[user.ID] = user
	}

	return userMap, nil
}
