package domain

import "context"

type CommentRepository interface {
	AddComment(ctx context.Context, comment Comment) error
	FetchCommentsByEntityID(ctx context.Context, entityID, userID, offset, limit int) ([]Comment, error)
	CountByEntityIDs(ctx context.Context, entityIDs []int) ([]CommentCount, error)
}
