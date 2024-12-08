package domain

type UserRepository interface {
	CreateUser(user *User) error
	GetUserByUsername(username string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	UpdateUser(user *User) error
}
