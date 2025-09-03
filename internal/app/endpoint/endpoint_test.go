package endpoint

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type MockService struct {
	DaysCt int64
}

func (ms *MockService) DaysCount() int64 {
	return ms.DaysCt
}

func TestStatus(t *testing.T) { //Mock testing
	endpoint := New(&MockService{DaysCt: 101})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/status", nil)
	contxt := echo.New().NewContext(req, rec)

	err := endpoint.Status(contxt) //TODO mock this
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Осталось дней: 101")
}
