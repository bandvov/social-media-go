package domain

import (
	"context"
)

type PostRepository interface {
	Create(ctx context.Context, post CreatePostRequest) error
	GetByID(ctx context.Context, id int) (*Post, error)
	Update(ctx context.Context, id int, post *Post) error
	Delete(ctx context.Context, id int) error
	GetByUserID(ctx context.Context, userID int, p Pagination) ([]Post, error)
	GetCountPostsByUser(ctx context.Context, userId int) (int, error)
}
