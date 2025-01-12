package application

import (
	"github.com/bandvov/social-media-go/domain"
)

type MockUserService struct {
	AuthenticateFunc      func(email, password string) (*domain.User, error)
	RegisterUserFunc      func(user domain.CreateUserRequest) error
	UpdateUserDataFunc    func(user *domain.User) error
	ChangeUserRoleFunc    func(userID int, newRole string, isAdmin bool) error
	FindByEmailFunc       func(email string) (*domain.User, error)
	GetUserByIDFunc       func(id int) (*domain.User, error)
	GetPublicProfilesFunc func(limit, offset int) ([]domain.User, error)
	GetAdminProfilesFunc  func(limit, offset int) ([]domain.User, error)

	GetUserProfileInfoFunc func(id, otherUser int) (*domain.User, error)
	GetAllUsersFunc        func(limit, offset int, sort, orderBy, search string) ([]*domain.User, error)
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

func (m *MockUserService) UpdateUserData(user *domain.User) error {
	return m.UpdateUserDataFunc(user)
}
func (m *MockUserService) FindByEmail(email string) (*domain.User, error) {
	return m.FindByEmailFunc(email)
}
func (m *MockUserService) GetUserByID(id int) (*domain.User, error) {
	return m.GetUserByIDFunc(id)
}
func (m *MockUserService) GetPublicProfiles(limit, offset int) ([]domain.User, error) {
	return m.GetPublicProfilesFunc(limit, offset)
}
func (m *MockUserService) GetAdminProfiles(limit, offset int) ([]domain.User, error) {
	return m.GetAdminProfilesFunc(limit, offset)
}

func (m *MockUserService) GetUserProfileInfo(id, otherUser int) (*domain.User, error) {
	return m.GetUserProfileInfoFunc(id, otherUser)
}
func (m *MockUserService) GetAllUsers(limit, offset int, sort, orderBy, search string) ([]*domain.User, error) {
	return m.GetAllUsersFunc(limit, offset, sort, orderBy, search)
}
