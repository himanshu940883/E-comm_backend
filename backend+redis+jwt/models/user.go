package models

type User struct {
	ID       int    `db:"id" json:"id"`
	UserID   string `db:"user_id" json:"user_id"`
	Password string `db:"password" json:"password"`
	Role     string `db:"role" json:"role"`
}
