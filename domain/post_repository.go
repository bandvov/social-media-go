package domain

type PostRepository interface {
	Create(post *CreatePostRequest) error
	GetByID(id int) (*Post, error)
	Update(post *Post) error
	Delete(id int) error
	// GetAll() ([]*Post, error)
}
