package application

import (
	"context"
	"fmt"
	"os"
	"posts/domain"

	"golang.org/x/sync/errgroup"
)

type PostServiceInterface interface {
	CreatePost(ctx context.Context, post *domain.CreatePostRequest) error
	DeletePost(ctx context.Context, id int) error
	UpdatePost(ctx context.Context, id int, post *domain.Post) error
	GetPostByID(ctx context.Context, id, targetUserId int) (*domain.Post, error)
	GetPostsByUser(ctx context.Context, userID, targetUserId int, p domain.Pagination) ([]domain.Post, int, error)
	GetCountPostsByUser(ctx context.Context, userID int) (int, error)
}

type PostService struct {
	postRepo domain.PostRepository
	fetcher  PostsFetcher
}

func NewPostService(
	repo domain.PostRepository,
	fetcher PostsFetcher,
) *PostService {
	return &PostService{
		postRepo: repo,
		fetcher:  fetcher,
	}
}

func (s *PostService) CreatePost(ctx context.Context, post *domain.CreatePostRequest) error {
	return s.postRepo.Create(ctx, post)
}

func (s *PostService) DeletePost(ctx context.Context, id int) error {
	return s.postRepo.Delete(ctx, id)
}

func (s *PostService) UpdatePost(ctx context.Context, id int, post *domain.Post) error {
	return s.postRepo.Update(ctx, id, post)
}

func (s *PostService) GetPostByID(ctx context.Context, id, targetUserId int) (*domain.Post, error) {
	post, err := s.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var (
		eg               errgroup.Group
		reactionStatsMap map[int]domain.ReactionStat
		reactionsMap     map[int]domain.Reaction
	)

	// Fetch user reactions
	eg.Go(func() error {
		var err error
		reactionStatsMap, err = s.fetcher.FetchReactionStats(ctx, []domain.Entity{{ID: id, Type: "post"}})
		return err
	})

	eg.Go(func() error {
		var err error
		reactionsMap, err = s.fetcher.FetchUserReactions(ctx, targetUserId, []domain.Entity{{ID: post.ID, Type: "post"}})
		return err
	})

	// Wait for all goroutines to complete
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	post.UserReaction = reactionsMap[post.ID].Type
	post.Reactions = reactionStatsMap[post.ID].Reactions
	return post, nil
}

func (s *PostService) GetCountPostsByUser(ctx context.Context, userID int) (int, error) {
	return s.postRepo.GetCountPostsByUser(ctx, userID)
}

func (s *PostService) GetPostsByUser(ctx context.Context, authorID, targetUserId int, p domain.Pagination) ([]domain.Post, int, error) {
	fmt.Fprintln(os.Stdout, "GetPostsByUser in service")
	posts, err := s.postRepo.GetByUserID(ctx, authorID, p)
	if err != nil {
		return nil, 0, err
	}
	var postIDs []int
	var entities []domain.Entity
	for _, v := range posts {
		postIDs = append(postIDs, v.ID)
		entities = append(entities, domain.Entity{ID: v.ID, Type: "post"})
	}

	var (
		reactionStatsMap map[int]domain.ReactionStat
		postsCount       int
		eg               errgroup.Group
	)
	commentsCountsMap := make(map[int]domain.Comment)

	eg.Go(func() error {
		var err error
		commentsCountsMap, err = s.fetcher.FetchCommentsCount(ctx, postIDs)
		return err
	})

	eg.Go(func() error {
		var err error
		reactionStatsMap, err = s.fetcher.FetchReactionStats(ctx, entities)
		return err
	})

	// Second task: Fetch posts count
	eg.Go(func() error {
		var err error
		postsCount, err = s.GetCountPostsByUser(ctx, authorID)
		if err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, 0, err
	}

	for i, post := range posts {
		posts[i].Reactions = reactionStatsMap[post.ID].Reactions
		posts[i].TotalCommentsCount = commentsCountsMap[post.ID].CommentCount
		posts[i].TotalReactionsCount = reactionStatsMap[post.ID].TotalCount
	}

	return posts, postsCount, nil
}
