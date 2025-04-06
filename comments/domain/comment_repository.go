package domain

import (
	"context"
)

type CommentRepository interface {
	AddComment(ctx context.Context, comment Comment) error
	FetchCommentsByEntityID(ctx context.Context, entityID int, pagination Pagination) ([]Comment, error)
	CountByEntityIDAndType(ctx context.Context, entities []Entity) ([]CommentCount, error)
}
