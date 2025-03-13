package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"posts/domain"
	"regexp"
	"strings"
)

type PostRepository struct {
	db    *sql.DB
	cache *RedisCache
}

func NewPostRepository(db *sql.DB, cache *RedisCache) *PostRepository {
	return &PostRepository{db: db, cache: cache}
}

func (r *PostRepository) Create(ctx context.Context, post domain.CreatePostRequest) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Prepare the insert query using a prepared statement with context
	stmt, err := r.db.PrepareContext(ctx, `
		INSERT INTO posts (author_id, content, visibility, pinned)
		VALUES ($1, $2, $3, $4) RETURNING id;
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the prepared statement with the context and parameters
	sqlResult, err := stmt.ExecContext(ctx, post.AuthorID, post.Content, post.Visibility, post.Pinned)
	if err != nil {
		return err
	}

	postID, err := sqlResult.LastInsertId()
	if err != nil {
		return err
	}

	keywords := extractKeywords(post.Content)
	tags := extractTags(post.Content)

	// Insert keywords and link them to the post
	for _, keyword := range keywords {
		var keywordID int
		err = tx.QueryRow("INSERT INTO keywords (keyword) VALUES ($1) ON CONFLICT (keyword) DO UPDATE SET keyword=EXCLUDED.keyword RETURNING id", keyword).Scan(&keywordID)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.Exec("INSERT INTO post_keywords (post_id, keyword_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", postID, keywordID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Insert tags
	for _, tag := range tags {
		var tagID int
		err = tx.QueryRow("INSERT INTO tags (tag) VALUES ($1) ON CONFLICT (tag) DO UPDATE SET tag=EXCLUDED.tag RETURNING id", tag).Scan(&tagID)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.Exec("INSERT INTO post_tags (post_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", postID, tagID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit transaction
	return tx.Commit()
}

func (r *PostRepository) Update(ctx context.Context, postId int, post *domain.Post) error {
	// Prepare the update query using a prepared statement with context
	stmt, err := r.db.PrepareContext(ctx, `
		UPDATE posts 
		SET content = $1, visibility = $2, pinned = $3
		WHERE id = $4;
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the prepared statement with the context and parameters
	_, err = stmt.ExecContext(ctx, post.Content, post.Visibility, post.Pinned, postId)
	return err
}

func (r *PostRepository) Delete(ctx context.Context, id int) error {
	// Prepare the delete query using a prepared statement with context
	stmt, err := r.db.PrepareContext(ctx,
		`UPDATE posts
		SET visibility = 5 
		WHERE id = $1;`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the prepared statement with the context and parameter
	_, err = stmt.ExecContext(ctx, id)
	return err
}
func (r *PostRepository) GetByID(ctx context.Context, id int) (*domain.Post, error) {
	var post domain.Post

	query := `
		SELECT
			id AS post_id,
			author_id,
			content,
			pinned,
			visibility,
			created_at,
			updated_at
		FROM posts
		WHERE id = $1;
	`
	// Prepare the select query using a prepared statement with context
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the prepared statement with context and parameters
	err = stmt.QueryRowContext(ctx, id).Scan(&post.ID, &post.AuthorID, &post.Content, &post.Pinned, &post.Visibility, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) GetByUserID(ctx context.Context, userID int, p domain.Pagination) ([]domain.Post, error) {

	query := `
		SELECT
			id AS post_id,
			author_id,
			content,
			visibility,
			pinned,
			created_at,
			updated_at
		FROM posts
		WHERE author_id = $1
		ORDER BY id DESC
		OFFSET $2
		LIMIT $3;`
	// Prepare the query using a prepared statement.
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the query with context and parameters
	rows, err := stmt.QueryContext(ctx, userID, p.Offset, p.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect the posts
	var posts []domain.Post
	for rows.Next() {
		var post domain.Post
		if err := rows.Scan(&post.ID, &post.AuthorID, &post.Content, &post.Visibility, &post.Pinned, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	// Check for any errors that may have occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) GetCountPostsByUser(ctx context.Context, authorID int) (int, error) {
	var postsCount int

	stmt, err := r.db.Prepare(`
		SELECT COUNT(*) AS posts_count
		FROM posts
		WHERE author_id = $1;
    `)

	if err != nil {
		return postsCount, fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(authorID).Scan(&postsCount)
	if err != nil {
		return postsCount, fmt.Errorf("Failed to execute query: %v", err)
	}
	return postsCount, nil
}

// Function to extract keywords from content
func extractKeywords(content string) []string {
	// Convert to lowercase
	content = strings.ToLower(content)

	// Remove punctuation using regex
	re := regexp.MustCompile(`[^\w\s]`)
	content = re.ReplaceAllString(content, "")

	// Split into words
	words := strings.Fields(content)

	// Define stopwords to ignore (extend as needed)
	stopwords := map[string]bool{
		"the": true, "is": true, "and": true, "or": true, "to": true, "a": true, "in": true, "of": true, "this": true,
	}

	// Filter out stopwords and short words
	var keywords []string
	seen := make(map[string]bool) // To avoid duplicates
	for _, word := range words {
		if len(word) > 2 && !stopwords[word] && !seen[word] { // Ignore short words and stopwords
			keywords = append(keywords, word)
			seen[word] = true
		}
	}
	return keywords
}

// Extracts hashtags (tags) from the content
func extractTags(content string) []string {
	re := regexp.MustCompile(`#\w+`) // Matches words prefixed with #
	matches := re.FindAllString(content, -1)

	// Convert to lowercase and remove duplicates
	tagSet := make(map[string]bool)
	var tags []string
	for _, tag := range matches {
		tag = strings.ToLower(tag)
		if !tagSet[tag] {
			tagSet[tag] = true
			tags = append(tags, tag)
		}
	}
	return tags
}
