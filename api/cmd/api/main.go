package main

import (
	"jan540/save-state/controllers"
	"jan540/save-state/db"
	"jan540/save-state/filesystem"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func main() {
	godotenv.Load()

	authSecret := os.Getenv("AUTH_SECRET")

	e := echo.New()

	// e.Use(middleware.Logger())

	db, err := db.InitDB(os.Getenv("DB_FILE"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	storage := filesystem.NewSaveStorage(os.Getenv("SAVE_DIRECTORY"))

	authController := controllers.NewAuthController(db, authSecret)
	saveController := controllers.NewSaveController(db, storage)

	e.POST("/login", authController.Login)
	e.POST("/register", authController.Register)

	p := e.Group("")

	p.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(authSecret),
	}))

	p.GET("/info", saveController.GetSaveInfos)
	p.GET("/sync/:game_code", saveController.GetSave)
	p.POST("/sync/:game_code", saveController.PostSave)
	p.GET("/test", func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		userId, _ := user.Claims.GetSubject()

		return c.JSON(200, userId)
	})

	e.Logger.Fatal(e.Start(":6969"))
}
