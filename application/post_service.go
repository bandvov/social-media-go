package application

import (
	"github.com/bandvov/social-media-go/domain"
)

type PostServiceInterface interface {
	Create(post *domain.CreatePostRequest) error
	Delete(id int) error
	Update(post *domain.Post) error
}

type PostService struct {
	postRepo domain.PostRepository
}

func NewPostService(repo domain.PostRepository) *PostService {
	return &PostService{postRepo: repo}
}

func (s *PostService) Create(post *domain.CreatePostRequest) error {
	return s.postRepo.Create(post)
}

func (s *PostService) Delete(id int) error {
	return s.Delete(id)
}
func (s *PostService) Update(post *domain.Post) error {
	return s.Update(post)
}
