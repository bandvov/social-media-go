package domain

type PostRepository interface {
	Create(post *CreatePostRequest) error
	// GetByID(id string) (*Post, error)
	Update(post *Post) error
	Delete(id int) error
	// GetAll() ([]*Post, error)
	// GetAllByUserId() ([]*Post, error)
}
