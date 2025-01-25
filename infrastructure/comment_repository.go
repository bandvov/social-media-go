package infrastructure

import (
	"database/sql"

	"github.com/bandvov/social-media-go/domain"
)

type PostgresCommentRepository struct {
	db *sql.DB
}

func NewPostgresCommentRepository(db *sql.DB) *PostgresCommentRepository {
	return &PostgresCommentRepository{db: db}
}

func (r *PostgresCommentRepository) AddComment(comment domain.Comment) error {
	_, err := r.db.Exec(
		"INSERT INTO comments (entity_id, entity_type, content, author_id) VALUES ($1, $2, $3, $4)",
		comment.EntityID, comment.EntityType, comment.Content, comment.AuthorID,
	)
	return err
}

func (r *PostgresCommentRepository) FetchCommentsByEntityID(entityID, offset, limit int) ([]domain.Comment, error) {
	rows, err := r.db.Query(`SELECT 
    c.id, 
	c.entity_id, 
    c.content, 
	c.author_id,
    c.created_at,   
    COALESCE(r.reply_count, 0) AS replies_count
	FROM comments c
	LEFT JOIN (
    SELECT 
        entity_id, 
        COUNT(*) AS reply_count
    FROM comments
    WHERE entity_type = 'reply'
    GROUP BY entity_id
	) r ON c.id = r.entity_id
	WHERE c.entity_id = $1 AND c.entity_type = 'comment'
	ORDER BY c.created_at DESC
	OFFSET $2 LIMIT $3;
	`, entityID, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var comment domain.Comment
		if err := rows.Scan(&comment.ID, &comment.EntityID, &comment.Content, &comment.AuthorID, &comment.CreatedAt, &comment.RepliesCount); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
