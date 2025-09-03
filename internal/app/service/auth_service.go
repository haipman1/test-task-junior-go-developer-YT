package service

import (
	"context"
	"fmt"
	"testtaskYT/internal/app/auth"
	"testtaskYT/internal/app/models"
	"testtaskYT/internal/app/repository"
)

type AuthService struct {
	userRepo   *repository.UserRepository
	jwtService *auth.JWTService
}

func NewAuthService(userRepo *repository.UserRepository, jwtService *auth.JWTService) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (*models.UserResponse, string, error) {
	existingUser, err := s.userRepo.GetUserByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, "", fmt.Errorf("user already exists")
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.userRepo.CreateUser(ctx, email, hashedPassword, "user")
	if err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	token, err := s.jwtService.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}
	return &models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, token, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.UserResponse, string, error) {
	fmt.Printf("DEBUG: Login attempt - email: %s, password: %s\n", email, password)

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		fmt.Printf("DEBUG: User not found - %v\n", err)
		return nil, "", fmt.Errorf("invalid credentials")
	}

	fmt.Printf("DEBUG: User found - ID: %d, Email: %s, PasswordHash: %s\n", user.ID, user.Email, user.Password)
	fmt.Printf("DEBUG: Input password: %s\n", password)
	fmt.Printf("DEBUG: Stored hash: %s\n", user.Password)

	// Проверяем пароль
	passwordMatch := auth.CheckPassword(password, user.Password)
	fmt.Printf("DEBUG: Password match: %v\n", passwordMatch)

	if !passwordMatch {
		fmt.Printf("DEBUG: Password mismatch!\n")
		return nil, "", fmt.Errorf("invalid credentials")
	}

	fmt.Printf("DEBUG: Login successful!\n")

	// Генерируем JWT токен
	token, err := s.jwtService.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, token, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*auth.Claims, error) {
	return s.jwtService.ValidateToken(tokenString)
}
