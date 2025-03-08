package domain

import (
	"context"
	"encoding/json"
)

type Entity struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}

type ReactionStat struct {
	EntityId   int             `json:"entity_id"`
	EntityType string          `json:"entity_type,omitempty"`
	Reactions  json.RawMessage `json:"reactions,omitempty"`
	TotalCount int             `json:"total_count,omitempty"`
}

type Reaction struct {
	UserId     int    `json:"user_id"`
	EntityId   int    `json:"entity_id"`
	EntityType string `json:"entity_type,omitempty"`
	Reaction   string `json:"reaction_type_id,omitempty"`
	Count      int    `json:"count,omitempty"`
	Type       string `json:"type,omitempty"`
}

type ReactionRepository interface {
	AddOrUpdateReaction(ctx context.Context, userId int, reaction Reaction) error
	RemoveReaction(ctx context.Context, userID, contentID string) error
	GetReactionsByEntityIDsAndTypes(ctx context.Context, entities []Entity) ([]Reaction, error)
	CountByEntityIDsAndTypes(ctx context.Context, entities []Entity) ([]Reaction, error)
	GetReacionStats(ctx context.Context, entities []Entity) ([]ReactionStat, error)
	GetUserReactions(ctx context.Context, userID int, entities []Entity) ([]Reaction, error)
}
