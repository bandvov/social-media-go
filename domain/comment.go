package domain

import "time"

type Comment struct {
	ID        int       `json:"id,omitempty"`
	EntityID  int       `json:"entity_id,omitempty"`
	Content   string    `json:"content,omitempty"`
	AuthorID  int       `json:"author_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
