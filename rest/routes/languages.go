package routes

import (
	"net/http"

	"github.com/florianwoelki/kira/pkg"
	"github.com/labstack/echo/v4"
)

type languagesResponse struct {
	Languages []pkg.Language `json:"languages"`
}

func Languages(c echo.Context) error {
	languages, err := pkg.GetLanguages()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{
			Message: "Could not get messages.",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, languagesResponse{
		Languages: languages,
	})
}
