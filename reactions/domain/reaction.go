package domain

import "context"

type Reaction struct {
	EntityId int    `json:"entity_id"`
	Reaction string `json:"reaction_type_id"`
	Count    int    `json:"count"`
}

type ReactionRepository interface {
	AddOrUpdateReaction(ctx context.Context, userId int, reaction Reaction) error
	RemoveReaction(ctx context.Context, userID, contentID string) error
	GetReactionsByEntityIDs(ctx context.Context, entityIDs []int) ([]Reaction, error)
	CountByEntityIDs(ctx context.Context, entityIDs []int) ([]Reaction, error)
}
