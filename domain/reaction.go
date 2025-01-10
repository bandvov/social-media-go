package domain

type Reaction struct {
	ReactionType string `json:"reaction_type"`
	Count int `json:"count"`
}