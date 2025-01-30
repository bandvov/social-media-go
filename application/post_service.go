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

func (s *PostService) GetPostsByUser(authorID, otherUserID, offset, limit int) ([]domain.Post, error) {
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

	// Create maps for data aggregation
	reactionMap := make(map[int][]domain.Reaction)
	commentMap := make(map[int][]domain.Comment)
	commentAuthorMap := make(map[int]domain.User)
	// Fetch comments and prepare auxiliary lists
	comments, err := s.commentRepo.GetCommentsByEntityIDs(postIDs)
	if err != nil {
		return nil, err
	}

	userIDList := make([]int, 0, len(comments))
	commentIDList := make([]int, 0, len(comments))
	for _, comment := range comments {
		commentMap[comment.EntityID] = append(commentMap[comment.EntityID], comment)
		userIDList = append(userIDList, comment.AuthorID)
		commentIDList = append(commentIDList, comment.EntityID)
	}

	// Create an error group for concurrent operations
	var eg errgroup.Group

	// Fetch reactions concurrently
	eg.Go(func() error {
		reactions, err := s.reactionRepo.GetReactionsByEntityIDs(append(postIDs, commentIDList...))
		if err != nil {
			return err
		}
		for _, reaction := range reactions {
			reactionMap[reaction.EntityId] = append(reactionMap[reaction.EntityId], reaction)
		}
		return nil
	})

	// Fetch user details concurrently
	eg.Go(func() error {
		userDetails, err := s.userRepo.GetUsersByID(context.Background(), userIDList)
		if err != nil {
			return err
		}
		for _, user := range userDetails {
			commentAuthorMap[user.ID] = user
		}
		return nil
	})

	// Wait for all concurrent tasks to complete
	if err := eg.Wait(); err != nil {
		return nil, err
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
