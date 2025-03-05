package domain

import (
	"encoding/json"
	"time"
)

type CommentType string

const (
	CommentTypeComment CommentType = "comment"
	CommentTypeReply   CommentType = "reply"
)

// PostVisibility represents the visibility of a post
type CommentStatus int

const (
	// Visibility constants
	Active  CommentStatus = iota // status = "active"
	Flagged                      // status = "flagged"
	Deleted                      // status = "deleted"
)

func (c CommentStatus) String() string {
	switch c {
	case Active:
		return "active"
	case Flagged:
		return "flagged"
	case Deleted:
		return "deleted"
	default:
		return "Unknown"
	}
}

type Comment struct {
	ID                  int             `json:"id,omitempty"`
	EntityID            int             `json:"entity_id,omitempty"`
	EntityType          CommentType     `json:"entity_type,omitempty"`
	Content             string          `json:"content,omitempty"`
	AuthorID            int             `json:"author_id,omitempty"`
	Username            string          `json:"username,omitempty"`
	ProfilePic          string          `json:"profile_pic,omitempty"`
	Status              CommentStatus   `json:"status,omitempty"`
	RepliesCount        int             `json:"replies_count,omitempty"`
	Reactions           json.RawMessage `json:"reactions,omitempty"`
	TotaReactionslCount int             `json:"total_reactions_count,omitempty"`
	UserReaction        string          `json:"user_reaction,omitempty"`
	CreatedAt           time.Time       `json:"created_at,omitempty"`
	UpdatedAt           time.Time       `json:"updated_at,omitempty"`
}

func (c *Comment) IsValidEntityId() bool {
	return c.EntityID > 0
}

func (c *Comment) IsValidAuthorId() bool {
	return c.AuthorID > 0

}

func (c *Comment) IsValidContent() bool {
	return c.Content != ""
}

type User struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	ProfilePic string `json:"profile_pic"`
}

type Reaction struct {
	EntityID int    `json:"entity_id"`
	Reaction string `json:"reaction_type"`
	Count    int    `json:"count"`
}

type CommentCount struct {
	EntityID     int `json:"entity_id"`
	CommentCount int `json:"comment_count,omitempty"`
	ReplyCount   int `json:"reply_count,omitempty"`
}
type Entity struct {
	ID   int    `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

type ReactionStat struct {
	EntityId   int             `json:"entity_id"`
	EntityType string          `json:"entity_type,omitempty"`
	Reactions  json.RawMessage `json:"reactions,omitempty"`
	TotalCount int             `json:"total_count,omitempty"`
}

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
