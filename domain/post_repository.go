package domain

type PostRepository interface {
	Create(post *CreatePostRequest) error
	GetByID(id int) (*Post, error)
	Update(id int, post *Post) error
	Delete(id int) error
	FindByUserID(userId, offset, limit int) ([]Post, error)
	GetCountPostsByUser(userId int) (int, error)
	// GetAll() ([]*Post, error)
}
