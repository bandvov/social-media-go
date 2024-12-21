package infrastructure

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/bandvov/social-media-go/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *domain.User) error {
	_, err := r.db.Exec("INSERT INTO users (password, email, status, role) VALUES ($1, $2, $3, $4)",
		user.Password, user.Email, user.Status, user.Role)
	return err
}

func (r *UserRepository) GetUserByUsername(username string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRow("SELECT id, username, password, email, status, role FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Status, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByID(id int) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow("SELECT id, username, password, email, status, role, profile_pic, created_at, updated_at  FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Status, &user.Role, &user.ProfilePic, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
func (r *UserRepository) GetUserByEmail(email string) (*domain.User, error) {
	fmt.Println("GetUserByEmail", email)
	var user domain.User
	err := r.db.QueryRow("SELECT id, username, password, email, status, role, profile_pic, created_at, updated_at FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Status, &user.Role, &user.ProfilePic, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(user *domain.User) error {
	query, err := r.buildUpdateQuery(user)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(query)
	return err
}

func (u *UserRepository) GetAllUsers(limit, offset int, sort string) ([]*domain.User, error) {
	// Validate and set default sorting
	order := "DESC"
	if sort == "asc" {
		order = "ASC"
	}

	query := fmt.Sprintf(`
        SELECT id, username, email, status, role, profile_pic, created_at
        FROM users
        ORDER BY created_at %s
        LIMIT $1 OFFSET $2
    `, order)

	rows, err := u.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Status, &user.Role, &user.ProfilePic, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserRepository) buildUpdateQuery(user *domain.User) (string, error) {
	var setClauses []string

	if user.FirstName != "" {
		setClauses = append(setClauses, fmt.Sprintf("first_name = '%s'", user.FirstName))
	}
	if user.Email != "" {
		setClauses = append(setClauses, fmt.Sprintf("email = '%s'", user.Email))
	}
	if user.LastName != "" {
		setClauses = append(setClauses, fmt.Sprintf("last_name = '%s'", user.LastName))
	}
	if user.Bio != "" {
		setClauses = append(setClauses, fmt.Sprintf("bio = '%s'", user.Bio))
	}
	if user.ProfilePic != "" {
		setClauses = append(setClauses, fmt.Sprintf("profile_pic = '%s'", user.ProfilePic))
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
	if user.Username != "" {
		setClauses = append(setClauses, fmt.Sprintf("username = '%s'", user.Username))
	}

	if len(setClauses) == 0 {
		return "", errors.New("No fields to update")
	}

	setClause := strings.Join(setClauses, ", ")
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = %d;", setClause, user.ID)
	fmt.Println("query: ", query)
	return query, nil
}
