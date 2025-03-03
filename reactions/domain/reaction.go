package domain

import (
	"context"
)

type Entity struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}

type Reaction struct {
	EntityId   int    `json:"entity_id"`
	EntityType string `json:"entity_type,omitempty"`
	Reaction   string `json:"reaction_type_id,omitempty"`
	Count      int    `json:"count"`
}

type ReactionRepository interface {
	AddOrUpdateReaction(ctx context.Context, userId int, reaction Reaction) error
	RemoveReaction(ctx context.Context, userID, contentID string) error
	GetReactionsByEntityIDs(ctx context.Context, entityIDs []int) ([]Reaction, error)
	CountByEntityIDsAndType(ctx context.Context, entities []Entity) ([]Reaction, error)
}
