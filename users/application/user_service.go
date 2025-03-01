package application

import (
	"context"
	"errors"
	"users/domain"

	"golang.org/x/crypto/bcrypt"
)

// UserServiceInterface defines methods for user-related operations.
type UserServiceInterface interface {
	Authenticate(ctx context.Context, email, password string) (*domain.User, error)
	RegisterUser(ctx context.Context, user domain.CreateUserRequest) error
	UpdateUserData(ctx context.Context, user *domain.User) error
	ChangeUserRole(ctx context.Context, userID int, newRole string, isAdmin bool) error
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
	GetPublicProfiles(ctx context.Context, limit, offset int) ([]domain.User, error)
	GetAdminProfiles(ctx context.Context, limit, offset int) ([]domain.User, error)
	GetUserProfileInfo(ctx context.Context, id int) (*domain.User, error)
	GetUsersByIDs(ctx context.Context, userIDs []int) (map[int]domain.User, error)
}
type UserService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) RegisterUser(ctx context.Context, u domain.CreateUserRequest) error {
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

	return s.userRepo.CreateUser(ctx, user)
}

func (s *UserService) Authenticate(ctx context.Context, email, password string) (*domain.User, error) {
	// Retrieve user by email
	user, err := s.userRepo.GetUserByEmail(ctx, email)
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

func (s *UserService) UpdateUserData(ctx context.Context, userData *domain.User) error {
	_, err := s.userRepo.GetUserByID(ctx, userData.ID)
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

	return s.userRepo.UpdateUser(ctx, userData)
}

func (s *UserService) ChangeUserRole(ctx context.Context, userID int, newRole string, isAdmin bool) error {
	return s.userRepo.UpdateUser(ctx, &domain.User{
		ID:   userID,
		Role: newRole,
	})
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

// GetPublicProfiles retrieves public profiles with pagination
func (s *UserService) GetPublicProfiles(ctx context.Context, limit, offset int) ([]domain.User, error) {
	return s.userRepo.GetPublicProfiles(ctx, limit, offset)
}

// GetAdminProfiles retrieves admin profiles with pagination
func (s *UserService) GetAdminProfiles(ctx context.Context, limit, offset int) ([]domain.User, error) {
	return s.userRepo.GetAdminProfiles(ctx, limit, offset)
}

func (s *UserService) GetUserProfileInfo(ctx context.Context, id int) (*domain.User, error) {
	return s.userRepo.GetUserProfileInfo(ctx, id)
}

func (s *UserService) GetUsersByIDs(ctx context.Context, userIDs []int) (map[int]domain.User, error) {
	userMap := make(map[int]domain.User)

	userDetails, err := s.userRepo.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	for _, user := range userDetails {
		userMap[user.ID] = user
	}

	return userMap, nil
}
