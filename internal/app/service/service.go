package service

//Вся бизнес логика
import (
	"context"
	"testtaskYT/internal/app/models"
	"testtaskYT/internal/app/repository"
	"time"
)

type Service struct {
	taskRepo *repository.TaskRepository
}

func New(rep *repository.TaskRepository) *Service {
	return &Service{
		taskRepo: rep,
	}
}

func (s *Service) DaysCount() int64 {
	dt := time.Date(2026, time.June, int(time.Monday), 0, 0, 0, 0, time.UTC)
	dt_until := int64(time.Until(dt).Hours()) / 24
	return dt_until
}

func (s *Service) CreateTask(ctx context.Context, title string) (*models.Task, error) {
	return s.taskRepo.CreateTask(ctx, title)
}
func (s *Service) GetTaskByID(ctx context.Context, id int) (*models.Task, error) {
	return s.taskRepo.GetById(ctx, id)
}

func (s *Service) DeleteTask(ctx context.Context, id int) error {
	return s.taskRepo.Delete(ctx, id)
}

func (s *Service) GetTasksPaginated(ctx context.Context, page, limit int) ([]models.Task, int, int, error) {
	// Устанавливаем значения по умолчанию
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// Получаем задачи
	tasks, err := s.taskRepo.GetTasksPaginated(ctx, page, limit)
	if err != nil {
		return nil, 0, 0, err
	}

	// Получаем общее количество
	total, err := s.taskRepo.GetTotalTasksCount(ctx)
	if err != nil {
		return nil, 0, 0, err
	}

	// Вычисляем общее количество страниц
	totalPages := (total + limit - 1) / limit

	return tasks, total, totalPages, nil
}

func (s *Service) UpdateTaskPartial(ctx context.Context, id int, updates map[string]interface{}) (*models.Task, error) {
	return s.taskRepo.UpdatePartial(ctx, id, updates)
}
