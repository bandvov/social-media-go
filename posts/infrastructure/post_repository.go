package infrastructure

import (
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

func (r *PostRepository) Create(post *domain.CreatePostRequest) error {
	_, err := r.db.Exec("INSERT INTO posts (author_id, content, visibility, pinned) VALUES ($1, $2, $3, $4);",
		post.AuthorID, post.Content, post.Visibility, post.Pinned)
	return err
}

func (r *PostRepository) Update(postId int, post *domain.Post) error {
	_, err := r.db.Exec("UPDATE posts SET content = $1, visibility = $2, pinned = $3, WHERE id = $4",
		post.Content, post.Visibility, post.Pinned, postId)
	return err
}
func (r *PostRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE from posts WHERE id = $1;", id)
	return err
}

func (r *PostRepository) GetByID(id int) (*domain.Post, error) {
	var post domain.Post
	err := r.db.QueryRow(`
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
	`, id).
		Scan(&post.ID, &post.AuthorID, &post.Content, &post.Pinned, &post.Visibility, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) GetByUserID(userID, offset, limit int) ([]domain.Post, error) {
	rows, err := r.db.Query(`
	SELECT
		id AS post_id,
		author_id,
		username AS author_name,
		content,
		visibility,
		pinned,
		created_at,
		updated_at	
	WHERE author_id = $1 -- Author ID
	ORDER BY id DESC
	OFFSET $2
	LIMIT $3;`, userID, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var post domain.Post
		if err := rows.Scan(&post.ID, &post.AuthorID, &post.AuthorName, &post.Content, &post.Visibility, &post.Pinned, &post.CreatedAt, &post.UpdatedAt, &post.Reactions); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepository) GetCountPostsByUser(authorID int) (int, error) {
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
