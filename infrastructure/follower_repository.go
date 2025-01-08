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

func (r *FollowerRepository) GetFollowers(userID, otherUser, limit, offset int, sort string, orderBy string, searchTerm string) ([]domain.User, error) {
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

	query := fmt.Sprintf(`
	SELECT 
    u.id,
    u.profile_pic,
    CASE 
        WHEN EXISTS (
            SELECT 1 
            FROM followers f2 
            WHERE f2.followee_id = f1.follower_id AND f2.follower_id = $2
        ) THEN TRUE
        ELSE FALSE
    END AS is_followed_by
	FROM followers f1
	JOIN users u ON f1.follower_id = u.id
	WHERE f1.followee_id = $1
	GROUP BY u.id, f1.follower_id, u.%v`, orderBy)

	if searchTerm != "" {
		query += fmt.Sprintf("\nWHERE position('%v' IN id) > 0 \n", searchTerm)
	}

	query += fmt.Sprintf("\nORDER BY %s %s\nLIMIT $3 OFFSET $4", orderBy, sort)
	fmt.Println(query)
	rows, err := r.db.Query(query, userID, otherUser, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %v", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.ProfilePic, &user.IsFollowedBy); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *FollowerRepository) GetFollowees(userID, otherUser, limit, offset int, sort string, orderBy string, searchTerm string) ([]domain.User, error) {
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

	query := fmt.Sprintf(`
	SELECT 
    u.id,
    u.profile_pic AS followee_profile_pic,
    CASE 
        WHEN EXISTS (
            SELECT 1 
            FROM followers f2 
            WHERE f2.followee_id = f.followee_id 
              AND f2.follower_id = $2
        ) THEN TRUE
        ELSE FALSE
    END AS is_followed_by
	FROM followers f
	JOIN users u ON f.followee_id = u.id
	WHERE f.follower_id = $1
	GROUP BY u.id, f.followee_id, u.%v`, orderBy)

	if searchTerm != "" {
		query += fmt.Sprintf("\nWHERE position('%v' IN id) > 0 \n", searchTerm)
	}

	query += fmt.Sprintf("\nORDER BY %s %s\nLIMIT $3 OFFSET $4", orderBy, sort)

	rows, err := r.db.Query(query, userID, otherUser, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get followees: %v", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.ProfilePic, &user.IsFollowedBy); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}
