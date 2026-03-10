package modules

import "time"

// User represents a user record in the system.
// swagger:model User
type User struct {
	ID        int        `db:"id" json:"id" example:"1"`
	Name      string     `db:"name" json:"name" example:"John Doe"`
	Email     string     `db:"email" json:"email" example:"john@example.com"`
	Age       int        `db:"age" json:"age" example:"25"`
	IsActive  bool       `db:"is_active" json:"is_active" example:"true"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
