package application

import (
	"errors"
	"fmt"

	"github.com/bandvov/social-media-go/domain"
	"golang.org/x/crypto/bcrypt"
)

// UserServiceInterface defines methods for user-related operations.
type UserServiceInterface interface {
	Authenticate(email, password string) (*domain.User, error)
	RegisterUser(user domain.CreateUserRequest) error
	UpdateUserData(domain.User) error
	ChangeUserRole(userID int, newRole string, isAdmin bool) error
	FindByEmail(email string) (*domain.User, error)
	GetUserByID(id int) (*domain.User, error)
	GetAllUsers(limit, offset int, sort string) ([]*domain.User, error)
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
	fmt.Println("here1")
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

func (s *UserService) UpdateUserData(userData domain.User) error {
	user, err := s.repo.GetUserByID(userData.ID)
	if err != nil {
		return err
	}

	if userData.Email != "" {
		user.UpdateEmail(userData.Email)
	}

	if userData.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.UpdatePassword(string(hashedPassword))
	}

	if userData.FirstName != "" {
		user.FirstName = userData.FirstName
	}

	if userData.LastName != "" {
		user.LastName = userData.LastName
	}

	if userData.Bio != "" {
		user.Bio = userData.Bio
	}

	if userData.ProfilePic != "" {
		user.ProfilePic = userData.ProfilePic
	}

	if userData.Username != "" {
		user.Username = userData.Username
	}
	if userData.Role != "" {
		user.Role = userData.Role
	}

	if userData.Status != "" {
		user.Status = userData.Status
	}

	return s.repo.UpdateUser(user)
}

func (s *UserService) ChangeUserRole(userID int, newRole string, isAdmin bool) error {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return err
	}

	return user.ChangeRole(newRole, isAdmin)
}

func (s *UserService) FindByEmail(email string) (*domain.User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s *UserService) GetUserByID(id int) (*domain.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserService) GetAllUsers(limit, offset int, sort string) ([]*domain.User, error) {
	return s.repo.GetAllUsers(limit, offset, sort)
}
