package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"testtaskYT/internal/app/endpoint"
	"testtaskYT/internal/app/service"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// Integration Tests

func TestApp_Initialization(t *testing.T) {
	a, err := New()

	assert.NoError(t, err)
	assert.NotNil(t, a)
	assert.NotNil(t, a.e)
	assert.NotNil(t, a.s)
	assert.NotNil(t, a.echo)
}

func TestApp_Routes(t *testing.T) {
	app, _ := New()

	rt := app.echo.Routes()
	existRoutes := make(map[string]string)

	for _, i := range rt {
		existRoutes[i.Path] = i.Name
	}

	supposedRoutes := []string{"/", "/status"}
	for _, i := range supposedRoutes {
		//if _, exist := existRoutes[i]; !exist {
		//	t.Errorf("Path \"%s\" doesn't exist", i)
		//	assert.Errorf("Path \"%s\" doesn't exist", i)
		//}

		_, exist2 := existRoutes[i]
		assert.True(t, exist2)
	}
}

func TestApp_CorrectShutdown(t *testing.T) {
	app, errApp := New()
	assert.NoError(t, errApp)

	go func() {
		errRun := app.Run() //blocking operation //TODO: how to use WG or someth here?
		if errRun != nil && errRun != http.ErrServerClosed {
			assert.Fail(t, "Unexpected app.Run() shutdown with error: ", errRun)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	errShutdown := app.echo.Shutdown(context.Background())
	assert.NoError(t, errShutdown, "Server shutdown with error!")
}

func TestApp_Endpoints_Requests(t *testing.T) {
	e := echo.New()
	service := service.New()
	endpoint := endpoint.New(service)

	e.GET("/", endpoint.HelloWorld)
	e.GET("/status", endpoint.Status)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, rec.Code, http.StatusOK)
	assert.Equal(t, rec.Body.String(), "Hello world!")

	req = httptest.NewRequest("GET", "/status", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, rec.Code, http.StatusOK)
	assert.Contains(t, rec.Body.String(), "Осталось дней: ")
}
