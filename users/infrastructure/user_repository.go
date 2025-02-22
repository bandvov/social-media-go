package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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
		return nil, fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).
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

func (r *UserRepository) GetUserProfileInfo(id int) (*domain.User, error) {
	var user domain.User

	cacheKey := fmt.Sprintf("user-profile-info:%v", id)
	ctx := context.Background()
	cachedUser, err := r.cache.Get(ctx, cacheKey)
	if err == nil && cachedUser != "" {
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			return &user, nil
		}
	}

	// Prepare the statement
	stmt, err := r.db.Prepare(`	
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
		return nil, fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).
		Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Role, &user.ProfilePic, &user.CreatedAt, &user.UpdatedAt)
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
	stmt, err := r.db.Prepare("SELECT password, email, FROM users WHERE email = $1;")
	if err != nil {
		return nil, fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(email).
		Scan(&user.Password, &user.Email)
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
