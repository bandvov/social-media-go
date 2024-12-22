package domain

type PostRepository interface {
	Create(post *CreatePostRequest) error
	GetByID(id int) (*Post, error)
	Update(id int, post *Post) error
	Delete(id int) error
	FindByUserID(id int) ([]Post, error)
	// GetAll() ([]*Post, error)
}
