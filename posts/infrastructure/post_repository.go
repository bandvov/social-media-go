package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"posts/domain"
)

type PostRepository struct {
	db    *sql.DB
	cache *RedisCache
}

func NewPostRepository(db *sql.DB, cache *RedisCache) *PostRepository {
	return &PostRepository{db: db, cache: cache}
}
func (r *PostRepository) Create(ctx context.Context, post *domain.CreatePostRequest) error {
	// Prepare the insert query using a prepared statement with context
	stmt, err := r.db.PrepareContext(ctx, `
		INSERT INTO posts (author_id, content, visibility, pinned)
		VALUES ($1, $2, $3, $4);
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the prepared statement with the context and parameters
	_, err = stmt.ExecContext(ctx, post.AuthorID, post.Content, post.Visibility, post.Pinned)
	return err
}

func (r *PostRepository) Update(ctx context.Context, postId int, post *domain.Post) error {
	// Prepare the update query using a prepared statement with context
	stmt, err := r.db.PrepareContext(ctx, `
		UPDATE posts 
		SET content = $1, visibility = $2, pinned = $3
		WHERE id = $4;
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the prepared statement with the context and parameters
	_, err = stmt.ExecContext(ctx, post.Content, post.Visibility, post.Pinned, postId)
	return err
}

func (r *PostRepository) Delete(ctx context.Context, id int) error {
	// Prepare the delete query using a prepared statement with context
	stmt, err := r.db.PrepareContext(ctx,
		`UPDATE posts
		SET visibility = 5 
		WHERE id = $1;`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the prepared statement with the context and parameter
	_, err = stmt.ExecContext(ctx, id)
	return err
}
func (r *PostRepository) GetByID(ctx context.Context, id int) (*domain.Post, error) {
	var post domain.Post

	query := `
		SELECT
			id AS post_id,
			author_id,
			content,
			pinned,
			visibility,
			created_at,
			updated_at
		FROM posts
		WHERE id = $1;
	`
	// Prepare the select query using a prepared statement with context
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the prepared statement with context and parameters
	err = stmt.QueryRowContext(ctx, id).Scan(&post.ID, &post.AuthorID, &post.Content, &post.Pinned, &post.Visibility, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) GetByUserID(ctx context.Context, userID int, p domain.Pagination) ([]domain.Post, error) {

	query := `
		SELECT
			id AS post_id,
			author_id,
			username AS author_name,
			content,
			visibility,
			pinned,
			created_at,
			updated_at
		FROM posts
		WHERE author_id = $1
		ORDER BY id DESC
		OFFSET $2
		LIMIT $3;`
	// Prepare the query using a prepared statement.
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the query with context and parameters
	rows, err := stmt.QueryContext(ctx, userID, p.Offset, p.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect the posts
	var posts []domain.Post
	for rows.Next() {
		var post domain.Post
		if err := rows.Scan(&post.ID, &post.AuthorID, &post.AuthorName, &post.Content, &post.Visibility, &post.Pinned, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	// Check for any errors that may have occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) GetCountPostsByUser(ctx context.Context, authorID int) (int, error) {
	var postsCount int

	stmt, err := r.db.Prepare(`
		SELECT COUNT(*) AS posts_count
		FROM posts
		WHERE author_id = $1;
    `)

	if err != nil {
		return postsCount, fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(authorID).Scan(&postsCount)
	if err != nil {
		return postsCount, fmt.Errorf("Failed to execute query: %v", err)
	}
	return postsCount, nil
}
