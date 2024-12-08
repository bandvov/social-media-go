package application

import (
	"context"
	"errors"

	"github.com/bandvov/social-media-go/domain"
	"github.com/google/uuid"
)

type PostServiceInterface interface {
	Create(context.Context, string, string) (*domain.Post, error)
}

type PostService struct {
	postRepo domain.PostRepository
}

func NewPostService(repo domain.PostRepository) *PostService {
	return &PostService{postRepo: repo}
}

func (s *PostService) Create(ctx context.Context, title, content string) (*domain.Post, error) {
	authorID, ok := ctx.Value("userID").(string)
	if !ok || authorID == "" {
		return nil, errors.New("unauthenticated")
	}

	post := &domain.Post{
		ID:       uuid.NewString(),
		Content:  content,
		AuthorID: authorID,
	}

	err := s.postRepo.Create(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}
