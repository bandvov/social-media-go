package domain

import "errors"

// Tag represents the domain model for a tag.
type Tag struct {
	ID   int
	Name string
}

// Validate checks if the tag is valid.
func (t *Tag) Validate() error {
	if t.Name == "" {
		return errors.New("tag name cannot be empty")
	}
	return nil
}
