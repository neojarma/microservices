package entity

import "time"

type User struct {
	ID         int       `json:"id"`
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name,omitempty"`
	LastName   string    `json:"last_name,omitempty"`
	Password   string    `json:"-"`
	UserActive int       `json:"user_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
