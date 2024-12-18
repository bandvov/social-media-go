package infrastructure

import (
	"database/sql"

	"github.com/bandvov/social-media-go/domain"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *domain.CreatePostRequest) error {
	_, err := r.db.Exec("INSERT INTO posts (author_id, content, visibility, tags, pinned) VALUES ($1, $2, $3, $4, $5);",
		post.AuthorID, post.Content, post.Visibility, post.Tags, post.Pinned)
	return err
}

func (r *PostRepository) Update(postId int, post *domain.Post) error {
	_, err := r.db.Exec("UPDATE posts SET content = $1, visibility = $2, pinned = $3, tags = $4 WHERE id = $5",
		post.Content, post.Visibility, post.Pinned, post.Tags, postId)
	return err
}
func (r *PostRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE from posts WHERE id = $1;", id)
	return err
}

func (r *PostRepository) GetByID(id int) (*domain.Post, error) {
	var post domain.Post
	err := r.db.QueryRow("SELECT id, author_id, content, pinned, visibility, tags, created_at, updated_at  FROM posts WHERE id = $1", id).
		Scan(&post.ID, &post.AuthorID, &post.Content, &post.Pinned, &post.Visibility, &post.Tags, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &post, nil
}
