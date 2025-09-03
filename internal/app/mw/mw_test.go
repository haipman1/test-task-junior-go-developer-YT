package mw

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// Проверка что Middleware вызывает next и возвращает его результат
// и что Header правильно извлекается и обрабатывается
func TestRoleCheck(t *testing.T) {
	testCases := []struct {
		name        string
		headerValue string
		expValue    string
	}{
		{"admin role", "admin", "admin"},
		{"user role", "user", "user"},
		{"empty role", "", "no-role"},
	}

	for _, cs := range testCases {
		t.Run(cs.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			req.Header.Set("User-Role", cs.expValue)

			contxt := echo.New().NewContext(req, rec)

			mw_next := func(ctx echo.Context) error {
				role := ctx.Request().Header.Get("User-Role")

				switch role {
				case "user":
					return ctx.String(http.StatusOK, "user")
				case "admin":
					return ctx.String(http.StatusOK, "admin")
				default:
					return ctx.String(http.StatusOK, "no-role")
				}
			}
			err := RoleCheck(mw_next)(contxt)

			assert.NoError(t, err)
			assert.Equal(t, rec.Code, http.StatusOK)
			assert.Equal(t, rec.Body.String(), cs.expValue)
		})
	}
}
