package infrastructure

import (
	"database/sql"
	"fmt"

	"github.com/bandvov/social-media-go/domain"
	"github.com/bandvov/social-media-go/utils"
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

// Fetch reactions by post IDs
func (r *ReactionRepository) GetReactionsByEntityIDs(postIDs []int) ([]domain.Reaction, error) {
	if len(postIDs) == 0 {
		return nil, nil
	}

	query := fmt.Sprintf(`
        SELECT r.entity_id, rt.name AS reaction, COUNT(r.id) AS count
        FROM reactions r
        JOIN reaction_types rt ON r.reaction_type_id = rt.id
        WHERE r.post_id IN (%s)
        GROUP BY r.post_id, rt.name`, utils.Placeholders(len(postIDs)))

	rows, err := r.db.Query(query, utils.ToInterface(postIDs)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reactions []domain.Reaction
	for rows.Next() {
		var reaction domain.Reaction
		if err := rows.Scan(&reaction.EntityId, &reaction.Reaction, &reaction.Count); err != nil {
			return nil, err
		}
		reactions = append(reactions, reaction)
	}

	return reactions, nil
}
