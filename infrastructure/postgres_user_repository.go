package infrastructure

import (
	"database/sql"

	"github.com/bandvov/social-media-go/domain"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) CreateUser(user *domain.User) error {
	_, err := r.db.Exec("INSERT INTO users (username, password, email, status, role) VALUES ($1, $2, $3, $4, $5)",
		user.Username, user.Password, user.Email, user.Status, user.Role)
	return err
}

func (r *PostgresUserRepository) GetUserByUsername(username string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRow("SELECT id, username, password, email, status, role FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Status, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *PostgresUserRepository) GetUserByID(id int64) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRow("SELECT id, username, password, email, status, role FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Status, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (r *PostgresUserRepository) GetUserByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRow("SELECT id, username, password, email, status, role FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Status, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *PostgresUserRepository) UpdateUser(user *domain.User) error {
	_, err := r.db.Exec("UPDATE users SET email = $1, password = $2, status = $3 WHERE id = $4",
		user.Email, user.Password, user.Status, user.ID)
	return err
}
