package application

import (
	"comments/domain"
	"context"
	"fmt"
	"os"

	"golang.org/x/sync/errgroup"
)

// CommentServiceInterface defines methods for tags-related operations.
type CommentServiceInterface interface {
	AddComment(ctx context.Context, c *domain.Comment) error
	GetCommentsByEntityID(ctx context.Context, entityId int, targetUserId int, pagination domain.Pagination) ([]domain.Comment, error)
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
func (s *CommentService) GetCommentsByEntityID(ctx context.Context, entityID int, userID int, p domain.Pagination) ([]domain.Comment, error) {
	comments, err := s.commentRepo.FetchCommentsByEntityID(ctx, entityID, p)
	if err != nil {
		return nil, err
	}

	if len(comments) == 0 {
		return comments, nil
	}

	// Collect user IDs and comment IDs
	userIDs := make(map[int]struct{})
	userIDs[userID] = struct{}{}
	entities := make([]domain.Entity, 0, len(comments))
	for _, comment := range comments {
		userIDs[comment.AuthorID] = struct{}{}
		entities = append(entities, domain.Entity{
			ID:   comment.ID,
			Type: string(comment.EntityType),
		})
	}

	// Maps for fetched data
	var (
		users         map[int]domain.User
		reactionStats map[int]domain.ReactionStat
	)

	var eg errgroup.Group

	// Fetch user details
	eg.Go(func() error {
		var err error
		users, err = s.commentFetcher.FetchUsersByID(ctx, getKeys(userIDs))
		return err
	})

	// Fetch total reactions
	eg.Go(func() error {
		var err error
		reactionStats, err = s.commentFetcher.FetchReactionStats(ctx, entities)
		return err
	})

	// Wait for all goroutines to complete
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	// Map fetched data to comments
	for i := range comments {
		comment := &comments[i]
		fmt.Fprintln(os.Stdout, "users", users)
		if user, exists := users[comment.AuthorID]; exists {
			comment.Username = user.Username
			comment.ProfilePic = user.ProfilePic
		}
		if reactionStat, exists := reactionStats[comment.ID]; exists {
			comment.Reactions = reactionStat.Reactions
		}
		if reactionStat, exists := reactionStats[comment.ID]; exists {
			comment.TotaReactionslCount = reactionStat.TotalCount
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
