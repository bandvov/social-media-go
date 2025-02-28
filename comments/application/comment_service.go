package application

import (
	"comments/domain"
	"context"
)

// CommentServiceInterface defines methods for tags-related operations.
type CommentServiceInterface interface {
	AddComment(ctx context.Context, c *domain.Comment) error
	GetCommentsByEntityID(ctx context.Context, entityID, userID, offset, limit int) ([]domain.Comment, error)
	GetCommentsAndRepliesCount(ctx context.Context, entityIDs []int) ([]domain.CommentCount, error)
}
type CommentService struct {
	commentRepo domain.CommentRepository
}

func NewCommentService(repo domain.CommentRepository) *CommentService {
	return &CommentService{
		commentRepo: repo,
	}
}

func (s *CommentService) AddComment(ctx context.Context, c *domain.Comment) error {
	comment := domain.Comment{
		EntityID:   c.EntityID,
		EntityType: c.EntityType,
		Content:    c.Content,
		AuthorID:   c.AuthorID,
		Status:     domain.Active,
	}
	return s.commentRepo.AddComment(ctx, comment)
}

func (s *CommentService) GetCommentsByEntityID(ctx context.Context, entityID, userID, offset, limit int) ([]domain.Comment, error) {
	return s.commentRepo.FetchCommentsByEntityID(ctx, entityID, userID, offset, limit)
}

func (s *CommentService) GetCommentsAndRepliesCount(ctx context.Context, entityIDs []int) ([]domain.CommentCount, error) {
	return s.commentRepo.CountByEntityIDs(ctx, entityIDs)
}
