package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"followers/domain"
)

type FollowRepository struct {
	db    *sql.DB
	cache *RedisCache
}

func NewFollowRepository(db *sql.DB, cache *RedisCache) *FollowRepository {
	return &FollowRepository{db: db, cache: cache}
}

func (r *FollowRepository) AddFollower(follower *domain.Follower) error {
	query := "INSERT INTO followers (follower_id, followee_id) VALUES ($1, $2)"
	_, err := r.db.Exec(query, follower.FollowerID, follower.FolloweeID)
	if err != nil {
		return fmt.Errorf("failed to add follower: %v", err)
	}
	return nil
}

func (r *FollowRepository) RemoveFollower(follower *domain.Follower) error {
	query := "DELETE FROM followers WHERE follower_id = $1 AND followee_id = $2"
	_, err := r.db.Exec(query, follower.FollowerID, follower.FolloweeID)
	if err != nil {
		return fmt.Errorf("failed to remove follower: %v", err)
	}
	return nil
}

func (r *FollowRepository) GetFollowers(ctx context.Context, userID, otherUser, limit, offset int, sort, orderBy, searchTerm string) ([]domain.Follow, error) {
	// Validate and set default sorting
	if sort == "" || sort == "desc" {
		sort = "DESC"
	}
	if sort == "asc" {
		sort = "ASC"
	}

	if limit == 0 {
		limit = 24
	}

	// Start building the query
	query := `
	SELECT 
		follower_id AS id,
		CASE 
			WHEN follower_id = $2 THEN TRUE      
			ELSE FALSE                            
		END AS follows_follower,
		CASE 
			WHEN followee_id = $2 THEN TRUE      
			ELSE FALSE                            
		END AS followed_by_follower
	FROM followers 
	WHERE followee_id = $1
	`

	// Add search term condition if searchTerm is provided
	if searchTerm != "" {
		query += " AND follower_id::text LIKE $5 "
	}

	// Add order by, limit, and offset clauses
	query += fmt.Sprintf(" ORDER BY %s %s LIMIT $3 OFFSET $4", orderBy, sort)

	// Prepare the query
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query: %v", err)
	}
	defer stmt.Close()

	// Set up query parameters: userID, otherUser, limit, offset, and searchTerm (if provided)
	params := []interface{}{userID, otherUser, limit, offset}
	if searchTerm != "" {
		params = append(params, "%"+searchTerm+"%")
	}

	// Execute the query
	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %v", err)
	}
	defer rows.Close()

	// Parse the results
	var followers []domain.Follow
	for rows.Next() {
		var follower domain.Follow
		if err := rows.Scan(&follower.ID, &follower.FollowsFollower, &follower.FollowedByFollower); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		followers = append(followers, follower)
	}

	// Return the followers
	return followers, nil
}

func (r *FollowRepository) GetFollowees(ctx context.Context, userID, otherUser, limit, offset int, sort, orderBy, searchTerm string) ([]domain.Follow, error) {
	// Validate and set default sorting
	if sort == "" || sort == "desc" {
		sort = "DESC"
	}
	if sort == "asc" {
		sort = "ASC"
	}
	if orderBy == "" {
		orderBy = "created_at"
	}
	if limit == 0 {
		limit = 24
	}

	// Base query
	query := `
	SELECT 
		followee_id AS id,
		CASE 
			WHEN follower_id = $2 THEN TRUE      
			ELSE FALSE                            
		END AS follows_follower,
		CASE 
			WHEN followee_id = $2 THEN TRUE      
			ELSE FALSE                            
		END AS followed_by_follower
	FROM followers 
	WHERE follower_id = $1
	`

	// Add search term condition if searchTerm is not empty
	if searchTerm != "" {
		query += " AND followee_id::text LIKE $5 "
	}

	// Add sorting, limit, and offset
	query += fmt.Sprintf(" ORDER BY %s %s LIMIT $3 OFFSET $4", orderBy, sort)

	// Prepare the query
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query: %v", err)
	}
	defer stmt.Close()

	// Parameters: userID, otherUser, limit, offset, searchTerm (if provided)
	params := []interface{}{userID, otherUser, limit, offset}
	if searchTerm != "" {
		params = append(params, "%"+searchTerm+"%")
	}

	// Execute the query
	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to get followees: %v", err)
	}
	defer rows.Close()

	// Parse the results
	var followees []domain.Follow
	for rows.Next() {
		var followee domain.Follow
		if err := rows.Scan(&followee.ID, &followee.FollowsFollower, &followee.FollowedByFollower); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		followees = append(followees, followee)
	}

	// Return the results
	return followees, nil
}

func (r *FollowRepository) GetFollowerStats(ctx context.Context, userID int) (int, int, error) {
	var followersCount, followeesCount int

	query := `
        SELECT 
            COUNT(DISTINCT follower_id) FILTER (WHERE followee_id = $1) AS followers_count,
            COUNT(DISTINCT followee_id) FILTER (WHERE follower_id = $1) AS followees_count
        FROM followers`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, 0, err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, userID).Scan(&followersCount, &followeesCount)
	if err != nil {
		return 0, 0, err
	}

	return followersCount, followeesCount, nil
}

func (r *FollowRepository) CheckFollowStatus(ctx context.Context, userID, targetUserID int64) (*domain.UserRelationship, error) {
	stmt, err := r.db.PrepareContext(ctx, `
    SELECT 
        EXISTS (SELECT 1 FROM followers WHERE follower_id = $1 AND followee_id = $2) AS is_follower,
        EXISTS (SELECT 1 FROM followers WHERE follower_id = $2 AND followee_id = $1) AS is_followed
`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close() // Close the statement when done

	var relationship domain.UserRelationship
	err = stmt.QueryRowContext(ctx, userID, targetUserID).Scan(&relationship.IsFollower, &relationship.IsFollowed)
	if err != nil {
		return nil, err
	}

	return &relationship, nil
}
