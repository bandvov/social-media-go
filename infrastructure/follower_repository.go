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

func (r *FollowerRepository) GetFollowers(userID int) ([]domain.User, error) {
	query := `
	SELECT 
		f.follower_id,
		JSON_AGG(JSON_BUILD_OBJECT(
			'id', f.followee_id,
			'profile_pic',uf.profile_pic
		)) AS followers
	FROM followers f
	LEFT JOIN users uf ON f.follower_id = uf.id
	WHERE f.followee_id = 1
	GROUP BY f.follower_id;`
	rows, err := r.db.Query(query, userID)
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

func (r *FollowerRepository) GetFollowees(userID int) ([]domain.User, error) {
	query := `
	SELECT 
		f.followee_id,
		JSON_AGG(JSON_BUILD_OBJECT(
			'id', f.follower_id,
			'profile_pic',uf.profile_pic
		)) AS followers
	FROM followers f
	LEFT JOIN users uf ON f.follower_id = uf.id
	WHERE f.follower_id = 1
	GROUP BY f.followee_id;`
	rows, err := r.db.Query(query, userID)
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
