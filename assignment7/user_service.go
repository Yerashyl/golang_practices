package assignment7

import (
	"errors"
	"fmt"
)

// UserService provides business logic for users.
type UserService struct {
	repo UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// RegisterUser checks if user exists, then creates it.
func (s *UserService) RegisterUser(user *User, email string) error {
	existing, err := s.repo.GetByEmail(email)
	if existing != nil {
		return fmt.Errorf("user with this email already exists")
	}

	if err != nil {
		// Expecting "not found" or similar, but assignment shows generic check
		// In a real app we'd check if err is specific.
		// For now following the snippet logic.
		return fmt.Errorf("error getting user with this email")
	}

	return s.repo.CreateUser(user)
}

// UpdateUserName checks name not empty, then updates.
func (s *UserService) UpdateUserName(id int, newName string) error {
	if newName == "" {
		return fmt.Errorf("name cannot be empty")
	}

	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return err
	}

	user.Name = newName
	return s.repo.UpdateUser(user)
}

// DeleteUser prevents deleting admin user.
func (s *UserService) DeleteUser(id int) error {
	if id == 1 {
		return errors.New("it is not allowed to delete admin user")
	}

	return s.repo.DeleteUser(id)
}
