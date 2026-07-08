package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

type User2FAData struct {
	IsEnabled bool
	Secret    string
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Get2FAData(ctx context.Context, userId string) (*User2FAData, error) {
	query := `
        SELECT two_factor_authentication, COALESCE(two_factor_secret, '') 
        FROM users
        WHERE id = $1;`

	var data User2FAData
	err := r.db.QueryRow(ctx, query, userId).Scan(&data.IsEnabled, &data.Secret)
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
