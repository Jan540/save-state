package filesystem

import (
	"fmt"
	"io"
	"jan540/save-state/models"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type SaveStorage struct {
	baseDir string
}

func NewSaveStorage(d string) *SaveStorage {
	return &SaveStorage{
		baseDir: d,
	}
}

func (ss *SaveStorage) ListFiles() ([]models.SaveFile, error) {
	filesRaw, err := os.ReadDir(ss.baseDir)

	if err != nil {
		return nil, err
	}

	files := make([]models.SaveFile, 0)

	for _, file := range filesRaw {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()

		if err != nil {
			return nil, err
		}

		file := models.SaveFile{
			Name:    info.Name(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		}

		files = append(files, file)
	}

	return files, nil
}

func (ss *SaveStorage) SaveSaves(userId string, saves []*multipart.FileHeader) error {
	for _, save := range saves {
		if err := ss.SaveSave(userId, save); err != nil {
			return err
		}
	}

	return nil
}

func (ss *SaveStorage) SaveSave(userId string, save *multipart.FileHeader) error {
	gameCode := save.Filename

	src, err := save.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	gameDir := filepath.Join(ss.baseDir, userId, gameCode)

	if err := os.MkdirAll(gameDir, os.ModePerm); err != nil {
		return err
	}

	dstPath := filepath.Join(gameDir, "current.sav")

	if _, err := os.Stat(dstPath); err == nil {
		backupFileName := fmt.Sprintf("backup%s.sav", time.Now().Format("060102-150405"))

		backupPath := filepath.Join(gameDir, backupFileName)

		if err := os.Rename(dstPath, backupPath); err != nil {
			return err
		}

		files, err := os.ReadDir(gameDir)
		if err != nil {
			return err
		}

		if len(files) > 5 {
			if err := ss.deleteOldestSave(gameDir, files); err != nil {
				return err
			}
		}
	}

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

func (ss *SaveStorage) deleteOldestSave(gameDir string, files []os.DirEntry) error {
	var oldest os.FileInfo

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			return err
		}

		if oldest == nil {
			oldest = info
			continue
		}

		fileInfo, _ := file.Info()

		if fileInfo.ModTime().Before(oldest.ModTime()) {
			oldest = info
		}
	}

	err := os.Remove(filepath.Join(gameDir, oldest.Name()))
	if err != nil {
		return err
	}

	return nil
}
