package repository

//Layer between App and DataBase

import (
	"context"
	"database/sql"
	"fmt"
	"testtaskYT/internal/app/models"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (r *TaskRepository) CreateTask(ctx context.Context, title string) (*models.Task, error) {
	query := `INSERT INTO tasks (title) VALUES ($1) 
              RETURNING id, title, completed, created_at, updated_at`

	task := &models.Task{}
	err := r.db.QueryRowContext(ctx, query, title).Scan(
		&task.ID, &task.Title, &task.Completed, &task.CreatedAt, &task.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *TaskRepository) GetTasksPaginated(ctx context.Context, page int, limit int) ([]models.Task, error) {
	offset := (page - 1) * limit
	query := `
		SELECT id, title, completed, created_at, updated_at
		FROM tasks
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID, &task.Title, &task.Completed, &task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
func (r *TaskRepository) GetTotalTasksCount(ctx context.Context) (int, error) {
	query := "SELECT COUNT(*) FROM tasks"

	var total int
	err := r.db.QueryRowContext(ctx, query).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *TaskRepository) GetById(ctx context.Context, id int) (*models.Task, error) {
	query := "SELECT id, title, completed, created_at, updated_at FROM tasks WHERE id = $1"

	task := &models.Task{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID, &task.Title, &task.Completed, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM tasks WHERE id=$1"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("Delete() - ExecContext error: %w", err)
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error in Delete() - RowsAffeced(): %w", err)
	}
	if rowsAff == 0 {
		return fmt.Errorf("Delete() - task not found")
	}
	return nil
}

func (r *TaskRepository) UpdatePartial(ctx context.Context, id int, updates map[string]interface{}) (*models.Task, error) {
	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	allowedFields := map[string]bool{
		"title":     true,
		"completed": true,
	}
	filteredUpdates := make(map[string]interface{})
	for key, value := range updates {
		if allowedFields[key] {
			filteredUpdates[key] = value
		}
	}
	if len(filteredUpdates) == 0 {
		return nil, fmt.Errorf("no valid fields to update")
	}

	query := "UPDATE tasks SET "
	values := []interface{}{}
	paramCount := 1

	for key, val := range updates {
		query += fmt.Sprintf("%s = $%d, ", key, paramCount)
		values = append(values, val)
		paramCount++
	}

	query += fmt.Sprintf("updated_at = NOW() WHERE id = $%d", paramCount)
	values = append(values, id)

	query += " RETURNING id, title, completed, created_at, updated_at"

	task := &models.Task{}
	err := r.db.QueryRowContext(ctx, query, values...).Scan(
		&task.ID, &task.Title, &task.Completed, &task.CreatedAt, &task.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return task, nil
}
