package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestHealth(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, CheckHealth(c))
	require.Equal(t, http.StatusOK, rec.Code)

	expected := `{"status":"ok"}`
	require.JSONEq(t, expected, rec.Body.String())
}
