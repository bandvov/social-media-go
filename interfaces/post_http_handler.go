package interfaces

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bandvov/social-media-go/application"
	"github.com/bandvov/social-media-go/domain"
)

type PostHTTPHandler struct {
	PostService application.PostServiceInterface
}

func NewPostHTTPHandler(postService application.PostServiceInterface) *PostHTTPHandler {
	return &PostHTTPHandler{PostService: postService}
}

func (p *PostHTTPHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	authorID, ok := r.Context().Value(userIDKey).(interface{}).(int)
	if !ok || authorID == 0 {
		http.Error(w, "unauthenticated", http.StatusBadRequest)
		return
	}

	var newPost domain.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if newPost.Content == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	newPost.AuthorID = authorID

	err := p.PostService.CreatePost(&newPost)
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

	err = p.PostService.UpdatePost(postID, &domain.Post{
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

	post, err := p.PostService.GetPostByID(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !isAdmin || (*post.Visibility == domain.Private && post.AuthorID != userID) {
		http.Error(w, "Access forbidden", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}

func (h *PostHTTPHandler) GetPostsByUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	posts, err := h.PostService.GetPostsByUser(userID)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No posts", http.StatusNotFound)
			return
		}
		http.Error(w, "could not retrieve posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
