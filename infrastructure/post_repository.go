package infrastructure

import (
	"database/sql"

	"github.com/bandvov/social-media-go/domain"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *domain.CreatePostRequest) error {
	_, err := r.db.Exec("INSERT INTO posts (author_id, content, visibility, tags, pinned) VALUES ($1, $2, $3, $4, $5);",
		post.AuthorID, post.Content, post.Visibility, post.Tags, post.Pinned)
	return err
}

func (r *PostRepository) Update(postId int, post *domain.Post) error {
	_, err := r.db.Exec("UPDATE posts SET content = $1, visibility = $2, pinned = $3, tags = $4 WHERE id = $5",
		post.Content, post.Visibility, post.Pinned, post.Tags, postId)
	return err
}
func (r *PostRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE from posts WHERE id = $1;", id)
	return err
}

func (r *PostRepository) GetByID(id int) (*domain.Post, error) {
	var post domain.Post
	err := r.db.QueryRow(`
	SELECT 
	p.id,
	p.author_id, 
	p.content, 
	p.pinned, 
	p.visibility, 
	p.tags, 
	p.created_at, 
	p.updated_at,
    json_agg(
        json_build_object(
            'id', r.id,
            'type', rt.name,
            'user', json_build_object(
                'id', u.id,
                'profile_pic', u.profile_pic
            )
        )
    ) AS reactions
	FROM posts p
	LEFT JOIN reactions r ON p.id = r.entity_id
	LEFT JOIN reaction_types rt ON r.reaction_type_id = rt.id
	LEFT JOIN users u ON r.user_id = u.id
	WHERE p.id = $1
	GROUP BY p.id
	ORDER BY p.id;`, id).
		Scan(&post.ID, &post.AuthorID, &post.Content, &post.Pinned, &post.Visibility, &post.Tags, &post.CreatedAt, &post.UpdatedAt, &post.Reactions)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) FindByUserID(userID int) ([]domain.Post, error) {
	rows, err := r.db.Query(`	
	SELECT 
	p.id,
	p.author_id, 
	p.content, 
	p.pinned, 
	p.visibility,
	p.created_at, 
	p.updated_at,
    json_agg(
        json_build_object(
            'id', r.id,
            'type', rt.name,
            'user', json_build_object(
                'id', u.id,
                'profile_pic', u.profile_pic
            )
        )
    ) AS reactions,
	 json_agg(
        json_build_object(
            'id', c.id,
            'author_id', c.author_id,
            'content', c.content,
            'replies', (
                SELECT json_agg(
                    json_build_object(
                        'id', nc.id,
                        'author_id', nc.author_id,
                        'content', nc.content
                    )
                )
                FROM comments nc
                WHERE nc.entity_id = c.id AND c.entity_type = 'comment'
            )
        )
    ) AS comments 
	FROM posts p
	LEFT JOIN reactions r ON p.id = r.entity_id
	LEFT JOIN reaction_types rt ON r.reaction_type_id = rt.id
	LEFT JOIN users u ON r.user_id = u.id
	LEFT JOIN comments c ON p.id = c.entity_id AND c.entity_type = 'post'
	WHERE user_id = $1
	GROUP BY p.id
	ORDER BY p.id;`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var post domain.Post
		if err := rows.Scan(&post.ID, &post.AuthorID, &post.Content, &post.Pinned, &post.Visibility, &post.Tags, &post.CreatedAt, &post.UpdatedAt, &post.Reactions); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
