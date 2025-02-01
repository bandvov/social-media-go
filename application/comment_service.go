package application

import "github.com/bandvov/social-media-go/domain"

// CommentServiceInterface defines methods for tags-related operations.
type CommentServiceInterface interface {
	AddComment(c *domain.Comment) error
	GetCommentsByEntityID(entityID, userID, offset, limit int) ([]domain.Comment, error)
	GetCommentsByEntityIDs(entityIDs []int) (map[int][]domain.Comment, []int, []int, error)
	GetCommentsAndRepliesCount(entityIDs []int) (int, int, error)
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

func (s *CommentService) GetCommentsByEntityIDs(entityIDs []int) (map[int][]domain.Comment, []int, []int, error) {
	commentMap := make(map[int][]domain.Comment)
	userIDList := []int{}
	commentIDList := []int{}

	comments, err := s.commentRepo.GetCommentsByEntityIDs(entityIDs)
	if err != nil {
		return nil, nil, nil, err
	}

	for _, comment := range comments {
		commentMap[comment.EntityID] = append(commentMap[comment.EntityID], comment)
		userIDList = append(userIDList, comment.AuthorID)
		commentIDList = append(commentIDList, comment.EntityID)
	}

	return commentMap, commentIDList, userIDList, nil
}

func (s *CommentService) GetCommentsAndRepliesCount(entityIDs []int) (int, int, error) {
	return s.commentRepo.CountByEntityIDs(entityIDs)
}
