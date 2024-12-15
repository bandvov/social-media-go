package interfaces

import (
	"encoding/json"
	"net/http"

	"github.com/bandvov/social-media-go/application"
	"github.com/bandvov/social-media-go/domain"
)

type PostHTTPHandler struct {
	PostService application.PostServiceInterface
}

func NewPostHTTPHandler(postService application.PostServiceInterface) *PostHTTPHandler {
	return &PostHTTPHandler{PostService: postService}
}

func (p *PostHTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	err := p.PostService.Create(&newPost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post created successfully"})
}

func (p *PostHTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode(map[string]string{"message": "post deleted successfully"})
}

func (p *PostHTTPHandler) Update(w http.ResponseWriter, r *http.Request) {

	var post *domain.CreatePostRequest

	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	err := p.PostService.Update(&domain.Post{
		Content: post.Content, Visibility: post.Visibility, Tags: post.Tags, Pinned: post.Pinned,
	})
	if err != nil {
		http.Error(w, "error updating post: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "post updated successfully"})
}
