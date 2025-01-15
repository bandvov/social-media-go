package application

import "github.com/bandvov/social-media-go/domain"

type MockPostService struct {
	CreatePostFunc   func(post *domain.CreatePostRequest) error
	DeletePostFunc   func(id int) error
	UpdatePostFunc   func(id int, post *domain.Post) error
	GetPostByIDFunc  func(id int) (*domain.Post, error)
	FindByUserIDFunc func(userID, otherUserId, offset, limit int) ([]domain.Post, error)
}

func (s *MockPostService) CreatePost(post *domain.CreatePostRequest) error {
	return s.CreatePostFunc(post)
}

func (s *MockPostService) DeletePost(id int) error {
	return s.DeletePostFunc(id)
}

func (s *MockPostService) UpdatePost(id int, post *domain.Post) error {
	return s.UpdatePostFunc(id, post)
}

func (s *MockPostService) GetPostByID(id int) (*domain.Post, error) {
	return s.GetPostByIDFunc(id)
}

func (s *MockPostService) GetPostsByUser(userID, otherUserId, offset, limit int) ([]domain.Post, error) {
	return s.FindByUserIDFunc(userID, otherUserId, offset, limit)
}
