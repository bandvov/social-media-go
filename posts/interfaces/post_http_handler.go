package interfaces

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"posts/application"
	"posts/domain"
	"posts/internal"
	"strconv"
	"time"
)

// Define keys for context
type contextKey string

const (
	userIDKey  contextKey = "userID"
	isAdminKey contextKey = "isAdmin"
)

type PostHTTPHandler struct {
	postService application.PostServiceInterface
	db          *sql.DB
	rdb         internal.RedisClient
}

func NewPostHTTPHandler(
	postService application.PostServiceInterface,
	rdb internal.RedisClient,

) *PostHTTPHandler {
	return &PostHTTPHandler{
		postService: postService,
		rdb:         rdb,
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	post, err := p.postService.GetPostByID(ctx, postID)
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
	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Move the logic to the service layer
	posts, postsCount, err := h.postService.GetPostsByUser(authorIDFromUrl, offset, limit)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"data":    posts,
		"hasMore": postsCount > offset+limit,
	}
	fmt.Println(time.Since(s))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *PostHTTPHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		slog.Warn("PostgreSQL health check failed", "error", err)
		http.Error(w, "Unhealthy", http.StatusServiceUnavailable)
		return
	}
	if err := h.rdb.Ping(ctx).Err(); err != nil {
		slog.Warn("Redis health check failed", "error", err)
		http.Error(w, "Unhealthy", http.StatusServiceUnavailable)
		return
	}

	slog.Info("Health check passed")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Ok")
}
