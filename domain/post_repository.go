package domain

type PostRepository interface {
	Create(post *Post) error
	// GetByID(id string) (*Post, error)
	// Update(post *Post) error
	// Delete(id string) error
	// GetAll() ([]*Post, error)
	// GetAllByUserId() ([]*Post, error)
}
