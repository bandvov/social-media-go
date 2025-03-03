package infrastructure

import (
	"comments/domain"
	"comments/utils"
	"context"
	"database/sql"
	"fmt"
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

func (r *PostgresCommentRepository) CountByEntityIDs(ctx context.Context, entityIDs []int) ([]domain.CommentCount, error) {
	// Generate the placeholders for the query based on the number of entityIDs
	placeholders := utils.Placeholders(len(entityIDs))

	// Prepare the query with placeholders
	query := fmt.Sprintf(`
        SELECT
			entity_id,
			COALESCE(COUNT(CASE WHEN entity_type = 'comment' THEN 1 END), 0) AS comment_count,
            COALESCE(COUNT(CASE WHEN entity_type = 'post' THEN 1 END), 0) AS reply_count
        FROM comments
        WHERE entity_id IN (%s)
		GROUP BY entity_id`, placeholders)

	// Prepare the statement using the context
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()

	// Execute the query using the provided entityIDs
	rows, err := stmt.QueryContext(ctx, utils.ToInterface(entityIDs)...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var counts []domain.CommentCount
	for rows.Next() {
		var count domain.CommentCount
		if err := rows.Scan(&count.EntityID, &count.CommentCount, &count.ReplyCount); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		counts = append(counts, count)
	}

	// Check for errors after the loop
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return counts, nil
}
