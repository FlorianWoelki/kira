package routes

import (
	"net/http"

	"github.com/florianwoelki/kira/internal"
	"github.com/labstack/echo/v4"
)

type languagesResponse struct {
	Languages []internal.Language `json:"languages"`
}

func Languages(c echo.Context) error {
	languages, err := internal.GetLanguages()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, languagesResponse{
		Languages: languages,
	})
}
