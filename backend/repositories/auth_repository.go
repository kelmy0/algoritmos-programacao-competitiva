package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT 
			u.id, u.name, u.username, u.email, u.password_hash, u.recovery_token, 
			u.recovery_token_expires_at, u.enable, u.two_factor_authentication, 
			u.two_factor_secret, u.role_id, u.failed_attempts, u.last_login, 
			u.blocked_until, u.created_at, u.updated_at, r.is_employee
		FROM users u
		INNER JOIN roles r ON u.role_id = r.id
		WHERE u.email = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.RecoveryToken,
		&user.RecoveryTokenExpiresAt,
		&user.Enable,
		&user.TwoFactorAuthentication,
		&user.TwoFactorSecret,
		&user.RoleId,
		&user.FailedAttempts,
		&user.LastLogin,
		&user.BlockedUntil,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Role.IsEmployee,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
