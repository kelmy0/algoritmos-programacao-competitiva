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

func (r *AuthRepository) SaveRefreshToken(ctx context.Context, tokenId, userId, familyId string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, family_id, expires_at) VALUES 
		($1, $2, $3, $4);
	`
	_, err := r.db.Exec(ctx, query, tokenId, userId, familyId, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to save refresh token (id: %s, user_id: %s): %w", tokenId, userId, err)
	}
	return nil
}

func (r *AuthRepository) GetRefreshTokenById(ctx context.Context, id string) (*models.RefreshToken, error) {
	query := `
        SELECT id, user_id, family_id, is_revoked, expires_at, created_at
        FROM refresh_tokens 
        WHERE id = $1;
    `

	var token models.RefreshToken
	err := r.db.QueryRow(ctx, query, id).Scan(
		&token.Id,
		&token.UserId,
		&token.FamilyId,
		&token.IsRevoked,
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

func (r *AuthRepository) RotateRefreshToken(ctx context.Context, oldTokenId, newTokenId, userId, familyId string, newExpiresAt time.Time) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction for token rotation: %w", err)
	}
	defer tx.Rollback(ctx)

	revokeQuery := `
		UPDATE refresh_tokens 
		SET is_revoked = TRUE 
		WHERE id = $1 AND user_id = $2;
	`

	res, err := tx.Exec(ctx, revokeQuery, oldTokenId, userId)
	if err != nil {
		return fmt.Errorf("failed to revoke old token (%s) in rotation: %w", oldTokenId, err)
	}
	if res.RowsAffected() == 0 {
		return models.ErrTokenNotFound
	}
	insertQuery := `
		INSERT INTO refresh_tokens (id, user_id, family_id, expires_at) 
		VALUES ($1, $2, $3, $4);
	`
	_, err = tx.Exec(ctx, insertQuery, newTokenId, userId, familyId, newExpiresAt)
	if err != nil {
		return fmt.Errorf("failed to insert new rotated token (%s): %w", newTokenId, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit token rotation transaction: %w", err)
	}

	return nil
}

func (r *AuthRepository) RevokeFamily(ctx context.Context, familyId string) error {
	query := `
		UPDATE refresh_tokens 
		SET is_revoked = TRUE 
		WHERE family_id = $1;
	`
	_, err := r.db.Exec(ctx, query, familyId)
	if err != nil {
		return fmt.Errorf("failed to revoke entire token family (%s): %w", familyId, err)
	}
	return nil
}

func (r *AuthRepository) DeleteFamily(ctx context.Context, familyId string) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE family_id = $1;
	`
	_, err := r.db.Exec(ctx, query, familyId)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token family (%s): %w", familyId, err)
	}
	return nil
}

func (r *AuthRepository) DeleteOtherFamilies(ctx context.Context, userId, currentFamilyId string) error {
	query := `
        DELETE FROM refresh_tokens
        WHERE user_id = $1 AND family_id != $2;
    `
	_, err := r.db.Exec(ctx, query, userId, currentFamilyId)
	if err != nil {
		return fmt.Errorf("failed to revoke other sessions for user (%s) keeping family (%s): %w", userId, currentFamilyId, err)
	}
	return nil
}

func (r *AuthRepository) DeleteAllUserRefreshTokens(ctx context.Context, userId string) error {
	query := `
        DELETE FROM refresh_tokens
        WHERE user_id = $1;
    `
	_, err := r.db.Exec(ctx, query, userId)
	if err != nil {
		return fmt.Errorf("failed to revoke other sessions for user (%s): %w", userId, err)
	}
	return nil
}
