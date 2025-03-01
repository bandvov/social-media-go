package infrastructure

import (
	"context"
	"users/domain"
)

type MockUserRepository struct {
	CreateUserFunc         func(ctx context.Context, user *domain.User) error
	GetUserByUsernameFunc  func(ctx context.Context, username string) (*domain.User, error)
	GetUserByEmailFunc     func(ctx context.Context, email string) (*domain.User, error)
	GetUserByIDFunc        func(ctx context.Context, id int) (*domain.User, error)
	GetPublicProfilesFunc  func(ctx context.Context, limit, offset int) ([]domain.User, error)
	GetAdminProfilesFunc   func(ctx context.Context, limit, offset int) ([]domain.User, error)
	GetUserProfileInfoFunc func(ctx context.Context, id int) (*domain.User, error)
	UpdateUserFunc         func(ctx context.Context, user *domain.User) error
	GetUsersByIDsFunc      func(ctx context.Context, userIDs []int) ([]domain.User, error)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.GetUserByEmailFunc != nil {
		return m.GetUserByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockUserRepository) GetPublicProfiles(ctx context.Context, limit, offset int) ([]domain.User, error) {
	if m.GetPublicProfilesFunc != nil {
		return m.GetPublicProfilesFunc(ctx, limit, offset)
	}
	return nil, nil
}
func (m *MockUserRepository) GetAdminProfiles(ctx context.Context, limit, offset int) ([]domain.User, error) {
	if m.GetAdminProfilesFunc != nil {
		return m.GetAdminProfilesFunc(ctx, limit, offset)
	}
	return nil, nil
}

func (m *MockUserRepository) GetUserProfileInfo(ctx context.Context, id int) (*domain.User, error) {
	if m.GetUserProfileInfoFunc != nil {
		return m.GetUserProfileInfoFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	if m.GetUserByUsernameFunc != nil {
		return m.GetUserByUsernameFunc(ctx, username)
	}
	return nil, nil
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	if m.UpdateUserFunc != nil {
		return m.UpdateUserFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) GetUsersByIDs(ctx context.Context, userIDs []int) ([]domain.User, error) {
	if m.GetUsersByIDsFunc != nil {
		return m.GetUsersByIDsFunc(ctx, userIDs)
	}
	return nil, nil
}
