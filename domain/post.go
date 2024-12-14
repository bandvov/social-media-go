package domain

import "time"

type CreatePostRequest struct {
	AuthorID   int       `json:"author_id"` // ID of the user who created the post
	Content    string    `json:"content"`
	Pinned     bool      `json:"pinned,omitempty"`
	Tags       string    `json:"tags,omitempty"`
	Visibility string    `json:"visibility,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

type Post struct {
	ID           string    `json:"id,omitempty"`
	AuthorID     int       `json:"author_id"` // ID of the user who created the post
	CommentCount int       `json:"comment_count,omitempty"`
	Content      string    `json:"content,omitempty"`
	LikeCount    int       `json:"like_count,omitempty"`
	Pinned       bool      `json:"pinned,omitempty"`
	ShareCount   int       `json:"share_count,omitempty"`
	Tags         string    `json:"tags,omitempty"`
	Visibility   string    `json:"visibility,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}
