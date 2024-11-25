package main

import (
	"jan540/save-state/auth"
	"jan540/save-state/controllers"
	"jan540/save-state/db"
	"jan540/save-state/filesystem"
	"net/http"
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
	e.Use(auth.ClerkMiddleware())

	db, err := db.InitDB(os.Getenv("DB_FILE"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	storage := filesystem.NewSaveStorage(os.Getenv("SAVE_DIRECTORY"))

	saveController := controllers.NewSaveController(db, storage)

	e.GET("/info", saveController.GetSaveInfos)
	e.GET("/sync/:game_code", saveController.GetSave)
	e.POST("/sync/:game_code", saveController.PostSave)

	e.GET("test", func(c echo.Context) error {
		userId := c.Get("userId").(string)

		return c.JSON(http.StatusOK, userId)
	})

	e.Logger.Fatal(e.Start(":6969"))
}
