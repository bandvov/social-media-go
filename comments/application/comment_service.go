package application

import (
	"comments/domain"
	"context"

	"golang.org/x/sync/errgroup"
)

// CommentServiceInterface defines methods for tags-related operations.
type CommentServiceInterface interface {
	AddComment(ctx context.Context, c *domain.Comment) error
	GetCommentsByEntityID(ctx context.Context, entityID, userID, offset, limit int) ([]domain.Comment, error)
	GetCommentsAndRepliesCount(ctx context.Context, entityIDs []int) ([]domain.CommentCount, error)
}
type CommentService struct {
	commentRepo    domain.CommentRepository
	commentFetcher *CommentFetcher
}

func NewCommentService(repo domain.CommentRepository, commentFetcher *CommentFetcher) *CommentService {
	return &CommentService{
		commentRepo:    repo,
		commentFetcher: commentFetcher,
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

// Fetch comments with all necessary details
func (s *CommentService) GetCommentsByEntityID(ctx context.Context, entityID, userID, offset, limit int) ([]domain.Comment, error) {
	comments, err := s.commentRepo.FetchCommentsByEntityID(ctx, entityID, offset, limit)
	if err != nil {
		return nil, err
	}

	if len(comments) == 0 {
		return comments, nil
	}

	// Collect user IDs and comment IDs
	userIDs := make(map[int]struct{})
	commentIDs := make([]int, len(comments))
	for i, comment := range comments {
		userIDs[comment.AuthorID] = struct{}{}
		commentIDs[i] = comment.ID
	}

	// Maps for fetched data
	var (
		users          map[int]domain.User
		totalReactions map[int]int
		// reactions      map[int][]domain.Reaction
	)

	var eg errgroup.Group

	// Fetch user details
	eg.Go(func() error {
		var err error
		users, err = s.commentFetcher.FetchUsersByID(ctx, getKeys(userIDs))
		return err
	})

	// // Fetch reactions
	// eg.Go(func() error {
	// 	var err error
	// 	reactions, err = s.commentFetcher.FetchReactions(ctx, commentIDs, userID)
	// 	return err
	// })

	// Fetch total reactions
	eg.Go(func() error {
		var err error
		totalReactions, err = s.commentFetcher.FetchTotalReactions(ctx, commentIDs)
		return err
	})

	// Wait for all goroutines to complete
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	// Map fetched data to comments
	for i := range comments {
		comment := &comments[i]
		if user, exists := users[comment.AuthorID]; exists {
			comment.Username = user.Username
			comment.ProfilePic = user.ProfilePic
		}
		// if reactionTypes, exists := reactions[comment.ID]; exists {
		// 	comment.Reactions = reactios
		// }
		if count, exists := totalReactions[comment.ID]; exists {
			comment.TotaReactionslCount = count
		}
	}

	return comments, nil
}

func (s *CommentService) GetCommentsAndRepliesCount(ctx context.Context, entityIDs []int) ([]domain.CommentCount, error) {
	return s.commentRepo.CountByEntityIDs(ctx, entityIDs)
}

// Helper function to extract keys from a map
func getKeys(m map[int]struct{}) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
