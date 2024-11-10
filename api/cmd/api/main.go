package main

import (
	"jan540/save-state/controllers"
	"jan540/save-state/filesystem"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	godotenv.Load()

	e := echo.New()

	// e.Use(middleware.Logger())

	clerk.SetKey(os.Getenv("CLERK_SECRET_KEY"))

	storage := filesystem.NewSaveStorage(os.Getenv("SAVE_DIRECTORY"))
	saveController := controllers.NewSaveController(*storage)

	e.GET("/", saveController.GetSaves)
	e.POST("/sync", saveController.PostSaves)

	e.Logger.Fatal(e.Start(":6969"))
}
