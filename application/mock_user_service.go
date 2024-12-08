package application

import (
	"github.com/bandvov/social-media-go/domain"
)

type MockUserService struct {
	AuthenticateFunc   func(email, password string) (*domain.User, error)
	RegisterUserFunc   func(user domain.CreateUserRequest) error
	UpdateUserDataFunc func(userID int, email, password, firstName, lastName, bio, profilePic string) error
	ChangeUserRoleFunc func(userID int, newRole string, isAdmin bool) error
	FindByEmailFunc    func(email string) (*domain.User, error)
	GetUserByIDFunc    func(id int) (*domain.User, error)
	GetAllUsersFunc    func(limit, offset int, sort string) ([]*domain.User, error)
}

func (m *MockUserService) Authenticate(email, password string) (*domain.User, error) {
	return m.AuthenticateFunc(email, password)
}
func (m *MockUserService) RegisterUser(user domain.CreateUserRequest) error {
	return m.RegisterUserFunc(user)
}
func (m *MockUserService) ChangeUserRole(userID int, newRole string, isAdmin bool) error {
	return m.ChangeUserRoleFunc(userID, newRole, isAdmin)
}

func (m *MockUserService) UpdateUserData(userID int, email, password, firstName, lastName, bio, profilePic string) error {
	return m.UpdateUserDataFunc(userID, email, password, firstName, lastName, bio, profilePic)
}
func (m *MockUserService) FindByEmail(email string) (*domain.User, error) {
	return m.FindByEmailFunc(email)
}
func (m *MockUserService) GetUserByID(id int) (*domain.User, error) {
	return m.GetUserByIDFunc(id)
}

func (m *MockUserService) GetAllUsers(limit, offset int, sort string) ([]*domain.User, error) {
	return m.GetAllUsersFunc(limit, offset, sort)
}
