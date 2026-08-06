package models

import "time"

type RefreshToken struct {
	Id        string    `db:"id"`
	UserId    string    `db:"user_id"`
	FamilyId  string    `db:"family_id"`
	IsRevoked bool      `db:"is_revoked"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}
