package application

import (
	"github.com/bandvov/social-media-go/domain"
)

type PostServiceInterface interface {
	CreatePost(post *domain.CreatePostRequest) error
	DeletePost(id int) error
	UpdatePost(post *domain.Post) error
	GetPostByID(id int) (*domain.Post, error)
}

type PostService struct {
	postRepo domain.PostRepository
}

func NewPostService(repo domain.PostRepository) *PostService {
	return &PostService{postRepo: repo}
}

func (s *PostService) CreatePost(post *domain.CreatePostRequest) error {
	return s.postRepo.Create(post)
}

func (s *PostService) DeletePost(id int) error {
	return s.postRepo.Delete(id)
}

func (s *PostService) UpdatePost(post *domain.Post) error {
	return s.postRepo.Update(post)
}

func (s *PostService) GetPostByID(id int) (*domain.Post, error) {
	return s.postRepo.GetByID(id)
}
