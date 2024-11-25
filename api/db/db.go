package db

import (
	"jan540/save-state/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type SaveDB struct {
	conn *sqlx.DB
}

func InitDB(dbFile string) (*SaveDB, error) {
	db, err := sqlx.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS saves(
			saveId INTEGER PRIMARY KEY AUTOINCREMENT,
			gameCode VARCHAR(4) NOT NULL,
			userId VARCHAR(255) NOT NULL,
			saveTime DATETIME NOT NULL,
			isBackup BOOLEAN NOT NULL DEFAULT FALSE,
			filename VARCHAR(255) NOT NULL DEFAULT 'current.sav',
			createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(gameCode, userId, saveTime)
		);`)

	return &SaveDB{conn: db}, err
}

func (sdb *SaveDB) Close() error {
	return sdb.conn.Close()
}

func (sdb *SaveDB) CreateSave(save *models.Save) error {
	_, err := sdb.conn.Exec(`
		INSERT INTO saves (gameCode, userId, saveTime) VALUES ($1, $2, $3);`,
		save.GameCode,
		save.UserId,
		save.SaveTime)

	if err != nil {
		return err
	}

	currentSave, err := sdb.GetCurrentSave(save.UserId, save.GameCode)
	if err != nil {
		return err
	}

	*save = currentSave

	return nil
}

func (sdb *SaveDB) DeleteSave(save models.Save) error {
	_, err := sdb.conn.Exec(`
		DELETE FROM saves
		WHERE saveId=$1;`,
		save.Id)

	return err
}

func (sdb *SaveDB) GetSaves(userId string) ([]models.Save, error) {
	var saves []models.Save

	err := sdb.conn.Select(&saves, `
		SELECT * FROM saves
		WHERE userId=$1
		AND isBackup=false;`,
		userId)

	return saves, err
}

func (sdb *SaveDB) GetCurrentSave(userId string, gameCode string) (models.Save, error) {
	var currentSave models.Save

	err := sdb.conn.Get(&currentSave, `
		SELECT * FROM saves
		WHERE userId=$1
		AND gameCode=$2
		AND isBackup=false
		LIMIT 1;`,
		userId,
		gameCode)

	return currentSave, err
}

func (sdb *SaveDB) GetOldestSave(userId string, gameCode string) (models.Save, error) {
	var oldestSave models.Save

	err := sdb.conn.Get(&oldestSave, `
		SELECT * FROM saves
		WHERE userId=$1
		AND gameCode=$2
		ORDER BY saveTime ASC
		LIMIT 1;`,
		userId,
		gameCode)

	return oldestSave, err
}

func (sdb *SaveDB) GetSaveCount(userId string, gameCode string) (int, error) {
	var count int

	err := sdb.conn.Get(&count, `
		SELECT count(*) FROM saves
		WHERE userId=$1
		AND gameCode=$2;`,
		userId,
		gameCode)

	return count, err
}

func (sdb *SaveDB) UpdateSave(save models.Save) error {
	_, err := sdb.conn.Exec(`
		UPDATE saves
		SET isBackup=$1, filename=$2
		WHERE saveId=$3`,
		save.IsBackup,
		save.Filename,
		save.Id)

	return err
}
