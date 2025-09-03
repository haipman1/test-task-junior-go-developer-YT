package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testtaskYT/internal/app/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, email, passwordHash, role string) (*models.User, error) {
	query := `INSERT INTO users (email, password, role) VALUES($1, $2, $3)
	RETURNING id, email, password, role, created_at`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email, passwordHash, role).Scan(
		&user.ID, &user.Email, &user.Password, &user.Role, &user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	query := `
        SELECT id, email, password, role, created_at, updated_at 
        FROM users 
        WHERE id = $1
    `

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.Role, &user.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
        SELECT id, email, password, role, created_at
        FROM users 
        WHERE email = $1
    `

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.Role, &user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
