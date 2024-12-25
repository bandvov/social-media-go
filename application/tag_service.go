package application

import "github.com/bandvov/social-media-go/domain"

// TagServiceInterface defines methods for tags-related operations.
type TagServiceInterface interface {
	CreateTag(name string) (*domain.Tag, error)
	GetAllTags() ([]*domain.Tag, error)
	DeleteTag(id int) error
	
}

// TagService provides use case methods for managing tags.
type TagService struct {
	repo domain.TagRepository
}

// NewTagService creates a new TagService.
func NewTagService(repository domain.TagRepository) *TagService {
	return &TagService{repo: repository}
}

// CreateTag validates and creates a new tag.
func (s *TagService) CreateTag(name string) (*domain.Tag, error) {
	tag := &domain.Tag{
		Name: name,
	}

	if err := tag.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// GetAllTags retrieves all tags.
func (s *TagService) GetAllTags() ([]*domain.Tag, error) {
	return s.repo.FindAll()
}

// DeleteTag removes tag from the table.
func (s *TagService) DeleteTag(id int) error {
	return s.repo.Delete(id)
}
