package domain

type Reaction struct {
	EntityId int    `json:"entity_id"`
	Reaction string `json:"reaction_type_id"`
	Count    int    `json:"count"`
}

type ReactionRepository interface {
	AddOrUpdateReaction(userId int, reaction Reaction) error
	RemoveReaction(userID, contentID string) error
	GetReactionsByEntityIDs(entityIDs []int) ([]Reaction, error)
	CountByEntityIDs(entityIDs []int) ([]Reaction, error)
}
