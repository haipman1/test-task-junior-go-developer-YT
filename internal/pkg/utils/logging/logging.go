package logging

//TODO доработать (опимизировать чтобы было меньше функий, объединить их как-то и тп)
import (
	"os"
	"strings"
	"testtaskYT/internal/pkg/config"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

func Init(config config.Config) {
	if logger != nil {
		return
	}
	logger = logrus.New()

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.00",
		ForceColors:     true,
		PadLevelText:    true,
	})
	if strings.ToLower(config.LogLevel) == "debug" {
		logger.SetLevel(logrus.DebugLevel)
	} else if strings.ToLower(config.LogLevel) == "info" {
		logger.SetLevel(logrus.InfoLevel)
	}

	logger.SetOutput(os.Stdout)
}

func MetodCalled(ctx echo.Context) {
	logger.WithFields(logrus.Fields{
		"client_ip": ctx.RealIP(),
	}).Infof("%s(%s) -> Method called", ctx.Request().Method, ctx.Path())
}

func StatusDaysCalculated(ctx echo.Context, daysCount int64) {
	logger.WithFields(logrus.Fields{
		"client_ip":  ctx.RealIP(),
		"days_count": daysCount,
	}).Debugf("%s(%s) -> Days count calculated", ctx.Request().Method, ctx.Path())
}

func StatusFailedSendResp(ctx echo.Context, err error, daysCount int64) {
	logger.WithFields(logrus.Fields{
		"error":      err.Error(),
		"client_ip":  ctx.RealIP(),
		"days_count": daysCount,
	}).Errorf("%s(%s) -> Failed to send response", ctx.Request().Method, ctx.Path())
}
func StatusReqCompleted(ctx echo.Context, daysCount int64) {
	logger.WithFields(logrus.Fields{
		"client_ip":  ctx.RealIP(),
		"days_count": daysCount,
	}).Infof("%s(%s) -> Request completed successfully", ctx.Request().Method, ctx.Path())
}

func CrTaskBindAndValidate(ctx echo.Context, title string) {
	logger.WithFields(logrus.Fields{
		"client_ip": ctx.RealIP(),
		"title":     title,
	}).Debugf("%s(%s) -> Binding and validation completed", ctx.Request().Method, ctx.Path())
}
func CrTaskFailedCrTaskInService(ctx echo.Context, err error, title string) {
	logger.WithFields(logrus.Fields{
		"client_ip": ctx.RealIP(),
		"title":     title,
		"error":     err.Error(),
	}).Errorf("%s(%s) -> Failed to create task in service layer", ctx.Request().Method, ctx.Path())
}
func CrTaskCreatedSucc(ctx echo.Context, title string) {
	logger.WithFields(logrus.Fields{
		"client_ip": ctx.RealIP(),
		"title":     title,
	}).Infof("%s(%s) -> Task created successfully", ctx.Request().Method, ctx.Path())
}

func GetTaskInvalidID(ctx echo.Context) {
	logger.WithFields(logrus.Fields{
		"client_ip":     ctx.RealIP(),
		"task_id_param": ctx.Param("id"),
	}).Errorf("%s(%s) -> Invalid task ID format", ctx.Request().Method, ctx.Path())
}
func GetTaskParcedID(ctx echo.Context, id int) {
	logger.WithFields(logrus.Fields{
		"client_ip": ctx.RealIP(),
		"task_id":   id,
	}).Debugf("%s(%s) -> Task ID parsed successfully", ctx.Request().Method, ctx.Path())
}
func GetTaskFromServiceFailed(ctx echo.Context, err error, id int) {
	logger.WithFields(logrus.Fields{
		"client_ip": ctx.RealIP(),
		"task_id":   id,
		"error":     err.Error(),
	}).Errorf("%s(%s) -> Failed to get task from service", ctx.Request().Method, ctx.Path())
}
func GetTaskSuccess(ctx echo.Context, id int, title string) {
	logger.WithFields(logrus.Fields{
		"client_ip": ctx.RealIP(),
		"task_id":   id,
		"title":     title,
	}).Infof("%s(%s) -> Task received successfully", ctx.Request().Method, ctx.Path())
}

func DelTaskInvalidID(ctx echo.Context, id int) {
	logger.WithFields(logrus.Fields{
		"client_ip":     ctx.RealIP(),
		"task_id_param": id,
	}).Errorf("%s(%s) -> Invalid task ID format", ctx.Request().Method, ctx.Path())
}
func DelTaskParsedID(ctx echo.Context, id int) {
	logger.WithFields(logrus.Fields{
		"client_ip": ctx.RealIP(),
		"task_id":   id,
	}).Debugf("%s(%s) -> Task ID parsed successfully", ctx.Request().Method, ctx.Path())
}
func DelTaskFromServiceFailed(ctx echo.Context, err error, id int) {
	logger.WithFields(logrus.Fields{
		"client_ip": ctx.RealIP(),
		"task_id":   id,
		"error":     err.Error(),
	}).Errorf("%s(%s) -> Failed to delete task in service layer", ctx.Request().Method, ctx.Path())
}
func DelTaskSuccess(ctx echo.Context, id int, start time.Time) {
	logger.WithFields(logrus.Fields{
		"client_ip":   ctx.RealIP(),
		"task_id":     id,
		"duration_ms": time.Since(start).Milliseconds(),
	}).Infof("%s(%s) -> Task deleted successfully", ctx.Request().Method, ctx.Path())
}

func GetAllTasksParamsParsed(ctx echo.Context, page, limit int) {
	logger.WithFields(logrus.Fields{
		"page":  page,
		"limit": limit,
	}).Debugf("%s(%s) -> Pagination params parsed", ctx.Request().Method, ctx.Path())
}

func GetAllTasksRetrievedSucc(ctx echo.Context, tasksCount, total, page, totalPages int, start time.Time) {
	logger.WithFields(logrus.Fields{
		"tasks_count": tasksCount,
		"total":       total,
		"page":        page,
		"total_pages": totalPages,
		"duration_ms": time.Since(start).Milliseconds(),
	}).Infof("%s(%s) -> Tasks retrieved successfully", ctx.Request().Method, ctx.Path())
}

func GetAllTasksFailedFromService(ctx echo.Context, err error, start time.Time) {
	logger.WithFields(logrus.Fields{
		"client_ip":   ctx.RealIP(),
		"error":       err.Error(),
		"duration_ms": time.Since(start).Milliseconds(),
	}).Errorf("%s(%s) -> Failed to get tasks from service", ctx.Request().Method, ctx.Path())
}

func UpdTaskPartInvalidID(ctx echo.Context, id int) {
	logger.WithFields(logrus.Fields{
		"client_ip":     ctx.RealIP(),
		"task_id_param": id,
	}).Errorf("%s(%s) -> Invalid task ID format", ctx.Request().Method, ctx.Path())
}
func UpdTaskPartParcedID(ctx echo.Context, id int) {
	logger.WithFields(logrus.Fields{
		"client_ip": ctx.RealIP(),
		"task_id":   id,
	}).Debugf("%s(%s) -> Task ID parsed successfully", ctx.Request().Method, ctx.Path())
}

func UpdTaskPartValidationFailed(ctx echo.Context, id int, start time.Time) {
	logger.WithFields(logrus.Fields{
		"client_ip":   ctx.RealIP(),
		"task_id":     id,
		"duration_ms": time.Since(start).Milliseconds(),
	}).Errorf("%s(%s) -> Validation failed", ctx.Request().Method, ctx.Path())
}

func UpdTaskPartNoFields(ctx echo.Context, id int, start time.Time) {
	logger.WithFields(logrus.Fields{
		"client_ip":   ctx.RealIP(),
		"task_id":     id,
		"duration_ms": time.Since(start).Milliseconds(),
	}).Warnf("%s(%s) -> No fields to update", ctx.Request().Method, ctx.Path())
}

func UpdTaskPartFailedFromService(ctx echo.Context, id int, updates map[string]interface{},
	err error, start time.Time) {
	logger.WithFields(logrus.Fields{
		"client_ip":   ctx.RealIP(),
		"task_id":     id,
		"updates":     updates,
		"error":       err.Error(),
		"duration_ms": time.Since(start).Milliseconds(),
	}).Errorf("%s(%s) -> Failed to update task in service", ctx.Request().Method, ctx.Path())
}
func UpdTaskPartSuccess(ctx echo.Context, id int, updates map[string]interface{}, start time.Time) {
	logger.WithFields(logrus.Fields{
		"client_ip":   ctx.RealIP(),
		"task_id":     id,
		"updates":     updates,
		"duration_ms": time.Since(start).Milliseconds(),
	}).Infof("%s(%s) -> Task updated successfully", ctx.Request().Method, ctx.Path())
}

// auth
func AuthFailed(ctx echo.Context, operation string, err error) {
	logger.WithFields(logrus.Fields{
		"operation": operation,
		"error":     err.Error(),
	}).Warnf("Auth operation failed: %s", operation)
}

func AuthSuccess(ctx echo.Context, operation string, userID int) {
	logger.WithFields(logrus.Fields{
		"operation": operation,
		"user_id":   userID,
	}).Infof("Auth operation successful: %s", operation)
}
