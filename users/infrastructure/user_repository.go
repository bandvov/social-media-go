package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
	"users/domain"
)

type UserRepository struct {
	db    *sql.DB
	cache Cache
}

func NewUserRepository(db *sql.DB, cache Cache) *UserRepository {
	return &UserRepository{db: db, cache: cache}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	// Prepare the statement
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO users (password, email, status, role) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute the statement
	_, err = stmt.ExecContext(ctx, user.Password, user.Email, user.Status, user.Role)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}
func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	user := &domain.User{}

	// Prepare the statement
	stmt, err := r.db.PrepareContext(ctx, "SELECT id, username, password, email, status, role FROM users WHERE username = $1")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute the statement
	err = stmt.QueryRowContext(ctx, username).
		Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Status, &user.Role)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User

	cacheKey := fmt.Sprintf("user:%d", id)

	// Try to get the user from cache
	cachedUser, err := r.cache.Get(ctx, cacheKey)
	if err == nil && cachedUser != "" {
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			return &user, nil
		}
	}

	// Prepare the statement for querying the database
	stmt, err := r.db.PrepareContext(ctx, `
	SELECT
		id,
		username,
		first_name,
		last_name,
		email,
		status,
		role,
		profile_pic,
		created_at,
		updated_at
	FROM public.users
	WHERE id = $1;
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute the query
	err = stmt.QueryRowContext(ctx, id).
		Scan(
			&user.ID,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Status,
			&user.Role,
			&user.ProfilePic,
			&user.CreatedAt,
			&user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Cache the result
	data, err := json.Marshal(user)
	if err == nil {
		r.cache.Set(ctx, cacheKey, string(data), 24*time.Hour)
	}

	return &user, nil
}

func (r *UserRepository) GetPublicProfiles(ctx context.Context, p domain.Pagination) ([]domain.User, error) {
	cacheKey := fmt.Sprintf("public_profiles:offset:%d:limit:%d", p.Offset, p.Limit)

	// Try to get the data from cache
	cachedData, err := r.cache.Get(ctx, cacheKey)
	if err == nil && cachedData != "" {
		var users []domain.User
		if err := json.Unmarshal([]byte(cachedData), &users); err == nil {
			return users, nil
		}
	}
	// Prepare the statement for querying the database
	stmt, err := r.db.PrepareContext(ctx,
		`SELECT id, username, profile_pic 
		 FROM users 
		 GROUP BY id 
		 OFFSET $1 LIMIT $2`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute the prepared statement with parameters
	rows, err := stmt.QueryContext(ctx, p.Offset, p.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch public profiles: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Username, &user.ProfilePic); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	// Cache the result
	data, err := json.Marshal(users)
	if err == nil {
		r.cache.Set(ctx, cacheKey, string(data), 24*time.Hour)
	}

	return users, nil
}
func (r *UserRepository) GetAdminProfiles(ctx context.Context, p domain.Pagination) ([]domain.User, error) {
	cacheKey := fmt.Sprintf("admin_profiles:offset:%d:limit:%d", p.Offset, p.Limit)

	// Try to get the data from the cache
	cachedData, err := r.cache.Get(ctx, cacheKey)
	if err == nil && cachedData != "" {
		var users []domain.User
		if err := json.Unmarshal([]byte(cachedData), &users); err == nil {
			return users, nil
		}
	}

	// Prepare the statement for querying the database
	stmt, err := r.db.PrepareContext(ctx,
		`SELECT id, username, email, role, status, created_at, updated_at 
		 FROM users
		 GROUP BY id 
		 OFFSET $1 LIMIT $2;
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute the prepared statement with parameters
	rows, err := stmt.QueryContext(ctx, p.Offset, p.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch admin profiles: %w", err)
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
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	// Cache the result
	data, err := json.Marshal(users)
	if err == nil {
		r.cache.Set(ctx, cacheKey, string(data), 24*time.Hour)
	}

	return users, nil
}

func (r *UserRepository) GetUserProfileInfo(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User

	cacheKey := fmt.Sprintf("user-profile-info:%v", id)

	// Try to get the data from the cache
	cachedUser, err := r.cache.Get(ctx, cacheKey)
	if err == nil && cachedUser != "" {
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			return &user, nil
		}
	}

	// Prepare the statement for querying the database
	stmt, err := r.db.PrepareContext(ctx, `	
		SELECT
			id,
			username,
			first_name,
			last_name,
			email,
			role,
			profile_pic,
			created_at,
			updated_at
		FROM public.users
		WHERE status != 'banned'
		AND id = $1;
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute the query with the context
	err = stmt.QueryRowContext(ctx, id).
		Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Role, &user.ProfilePic, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Cache the result
	data, err := json.Marshal(user)
	if err == nil {
		r.cache.Set(ctx, cacheKey, string(data), 24*time.Hour)
	}

	return &user, nil
}
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	cacheKey := fmt.Sprintf("user:%v", email)

	// Try to get the data from the cache
	cachedUser, err := r.cache.Get(ctx, cacheKey)
	if err == nil && cachedUser != "" {
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			return &user, nil
		}
	}

	// Prepare the statement for querying the database
	stmt, err := r.db.PrepareContext(ctx, "SELECT id, password, email FROM users WHERE email = $1;")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute the query with the context
	err = stmt.QueryRowContext(ctx, email).
		Scan(&user.ID, &user.Password, &user.Email)
	if err != nil {
		return nil, err
	}

	// Cache the result
	data, err := json.Marshal(user)
	if err == nil {
		r.cache.Set(ctx, cacheKey, string(data), 24*time.Hour)
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	query, err := r.buildUpdateQuery(user)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, query)
	// Prepare the statement with context
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute the statement with context
	_, err = stmt.ExecContext(ctx)
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
func (r *UserRepository) GetUsersByIDs(ctx context.Context, userIDs []int) ([]domain.User, error) {
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
		return nil, fmt.Errorf("failed to execute query: %v", err)
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

	return users, nil
}
