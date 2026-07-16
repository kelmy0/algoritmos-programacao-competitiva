package repositories

import (
	"context"
	"errors"
	"fmt"
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

func (r *AuthRepository) SaveRefreshToken(ctx context.Context, tokenId, userId string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, expires_at) VALUES 
		($1, $2, $3);
	`
	_, err := r.db.Exec(ctx, query, tokenId, userId, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to save refresh token (id: %s, user_id: %s): %w", tokenId, userId, err)
	}
	return nil
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
		return nil, fmt.Errorf("failed to query refresh token by id (%s): %w", id, err)
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
		return fmt.Errorf("failed to execute delete query for refresh token (id: %s, user_id: %s): %w", tokenId, userId, err)
	}

	if result.RowsAffected() == 0 {
		return models.ErrTokenNotFound
	}
	return nil
}

func (r *AuthRepository) DeleteAllRefreshToken(ctx context.Context, userId, tokenId string) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE user_id = $1 AND id != $2;
	`
	_, err := r.db.Exec(ctx, query, userId, tokenId)
	if err != nil {
		return fmt.Errorf("failed to revoke other refresh tokens for user (%s) keeping current token (%s): %w", userId, tokenId, err)
	}
	return nil
}
