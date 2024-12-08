package infrastructure

import (
	"database/sql"
	"fmt"

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
	user := &domain.User{}
	err := r.db.QueryRow("SELECT id, username, password, email, status, role FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Status, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (r *UserRepository) GetUserByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRow("SELECT id,  password, email, status, role FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Password, &user.Email, &user.Status, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateUser(user *domain.User) error {
	_, err := r.db.Exec("UPDATE users SET email = $1, password = $2, status = $3 WHERE id = $4",
		user.Email, user.Password, user.Status, user.ID)
	return err
}

func (u *UserRepository) GetAllUsers(limit, offset int, sort string) ([]*domain.User, error) {
	// Validate and set default sorting
	order := "DESC"
	if sort == "asc" {
		order = "ASC"
	}

	query := fmt.Sprintf(`
        SELECT id, username, email, status, role, created_at
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
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Status, &user.Role, &user.CreatedAt)
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
