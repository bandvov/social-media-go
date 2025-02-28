package domain

type MockCommentRepository struct {
	AddCommentFunc              func(comment Comment) error
	FetchCommentsByEntityIDFunc func(entityID, userID, offset, limit int) ([]Comment, error)
	CountByEntityIDsFunc        func(entityIDs []int) ([]CommentCount, error)
}

func (m *MockCommentRepository) AddComment(comment Comment) error {
	return m.AddCommentFunc(comment)
}
