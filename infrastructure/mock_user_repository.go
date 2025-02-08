package infrastructure

import (
	"context"

	"github.com/bandvov/social-media-go/domain"
)

type MockUserRepository struct {
	CreateUserFunc         func(user *domain.User) error
	GetUserByUsernameFunc  func(username string) (*domain.User, error)
	GetUserByEmailFunc     func(email string) (*domain.User, error)
	GetUserByIDFunc        func(id int) (*domain.User, error)
	GetPublicProfilesFunc  func(limit, offset int) ([]domain.User, error)
	GetAdminProfilesFunc   func(limit, offset int) ([]domain.User, error)
	GetUserProfileInfoFunc func(id, authenticatedUser int) (*domain.User, error)
	UpdateUserFunc         func(user *domain.User) error
	GetUsersByIDFunc       func(ctx context.Context, userIDs []int) ([]domain.User, error)
}

func (m *MockUserRepository) CreateUser(user *domain.User) error {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(user)
	}
	return nil
}

func (m *MockUserRepository) GetUserByEmail(email string) (*domain.User, error) {
	if m.GetUserByEmailFunc != nil {
		return m.GetUserByEmailFunc(email)
	}
	return nil, nil
}

func (m *MockUserRepository) GetUserByID(id int) (*domain.User, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(id)
	}
	return nil, nil
}

func (m *MockUserRepository) GetPublicProfiles(limit, offset int) ([]domain.User, error) {
	if m.GetPublicProfilesFunc != nil {
		return m.GetPublicProfilesFunc(limit, offset)
	}
	return nil, nil
}
func (m *MockUserRepository) GetAdminProfiles(limit, offset int) ([]domain.User, error) {
	if m.GetAdminProfilesFunc != nil {
		return m.GetAdminProfilesFunc(limit, offset)
	}
	return nil, nil
}

func (m *MockUserRepository) GetUserProfileInfo(id, authenticatedUser int) (*domain.User, error) {
	if m.GetUserProfileInfoFunc != nil {
		return m.GetUserProfileInfoFunc(id, authenticatedUser)
	}
	return nil, nil
}

func (m *MockUserRepository) GetUserByUsername(username string) (*domain.User, error) {
	if m.GetUserByUsernameFunc != nil {
		return m.GetUserByUsernameFunc(username)
	}
	return nil, nil
}

func (m *MockUserRepository) UpdateUser(user *domain.User) error {
	if m.UpdateUserFunc != nil {
		return m.UpdateUserFunc(user)
	}
	return nil
}

func (m *MockUserRepository) GetUsersByID(ctx context.Context, userIDs []int) ([]domain.User, error) {
	if m.GetUsersByIDFunc != nil {
		return m.GetUsersByIDFunc(ctx, userIDs)
	}
	return nil, nil
}
