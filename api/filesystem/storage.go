package filesystem

import (
	"jan540/save-state/models"
	"os"
)

type SaveStorage struct {
	directory string
}

func NewSaveStorage(d string) *SaveStorage {
	return &SaveStorage{
		directory: d,
	}
}

func (ss *SaveStorage) ListFiles() ([]models.SaveFile, error) {
	filesRaw, err := os.ReadDir(ss.directory)

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
