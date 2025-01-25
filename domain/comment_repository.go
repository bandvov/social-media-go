package domain

type CommentRepository interface {
	AddComment(comment Comment) error
	FetchCommentsByEntityID(entityID, offset, limit int) ([]Comment, error)
}
