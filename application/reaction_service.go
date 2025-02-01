package application

import "github.com/bandvov/social-media-go/domain"

type ReactionServiceInterface interface {
	AddOrUpdateReaction(userID int, reaction domain.Reaction) error
	RemoveReaction(userID, contentID string) error
	GetReactions(entityIDs []int) (map[int][]domain.Reaction, error)
}
type ReactionService struct {
	reactionRepo domain.ReactionRepository
}

func NewReactionService(reactionRepo domain.ReactionRepository) *ReactionService {
	return &ReactionService{reactionRepo: reactionRepo}
}

func (s *ReactionService) AddOrUpdateReaction(userID int, reaction domain.Reaction) error {
	return s.reactionRepo.AddOrUpdateReaction(userID, reaction)
}

func (s *ReactionService) RemoveReaction(userID, contentID string) error {
	return s.reactionRepo.RemoveReaction(userID, contentID)
}

func (s *ReactionService) GetReactions(entityIDs []int) (map[int][]domain.Reaction, error) {
	reactionMap := make(map[int][]domain.Reaction)

	reactions, err := s.reactionRepo.GetReactionsByEntityIDs(entityIDs)
	if err != nil {
		return nil, err
	}
	for _, reaction := range reactions {
		reactionMap[reaction.EntityId] = append(reactionMap[reaction.EntityId], reaction)
	}

	return reactionMap, nil
}
