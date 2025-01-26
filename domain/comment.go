package domain

import (
	"encoding/json"
	"time"
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
	EntityType          string          `json:"entity_type,omitempty"`
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
