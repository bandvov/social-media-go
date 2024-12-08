package interfaces

import (
	"net/http"

	"github.com/bandvov/social-media-go/application"
)

type PostHTTPHandler struct {
	PostService application.PostServiceInterface
}

func NewPostHTTPHandler(postService application.PostServiceInterface) *PostHTTPHandler {
	return &PostHTTPHandler{PostService: postService}
}

func (p *PostHTTPHandler) Create(w http.ResponseWriter, r *http.Request) {}
