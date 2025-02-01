package interfaces

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bandvov/social-media-go/application"
	"github.com/bandvov/social-media-go/domain"
	"golang.org/x/sync/errgroup"
)

type PostHTTPHandler struct {
	postService     application.PostServiceInterface
	commentService  application.CommentServiceInterface
	userService     application.UserServiceInterface
	reactionService application.ReactionServiceInterface
}

func NewPostHTTPHandler(
	postService application.PostServiceInterface,
	commentService application.CommentServiceInterface,
	userService application.UserServiceInterface,
	reactionService application.ReactionServiceInterface,
) *PostHTTPHandler {
	return &PostHTTPHandler{
		postService:     postService,
		commentService:  commentService,
		userService:     userService,
		reactionService: reactionService}
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
	userID, ok := r.Context().Value(userIDKey).(interface{}).(int)
	if !ok || userID == 0 {
		http.Error(w, "unauthenticated", http.StatusBadRequest)
		return
	}

	isAdmin := r.Context().Value(isAdminKey).(bool)

	id := r.PathValue("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := p.postService.GetPostByID(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("here")
	if !isAdmin || (*post.Visibility == domain.Private && post.AuthorID != userID) {
		fmt.Println("here1")
		http.Error(w, "Access forbidden", http.StatusForbidden)
		return
	}
	fmt.Println("here2")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}

// func (h *PostHTTPHandler) GetPostsByUser(w http.ResponseWriter, r *http.Request) {
// 	userID, ok := r.Context().Value(userIDKey).(interface{}).(int)
// 	if !ok || userID == 0 {
// 		http.Error(w, "unauthenticated", http.StatusBadRequest)
// 		return
// 	}

// 	idStr := r.PathValue("id")
// 	userIDFromUrl, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "invalid user ID", http.StatusBadRequest)
// 		return
// 	}

// 	query := r.URL.Query()

// 	// Parse `limit` and `offset` with default values
// 	page, err := strconv.Atoi(query.Get("page"))
// 	if err != nil || page < 1 {
// 		page = 1 // Default offset
// 	}

// 	limit, err := strconv.Atoi(query.Get("limit"))
// 	if err != nil || limit <= 0 {
// 		limit = 10 // Default limit
// 	}

// 	offset := (page - 1) * limit

// 	var posts []domain.Post
// 	var postsCount int

// 	// Create a new errgroup
// 	var g errgroup.Group
// 	// First task: Fetch posts
// 	g.Go(func() error {
// 		var err error
// 		posts, err = h.postService.GetPostsByUser(userIDFromUrl, userID, offset, limit)
// 		if err != nil {
// 			if errors.Is(err, sql.ErrNoRows) {
// 				http.Error(w, "No posts", http.StatusNotFound)
// 				return nil // No posts is not an error; early return
// 			}
// 			return err // Return other errors
// 		}
// 		return nil
// 	})

// 	// Second task: Fetch posts count
// 	g.Go(func() error {
// 		var err error
// 		postsCount, err = h.postService.GetCountPostsByUser(userIDFromUrl)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})

// 	// Wait for both goroutines to finish
// 	if err := g.Wait(); err != nil {
// 		http.Error(w, "could not complete request", http.StatusInternalServerError)
// 		return
// 	}

// 	response := map[string]interface{}{
// 		"data":    posts,
// 		"hasMore": postsCount > offset+limit,
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

func (h *PostHTTPHandler) GetPostsByUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(interface{}).(int)
	if !ok || userID == 0 {
		http.Error(w, "unauthenticated", http.StatusBadRequest)
		return
	}

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
	countsMap := make(map[int]domain.CommentCount)

	eg.Go(func() error {
		counts, err := h.commentService.GetCommentsAndRepliesCount(postIDs)
		for _, count := range counts {
			countsMap[count.EntityID] = count
		}
		return err
	})

	eg.Go(func() error {
		reactionMap, err = h.reactionService.GetReactions(postIDs)
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
		posts[i].TotalCommentsCount = countsMap[post.ID].CommentCount + countsMap[post.ID].CommentCount
	}

	response := map[string]interface{}{
		"data":    posts,
		"hasMore": postsCount > offset+limit,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
