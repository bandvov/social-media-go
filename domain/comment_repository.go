package domain

type CommentRepository interface {
	AddComment(comment Comment) error
	GetCommentsByEntityID(entityID int) ([]Comment, error)
}
