package models

import (
	"time"
)

type User struct {
	Id                      string     `db:"id"`
	Name                    string     `db:"name"`
	Username                string     `db:"username"`
	Email                   string     `db:"email"`
	PasswordHash            *string    `db:"password_hash"`
	SsoProvider             *string    `db:"sso_provider"`
	SsoUserId               *string    `db:"sso_user_id"`
	RecoveryTokenHash       *string    `db:"recovery_token_hash"`
	RecoveryTokenExpiresAt  *time.Time `db:"recovery_token_expires_at"`
	Enable                  bool       `db:"enable"`
	TwoFactorAuthentication bool       `db:"two_factor_authentication"`
	TwoFactorSecret         *string    `db:"two_factor_secret"`
	RoleId                  int        `db:"role_id"`
	FailedAttempts          int        `db:"failed_attempts"`
	LastLogin               *time.Time `db:"last_login"`
	BlockedUntil            *time.Time `db:"blocked_until"`
	CreatedAt               time.Time  `db:"created_at"`
	UpdatedAt               time.Time  `db:"updated_at"`
	Role                    Role       `db:"role"`
	Permissions             []string   `db:"permissions"`
}
