package app

import (
	"database/sql"
	"fmt"
	"testtaskYT/internal/app/auth"
	"testtaskYT/internal/app/endpoint"
	"testtaskYT/internal/app/mw"
	"testtaskYT/internal/app/repository"
	service "testtaskYT/internal/app/service"
	"testtaskYT/internal/pkg/config"
	"testtaskYT/internal/pkg/utils/logging"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type App struct {
	e    *endpoint.Endpoint
	s    *service.Service
	echo *echo.Echo
	db   *sql.DB
}

func New() (*App, error) {
	a := &App{}

	cfg, err := config.LoadConfig()
	fmt.Printf("Config loaded: %+v\n", cfg) //!!!
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	logging.Init(cfg)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	db, err := repository.NewDB(dsn)
	if err != nil {
		return nil, err
	}

	a.db = db

	jwtService := auth.NewJWTService("supersecret")
	userRepo := repository.NewUserRepository(db) //TODO new db
	taskRepo := repository.NewTaskRepository(db)

	// Инициализируем сервисы
	authService := service.NewAuthService(userRepo, jwtService)
	taskService := service.New(taskRepo)

	// Инициализируем эндпоинты
	authEndpoint := endpoint.NewAuthEndpoint(authService)
	taskEndpoint := endpoint.New(taskService)

	a.s = service.New(taskRepo)
	a.e = endpoint.New(a.s)

	//a.logger.Info("App.go/ Layers was initialized")

	a.echo = echo.New()
	a.echo.POST("/register", authEndpoint.Register)
	a.echo.POST("/login", authEndpoint.Login)

	// Защищенные роуты
	authGroup := a.echo.Group("/tasks")
	authGroup.Use(mw.AuthMiddleware(authService)) // JWT аутентификация

	authGroup.POST("", taskEndpoint.CreateTask)
	authGroup.GET("", taskEndpoint.GetAllTasks)
	authGroup.GET("/:id", taskEndpoint.GetTask)
	authGroup.DELETE("/:id", taskEndpoint.DeleteTask)
	authGroup.PATCH("/:id", taskEndpoint.UpdateTaskPartial)

	// Админские роуты
	adminGroup := a.echo.Group("/admin")
	adminGroup.Use(mw.AuthMiddleware(authService))
	adminGroup.Use(mw.RoleMiddleware("admin")) // Только для админов

	//a.echo.Use(mw.RequestLogger)

	a.echo.GET("/", a.e.HelloWorld)
	a.echo.GET("/status", a.e.Status)
	a.echo.POST("/tasks", a.e.CreateTask)
	a.echo.GET("/tasks", a.e.GetAllTasks)
	a.echo.GET("/tasks/:id", a.e.GetTask)
	a.echo.DELETE("/tasks/:id", a.e.DeleteTask)
	a.echo.PATCH("/tasks/:id", a.e.UpdateTaskPartial)

	//a.logger.Info("App.go/ Routes was initialized")
	return a, nil
}

func (a *App) Run() error {
	cfg, _ := config.LoadConfig()

	//Таймауты сервера
	a.echo.Server.ReadTimeout = cfg.ReadTimeout
	a.echo.Server.WriteTimeout = cfg.WriteTimeout

	fmt.Printf("Server starting on %s...", cfg.ServerPort)

	err := a.echo.Start(cfg.ServerPort)
	if err != nil {
		return fmt.Errorf("Error in starting server: %w", err)
	}
	return nil
}
