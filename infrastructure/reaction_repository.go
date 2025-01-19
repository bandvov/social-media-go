package infrastructure

import (
	"database/sql"

	"github.com/bandvov/social-media-go/domain"
)

type ReactionRepository struct {
	db *sql.DB
}

func NewReactionRepository(db *sql.DB) *ReactionRepository {
	return &ReactionRepository{db: db}
}

func (r *ReactionRepository) AddOrUpdateReaction(userID int, reaction domain.Reaction) error {
	query := `
        INSERT INTO reactions (user_id, entity_id, reaction_type_id)
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id, entity_id)
        DO UPDATE SET reaction_type_id = $3
    `
	_, err := r.db.Exec(query, userID, reaction.EntityId, reaction.Reaction)
	return err
}

func (r *ReactionRepository) RemoveReaction(userID, entityID string) error {
	query := `DELETE FROM reactions WHERE user_id = $1 AND entity_id = $2`
	_, err := r.db.Exec(query, userID, entityID)
	return err
}
