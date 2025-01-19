package application

import "github.com/bandvov/social-media-go/domain"

type ReactionServiceInterface interface {
	AddOrUpdateReaction(userID int, reaction domain.Reaction) error
	RemoveReaction(userID, contentID string) error
}
type ReactionService struct {
	repo domain.ReactionRepository
}

func NewReactionService(repo domain.ReactionRepository) *ReactionService {
	return &ReactionService{repo: repo}
}

func (s *ReactionService) AddOrUpdateReaction(userID int, reaction domain.Reaction) error {
	return s.repo.AddOrUpdateReaction(userID, reaction)
}

func (s *ReactionService) RemoveReaction(userID, contentID string) error {
	return s.repo.RemoveReaction(userID, contentID)
}
