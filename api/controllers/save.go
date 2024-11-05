package controllers

import (
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

func (sc *SaveController) GetSavedFiles(c echo.Context) error {
	files, err := sc.storage.ListFiles()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, files)
}
