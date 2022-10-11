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
	Output string `json:"output"`
	Error  string `json:"error"`
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
		Output: output.Result,
		Error:  output.Error,
	})
	return nil
}
