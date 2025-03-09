package domain

import (
	"time"
)

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type Reaction struct {
	EntityId int    `json:"entity_id"`
	Reaction string `json:"reaction_type_id"`
	Count    int    `json:"count,omitempty"`
	Type     string `json:"type",omitempty`
}

type Comment struct {
	EntityID     int `json:"entity_id"`
	CommentCount int `json:"comment_count"`
	ReplyCount   int `json:"reply_count"`
}

type CreatePostRequest struct {
	AuthorID   int            `json:"author_id,omitempty"` // ID of the user who created the post
	Content    string         `json:"content,omitempty"`
	Pinned     bool           `json:"pinned,omitempty"`
	Tags       string         `json:"tags,omitempty"`
	Visibility PostVisibility `json:"visibility,omitempty"`
}

type Post struct {
	ID                  int             `json:"id,omitempty"`
	AuthorID            int             `json:"author_id,omitempty"` // ID of the user who created the post
	Content             string          `json:"content,omitempty"`
	AuthorName          string          `json:"author_name,omitempty"`
	Pinned              bool            `json:"pinned,omitempty"`
	Tags                string          `json:"tags,omitempty"`
	Visibility          *PostVisibility `json:"visibility,omitempty"`
	Reactions           []Reaction      `json:"reactions,omitempty"`
	TotalReactionsCount int             `json:"total_reactions_count,omitempty"`
	TotalCommentsCount  int             `json:"total_comments_count,omitempty"`
	UserReaction        string          `json:"user_reaction,omitempty"`
	CreatedAt           time.Time       `json:"created_at,omitempty"`
	UpdatedAt           time.Time       `json:"updated_at,omitempty"`
}

// PostVisibility represents the visibility of a post
type PostVisibility int

const (
	// Visibility constants
	Public    PostVisibility = iota // Public visibility
	Private                         // Private visibility
	Unlisted                        // Unlisted visibility
	Followers                       // Followers visibility
	Hidden                          // Hidden visibility
)

func (v PostVisibility) String() string {
	switch v {
	case Public:
		return "Public"
	case Private:
		return "Private"
	case Followers:
		return "Followers"
	case Unlisted:
		return "Unlisted"
	default:
		return "Unknown"
	}
}

type PostSearchOptions struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Sort   string `json:"sort"`
	Search string `json:"search"`
}
type Entity struct {
	ID     int    `json:"id,omitempty"`
	userID int    `json:"user_id,omitempty"`
	Type   string `json:"type,omitempty"`
}

type Request[T any] struct {
	Data T `json:"data"`
}

type Response[T any] struct {
	Data    T      `json:"data"`
	Message string `json:"message"`
}
