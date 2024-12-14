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

func (r *PostRepository) Create(post *domain.Post) error {
	_, err := r.db.Exec("INSERT INTO posts (author_id, content, visibility, tags, like_count, comment_count, share_count, pinned) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);",
		post.AuthorID, post.Content, post.Visibility, post.Tags, post.LikeCount, post.CommentCount, post.ShareCount, post.Pinned)
	return err
}

func (r *PostRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE from posts WHERE id = $1;", id)
	return err
}
