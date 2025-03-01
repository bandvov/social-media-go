package application

import (
	"context"
	"users/domain"
)

type MockUserService struct {
	AuthenticateFunc       func(ctx context.Context, email, password string) (*domain.User, error)
	RegisterUserFunc       func(ctx context.Context, user domain.CreateUserRequest) error
	UpdateUserDataFunc     func(ctx context.Context, user *domain.User) error
	ChangeUserRoleFunc     func(ctx context.Context, userID int, newRole string, isAdmin bool) error
	FindByEmailFunc        func(ctx context.Context, email string) (*domain.User, error)
	GetUserByIDFunc        func(ctx context.Context, id int) (*domain.User, error)
	GetPublicProfilesFunc  func(ctx context.Context, limit, offset int) ([]domain.User, error)
	GetAdminProfilesFunc   func(ctx context.Context, limit, offset int) ([]domain.User, error)
	GetUserProfileInfoFunc func(ctx context.Context, id int) (*domain.User, error)
	GetUsersByIDsFunc      func(ctx context.Context, userIDs []int) (map[int]domain.User, error)
}

func (m *MockUserService) Authenticate(ctx context.Context, email, password string) (*domain.User, error) {
	return m.AuthenticateFunc(ctx, email, password)
}
func (m *MockUserService) RegisterUser(ctx context.Context, user domain.CreateUserRequest) error {
	return m.RegisterUserFunc(ctx, user)
}
func (m *MockUserService) ChangeUserRole(ctx context.Context, userID int, newRole string, isAdmin bool) error {
	return m.ChangeUserRoleFunc(ctx, userID, newRole, isAdmin)
}

func (m *MockUserService) UpdateUserData(ctx context.Context, user *domain.User) error {
	return m.UpdateUserDataFunc(ctx, user)
}
func (m *MockUserService) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return m.FindByEmailFunc(ctx, email)
}
func (m *MockUserService) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	return m.GetUserByIDFunc(ctx, id)
}
func (m *MockUserService) GetPublicProfiles(ctx context.Context, limit, offset int) ([]domain.User, error) {
	return m.GetPublicProfilesFunc(ctx, limit, offset)
}
func (m *MockUserService) GetAdminProfiles(ctx context.Context, limit, offset int) ([]domain.User, error) {
	return m.GetAdminProfilesFunc(ctx, limit, offset)
}

func (m *MockUserService) GetUserProfileInfo(ctx context.Context, id int) (*domain.User, error) {
	return m.GetUserProfileInfoFunc(ctx, id)
}

func (m *MockUserService) GetUsersByIDs(ctx context.Context, userIDs []int) (map[int]domain.User, error) {
	return m.GetUsersByIDsFunc(ctx, userIDs)
}
