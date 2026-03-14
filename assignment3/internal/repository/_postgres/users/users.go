package users

import (
	"assignment2/internal/repository/_postgres"
	"assignment2/pkg/modules"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Repository struct {
	db               *_postgres.Dialect
	executionTimeout time.Duration
}

// NewUserRepository создает новый экземпляр репозитория
func NewUserRepository(db *_postgres.Dialect) *Repository {
	return &Repository{
		db:               db,
		executionTimeout: time.Second * 5,
	}
}

// CreateUser создает пользователя и возвращает его ID
// CreateUser с поддержкой транзакции (Medium)
func (r *Repository) CreateUser(user modules.User) (int, error) {
	// Начинаем транзакцию
	tx, err := r.db.DB.Beginx()
	if err != nil {
		return 0, err
	}

	var id int
	query := `INSERT INTO users (name, email, age, is_active) VALUES ($1, $2, $3, $4) RETURNING id`
	err = tx.QueryRow(query, user.Name, user.Email, user.Age, user.IsActive).Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Запись в аудит-лог (требование Medium)
	auditQuery := `INSERT INTO audit_logs (user_id, action) VALUES ($1, $2)`
	_, err = tx.Exec(auditQuery, id, "User Created")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

// GetUsers с фильтрацией удаленных (Medium)
func (r *Repository) GetUsers(limit, offset int) ([]modules.User, error) {
	var users []modules.User
	// Только те, у кого deleted_at IS NULL (Soft Delete)
	query := `SELECT id, name, email, age, is_active FROM users 
              WHERE deleted_at IS NULL 
              LIMIT $1 OFFSET $2`

	err := r.db.DB.Select(&users, query, limit, offset)
	return users, err
}

func (r *Repository) DeleteUser(id int) error {
	// Вместо DELETE делаем UPDATE (Soft Delete)
	query := `UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	res, err := r.db.DB.Exec(query, id)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return fmt.Errorf("user not found or already deleted")
	}
	return nil
}

// GetUserByID находит одного пользователя по ID
func (r *Repository) GetUserByID(id int) (*modules.User, error) {
	var user modules.User
	query := `SELECT id, name, email, age, is_active FROM users WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.DB.Get(&user, query, id)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// UpdateUser обновляет данные пользователя
func (r *Repository) UpdateUser(id int, user modules.User) error {
	query := `UPDATE users SET name=$1, email=$2, age=$3, is_active=$4 WHERE id=$5`

	result, err := r.db.DB.Exec(query, user.Name, user.Email, user.Age, user.IsActive, id)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("no user found with given ID")
	}
	return nil
}
