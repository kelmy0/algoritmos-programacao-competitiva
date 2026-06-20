package models

import (
	"time"
)

type Administrator struct {
	Id                     string     `db:"id"`
	Name                   string     `db:"name"`
	Email                  string     `db:"email"`
	PasswordHash           string     `db:"password_hash"`
	RecoveryToken          *string    `db:"recovery_token"`
	RecoveryTokenExpiresAt *time.Time `db:"recovery_token_expires_at"`
	//RefreshToken            string     `db:"refresh_token"`
	//refreshTokenExpiresAt   *time.Time `db:"refresh_token_expires_at"`
	Enable                  bool       `db:"enable"`
	TwoFactorAuthentication bool       `db:"two_factor_authentication"`
	TwoFactorSecret         *string    `db:"two_factor_secret"`
	RoleId                  int        `db:"role_id"`
	FailedAttempts          int        `db:"failed_attempts"`
	LastLogin               *time.Time `db:"last_login"`
	BlockedUntil            *time.Time `db:"blocked_until"`
	CreatedAt               time.Time  `db:"created_at"`
	UpdatedAt               time.Time  `db:"updated_at"`
}
