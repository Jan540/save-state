package controllers

import (
	"encoding/json"
	"jan540/save-state/db"
	"jan540/save-state/filesystem"
	"jan540/save-state/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SaveController struct {
	db      *db.SaveDB
	storage *filesystem.SaveStorage
}

func NewSaveController(db *db.SaveDB, s *filesystem.SaveStorage) *SaveController {
	return &SaveController{
		db:      db,
		storage: s,
	}
}

func (sc *SaveController) GetSaveInfos(c echo.Context) error {
	userId := "1234-1234-1234-1234"

	saves, err := sc.db.GetSaves(userId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get saves 🫣")
	}

	return c.JSON(http.StatusOK, saves)
}

type PostSaveReq struct {
	GameCode string `param:"game_code"`
}

type PostSaveRes struct {
	Message string      `json:"message"`
	Save    models.Save `json:"save"`
}

func (sc *SaveController) PostSave(c echo.Context) error {
	userId := "1234-1234-1234-1234"

	var req PostSaveReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request 😥: "+err.Error())
	}

	saveFile, err := c.FormFile("save")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Couldn't read save file from form 🤷: "+err.Error())
	}

	rawMetadata := c.FormValue("metadata")
	if len(rawMetadata) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "No metadata found 😿")
	}

	var metadata models.SaveMetadata

	err = json.Unmarshal([]byte(rawMetadata), &metadata)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse metadata 🏷️: "+err.Error())
	}

	saveCount, err := sc.db.GetSaveCount(userId, req.GameCode)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get save count 🤷: "+err.Error())
	}

	if saveCount >= 10 {
		oldestSave, err := sc.db.GetOldestSave(userId, req.GameCode)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get oldest save 🤷: "+err.Error())
		}

		if err = sc.db.DeleteSave(oldestSave); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete oldest save from db 😢: "+err.Error())
		}

		if err = sc.storage.DeleteSave(oldestSave); err != nil {
			// TODO: revert db changes
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete oldest save from storage 😢: "+err.Error())
		}
	}

	if saveCount > 0 {
		currentSave, err := sc.db.GetCurrentSave(userId, req.GameCode)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get current save 🤷: "+err.Error())
		}

		backup, err := sc.storage.CreateBackup(currentSave)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create backup 😢: "+err.Error())
		}

		if err = sc.db.UpdateSave(backup); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update save in db 😢: "+err.Error())
		}
	}

	newSave := &models.Save{
		GameCode: req.GameCode,
		UserId:   userId,
		SaveTime: metadata.SaveTime,
	}

	if err = sc.db.CreateSave(newSave); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create save in db 🤷: "+err.Error())
	}

	if err = sc.storage.SaveSave(*newSave, saveFile); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save file in storage 🤷: "+err.Error())
	}

	return c.JSON(http.StatusOK, &PostSaveRes{
		Message: "Save successful 🎉",
		Save:    *newSave,
	})
}
