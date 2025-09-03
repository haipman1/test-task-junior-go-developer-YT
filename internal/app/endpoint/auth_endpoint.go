package endpoint

import (
	"context"
	"net/http"
	"testtaskYT/internal/app/auth"
	"testtaskYT/internal/app/models"
	"testtaskYT/internal/pkg/utils/logging"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// AuthService - интерфейс для сервиса аутентификации
type AuthService interface {
	Register(ctx context.Context, email, password string) (*models.UserResponse, string, error)
	Login(ctx context.Context, email, password string) (*models.UserResponse, string, error)
	ValidateToken(tokenString string) (*auth.Claims, error)
}

// AuthEndpoint - эндпоинты аутентификации
type AuthEndpoint struct {
	authService AuthService
	v           *validator.Validate
}

func NewAuthEndpoint(authService AuthService) *AuthEndpoint {
	return &AuthEndpoint{
		authService: authService,
		v:           validator.New(),
	}
}

// RegisterRequest - запрос на регистрацию
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginRequest - запрос на логин
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse - ответ с токеном
type AuthResponse struct {
	User  *models.UserResponse `json:"user"`
	Token string               `json:"token"`
}

func (e *AuthEndpoint) bindAndValidate(ctx echo.Context, structurePtr interface{}) error {
	// 1. Парсим JSON в структуру
	if err := ctx.Bind(structurePtr); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "Invalid JSON")
	}

	// 2. Валидируем по правилам из тегов validate
	if err := e.v.Struct(structurePtr); err != nil {
		// Преобразуем ошибки валидации в красивый формат
		if errs, ok := err.(validator.ValidationErrors); ok {
			errsMap := map[string]string{}
			for _, e := range errs {
				errsMap[e.Field()] = getValidationMessage(e)
			}
			return ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
				"error":   "Validation failed",
				"details": errsMap,
			})
		}
		return errorResponse(ctx, http.StatusInternalServerError, "Validation error")
	}

	return nil
}

// Register - регистрация пользователя
func (e *AuthEndpoint) Register(ctx echo.Context) error {
	logging.MetodCalled(ctx)

	var req RegisterRequest
	if err := e.bindAndValidate(ctx, &req); err != nil {
		return err
	}

	user, token, err := e.authService.Register(ctx.Request().Context(), req.Email, req.Password)
	if err != nil {
		logging.AuthFailed(ctx, "register", err)
		return errorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	logging.AuthSuccess(ctx, "register", user.ID)
	return ctx.JSON(http.StatusCreated, AuthResponse{
		User:  user,
		Token: token,
	})
}

// Login - аутентификация пользователя
func (e *AuthEndpoint) Login(ctx echo.Context) error {
	logging.MetodCalled(ctx)

	var req LoginRequest
	if err := e.bindAndValidate(ctx, &req); err != nil {
		return err
	}

	user, token, err := e.authService.Login(ctx.Request().Context(), req.Email, req.Password)
	if err != nil {
		logging.AuthFailed(ctx, "login", err)
		return errorResponse(ctx, http.StatusUnauthorized, "Invalid credentials")
	}

	logging.AuthSuccess(ctx, "login", user.ID)
	return ctx.JSON(http.StatusOK, AuthResponse{
		User:  user,
		Token: token,
	})
}
