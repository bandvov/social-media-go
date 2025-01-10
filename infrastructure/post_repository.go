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
	_, err := r.db.Exec("INSERT INTO posts (author_id, content, visibility, pinned) VALUES ($1, $2, $3, $4);",
		post.AuthorID, post.Content, post.Visibility, post.Pinned)
	return err
}

func (r *PostRepository) Update(postId int, post *domain.Post) error {
	_, err := r.db.Exec("UPDATE posts SET content = $1, visibility = $2, pinned = $3, WHERE id = $4",
		post.Content, post.Visibility, post.Pinned, postId)
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
    p.id AS post_id,
    p.author_id,
    u.username AS author_name,
    p.content,
    p.visibility,
    p.pinned,
	p.created_at,
	p.updated_at,
    COALESCE(
        json_agg(
            json_build_object(
                'reaction_type', grouped_reactions.reaction_type,
                'count', grouped_reactions.reaction_count
            )
        ),
        '[]'
    ) AS reactions,
	COALESCE(SUM(grouped_reactions.reaction_count), 0) AS total_count
	FROM 
		posts p
	LEFT JOIN 
    users u ON p.author_id = u.id
	LEFT JOIN (
		SELECT 
			r.entity_id AS post_id,
			rt.name AS reaction_type,
			COUNT(r.id) AS reaction_count
		FROM 
			reactions r
		LEFT JOIN 
			reaction_types rt ON r.reaction_type_id = rt.id
		GROUP BY 
			r.entity_id, rt.name
	) grouped_reactions ON p.id = grouped_reactions.post_id
	WHERE
		p.id = $1 -- Replace with the post ID you want to query
	GROUP BY 
		p.id, u.username
	ORDER BY 
    p.id;
;
;
`, id).
		Scan(&post.ID, &post.AuthorID, post.AuthorName, &post.Content, &post.Pinned, &post.Visibility, &post.CreatedAt, &post.UpdatedAt, &post.Reactions, &post.TotalCount)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) FindByUserID(userID int) ([]domain.Post, error) {
	rows, err := r.db.Query(`	
	SELECT 
    p.id AS post_id,
    p.author_id,
    u.username AS author_name,
    p.content,
    p.visibility,
    p.pinned,
	p.created_at,
	p.updated_at,
    COALESCE(
        json_agg(
            json_build_object(
                'reaction_type', grouped_reactions.reaction_type,
                'count', grouped_reactions.reaction_count
            )
        ),
        '[]'
    ) AS reactions,
	COALESCE(SUM(grouped_reactions.reaction_count), 0) AS total_count
	FROM 
		posts p
	LEFT JOIN 
		users u ON p.author_id = u.id
	LEFT JOIN (
		SELECT 
			r.entity_id AS post_id,
			rt.name AS reaction_type,
			COUNT(r.id) AS reaction_count
		FROM 
			reactions r
		LEFT JOIN 
			reaction_types rt ON r.reaction_type_id = rt.id
		GROUP BY 
			r.entity_id, rt.name
	) grouped_reactions ON p.id = grouped_reactions.post_id
	WHERE
		p.author_id = $1 -- Replace with the post ID you want to query
	GROUP BY 
		p.id, u.username
	ORDER BY 
    p.id;

	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var post domain.Post
		if err := rows.Scan(&post.ID,&post.AuthorID, &post.AuthorName, &post.Content, &post.Visibility, &post.Pinned, &post.CreatedAt, &post.UpdatedAt,&post.Reactions, &post.TotalCount); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
