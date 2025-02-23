package interfaces

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"posts/application"
	"posts/domain"
	"posts/internal"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"
)

// Define keys for context
type contextKey string

const (
	userIDKey  contextKey = "userID"
	isAdminKey contextKey = "isAdmin"
)

type PostHTTPHandler struct {
	postService     application.PostServiceInterface
	commentsClient  internal.ClientInterface
	reactionsClient internal.ClientInterface
}

func NewPostHTTPHandler(
	postService application.PostServiceInterface,
	commentsClient internal.ClientInterface,
	reactionsClient internal.ClientInterface,

) *PostHTTPHandler {
	return &PostHTTPHandler{
		postService:     postService,
		commentsClient:  commentsClient,
		reactionsClient: reactionsClient,
	}
}

func (p *PostHTTPHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	authorID, ok := r.Context().Value(userIDKey).(interface{}).(int)
	if !ok || authorID == 0 {
		http.Error(w, "unauthenticated", http.StatusBadRequest)
		return
	}

	var newPost struct {
		Data domain.CreatePostRequest `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if newPost.Data.Content == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	newPost.Data.AuthorID = authorID

	err := p.postService.CreatePost(&newPost.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post created successfully"})
}

func (p *PostHTTPHandler) DeletePost(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode(map[string]string{"message": "post deleted successfully"})
}

func (p *PostHTTPHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	var post *domain.CreatePostRequest

	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = p.postService.UpdatePost(postID, &domain.Post{
		Content: post.Content, Visibility: &post.Visibility, Tags: post.Tags, Pinned: post.Pinned,
	})

	if err != nil {
		http.Error(w, "error updating post: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "post updated successfully"})
}

func (p *PostHTTPHandler) GetPost(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := p.postService.GetPostByID(postID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}

func (h *PostHTTPHandler) GetPostsByUser(w http.ResponseWriter, r *http.Request) {

	s := time.Now()
	idStr := r.PathValue("id")
	authorIDFromUrl, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	query := r.URL.Query()

	// Parse `limit` and `offset` with default values
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page < 1 {
		page = 1 // Default offset
	}

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}

	offset := (page - 1) * limit

	posts, postIDs, err := h.postService.GetPostsByUser(authorIDFromUrl, offset, limit)
	if err != nil || len(posts) == 0 {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
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
		err = h.commentsClient.GetJSON(req, counts)
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
		err = h.reactionsClient.GetJSON(req, counts)
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
		err = h.reactionsClient.GetJSON(req, reactionMap)
		return err
	})

	// Second task: Fetch posts count
	eg.Go(func() error {
		var err error
		postsCount, err = h.postService.GetCountPostsByUser(authorIDFromUrl)
		if err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to fetch posts", http.StatusBadRequest)
		return
	}

	for i, post := range posts {
		posts[i].Reactions = reactionMap[post.ID]
		posts[i].TotalCommentsCount = commentsCountsMap[post.ID].CommentCount + commentsCountsMap[post.ID].CommentCount
		posts[i].TotalReactionsCount = reactionsCountsMap[post.ID].Count
	}

	response := map[string]interface{}{
		"data":    posts,
		"hasMore": postsCount > offset+limit,
	}
	fmt.Println(time.Since(s))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
