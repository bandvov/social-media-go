package domain

import "context"

type MockCommentRepository struct {
	AddCommentFunc              func(ctx context.Context, comment Comment) error
	FetchCommentsByEntityIDFunc func(ctx context.Context, entityID, offset, limit int) ([]Comment, error)
	CountByEntityIDsFunc        func(ctx context.Context, entityIDs []int) ([]CommentCount, error)
}

func (m *MockCommentRepository) AddComment(ctx context.Context, comment Comment) error {
	return m.AddCommentFunc(ctx, comment)
}

func (m *MockCommentRepository) FetchCommentsByEntityID(ctx context.Context, entityID, offset, limit int) ([]Comment, error) {
	return m.FetchCommentsByEntityIDFunc(ctx, entityID, offset, limit)
}
