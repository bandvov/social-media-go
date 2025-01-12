package application

import (
	"errors"
	"strconv"
	"time"

	"github.com/bandvov/social-media-go/domain"
	"github.com/bandvov/social-media-go/infrastructure"
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
	GetAllUsers(limit, offset int, sort, orderBy, search string) ([]*domain.User, error)
}
type UserService struct {
	repo  domain.UserRepository
	cache infrastructure.Cache
}

func NewUserService(repo domain.UserRepository, cache infrastructure.Cache) *UserService {
	return &UserService{repo: repo, cache: cache}
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

	return s.repo.CreateUser(user)
}

func (s *UserService) Authenticate(email, password string) (*domain.User, error) {
	// Retrieve user by email
	user, err := s.repo.GetUserByEmail(email)
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
	_, err := s.repo.GetUserByID(userData.ID)
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

	return s.repo.UpdateUser(userData)
}

func (s *UserService) ChangeUserRole(userID int, newRole string, isAdmin bool) error {
	return s.repo.UpdateUser(&domain.User{
		ID:   userID,
		Role: newRole,
	})
}

func (s *UserService) GetUserByID(id int) (*domain.User, error) {
	// Try to fetch from cache
	strId := strconv.Itoa(id)
	cachedUser, err := s.cache.Get(strId)
	if err == nil && cachedUser != nil {
		if user, ok := cachedUser.(*domain.User); ok {
			return user, nil
		}
	}

	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	// Store in cache

	_ = s.cache.Set(strId, user, 24*time.Hour)
	return user, nil
}

// GetPublicProfiles retrieves public profiles with pagination
func (s *UserService) GetPublicProfiles(limit, offset int) ([]domain.User, error) {
	return s.repo.GetPublicProfiles(limit, offset)
}

// GetAdminProfiles retrieves admin profiles with pagination
func (s *UserService) GetAdminProfiles(limit, offset int) ([]domain.User, error) {
	return s.repo.GetAdminProfiles(limit, offset)
}

func (s *UserService) GetUserProfileInfo(id, otherUser int) (*domain.User, error) {
	// Try to fetch from cache
	strId := strconv.Itoa(id)
	cachedUser, err := s.cache.Get(strId)
	if err == nil && cachedUser != nil {
		if user, ok := cachedUser.(*domain.User); ok {
			return user, nil
		}
	}

	user, err := s.repo.GetUserProfileInfo(id, otherUser)
	if err != nil {
		return nil, err
	}
	// Store in cache

	_ = s.cache.Set(strId, user, 24*time.Hour)
	return user, nil
}

func (s *UserService) GetAllUsers(limit, offset int, sort, orderBy, search string) ([]*domain.User, error) {
	return s.repo.GetAllUsers(limit, offset, sort, orderBy, search)
}
