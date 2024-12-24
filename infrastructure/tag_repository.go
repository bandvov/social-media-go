package infrastructure

import (
	"database/sql"

	"github.com/bandvov/social-media-go/domain"
)

type TagRepository struct {
	db *sql.DB
}

// NewTagRepository creates a new repository instance.
func NewTagRepository(db *sql.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) Create(tag *domain.Tag) error {
	_, err := r.db.Exec("INSERT INTO tags (id, name) VALUES ($1, $2)", tag.ID, tag.Name)
	return err
}

func (r *TagRepository) FindByID(id string) (*domain.Tag, error) {
	row := r.db.QueryRow("SELECT id, name FROM tags WHERE id = $1", id)
	tag := &domain.Tag{}
	if err := row.Scan(&tag.ID, &tag.Name); err != nil {
		return nil, err
	}
	return tag, nil
}

func (r *TagRepository) FindAll() ([]*domain.Tag, error) {
	rows, err := r.db.Query("SELECT id, name FROM tags")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*domain.Tag
	for rows.Next() {
		tag := &domain.Tag{}
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (r *TagRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM tags WHERE id = $1", id)
	return err
}
