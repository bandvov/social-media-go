package infrastructure

import (
	"database/sql"
	"fmt"

	"github.com/bandvov/social-media-go/domain"
)

type FollowerRepository struct {
	db *sql.DB
}

func NewFollowerRepository(db *sql.DB) *FollowerRepository {
	return &FollowerRepository{db: db}
}

func (r *FollowerRepository) AddFollower(follower *domain.Follower) error {
	query := "INSERT INTO followers (follower_id, followee_id) VALUES ($1, $2)"
	_, err := r.db.Exec(query, follower.FollowerID, follower.FolloweeID)
	if err != nil {
		return fmt.Errorf("failed to add follower: %v", err)
	}
	return nil
}

func (r *FollowerRepository) RemoveFollower(follower *domain.Follower) error {
	query := "DELETE FROM followers WHERE follower_id = $1 AND followee_id = $2"
	_, err := r.db.Exec(query, follower.FollowerID, follower.FolloweeID)
	if err != nil {
		return fmt.Errorf("failed to remove follower: %v", err)
	}
	return nil
}

func (r *FollowerRepository) GetFollowers(userID int, limit, offset int, sort string, orderBy string, searchTerm string) ([]domain.User, error) {
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

	query := `
	SELECT 
		f.follower_id,
		JSON_AGG(JSON_BUILD_OBJECT(
			'id', f.followee_id,
			'profile_pic',uf.profile_pic
		)) AS followers
	FROM followers f
	LEFT JOIN users uf ON f.follower_id = uf.id
	WHERE f.followee_id = $1
	GROUP BY f.follower_id;`

	if searchTerm != "" {
		query += fmt.Sprintf("\nWHERE position('%v' IN id) > 0 \n", searchTerm)
	}

	query += fmt.Sprintf("\nORDER BY %s %s\nLIMIT $2 OFFSET $3", orderBy, sort)

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %v", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Email, &user.ProfilePic); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *FollowerRepository) GetFollowees(userID int, limit, offset int, sort string, orderBy string, searchTerm string) ([]domain.User, error) {
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

	query := `
	SELECT 
		f.followee_id,
		JSON_AGG(JSON_BUILD_OBJECT(
			'id', f.follower_id,
			'profile_pic',uf.profile_pic
		)) AS followers
	FROM followers f
	LEFT JOIN users uf ON f.follower_id = uf.id
	WHERE f.follower_id = $1
	GROUP BY f.followee_id;`

	if searchTerm != "" {
		query += fmt.Sprintf("\nWHERE position('%v' IN id) > 0 \n", searchTerm)
	}

	query += fmt.Sprintf("\nORDER BY %s %s\nLIMIT $2 OFFSET $3", orderBy, sort)

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get followees: %v", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Email, &user.ProfilePic); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}
