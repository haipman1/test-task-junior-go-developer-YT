package models

//Это модель бизнес-сущности "Задача"

import "time"

type Task struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title" binding:"required"`
	Completed bool      `json:"completed" db:"completed"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
