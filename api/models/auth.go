package models

type User struct {
	UserId    string `db:"userId" json:"user_id"`
	CreatedAt string `db:"createdAt" json:"created_at"`
	Username  string `db:"username" json:"username"`
	Password  string `db:"password" json:"password"`
}
