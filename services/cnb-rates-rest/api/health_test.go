package api

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	t.Log("healthy")
	{
		router := echo.New()
		router.HEAD("/health", HealtCheck())

		req := httptest.NewRequest(http.MethodHead, "/health", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	}

}
