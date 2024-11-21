package models

import "time"

type SaveMetadata struct {
	FileName string    `json:"filename,omitempty"`
	GameCode string    `json:"game_code"`
	SaveTime time.Time `json:"save_time"`
}

type Save struct {
	Id       string    `db:"saveId" json:"id"`
	GameCode string    `db:"gameCode" json:"game_code"`
	UserId   string    `db:"userId" json:"user_id"`
	SaveTime time.Time `db:"saveTime" json:"save_time"`
	IsBackup bool      `db:"isBackup" json:"is_backup"`
	Filename string    `db:"filename" json:"filename"`
}
