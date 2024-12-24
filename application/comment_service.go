package application

import "github.com/bandvov/social-media-go/domain"

type CommentService struct {
	commentRepo domain.CommentRepository
}

func NewCommentService(repo domain.CommentRepository) *CommentService {
	return &CommentService{
		commentRepo: repo,
	}
}

func (s *CommentService) AddComment(entityID int, content string, authorID int) error {
	comment := domain.Comment{
		EntityID: entityID,
		Content:  content,
		AuthorID: authorID,
	}
	return s.commentRepo.AddComment(comment)
}

func (s *CommentService) GetComments(entityID int) ([]domain.Comment, error) {
	return s.commentRepo.GetCommentsByEntityID(entityID)
}
