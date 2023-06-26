package routes

import "github.com/labstack/echo/v4"

func CheckHealth(c echo.Context) error {
	status := map[string]string{
		"status": "ok",
	}
	return c.JSON(200, status)
}
