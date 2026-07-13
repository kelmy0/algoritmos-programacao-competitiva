package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (r *UserRepository) GetUserByEmailForAuth(ctx context.Context, email string) (*models.User, error) {
	return r.getForAuth(ctx, email, "email")
}

func (r *UserRepository) GetUserByIdForAuth(ctx context.Context, id string) (*models.User, error) {
	return r.getForAuth(ctx, id, "id")
}

func (r *UserRepository) getForAuth(ctx context.Context, value, field string) (*models.User, error) {
	query := fmt.Sprintf(`
        SELECT 
            u.id, u.username, u.email, u.password_hash, u.enable, 
            u.two_factor_authentication, u.two_factor_secret, r.is_employee,
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
		&user.Id, &user.Username, &user.Email, &user.PasswordHash, &user.Enable,
		&user.TwoFactorAuthentication, &user.TwoFactorSecret, &user.Role.IsEmployee,
		&user.Permissions,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, fmt.Errorf("database query failed: %w", err)
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get auth data: %w", err)
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
		return fmt.Errorf("failed to save 2fa secret: %w", err)
	}

	if res.RowsAffected() == 0 {
		return models.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) Enable2FA(ctx context.Context, userId string) error {
	query := `
		UPDATE users 
		SET two_factor_authentication = TRUE
		WHERE id = $1;
	`
	res, err := r.db.Exec(ctx, query, userId)
	if err != nil {
		return fmt.Errorf("failed to enable 2fa: %w", err)
	}

	if res.RowsAffected() == 0 {
		return models.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) Disable2FA(ctx context.Context, userId string) error {
	query := `
		UPDATE users 
		SET two_factor_authentication = FALSE,
			two_factor_secret = NULL
		WHERE id = $1;
	`
	res, err := r.db.Exec(ctx, query, userId)
	if err != nil {
		return fmt.Errorf("failed to disable 2fa: %w", err)
	}
	if res.RowsAffected() == 0 {
		return models.ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) CheckUserExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, data models.NewUser) (string, error) {
	query := `
        INSERT INTO users(name, username, email, password_hash, role_id) VALUES
        ($1, $2, $3, $4, 1)
        RETURNING id;
    `
	var insertedId string

	err := r.db.QueryRow(ctx, query, data.Name, data.Username, data.Email, data.PasswordHash).Scan(&insertedId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return "", models.ErrUserAlreadyExists
			}
		}

		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return insertedId, nil
}

func (r *UserRepository) GetUserBySocialID(ctx context.Context, provider, socialId string) (*models.User, error) {
	query := `
        SELECT 
            u.id, u.username, u.email, u.password_hash, u.enable, 
            u.two_factor_authentication, u.two_factor_secret, r.is_employee,
            COALESCE(array_agg(p.slug) FILTER (WHERE p.slug IS NOT NULL), '{}') as permissions
        FROM users u
        INNER JOIN user_social_accounts usa ON u.id = usa.user_id
        INNER JOIN roles r ON u.role_id = r.id
        LEFT JOIN role_permissions rp ON r.id = rp.role_id
        LEFT JOIN permissions p ON rp.permission_id = p.id
        WHERE usa.provider = $1 AND usa.social_user_id = $2
        GROUP BY u.id, r.id
    `
	var user models.User
	err := r.db.QueryRow(ctx, query, provider, socialId).Scan(
		&user.Id, &user.Username, &user.Email, &user.PasswordHash, &user.Enable,
		&user.TwoFactorAuthentication, &user.TwoFactorSecret, &user.Role.IsEmployee,
		&user.Permissions,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to query social user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) CreateSocialUser(ctx context.Context, newUser models.NewUserGoogle, provider, socialId string) (*models.User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	queryUser := `
        INSERT INTO users (name, username, email, role_id) 
        VALUES ($1, $2, $3, 1) 
        RETURNING id
    `
	var userID string
	err = tx.QueryRow(ctx, queryUser, newUser.Name, newUser.Username, newUser.Email).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("error inserting user within transaction: %w", err)
	}

	querySocial := `
        INSERT INTO user_social_accounts (user_id, provider, social_user_id) 
        VALUES ($1, $2, $3)
    `
	_, err = tx.Exec(ctx, querySocial, userID, provider, socialId)
	if err != nil {
		return nil, fmt.Errorf("error inserting social account within transaction: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("error committing social transaction: %w", err)
	}

	return r.GetUserByIdForAuth(ctx, userID)
}

func (r *UserRepository) CreateSocialLink(ctx context.Context, id, provider, socialUserId string) error {
	query := `
        INSERT INTO user_social_accounts (user_id, provider, social_user_id)
        VALUES ($1, $2, $3)
    `
	_, err := r.db.Exec(ctx, query, id, provider, socialUserId)
	if err != nil {
		return fmt.Errorf("failed to link social account: %w", err)
	}

	return nil
}

func (r *UserRepository) ChangePassword(ctx context.Context, id, newPassword string) error {
	query := `
        UPDATE users
        SET password_hash = $1
        WHERE id = $2
    `

	res, err := r.db.Exec(ctx, query, newPassword, id)
	if err != nil {
		return fmt.Errorf("failed to change password: %w", err)
	}
	if res.RowsAffected() == 0 {
		return models.ErrUserNotFound
	}
	return nil
}
