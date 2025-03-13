package userschema

import (
	"time"
)

type UserModel struct {
	ID        int64     `json:"id" db:"id"`
	Login     string    `json:"login" db:"login"`
	Password  string    `json:"password_hash" db:"password_hash"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
