package application

import (
	"errors"

	"github.com/bandvov/social-media-go/domain"
	"golang.org/x/crypto/bcrypt"
)

// UserServiceInterface defines methods for user-related operations.
type UserServiceInterface interface {
	Authenticate(email, password string) (*domain.User, error)
	RegisterUser(user domain.CreateUserRequest) error
	UpdateUserData(userID int64, email, password, firstName, lastName, bio, profilePic string) error
	ChangeUserRole(userID int64, newRole string, isAdmin bool) error
	FindByEmail(email string) (*domain.User, error)
	GetUserByID(id int64) (*domain.User, error)
}
type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
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
		return nil, errors.New("invalid email or password")
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

func (s *UserService) UpdateUserData(userID int64, email, password, firstName, lastName, bio, profilePic string) error {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return err
	}

	if email != "" {
		user.UpdateEmail(email)
	}

	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.UpdatePassword(string(hashedPassword))
	}

	if firstName != "" {
		user.FirstName = firstName
	}

	if lastName != "" {
		user.LastName = lastName
	}

	if bio != "" {
		user.Bio = bio
	}

	if profilePic != "" {
		user.ProfilePic = profilePic
	}

	return s.repo.UpdateUser(user)
}

func (s *UserService) ChangeUserRole(userID int64, newRole string, isAdmin bool) error {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return err
	}

	return user.ChangeRole(newRole, isAdmin)
}
func (s *UserService) FindByEmail(email string) (*domain.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByID(id int64) (*domain.User, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
