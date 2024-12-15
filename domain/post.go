package domain

import "time"

type CreatePostRequest struct {
	AuthorID   int       `json:"author_id,omitempty"` // ID of the user who created the post
	Content    string    `json:"content,omitempty"`
	Pinned     bool      `json:"pinned,omitempty"`
	Tags       string    `json:"tags,omitempty"`
	Visibility string    `json:"visibility,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

type Post struct {
	ID         string    `json:"id,omitempty"`
	AuthorID   int       `json:"author_id,omitempty"` // ID of the user who created the post
	Content    string    `json:"content,omitempty"`
	Pinned     bool      `json:"pinned,omitempty"`
	Tags       string    `json:"tags,omitempty"`
	Visibility string    `json:"visibility,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}
