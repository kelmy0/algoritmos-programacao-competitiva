package models

import "time"

type RefreshToken struct {
	Id        string    `db:"id"`
	UserId    string    `db:"user_id"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}
