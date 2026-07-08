package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Save2FASecret(ctx context.Context, userId, secret string) error {
	query := `
		UPDATE users 
		SET two_factor_secret = $1
		WHERE id = $2;
	`
	_, err := r.db.Exec(ctx, query, secret, userId)
	return err
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
