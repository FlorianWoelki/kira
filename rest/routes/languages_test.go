package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/florianwoelki/kira/pkg"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestLanguages(t *testing.T) {
	// Mock loaded languages for accepting valid return value for `pkg.GetLanguages()`.
	pkg.LoadedLanguages = map[string]pkg.Language{
		"go":  {Name: "Go", Version: "1.15.6", Extension: ".go", Timeout: 5},
		"cpp": {Name: "Cpp", Version: "10.2.1", Extension: ".cpp", Timeout: 5, Compiled: true},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/languages", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, Languages(c))
	require.Equal(t, http.StatusOK, rec.Code)

	expected := `{"languages":[{"compiled":false,"name":"Go","version":"1.15.6","extension":".go","timeout":5},{"compiled":true,"name":"Cpp","version":"10.2.1","extension":".cpp","timeout":5}]}`
	require.JSONEq(t, expected, rec.Body.String())
}
