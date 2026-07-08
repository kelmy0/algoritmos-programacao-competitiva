package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

type UserAuthData struct {
	IsEnabled    bool
	Secret       string
	PasswordHash string
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return r.getUserBy(ctx, email, "email")
}

func (r *UserRepository) GetUserById(ctx context.Context, id string) (*models.User, error) {
	return r.getUserBy(ctx, id, "id")
}

func (r *UserRepository) getUserBy(ctx context.Context, value, field string) (*models.User, error) {
	query := fmt.Sprintf(`
        SELECT 
            u.id, u.name, u.username, u.email, u.password_hash, u.recovery_token_hash, 
            u.recovery_token_expires_at, u.enable, u.two_factor_authentication, 
            u.two_factor_secret, u.role_id, u.failed_attempts, u.last_login, 
            u.blocked_until, u.created_at, u.updated_at, r.is_employee,
            COALESCE(array_agg(p.slug) FILTER (WHERE p.slug IS NOT NULL), '{}') as permissions
        FROM users u
        INNER JOIN roles r ON u.role_id = r.id
        LEFT JOIN role_permissions rp ON r.id = rp.role_id
        LEFT JOIN permissions p ON rp.permission_id = p.id
        WHERE u.%s = $1
        GROUP BY u.id, r.id
    `, field)

	var user models.User
	err := r.db.QueryRow(ctx, query, value).Scan(
		&user.Id, &user.Name, &user.Username, &user.Email, &user.PasswordHash, &user.RecoveryTokenHash,
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

func (r *UserRepository) GetAuthData(ctx context.Context, userId string) (*UserAuthData, error) {
	query := `
        SELECT two_factor_authentication, COALESCE(two_factor_secret, ''), password_hash
        FROM users
        WHERE id = $1;`

	var data UserAuthData
	err := r.db.QueryRow(ctx, query, userId).Scan(&data.IsEnabled, &data.Secret, &data.PasswordHash)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *UserRepository) Save2FASecret(ctx context.Context, userId, secret string) error {
	query := `
		UPDATE users 
		SET two_factor_secret = $1
		WHERE id = $2;
	`
	res, err := r.db.Exec(ctx, query, secret, userId)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("User not found")
	}

	return nil
}

func (r *UserRepository) Enable2FA(ctx context.Context, userId string) error {
	query := `
		UPDATE users 
		SET two_factor_authentication = TRUE
		WHERE id = $1;
	`
	_, err := r.db.Exec(ctx, query, userId)
	return err
}

func (r *UserRepository) Disable2FA(ctx context.Context, userId string) error {
	query := `
		UPDATE users 
		SET two_factor_authentication = FALSE,
			two_factor_secret = NULL
		WHERE id = $1;
	`
	_, err := r.db.Exec(ctx, query, userId)
	return err
}
