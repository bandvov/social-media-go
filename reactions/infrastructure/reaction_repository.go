package infrastructure

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"reactions/domain"
	"time"

	"github.com/lib/pq"
	pg "github.com/lib/pq"
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

func (r *ReactionRepository) GetReactionsByEntityIDsAndTypes(ctx context.Context, entities []domain.Entity) ([]domain.Reaction, error) {
	if len(entities) == 0 {
		return nil, nil
	}

	// Extract entity IDs and types
	var ids []int
	var types []string
	for _, e := range entities {
		ids = append(ids, e.ID)
		types = append(types, e.Type)
	}

	query := `
        SELECT r.entity_id, rt.name AS reaction, COUNT(r.id) AS count
        FROM reactions r
        JOIN reaction_types rt ON r.reaction_type_id = rt.id
        WHERE (entity_id, entity_type) IN (SELECT * FROM UNNEST($1::int[], $2::entity_type[]))
        GROUP BY r.entity_id, rt.name`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := r.db.QueryContext(ctx, query, pg.Array(ids), pg.Array(types))
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

func (r *ReactionRepository) CountByEntityIDsAndTypes(ctx context.Context, entities []domain.Entity) ([]domain.Reaction, error) {
	if len(entities) == 0 {
		return nil, nil
	}

	// Extract entity IDs and types
	var ids []int
	var types []string
	for _, e := range entities {
		ids = append(ids, e.ID)
		types = append(types, e.Type)
	}

	query := `
        SELECT entity_id, entity_type, COUNT(*) AS count
        FROM reactions 
        WHERE (entity_id, entity_type) IN (SELECT * FROM UNNEST($1::int[], $2::entity_type[]))
        GROUP BY entity_id, entity_type;
    `

	// Execute query
	rows, err := r.db.QueryContext(ctx, query, pg.Array(ids), pg.Array(types))
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

	return counts, rows.Err() // Ensure any iteration errors are returned
}

func (r *ReactionRepository) GetReacionStats(ctx context.Context, entities []domain.Entity) ([]domain.ReactionStat, error) {
	if len(entities) == 0 {
		return nil, nil
	}

	// Extract entity IDs and types
	var ids []int
	var types []string
	for _, e := range entities {
		ids = append(ids, e.ID)
		types = append(types, e.Type)
	}

	query := `
	WITH reactions_aggregated AS (
		SELECT 
		entity_id,
		entity_type,
		user_id,
		reaction_type_id,
		COUNT(*) AS count
		FROM reactions
		WHERE (entity_id, entity_type) IN (SELECT * FROM UNNEST($1::int[], $2::entity_type[]))  -- Filtering by entity_id and entity_type
		GROUP BY entity_id, entity_type, user_id, reaction_type_id
		), reactions_grouped AS (
			SELECT
			entity_id,
			entity_type,
			jsonb_object_agg(rt.name, r.count) AS reactions  -- Aggregating reaction counts per type
			FROM reactions_aggregated r
			JOIN reaction_types rt ON r.reaction_type_id = rt.id
			GROUP BY entity_id, entity_type
			)
			SELECT 
			rg.entity_id,
			rg.entity_type,
			SUM(
				(SELECT SUM((kv.value)::int) FROM jsonb_each_text(rg.reactions) AS kv(key, value))  -- Summing reaction counts across all types
				) AS total_reactions,
				rg.reactions  -- Returning the grouped reactions as is
				FROM reactions_grouped rg
				GROUP BY rg.entity_id, rg.entity_type, rg.reactions;`

	// Execute query
	rows, err := r.db.QueryContext(ctx, query, pg.Array(ids), pg.Array(types))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []domain.ReactionStat
	for rows.Next() {
		var stat domain.ReactionStat
		if err := rows.Scan(&stat.EntityId, &stat.EntityType, &stat.TotalCount, &stat.Reactions); err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}

	return stats, rows.Err() // Ensure any iteration errors are returned
}

func (r *ReactionRepository) GetUserReactions(ctx context.Context, userID int, entities []domain.Entity) ([]domain.Reaction, error) {

	// Generate cache key using a sorted JSON format for consistency
	entityKey, _ := json.Marshal(entities)
	cacheKey := fmt.Sprintf("user:%d:entities:%x", userID, sha256.Sum256(entityKey))

	// Check Redis cache
	cachedData, err := r.cache.Get(ctx, cacheKey)
	if err == nil {
		var reactions []domain.Reaction
		if err := json.Unmarshal([]byte(cachedData), &reactions); err == nil {
			return reactions, nil
		}
	}

	if len(entities) == 0 {
		return nil, nil
	}

	// Extract entity IDs and types
	var ids []int
	var types []string
	for _, e := range entities {
		ids = append(ids, e.ID)
		types = append(types, e.Type)
	}

	// Query reactions
	query := `SELECT user_id, entity_id, entity_type, rt.name 
		FROM reactions
		LEFT JOIN reaction_types rt 
		ON  reaction_type_id = rt.id
		WHERE user_id = $1 
		AND entity_id = ANY($2::int[]) 
		AND entity_type = ANY($3::entity_type[]);`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	// Use prepared statement
	rows, err := stmt.QueryContext(ctx, userID, pq.Array(ids), pg.Array(types))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var reactions []domain.Reaction
	for rows.Next() {
		var reaction domain.Reaction
		if err := rows.Scan(&reaction.UserId, &reaction.EntityId, &reaction.EntityType, &reaction.Name); err != nil {
			return nil, err
		}
		reactions = append(reactions, reaction)
	}

	// Store in Redis for future requests
	if len(reactions) > 0 {
		data, _ := json.Marshal(reactions)
		r.cache.Set(ctx, cacheKey, string(data), 10*time.Minute)
	}

	return reactions, nil
}
