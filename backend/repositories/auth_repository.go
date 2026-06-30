package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
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
			u.blocked_until, u.created_at, u.updated_at, r.is_employee,
			COALESCE(array_agg(p.slug) FILTER (WHERE p.slug IS NOT NULL), '{}') as permissions
		FROM users u
		INNER JOIN roles r ON u.role_id = r.id
		LEFT JOIN role_permissions rp ON r.id = rp.role_id
        LEFT JOIN permissions p ON rp.permission_id = p.id
		WHERE u.email = $1
		GROUP BY u.id, r.id
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.Id, &user.Name, &user.Username, &user.Email, &user.PasswordHash, &user.RecoveryToken,
		&user.RecoveryTokenExpiresAt, &user.Enable, &user.TwoFactorAuthentication,
		&user.TwoFactorSecret, &user.Role.Id, &user.FailedAttempts, &user.LastLogin,
		&user.BlockedUntil, &user.CreatedAt, &user.UpdatedAt,
		&user.Role.IsEmployee,
		&user.Permissions,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) SaveRefreshToken(ctx context.Context, tokenId, userId string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, expires_at) VALUES 
		($1, $2, $3);
	`
	_, err := r.db.Exec(ctx, query, tokenId, userId, expiresAt)
	return err
}

func (r *AuthRepository) GetRefreshTokenById(ctx context.Context, id string) (*models.RefreshToken, error) {
	query := `
        SELECT id, user_id, expires_at, created_at 
        FROM refresh_tokens 
        WHERE id = $1;
    `

	var token models.RefreshToken
	err := r.db.QueryRow(ctx, query, id).Scan(
		&token.Id,
		&token.UserId,
		&token.ExpiresAt,
		&token.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &token, nil
}

func (r *AuthRepository) DeleteRefreshTokenById(ctx context.Context, userId, tokenId string) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE id = $1 AND user_id = $2;
	`
	result, err := r.db.Exec(ctx, query, tokenId, userId)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("refresh token not found or already revoked")
	}
	return nil
}

func (r *AuthRepository) DeleteAllRefreshToken(ctx context.Context, userId, tokenId string) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE user_id = $1 AND id != $2;
	`
	_, err := r.db.Exec(ctx, query, userId, tokenId)
	return err
}
