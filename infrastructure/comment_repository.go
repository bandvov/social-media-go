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
		"INSERT INTO comments (id, entity_id, content, author_id) VALUES ($1, $2, $3, $4)",
		comment.ID, comment.EntityID, comment.Content, comment.AuthorID,
	)
	return err
}

func (r *PostgresCommentRepository) GetCommentsByEntityID(entityID int) ([]domain.Comment, error) {
	rows, err := r.db.Query("SELECT id, entity_id, content, author_id FROM comments WHERE post_id = $1", entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var comment domain.Comment
		if err := rows.Scan(&comment.ID, &comment.EntityID, &comment.Content, &comment.AuthorID); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
