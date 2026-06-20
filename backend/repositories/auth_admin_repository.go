package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
)

type AuthAdminRepository struct {
	db *pgxpool.Pool
}

func NewAuthAdminRepository(db *pgxpool.Pool) *AuthAdminRepository {
	return &AuthAdminRepository{db: db}
}

func (r *AuthAdminRepository) GetAdminByEmail(ctx context.Context, email string) (*models.Administrator, error) {
	query := `
		SELECT id, name, email, password_hash, recovery_token,  recovery_token_expires_at, enable, 
            two_factor_authentication, two_factor_secret, role_id, failed_attempts, 
            last_login, blocked_until, created_at, updated_at
		FROM administrators
		WHERE email = $1
	`

	var admin models.Administrator
	err := r.db.QueryRow(ctx, query, email).Scan(
		&admin.Id,
		&admin.Name,
		&admin.Email,
		&admin.PasswordHash,
		&admin.RecoveryToken,
		&admin.RecoveryTokenExpiresAt,
		&admin.Enable,
		&admin.TwoFactorAuthentication,
		&admin.TwoFactorSecret,
		&admin.RoleId,
		&admin.FailedAttempts,
		&admin.LastLogin,
		&admin.BlockedUntil,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &admin, nil
}
