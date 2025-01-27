package application

import (
	"context"

	"github.com/bandvov/social-media-go/domain"
	"golang.org/x/sync/errgroup"
)

type PostServiceInterface interface {
	CreatePost(post *domain.CreatePostRequest) error
	DeletePost(id int) error
	UpdatePost(id int, post *domain.Post) error
	GetPostByID(id int) (*domain.Post, error)
	GetPostsByUser(userID, otherUserId, offset, limit int) ([]domain.Post, error)
	GetCountPostsByUser(userID int) (int, error)
}

type PostService struct {
	reactionRepo domain.ReactionRepository
	postRepo     domain.PostRepository
	commentRepo  domain.CommentRepository
	userRepo     domain.UserRepository
}

func NewPostService(repo domain.PostRepository) *PostService {
	return &PostService{postRepo: repo}
}

func (s *PostService) CreatePost(post *domain.CreatePostRequest) error {
	return s.postRepo.Create(post)
}

func (s *PostService) DeletePost(id int) error {
	return s.postRepo.Delete(id)
}

func (s *PostService) UpdatePost(id int, post *domain.Post) error {
	return s.postRepo.Update(id, post)
}

func (s *PostService) GetPostByID(id int) (*domain.Post, error) {
	return s.postRepo.GetByID(id)
}

func (s *PostService) GetPostsByUser(authorID, otherUserId, offset, limit int) ([]domain.Post, error) {
	// Fetch posts
	posts, err := s.postRepo.GetPosts(authorID, offset, limit)
	if err != nil {
		return nil, err
	}

	if len(posts) == 0 {
		return posts, nil
	}

	// Extract post IDs for batch operations
	postIDs := make([]int, len(posts))
	for i, post := range posts {
		postIDs[i] = post.ID
	}

	// Create error group for concurrent fetching
	var eg errgroup.Group
	var userIDList []int
	reactionMap := make(map[int][]domain.Reaction)
	commentMap := make(map[int][]domain.Comment)
	commentAuthorMap := make(map[int]domain.User)

	// Fetch reactions concurrently
	eg.Go(func() error {
		reactions, err := s.reactionRepo.GetReactionsByPostIDs(postIDs)
		if err != nil {
			return err
		}
		// Map reactions to post IDs
		for _, reaction := range reactions {
			reactionMap[reaction.EntityId] = append(reactionMap[reaction.EntityId], reaction)
		}
		return nil
	})

	// Fetch comments
	eg.Go(func() error {
		comments, err := s.commentRepo.GetCommentsByPostIDs(postIDs)
		if err != nil {
			return err
		}

		for _, comment := range comments {
			commentMap[comment.EntityID] = append(commentMap[comment.EntityID], comment)
			userIDList = append(userIDList, comment.AuthorID)

		}

		return err
	})

	// Wait for both operations to complete
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	userDetails, err := s.userRepo.GetUsers(context.Background(), userIDList)
	if err != nil {
		return nil, err
	}

	// Create a map of user details for quick lookup
	for _, user := range userDetails {
		commentAuthorMap[user.ID] = user
	}

	// Populate posts with comments and reactions
	for i, post := range posts {
		posts[i].Reactions = reactionMap[post.ID]
		comments := commentMap[post.ID]
		for j, comment := range comments {
			if user, exists := commentAuthorMap[comment.AuthorID]; exists {
				comments[j].ProfilePic = *user.ProfilePic
				comments[j].Username = *user.Username
			}
		}
		posts[i].Comments = comments
	}

	return posts, nil
}

func (s *PostService) GetCountPostsByUser(userID int) (int, error) {
	return s.postRepo.GetCountPostsByUser(userID)
}
