package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/florianwoelki/kira/internal"
	"github.com/florianwoelki/kira/internal/pool"
	"github.com/labstack/echo/v4"
)

type executeBody struct {
	Language string            `json:"language" binding:"required"`
	Content  string            `json:"content" binding:"required"`
	Stdin    []string          `json:"stdin,omitempty"`
	Tests    []pool.TestResult `json:"tests,omitempty"`
}

type executeResponse struct {
	CompileOutput pool.Output     `json:"compileOutput"`
	RunOutput     pool.Output     `json:"runOutput"`
	TestOutput    pool.TestOutput `json:"testOutput,omitempty"`
}

func Execute(c echo.Context, rceEngine *internal.RceEngine) error {
	// Setting default values so that the optional fields are not empty.
	body := executeBody{
		Tests: []pool.TestResult{},
	}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	bypassCacheStr := c.QueryParam("bypass_cache")
	bypassCache, err := strconv.ParseBool(bypassCacheStr)
	if len(bypassCacheStr) == 0 {
		bypassCache = false
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	output, err := rceEngine.Dispatch(body.Language, body.Content, body.Stdin, body.Tests, bypassCache)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, executeResponse{
		CompileOutput: output.CompileOutput,
		RunOutput:     output.RunOutput,
		TestOutput:    output.TestOutput,
	})
	return nil
}
