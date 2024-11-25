package filesystem

import (
	"fmt"
	"io"
	"jan540/save-state/models"
	"mime/multipart"
	"os"
	"path/filepath"
)

type SaveStorage struct {
	baseDir string
}

func NewSaveStorage(d string) *SaveStorage {
	return &SaveStorage{
		baseDir: d,
	}
}

func (ss *SaveStorage) GetSavePath(userId string, gameCode string) (string, error) {
	path := filepath.Join(ss.baseDir, userId, gameCode, "current.sav")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("Save not found")
	}

	return path, nil
}

func (ss *SaveStorage) SaveSave(data models.Save, file *multipart.FileHeader) error {
	gameDir := filepath.Join(ss.baseDir, data.UserId, data.GameCode)

	if err := os.MkdirAll(gameDir, os.ModePerm); err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dstPath := filepath.Join(gameDir, "current.sav")

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

func (ss *SaveStorage) DeleteSave(save models.Save) error {
	path := filepath.Join(ss.baseDir, save.UserId, save.GameCode, save.Filename)
	err := os.Remove(path)
	return err
}

func (ss *SaveStorage) CreateBackup(save models.Save) (models.Save, error) {
	gameDir := filepath.Join(ss.baseDir, save.UserId, save.GameCode)
	backupFilename := fmt.Sprintf("backup_%d", save.SaveTime.Unix())

	oldPath := filepath.Join(gameDir, save.Filename)
	newPath := filepath.Join(gameDir, backupFilename)
	save.Filename = backupFilename
	save.IsBackup = true

	err := os.Rename(oldPath, newPath)
	return save, err
}
