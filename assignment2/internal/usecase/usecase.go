package usecase

import (
	"assignment2/internal/repository"
	"assignment2/pkg/modules"
	"fmt"
)

// Описываем интерфейс Usecase
type UserUsecase interface {
	CreateUser(user modules.User) (int, error)
	GetAllUsers(limit int, offset int) ([]modules.User, error)
	GetUserByID(id int) (*modules.User, error)
	UpdateUser(id int, user modules.User) error
	DeleteUser(id int) error
}

type Usecases struct {
	UserUsecase
}

// Конструктор для всех Usecase
func NewUsecases(repos *repository.Repositories) *Usecases {
	return &Usecases{
		UserUsecase: NewUserLogic(repos.UserRepository),
	}
}

// Реализация логики
type userLogic struct {
	repo repository.UserRepository
}

func NewUserLogic(repo repository.UserRepository) UserUsecase {
	return &userLogic{repo: repo}
}

func (u *userLogic) CreateUser(user modules.User) (int, error) {
	// Валидация (Easy)
	if user.Name == "" {
		return 0, fmt.Errorf("name cannot be empty")
	}
	if user.Age < 0 || user.Age > 130 {
		return 0, fmt.Errorf("invalid age")
	}
	if user.Email == "" {
		return 0, fmt.Errorf("email is required")
	}
	return u.repo.CreateUser(user)
}

func (u *userLogic) GetAllUsers(limit int, offset int) ([]modules.User, error) {
	if limit <= 0 {
		limit = 10
	}
	return u.repo.GetUsers(limit, offset)
}

func (u *userLogic) GetUserByID(id int) (*modules.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *userLogic) UpdateUser(id int, user modules.User) error {
	return u.repo.UpdateUser(id, user)
}

func (u *userLogic) DeleteUser(id int) error {
	return u.repo.DeleteUser(id)
}
