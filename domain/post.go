package domain

import "time"

type CreatePostRequest struct {
	AuthorID   int            `json:"author_id,omitempty"` // ID of the user who created the post
	Content    string         `json:"content,omitempty"`
	Pinned     bool           `json:"pinned,omitempty"`
	Tags       string         `json:"tags,omitempty"`
	Visibility PostVisibility `json:"visibility,omitempty"`
	CreatedAt  time.Time      `json:"created_at,omitempty"`
	UpdatedAt  time.Time      `json:"updated_at,omitempty"`
}

type Post struct {
	ID         string         `json:"id,omitempty"`
	AuthorID   int            `json:"author_id,omitempty"` // ID of the user who created the post
	Content    string         `json:"content,omitempty"`
	Pinned     bool           `json:"pinned,omitempty"`
	Tags       string         `json:"tags,omitempty"`
	Visibility PostVisibility `json:"visibility,omitempty"`
	CreatedAt  time.Time      `json:"created_at,omitempty"`
	UpdatedAt  time.Time      `json:"updated_at,omitempty"`
}

// PostVisibility represents the visibility of a post
type PostVisibility int

const (
	// Visibility constants
	Public   PostVisibility = iota // Public visibility
	Private                        // Private visibility
	Unlisted                       // Unlisted visibility
)

func (v PostVisibility) String() string {
	switch v {
	case Public:
		return "Public"
	case Private:
		return "Private"
	case Unlisted:
		return "Unlisted"
	default:
		return "Unknown"
	}
}
