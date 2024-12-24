package domain

// TagRepository defines the interface for tag persistence.
type TagRepository interface {
	Create(tag *Tag) error
	FindByID(id string) (*Tag, error)
	FindAll() ([]*Tag, error)
	Delete(id int) error
}
