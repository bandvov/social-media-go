package domain

import "context"

type PostRepository interface {
	Create(post *CreatePostRequest) error
	GetByID(id int) (*Post, error)
	Update(id int, post *Post) error
	Delete(id int) error
	GetByUserID(ctx context.Context, userID, offset, limit int) ([]Post, error)
	GetCountPostsByUser(userId int) (int, error)
}
