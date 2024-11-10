package controllers

import (
	"fmt"
	"jan540/save-state/filesystem"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SaveController struct {
	storage filesystem.SaveStorage
}

func NewSaveController(s filesystem.SaveStorage) *SaveController {
	return &SaveController{
		storage: s,
	}
}

func (sc *SaveController) GetSaves(c echo.Context) error {
	files, err := sc.storage.ListFiles()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, files)
}

type PostSaveRes struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// gameCode is the name of the save file
func (sc *SaveController) PostSaves(c echo.Context) error {
	userId := "1234-1234-1234-1234"

	form, formErr := c.MultipartForm()
	if formErr != nil {
		return c.JSON(http.StatusInternalServerError, &PostSaveRes{
			Success: false,
			Message: "Failed to parse form 😥",
		})
	}

	saves := form.File["saves"]

	storageErr := sc.storage.SaveSaves(userId, saves)
	if storageErr != nil {
		return c.JSON(http.StatusInternalServerError, &PostSaveRes{
			Success: false,
			Message: fmt.Sprintf("Failed to save files: %s", storageErr.Error()),
		})
	}

	return c.JSON(http.StatusOK, &PostSaveRes{
		Success: true,
		Message: "Synced Saves ✅",
	})
}
