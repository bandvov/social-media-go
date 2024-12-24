package interfaces

import (
	"encoding/json"
	"net/http"

	"github.com/bandvov/social-media-go/application"
	"github.com/bandvov/social-media-go/domain"
)

type TagHandler struct {
	TagService application.TagServiceInterface
}

// NewTagHandler creates a new HTTP handler for tags.
func NewTagHandler(service application.TagServiceInterface) *TagHandler {
	return &TagHandler{TagService: service}
}

// CreateTag handles creating a new tag.
func (h *TagHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	var req domain.Tag

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.TagService.DeleteTag(req.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "tag created successfully"})
}

// GetTags handles retrieving all tags.
func (h *TagHandler) GetTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.TagService.GetAllTags()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tags)
}

// GetTags handles retrieving all tags.
func (h *TagHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "tag deleted successfully"})
}
