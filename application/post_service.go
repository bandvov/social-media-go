package application

import (
	"github.com/bandvov/social-media-go/domain"
)

type PostServiceInterface interface {
	Create(post *domain.Post) error
}

type PostService struct {
	postRepo domain.PostRepository
}

func NewPostService(repo domain.PostRepository) *PostService {
	return &PostService{postRepo: repo}
}

func (s *PostService) Create(post *domain.Post) error {
	return s.postRepo.Create(post)
}
