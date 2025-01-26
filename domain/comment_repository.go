package domain

type CommentRepository interface {
	AddComment(comment Comment) error
	FetchCommentsByEntityID(entityID, userID, offset, limit int) ([]Comment, error)
}
