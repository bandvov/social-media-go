package infrastructure

import (
	"comments/domain"
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type PostgresCommentRepository struct {
	db    *sql.DB
	cache Cache
}

func NewPostgresCommentRepository(db *sql.DB, cache Cache) *PostgresCommentRepository {
	return &PostgresCommentRepository{db: db, cache: cache}
}

func (r *PostgresCommentRepository) AddComment(ctx context.Context, comment domain.Comment) error {
	// Prepare the SQL statement with the context
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO comments (entity_id, entity_type, content, author_id) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute the prepared statement with the provided comment values
	_, err = stmt.ExecContext(ctx, comment.EntityID, comment.EntityType, comment.Content, comment.AuthorID)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

func (r *PostgresCommentRepository) FetchCommentsByEntityID(ctx context.Context, entityID int, p domain.Pagination) ([]domain.Comment, error) {

	// Prepare the SQL query
	stmt, err := r.db.PrepareContext(ctx, `
	SELECT 
		c.id, 
		c.author_id, 
		c.entity_id, 
		c.entity_type,
		c.content, 
		c.created_at,
    COALESCE(r.reply_count, 0) AS reply_count
	FROM comments c
	LEFT JOIN (
		SELECT 
			entity_id, 
			COUNT(*) AS reply_count
		FROM comments
		WHERE entity_type = 'comment'
		GROUP BY entity_id
	) r ON c.id = r.entity_id
	WHERE c.entity_id = $1 AND c.entity_type = 'post'
	ORDER BY c.created_at DESC
	OFFSET $2 LIMIT $3;
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the prepared statement
	rows, err := stmt.QueryContext(ctx, entityID, p.Offset, p.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var comment domain.Comment
		if err := rows.Scan(
			&comment.ID,
			&comment.AuthorID,
			&comment.EntityID,
			&comment.EntityType,
			&comment.Content,
			&comment.CreatedAt,
			&comment.RepliesCount,
		); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (r *PostgresCommentRepository) CountByEntityIDAndType(ctx context.Context, entities []domain.Entity) ([]domain.CommentCount, error) {
	if len(entities) == 0 {
		return nil, nil
	}

	// Split into slices of IDs and Types
	ids := make([]int, 0, len(entities))
	types := make([]string, 0, len(entities))
	for _, e := range entities {
		ids = append(ids, e.ID)
		types = append(types, e.Type)
	}

	query := `
		WITH entity_input AS (
			SELECT UNNEST($1::int[]) AS entity_id, UNNEST($2::entity_type[]) AS entity_type
		)
		SELECT 
			e.entity_id,
			e.entity_type,
			COALESCE(COUNT(CASE WHEN c.entity_type = 'comment' THEN 1 END), 0) AS comment_count,
			COALESCE(COUNT(CASE WHEN c.entity_type = 'post' THEN 1 END), 0) AS reply_count
		FROM entity_input e
		LEFT JOIN comments c 
			ON e.entity_id = c.entity_id AND e.entity_type = c.entity_type
		GROUP BY e.entity_id, e.entity_type;
	`

	rows, err := r.db.QueryContext(ctx, query, pq.Array(ids), pq.Array(types))
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var counts []domain.CommentCount
	for rows.Next() {
		var count domain.CommentCount
		if err := rows.Scan(&count.EntityID, &count.EntityType, &count.CommentCount, &count.ReplyCount); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		counts = append(counts, count)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return counts, nil
}
