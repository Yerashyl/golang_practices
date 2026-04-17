package assignment7

// User represents a user in the system.
type User struct {
	ID    int
	Name  string
	Email string
}

// UserRepository defines the persistence layer for users.
type UserRepository interface {
	GetByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(user *User) error
	UpdateUser(user *User) error
	DeleteUser(id int) error
}
