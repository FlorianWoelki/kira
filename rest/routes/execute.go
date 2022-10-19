package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/florianwoelki/kira/internal"
	"github.com/labstack/echo/v4"
)

type executeBody struct {
	Language string `json:"language" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type executeResponse struct {
	CompileOutput string `json:"compileOutput"`
	CompileError  string `json:"compileError"`
	CompileTime   int64  `json:"compileTime"`
	RunOutput     string `json:"runOutput"`
	RunError      string `json:"runError"`
	RunTime       int64  `json:"runTime"`
}

func Execute(c echo.Context, rceEngine *internal.RceEngine) error {
	body := executeBody{}

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

	output, err := rceEngine.Dispatch(body.Language, body.Content, bypassCache)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, executeResponse{
		CompileOutput: output.CompileResult,
		CompileError:  output.CompileError,
		CompileTime:   output.CompileTime.Milliseconds(),
		RunOutput:     output.RunResult,
		RunError:      output.RunError,
		RunTime:       output.RunTime.Milliseconds(),
	})
	return nil
}
