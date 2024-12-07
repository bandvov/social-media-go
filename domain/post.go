package domain

import "time"

type Post struct {
	ID           string `json:"id"`
	AuthorID     string `json:"author_id"` // ID of the user who created the post
	Content      string `json:"content"`
	LikeCount    int    `json:"like_count"`
	CommentCount int    `json:"comment_count"`
	ShareCount   int    `json:"share_count"`
	Visibility   string `json:"visibility"`
	Tags         []string
	Pinned       bool `json:"pinned"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	// Status       string // "draft", "published", etc.
}
