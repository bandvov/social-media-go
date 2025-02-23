package application

import (
	"bytes"
	"encoding/json"
	"net/http"
	"posts/domain"
	"posts/internal"

	"golang.org/x/sync/errgroup"
)

type PostServiceInterface interface {
	CreatePost(post *domain.CreatePostRequest) error
	DeletePost(id int) error
	UpdatePost(id int, post *domain.Post) error
	GetPostByID(id int) (*domain.Post, error)
	GetPostsByUser(userID, offset, limit int) ([]domain.Post, int, error)
	GetCountPostsByUser(userID int) (int, error)
}

type PostService struct {
	postRepo        domain.PostRepository
	commentsClient  internal.ClientInterface
	reactionsClient internal.ClientInterface
}

// GetPostsByUser implements PostServiceInterface.
func (s *PostService) GetPostsByUser(userID int, offset int, limit int) ([]domain.Post, int, error) {
	panic("unimplemented")
}

func NewPostService(
	repo domain.PostRepository,
	commentsClient internal.ClientInterface,
	reactionsClient internal.ClientInterface,
) *PostService {
	return &PostService{postRepo: repo,
		commentsClient:  commentsClient,
		reactionsClient: reactionsClient,
	}
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

func (s *PostService) GetCountPostsByUser(userID int) (int, error) {
	return s.postRepo.GetCountPostsByUser(userID)
}

func (s *PostService) FetchPostsByUser(authorID int, offset, limit int) ([]domain.Post, int, error) {
	posts, err := s.postRepo.GetByUserID(authorID, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	var postIDs []int
	for _, v := range posts {
		postIDs = append(postIDs, v.ID)
	}

	var (
		reactionMap map[int][]domain.Reaction
		postsCount  int
		eg          errgroup.Group
	)
	commentsCountsMap := make(map[int]domain.Comment)
	reactionsCountsMap := make(map[int]domain.Reaction)

	eg.Go(func() error {
		// figure out if it is better to move it into outer scope
		jsonData, err := json.Marshal(postIDs)
		if err != nil {
			return err
		}

		req, err := http.NewRequest("GET", "/count", bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}

		var counts []domain.Comment
		err = s.commentsClient.GetJSON(req, counts)
		if err != nil {
			return err
		}

		for _, count := range counts {
			commentsCountsMap[count.EntityID] = count
		}
		return nil
	})

	eg.Go(func() error {
		jsonData, err := json.Marshal(postIDs)
		if err != nil {
			return err
		}

		req, err := http.NewRequest("GET", "/counts", bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}

		var counts []domain.Reaction
		err = s.reactionsClient.GetJSON(req, counts)
		if err != nil {
			return err
		}

		for _, count := range counts {
			reactionsCountsMap[count.EntityId] = count
		}
		return nil
	})

	eg.Go(func() error {
		jsonData, err := json.Marshal(postIDs)
		if err != nil {
			return err
		}
		req, err := http.NewRequest("GET", "/", bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}
		err = s.reactionsClient.GetJSON(req, reactionMap)
		return err
	})

	// Second task: Fetch posts count
	eg.Go(func() error {
		var err error
		postsCount, err = s.GetCountPostsByUser(authorID)
		if err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, 0, err
	}

	for i, post := range posts {
		posts[i].Reactions = reactionMap[post.ID]
		posts[i].TotalCommentsCount = commentsCountsMap[post.ID].CommentCount
		posts[i].TotalReactionsCount = reactionsCountsMap[post.ID].Count
	}

	return posts, postsCount, nil
}
