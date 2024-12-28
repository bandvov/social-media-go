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

	// Prepare the statement
	stmt, err := r.db.Prepare(`
	WITH followers_agg AS (
    SELECT
        f.followee_id,
        JSON_AGG(JSON_BUILD_OBJECT(
            'id', f.follower_id,
			'profile_pic',uf.profile_pic
        )) AS followers
    FROM followers f
    LEFT JOIN users uf ON f.follower_id = uf.id
    GROUP BY f.followee_id
	),
	followees_agg AS (
		SELECT
			fw.follower_id,
			JSON_AGG(JSON_BUILD_OBJECT(
				'id', fw.followee_id,
				'profile_pic',uw.profile_pic
			)) AS followees
		FROM followers fw
		LEFT JOIN users uw ON fw.followee_id = uw.id
		GROUP BY fw.follower_id
	),
	post_counts AS (
		SELECT
			p.author_id,
			COUNT(*) AS post_count
		FROM posts p
		GROUP BY p.author_id
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
		COALESCE(fa.followers, '[]') AS followers,
		COALESCE(fe.followees, '[]') AS followees
	FROM users u
	LEFT JOIN followers_agg fa ON u.id = fa.followee_id
	LEFT JOIN followees_agg fe ON u.id = fe.follower_id
	LEFT JOIN post_counts pc ON u.id = pc.author_id
	WHERE u.id = $1;`)
	if err != nil {
		return nil, fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).
		Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Status, &user.Role, &user.ProfilePic, &user.CreatedAt, &user.UpdatedAt, &user.PostsCount, &user.Followers, &user.Followeees)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
func (r *UserRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User

	// Prepare the statement
	stmt, err := r.db.Prepare("SELECT id, username, password, email, status, role, profile_pic, created_at, updated_at FROM users WHERE email = $1")
	if err != nil {
		return nil, fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(email).
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

	// Prepare the statement
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	return err
}

func (r *UserRepository) GetAllUsers(limit, offset int, sort string, orderBy string, search string) ([]*domain.User, error) {
	// Validate and set default sorting
	if sort == "desc" {
		sort = "DESC"
	}
	if sort == "asc" {
		sort = "ASC"
	}
	baseQuery := `SELECT id, username, email, status, role, profile_pic, created_at FROM users`
	if search != "" {
		baseQuery += " AND email ILIKE " + "%" + search + "%"
	}
	query := fmt.Sprintf(`        
        ORDER BY %s %s
        LIMIT $1 OFFSET $2
    `, orderBy, sort)

	// Prepare the statement
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(query, limit, offset)
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
