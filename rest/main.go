package main

import (
	"log"
	"os"

	"github.com/florianwoelki/kira/pkg"
	"github.com/florianwoelki/kira/rest/routes"
	"github.com/labstack/echo/v4"
)

var logger *log.Logger = log.New(os.Stdout, "api: ", log.LstdFlags|log.Lshortfile)
var local = false // TODO: Change to env variable.

func main() {
	if !local {
		err := pkg.CreateRunners()
		if err != nil {
			logger.Fatalf("Error while trying to create runners: %v+", err)
		}

		err = pkg.CreateUsers()
		if err != nil {
			logger.Fatalf("Error while trying to create users: %v+", err)
		}
	}

	err := pkg.LoadLanguages()
	if err != nil {
		logger.Fatalf("Error while loading languages: %v+", err)
	}

	e := echo.New()
	e.GET("/languages", routes.Languages)
	e.POST("/execute", routes.Execute)

	e.Logger.Fatal(e.Start(":9090"))
}
