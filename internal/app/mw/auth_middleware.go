package mw

import (
	"net/http"
	"strings"
	"testtaskYT/internal/app/service"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(authService *service.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Authorization header required",
				})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid authorization format. Expected: Bearer <token>",
				})
			}

			tokenString := parts[1]

			claims, err := authService.ValidateToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid token",
				})
			}

			c.Set("userID", claims.UserID)
			c.Set("userEmail", claims.Email)
			c.Set("userRole", claims.Role)

			return next(c)
		}
	}
}

func RoleMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole := c.Get("userRole").(string)

			for _, allowedRole := range allowedRoles {
				if userRole == allowedRole {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Access denied. Insufficient permissions.",
			})
		}
	}
}
