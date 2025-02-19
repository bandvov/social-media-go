package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bandvov/social-media-go/domain"
)

type UserRepository struct {
	db    *sql.DB
	cache Cache
}

func NewUserRepository(db *sql.DB, cache Cache) *UserRepository {
	return &UserRepository{db: db, cache: cache}
}

func (r *UserRepository) CreateUser(user *domain.User) error {
	// Prepare the statement
	stmt, err := r.db.Prepare("INSERT INTO users (password, email, status, role) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Password, user.Email, user.Status, user.Role)
	return err
}

func (r *UserRepository) GetUserByUsername(username string) (*domain.User, error) {
	user := &domain.User{}

	// Prepare the statement
	stmt, err := r.db.Prepare("SELECT id, username, password, email, status, role FROM users WHERE username = $1")
	if err != nil {
		return nil, fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).
		Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Status, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByID(id int) (*domain.User, error) {
	var user domain.User

	cacheKey := fmt.Sprintf("user:%d", id)

	ctx := context.Background()
	cachedUser, err := r.cache.Get(ctx, cacheKey)
	if err == nil && cachedUser != "" {
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			return &user, nil
		}
	}

	stmt, err := r.db.Prepare(`	
	WITH 
	post_counts AS (
		SELECT
			p.author_id,
			COUNT(*) AS post_count
		FROM posts p
		GROUP BY p.author_id
	),
	follower_stats AS (
    SELECT
        f.followee_id,
        COUNT(f.follower_id) AS follower_count,
        COUNT(DISTINCT f.follower_id) FILTER (WHERE f.followee_id IS NOT NULL) AS followee_count
    FROM followers f
    GROUP BY f.followee_id
	)
	SELECT
		u.id,
		u.username,
		u.password,
		u.email,
		u.status,
		u.role,
		u.profile_pic,
		u.created_at,
		u.updated_at,
		COALESCE(pc.post_count, 0) AS post_count,
		COALESCE(fs.follower_count, 0) AS followers_count,
		COALESCE(fs.followee_count, 0) AS followees_count
	FROM users u
	LEFT JOIN post_counts pc ON u.id = pc.author_id
	LEFT JOIN follower_stats fs ON u.id = fs.followee_id
	WHERE u.id = $1;`)
	if err != nil {
		return nil, fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).
		Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Status, &user.Role, &user.ProfilePic, &user.CreatedAt, &user.UpdatedAt, &user.PostsCount, &user.FollowersCount, &user.FolloweesCount)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(user)
	if err == nil {
		r.cache.Set(ctx, cacheKey, string(data), 24*time.Hour)
	}
	return &user, nil
}

func (r *UserRepository) GetPublicProfiles(offset, limit int) ([]domain.User, error) {
	cacheKey := fmt.Sprintf("public_profiles:limit:%d:offset:%d", limit, offset)

	ctx := context.Background()
	cachedData, err := r.cache.Get(ctx, cacheKey)
	if err == nil && cachedData != "" {
		var users []domain.User
		if err := json.Unmarshal([]byte(cachedData), &users); err == nil {
			return users, nil
		}
	}

	stmt, err := r.db.Prepare(`SELECT id, username, profile_pic FROM users OFFSET $1 LIMIT $2`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	// Execute the prepared statement with parameters
	rows, err := stmt.Query(offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch public profiles: %v", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Username, &user.ProfilePic); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	data, err := json.Marshal(users)
	if err == nil {
		r.cache.Set(ctx, cacheKey, string(data), 24*time.Hour)
	}

	return users, nil
}

func (r *UserRepository) GetAdminProfiles(limit, offset int) ([]domain.User, error) {
	cacheKey := fmt.Sprintf("admin_profiles:limit:%d:offset:%d", limit, offset)

	ctx := context.Background()
	cachedData, err := r.cache.Get(ctx, cacheKey)
	if err == nil && cachedData != "" {
		var users []domain.User
		if err := json.Unmarshal([]byte(cachedData), &users); err == nil {
			return users, nil
		}
	}

	stmt, err := r.db.Prepare(`
		SELECT id, username, email, role, status, created_at, updated_at 
		FROM users 
		LIMIT $1 OFFSET $2
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch admin profiles: %v", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Role,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	data, err := json.Marshal(users)
	if err == nil {
		r.cache.Set(ctx, cacheKey, string(data), 24*time.Hour)
	}

	return users, nil
}

func (r *UserRepository) GetUserProfileInfo(id, authenticatedUser int) (*domain.User, error) {
	var user domain.User

	cacheKey := fmt.Sprintf("public-user:%v:authenticatedUser:%v", id, authenticatedUser)
	ctx := context.Background()
	cachedUser, err := r.cache.Get(ctx, cacheKey)
	if err == nil && cachedUser != "" {
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			return &user, nil
		}
	}

	// Prepare the statement
	stmt, err := r.db.Prepare(`	
	WITH 
    post_counts AS (
        SELECT
            author_id,
            COUNT(*) AS post_count
        FROM posts
        WHERE author_id = $1 -- Filter early to reduce computation
        GROUP BY author_id
    ),
    follower_stats AS (
        SELECT
            $1 AS user_id,
            COUNT(DISTINCT follower_id) FILTER (WHERE followee_id = $1) AS followers_count, -- Count of users who follow the current user
            COUNT(DISTINCT followee_id) FILTER (WHERE follower_id = $1) AS followees_count -- Count of users the current user follows
        FROM followers
    ),
    relationship_flags AS (
        SELECT 
            MAX(CASE WHEN follower_id = $1 AND followee_id = $2 THEN 1 ELSE 0 END)::BOOLEAN AS is_followee,
            MAX(CASE WHEN followee_id = $1 AND follower_id = $2 THEN 1 ELSE 0 END)::BOOLEAN AS is_follower
        FROM followers
        WHERE 
            (follower_id = $1 AND followee_id = $2)
            OR (followee_id = $1 AND follower_id = $2)
    )
	SELECT
    u.id,
    u.username,
    u.first_name,
    u.last_name,
    u.email,
    u.status,
    u.role,
    u.profile_pic,
    u.created_at,
    u.updated_at,
    COALESCE(pc.post_count, 0) AS post_count,
    COALESCE(fs.followers_count, 0) AS followers_count,
    COALESCE(fs.followees_count, 0) AS followees_count,
    COALESCE(rf.is_follower, FALSE) AS is_follower,
    COALESCE(rf.is_followee, FALSE) AS is_followee
	FROM users u
	LEFT JOIN post_counts pc ON u.id = pc.author_id
	LEFT JOIN follower_stats fs ON u.id = fs.user_id
	LEFT JOIN relationship_flags rf ON TRUE
	WHERE u.id = $1;
;
`)
	if err != nil {
		return nil, fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(id, authenticatedUser).
		Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Status, &user.Role, &user.ProfilePic, &user.CreatedAt, &user.UpdatedAt, &user.PostsCount, &user.FollowersCount, &user.FolloweesCount, &user.IsFollower, &user.IsFollowee)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(user)
	if err == nil {
		r.cache.Set(ctx, cacheKey, string(data), 24*time.Hour)
	}
	return &user, nil
}
func (r *UserRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User

	cacheKey := fmt.Sprintf("user:%v", email)

	ctx := context.Background()
	cachedUser, err := r.cache.Get(ctx, cacheKey)
	if err == nil && cachedUser != "" {
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			return &user, nil
		}
	}

	// Prepare the statement
	stmt, err := r.db.Prepare("SELECT id, username, password, email, status, role, profile_pic, created_at, updated_at FROM users WHERE email = $1;")
	if err != nil {
		return nil, fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(email).
		Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Status, &user.Role, &user.ProfilePic, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(user)
	if err == nil {
		r.cache.Set(ctx, cacheKey, string(data), 24*time.Hour)
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(user *domain.User) error {
	query, err := r.buildUpdateQuery(user)
	if err != nil {
		return err
	}

	// Prepare the statement
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	return err
}

func (u *UserRepository) buildUpdateQuery(user *domain.User) (string, error) {
	var setClauses []string

	if user.FirstName != nil {
		setClauses = append(setClauses, fmt.Sprintf("first_name = '%s'", *user.FirstName))
	}
	if user.Email != "" {
		setClauses = append(setClauses, fmt.Sprintf("email = '%s'", user.Email))
	}
	if user.LastName != nil {
		setClauses = append(setClauses, fmt.Sprintf("last_name = '%s'", *user.LastName))
	}
	if user.Bio != nil {
		setClauses = append(setClauses, fmt.Sprintf("bio = '%s'", *user.Bio))
	}
	if user.ProfilePic != nil {
		setClauses = append(setClauses, fmt.Sprintf("profile_pic = '%s'", *user.ProfilePic))
	}
	if user.Password != "" {
		setClauses = append(setClauses, fmt.Sprintf("password = '%s'", user.Password))
	}
	if user.Status != "" {
		setClauses = append(setClauses, fmt.Sprintf("status = '%s'", user.Status))
	}
	if user.Role != "" {
		setClauses = append(setClauses, fmt.Sprintf("role = '%s'", user.Role))
	}
	if user.Username != nil {
		setClauses = append(setClauses, fmt.Sprintf("username = '%s'", *user.Username))
	}

	if len(setClauses) == 0 {
		return "", errors.New("No fields to update")
	}

	setClause := strings.Join(setClauses, ", ")
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = %d;", setClause, user.ID)
	return query, nil
}

// GetUsersByID fetches user details for a given set of user IDs.
func (r *UserRepository) GetUsersByID(ctx context.Context, userIDs []int) ([]domain.User, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}

	// Generate placeholders for the SQL IN clause
	placeholders := make([]string, len(userIDs))
	args := make([]interface{}, len(userIDs))
	for i, id := range userIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT id, username, profile_pic
		FROM users
		WHERE id IN (%s)
	`, strings.Join(placeholders, ", "))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Username, &user.ProfilePic); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
