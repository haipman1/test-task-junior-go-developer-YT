package models

import "time"

type User struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email" binding:"required"`
	Password  string    `json:"password" db:"password" binding:"required"`
	Role      string    `json:"role" db:"role"` //user, admin
	CreatedAt time.Time `json:"registered" db:"registered"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
