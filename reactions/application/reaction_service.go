package application

import (
	"context"
	"reactions/domain"
)

type ReactionServiceInterface interface {
	AddOrUpdateReaction(ctx context.Context, userID int, reaction domain.Reaction) error
	RemoveReaction(ctx context.Context, userID, contentID string) error
	GetReactions(ctx context.Context, entities []domain.Entity) ([]domain.Reaction, error)
	GetReactionsCount(ctx context.Context, entities []domain.Entity) ([]domain.Reaction, error)
	GetReactionStats(ctx context.Context, entities []domain.Entity) ([]domain.ReactionStat, error)
	GetUserReactions(ctx context.Context, userId int, entities []domain.Entity) ([]domain.Reaction, error)
}
type ReactionService struct {
	reactionRepo domain.ReactionRepository
}

func NewReactionService(reactionRepo domain.ReactionRepository) *ReactionService {
	return &ReactionService{reactionRepo: reactionRepo}
}

func (s *ReactionService) AddOrUpdateReaction(ctx context.Context, userID int, reaction domain.Reaction) error {
	return s.reactionRepo.AddOrUpdateReaction(ctx, userID, reaction)
}

func (s *ReactionService) RemoveReaction(ctx context.Context, userID, contentID string) error {
	return s.reactionRepo.RemoveReaction(ctx, userID, contentID)
}

func (s *ReactionService) GetReactions(ctx context.Context, entities []domain.Entity) ([]domain.Reaction, error) {
	return s.reactionRepo.GetReactionsByEntityIDsAndTypes(ctx, entities)
}

func (s *ReactionService) GetReactionsCount(ctx context.Context, entities []domain.Entity) ([]domain.Reaction, error) {
	return s.reactionRepo.CountByEntityIDsAndTypes(ctx, entities)
}

func (s *ReactionService) GetReactionStats(ctx context.Context, entities []domain.Entity) ([]domain.ReactionStat, error) {
	return s.reactionRepo.GetReacionStats(ctx, entities)
}

func (s *ReactionService) GetUserReactions(ctx context.Context, userID int, entities []domain.Entity) ([]domain.Reaction, error) {
	return s.reactionRepo.GetUserReactions(ctx, userID, entities)
}
