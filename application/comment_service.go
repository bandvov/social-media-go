package application

import "github.com/bandvov/social-media-go/domain"

// CommentServiceInterface defines methods for tags-related operations.
type CommentServiceInterface interface {
	AddComment(c *domain.Comment) error
	GetCommentsByEntityID(entityID, userID, offset, limit int) ([]domain.Comment, error)
}
type CommentService struct {
	commentRepo domain.CommentRepository
}

func NewCommentService(repo domain.CommentRepository) *CommentService {
	return &CommentService{
		commentRepo: repo,
	}
}

func (s *CommentService) AddComment(c *domain.Comment) error {
	comment := domain.Comment{
		EntityID:   c.EntityID,
		EntityType: c.EntityType,
		Content:    c.Content,
		AuthorID:   c.AuthorID,
		Status:     domain.Active,
	}
	return s.commentRepo.AddComment(comment)
}

func (s *CommentService) GetCommentsByEntityID(entityID, userID, offset, limit int) ([]domain.Comment, error) {
	return s.commentRepo.FetchCommentsByEntityID(entityID, userID, offset, limit)
}
