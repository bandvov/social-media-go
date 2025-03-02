package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"reactions/domain"
	"reactions/utils"
)

type ReactionRepository struct {
	db    *sql.DB
	cache Cache
}

func NewReactionRepository(db *sql.DB, cache Cache) *ReactionRepository {
	return &ReactionRepository{db: db, cache: cache}
}

func (r *ReactionRepository) AddOrUpdateReaction(ctx context.Context, userID int, reaction domain.Reaction) error {
	query := `
        INSERT INTO reactions (user_id, entity_id, reaction_type_id)
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id, entity_id)
        DO UPDATE SET reaction_type_id = EXCLUDED.reaction_type_id
    `

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, userID, reaction.EntityId, reaction.Reaction)
	return err
}
func (r *ReactionRepository) RemoveReaction(ctx context.Context, userID, entityID string) error {
	query := `DELETE FROM reactions WHERE user_id = $1 AND entity_id = $2`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, userID, entityID)
	return err
}

func (r *ReactionRepository) GetReactionsByEntityIDs(ctx context.Context, postIDs []int) ([]domain.Reaction, error) {
	if len(postIDs) == 0 {
		return nil, nil
	}

	query := fmt.Sprintf(`
        SELECT r.entity_id, rt.name AS reaction, COUNT(r.id) AS count
        FROM reactions r
        JOIN reaction_types rt ON r.reaction_type_id = rt.id
        WHERE r.entity_id IN (%s)
        GROUP BY r.entity_id, rt.name`, utils.Placeholders(len(postIDs)))

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, utils.ToInterface(postIDs)...)
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

	return reactions, rows.Err() // Ensure any potential iteration errors are returned
}

func (r *ReactionRepository) CountByEntityIDs(ctx context.Context, entityIDs []int) ([]domain.Reaction, error) {
	if len(entityIDs) == 0 {
		return nil, nil
	}

	query := fmt.Sprintf(`
        SELECT entity_id,entity_type, COUNT(*) AS count
		FROM reactions 
		WHERE entity_id IN (%s)
		GROUP BY entity_id,entity_type;`,
		utils.Placeholders(len(entityIDs)))

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, utils.ToInterface(entityIDs)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var counts []domain.Reaction
	for rows.Next() {
		var count domain.Reaction
		if err := rows.Scan(&count.EntityId, &count.EntityType, &count.Count); err != nil {
			return nil, err
		}
		counts = append(counts, count)
	}

	return counts, rows.Err() // Ensure any potential iteration errors are returned
}
