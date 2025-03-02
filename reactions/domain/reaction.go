package domain

import "context"

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
	CountByEntityIDs(ctx context.Context, entityIDs []int) ([]Reaction, error)
}
