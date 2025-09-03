package endpoint

//TODO: Выносим логирование в хелперы
//HTTP логика: URL + HTTP метод
import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"testtaskYT/internal/app/models"
	"testtaskYT/internal/pkg/utils/logging"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func errorResponse(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
		Code:    statusCode,
	})
}
func (e *Endpoint) bindAndValidate(ctx echo.Context, structurePtr interface{}) error {
	if err := ctx.Bind(structurePtr); err != nil {
		return fmt.Errorf("Indalid JSON: %w", err)
	}
	//Validation
	if err := e.v.Struct(structurePtr); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			errs_map := map[string]string{}
			for _, e := range errs {
				errs_map[e.Field()] = getValidationMessage(e)
			}
			return ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
				"error":   "JSON request validation failed",
				"details": errs_map,
				"code":    http.StatusUnprocessableEntity,
			})
		}
		return fmt.Errorf("error in Validator JSON: %w", err)
	}
	return nil
}

func getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return fmt.Sprintf("Minimum length is %s characters", fe.Param())
	case "max":
		return fmt.Sprintf("Maximum length is %s characters", fe.Param())
	case "email":
		return "Invalid email format"
	case "omitempty":
		return "Field has invalid format"
	default:
		return fe.Error()
	}
}

type Service interface {
	DaysCount() int64
	CreateTask(ctx context.Context, title string) (*models.Task, error)
	GetTaskByID(ctx context.Context, id int) (*models.Task, error)
	DeleteTask(ctx context.Context, id int) error
	UpdateTaskPartial(ctx context.Context, id int, updates map[string]interface{}) (*models.Task, error)
	GetTasksPaginated(ctx context.Context, page, limit int) ([]models.Task, int, int, error)
}
type Endpoint struct {
	s Service
	v *validator.Validate
}

func New(s Service) *Endpoint {
	return &Endpoint{
		s: s,
		v: validator.New(), //TODO переместить этот валидатор и в auth_endpoint.go
	}
}

func (e *Endpoint) Status(ctx echo.Context) error {
	logging.MetodCalled(ctx)

	d := e.s.DaysCount()

	logging.StatusDaysCalculated(ctx, d)

	dt_str := fmt.Sprintf("Осталось дней: %d", d)
	err := ctx.String(http.StatusOK, dt_str)

	if err != nil {
		logging.StatusFailedSendResp(ctx, err, d)
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	logging.StatusReqCompleted(ctx, d)
	return nil
}

func (e *Endpoint) HelloWorld(ctx echo.Context) error {
	logging.MetodCalled(ctx)
	return ctx.String(http.StatusOK, "Hello world!")
}

func (e *Endpoint) CreateTask(ctx echo.Context) error {
	logging.MetodCalled(ctx)

	request := CreateTaskRequest{}
	if err := e.bindAndValidate(ctx, &request); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	logging.CrTaskBindAndValidate(ctx, request.Title)

	task, err := e.s.CreateTask(ctx.Request().Context(), request.Title)
	if err != nil {
		logging.CrTaskFailedCrTaskInService(ctx, err, request.Title)
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	logging.CrTaskCreatedSucc(ctx, request.Title)

	return ctx.JSON(http.StatusCreated, task)
}

func (e *Endpoint) GetTask(ctx echo.Context) error {
	logging.MetodCalled(ctx)

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logging.GetTaskInvalidID(ctx)

		return errorResponse(ctx, http.StatusBadRequest, "Invalid task ID")
	}
	logging.GetTaskParcedID(ctx, id)

	task, err := e.s.GetTaskByID(ctx.Request().Context(), id)
	if err != nil {
		logging.GetTaskFromServiceFailed(ctx, err, id)
		return errorResponse(ctx, http.StatusNotFound, err.Error())
	}
	logging.GetTaskSuccess(ctx, task.ID, task.Title)
	return ctx.JSON(http.StatusOK, task)
}

func (e *Endpoint) DeleteTask(ctx echo.Context) error {
	start := time.Now()

	logging.MetodCalled(ctx)

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logging.DelTaskInvalidID(ctx, id)
		return errorResponse(ctx, http.StatusBadRequest, "Invalid task ID")
	}
	logging.DelTaskParsedID(ctx, id)

	err = e.s.DeleteTask(ctx.Request().Context(), id)
	if err != nil {
		logging.DelTaskFromServiceFailed(ctx, err, id)
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	logging.DelTaskSuccess(ctx, id, start)
	return ctx.JSON(http.StatusOK, map[string]string{"message": "Task deleted successfully"})
}

func (e *Endpoint) GetAllTasks(ctx echo.Context) error {
	start := time.Now()
	logging.MetodCalled(ctx)

	// Парсим параметры пагинации
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	// Ограничиваем максимальный лимит
	if limit > 100 {
		limit = 100
	}

	logging.GetAllTasksParamsParsed(ctx, page, limit)

	// Получаем данные с пагинацией
	tasks, total, totalPages, err := e.s.GetTasksPaginated(ctx.Request().Context(), page, limit)
	if err != nil {
		logging.GetAllTasksFailedFromService(ctx, err, start)
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	// Формируем ответ
	response := GetAllTasksResponse{
		Tasks:      tasks,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	logging.GetAllTasksRetrievedSucc(ctx, len(tasks), total, page, totalPages, start)
	return ctx.JSON(http.StatusOK, response)
}

func (e *Endpoint) UpdateTaskPartial(ctx echo.Context) error {
	start := time.Now()

	logging.MetodCalled(ctx)

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logging.UpdTaskPartInvalidID(ctx, id)
		return errorResponse(ctx, http.StatusBadRequest, "Invalid task ID")
	}
	logging.UpdTaskPartParcedID(ctx, id)

	updates := UpdateTaskPartialRequest{}
	if err := e.bindAndValidate(ctx, &updates); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	if !updates.hasUpdate() {
		logging.UpdTaskPartValidationFailed(ctx, id, start)
		return errorResponse(ctx, http.StatusBadRequest,
			"No data or some error data has been sent for modification",
		)
	}
	logging.UpdTaskPartNoFields(ctx, id, start)

	task, err := e.s.UpdateTaskPartial(ctx.Request().Context(), id, updates.toMap())
	if err != nil {
		logging.UpdTaskPartFailedFromService(ctx, id, updates.toMap(), err, start)
		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	logging.UpdTaskPartSuccess(ctx, id, updates.toMap(), start)
	return ctx.JSON(http.StatusOK, task)
}
