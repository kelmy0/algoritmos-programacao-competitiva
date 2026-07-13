package models

import (
	"time"
)

type User struct {
	Id                      string              `db:"id"`
	Name                    string              `db:"name"`
	Username                string              `db:"username"`
	Email                   string              `db:"email"`
	PasswordHash            *string             `db:"password_hash"`
	RecoveryTokenHash       *string             `db:"recovery_token_hash"`
	RecoveryTokenExpiresAt  *time.Time          `db:"recovery_token_expires_at"`
	Enable                  bool                `db:"enable"`
	TwoFactorAuthentication bool                `db:"two_factor_authentication"`
	TwoFactorSecret         *string             `db:"two_factor_secret"`
	RoleId                  int                 `db:"role_id"`
	FailedAttempts          int                 `db:"failed_attempts"`
	LastLogin               *time.Time          `db:"last_login"`
	BlockedUntil            *time.Time          `db:"blocked_until"`
	CreatedAt               time.Time           `db:"created_at"`
	UpdatedAt               time.Time           `db:"updated_at"`
	Role                    Role                `db:"role"`
	Permissions             []string            `db:"permissions"`
	SocialAccounts          []UserSocialAccount `db:"-"`
}

type NewUser struct {
	Name         string `db:"name"`
	Username     string `db:"username"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
}

type UserSocialAccount struct {
	Id           string    `db:"id"`
	UserId       string    `db:"user_id"`
	Provider     string    `db:"provider"`
	SocialUserId string    `db:"social_user_id"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type NewUserGoogle struct {
	Name         string `db:"name"`
	Username     string `db:"username"`
	Email        string `db:"email"`
	Provider     string `db:"provider"`
	SocialUserId string `db:"social_user_id"`
}
