package domain

type Reaction struct {
	EntityId string `json:"entity_id"`
	Reaction string `json:"reaction_type_id"`
}

type ReactionRepository interface {
	AddOrUpdateReaction(userId int, reaction Reaction) error
	RemoveReaction(userID, contentID string) error
}
