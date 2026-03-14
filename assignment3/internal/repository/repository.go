package repository

import (
	"assignment2/internal/repository/_postgres"
	"assignment2/internal/repository/_postgres/users" // Важно: импортируем подпакет users
	"assignment2/pkg/modules"
)

type UserRepository interface {
	CreateUser(user modules.User) (int, error)
	// Добавляем int, int для limit и offset
	GetUsers(limit int, offset int) ([]modules.User, error)
	GetUserByID(id int) (*modules.User, error)
	UpdateUser(id int, user modules.User) error
	DeleteUser(id int) error
}

type Repositories struct {
	UserRepository
}

func NewRepositories(db *_postgres.Dialect) *Repositories {
	return &Repositories{
		// Мы обращаемся к пакету users, который импортировали в строке 5
		UserRepository: users.NewUserRepository(db),
	}
}
