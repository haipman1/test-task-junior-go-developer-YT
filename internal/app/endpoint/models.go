package endpoint

//Это модели для ОБМЕНА ДАННЫМИ с клиентом

import (
	"testtaskYT/internal/app/models"
	"time"
)

// МОДЕЛИ ЗАПРОСОВ (что приходит от клиента)
type CreateTaskRequest struct {
	Title string `json:"title" validate:"required,min=2,max=100"`
}

type UpdateTaskPartialRequest struct {
	Title     *string `json:"title" validate:"omitempty,min=2,max=100"`
	Completed *bool   `json:"completed" validate:"omitempty"`
}

func (st *UpdateTaskPartialRequest) toMap() map[string]interface{} {
	updates_map := map[string]interface{}{}
	if st.Title != nil {
		updates_map["title"] = *st.Title
	}
	if st.Completed != nil {
		updates_map["completed"] = *st.Completed
	}
	return updates_map
}
func (st *UpdateTaskPartialRequest) hasUpdate() bool {
	if st.Title != nil || st.Completed != nil {
		return true
	}
	return false
}

type GetAllTasksRequest struct {
	Page  int `query:"page" validate:"omitempty,min=1"`
	Limit int `query:"limit" validate:"omitempty,min=1,max=100"`
}

// МОДЕЛИ ОТВЕТОВ (что возвращаем клиенту)
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}

type TaskResponse struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type GetAllTasksResponse struct {
	Tasks      []models.Task `json:"tasks"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int           `json:"total_pages"`
}
