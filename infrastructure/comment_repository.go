package infrastructure

import (
	"database/sql"
	"fmt"

	"github.com/bandvov/social-media-go/domain"
	"github.com/bandvov/social-media-go/utils"
)

type PostgresCommentRepository struct {
	db *sql.DB
}

func NewPostgresCommentRepository(db *sql.DB) *PostgresCommentRepository {
	return &PostgresCommentRepository{db: db}
}

func (r *PostgresCommentRepository) AddComment(comment domain.Comment) error {
	_, err := r.db.Exec(
		"INSERT INTO comments (entity_id, entity_type, content, author_id) VALUES ($1, $2, $3, $4)",
		comment.EntityID, comment.EntityType, comment.Content, comment.AuthorID,
	)
	return err
}

func (r *PostgresCommentRepository) FetchCommentsByEntityID(entityID, userID, offset, limit int) ([]domain.Comment, error) {

	// Prepare the SQL query
	stmt, err := r.db.Prepare(`
	SELECT 
    c.id, 
    c.entity_id, 
    c.content, 
    c.author_id, 
    u.username, 
    u.profile_pic, 
    c.created_at,   
    COALESCE((
        SELECT rt.name
        FROM reactions rct2
        JOIN reaction_types rt ON rct2.reaction_type_id = rt.id
        WHERE rct2.entity_id = c.id AND rct2.user_id = $2
        LIMIT 1
    ),'') AS user_reaction,
    COALESCE(SUM(grouped_reactions.reaction_count), 0) AS total_reactions_count,
    COALESCE(
        json_agg(
            json_build_object(
                'reaction_type', grouped_reactions.reaction_type,
                'count', grouped_reactions.reaction_count
            )
        ) FILTER (WHERE grouped_reactions.reaction_type IS NOT NULL),
        '[]'
    ) AS reactions,
    COALESCE(r.reply_count, 0) AS replies_count  -- Added this line
FROM comments c
LEFT JOIN (
    SELECT 
        r.entity_id,
        rt.name AS reaction_type,
        COUNT(r.id) AS reaction_count
    FROM 
        reactions r
    LEFT JOIN 
        reaction_types rt ON r.reaction_type_id = rt.id
    GROUP BY 
        r.entity_id, rt.name
	) grouped_reactions ON c.id = grouped_reactions.entity_id
	LEFT JOIN (
    SELECT 
        entity_id, 
        COUNT(*) AS reply_count
    FROM comments
    WHERE entity_type = 'reply'
    GROUP BY entity_id
	) r ON c.id = r.entity_id
	LEFT JOIN users u ON c.author_id = u.id
	WHERE c.entity_id = $1 AND c.entity_type = 'comment'
	GROUP BY c.id, u.username, u.profile_pic, r.reply_count
	ORDER BY c.created_at DESC
	OFFSET $3 LIMIT $4;

	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the prepared statement
	rows, err := stmt.Query(entityID, userID, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var comment domain.Comment
		if err := rows.Scan(
			&comment.ID,
			&comment.EntityID,
			&comment.Content,
			&comment.AuthorID,
			&comment.Username,
			&comment.ProfilePic,
			&comment.CreatedAt,
			&comment.UserReaction,
			&comment.TotaReactionslCount,
			&comment.Reactions,
			&comment.RepliesCount,
		); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

// Fetch comments by post IDs
func (r *PostgresCommentRepository) GetCommentsByEntityIDs(entityIDs []int) ([]domain.Comment, error) {
	if len(entityIDs) == 0 {
		return nil, nil
	}

	// Prepare query with IN clause
	query := fmt.Sprintf(`
	SELECT id, entity_id, content, author_id, created_at
	FROM comments
	WHERE entity_id IN (%s)`, utils.Placeholders(len(entityIDs)))
	
	rows, err := r.db.Query(query, utils.ToInterface(entityIDs)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	fmt.Println("here1")
	
	var comments []domain.Comment
	for rows.Next() {
		var comment domain.Comment
		if err := rows.Scan(&comment.ID, &comment.EntityID, &comment.Content, &comment.AuthorID, &comment.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	fmt.Println("here2")
	
	return comments, nil
}
