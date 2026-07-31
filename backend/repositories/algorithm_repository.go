package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
)

type AlgorithmRepository struct {
	db *pgxpool.Pool
}

func NewAlgorithmRepository(db *pgxpool.Pool) *AlgorithmRepository {
	return &AlgorithmRepository{db: db}
}

func (r *AlgorithmRepository) List(ctx context.Context, limit, offset int) ([]dto.AlgorithmDTO, error) {
	query := `
		SELECT public_id, slug, name, category, difficulty
		FROM algorithms
		WHERE status = 'approved'
		ORDER BY updated_at DESC 
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]dto.AlgorithmDTO, 0, limit)

	for rows.Next() {
		var algo dto.AlgorithmDTO

		err := rows.Scan(
			&algo.PublicId,
			&algo.Slug,
			&algo.Name,
			&algo.Category,
			&algo.Difficulty,
		)

		if err != nil {
			return nil, err
		}
		list = append(list, algo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *AlgorithmRepository) ListAdmin(ctx context.Context, limit, offset int, userId, status string) ([]dto.AlgorithmDTO, error) {
	query := `
		SELECT public_id, slug, name, category, difficulty, status
		FROM algorithms
		WHERE author_id = $1
	`

	args := []any{userId}
	if status != "" {
		args = append(args, status)
		query += fmt.Sprintf(" AND status = $%d", len(args))
	}

	args = append(args, limit, offset)
	query += fmt.Sprintf(" ORDER BY updated_at DESC LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrAlgorithmsNotFound
		}
		return nil, err
	}
	defer rows.Close()

	var list []dto.AlgorithmDTO
	for rows.Next() {
		var algo dto.AlgorithmDTO

		err := rows.Scan(
			&algo.PublicId, &algo.Slug, &algo.Name, &algo.Category,
			&algo.Difficulty, &algo.Status,
		)

		if err != nil {
			return nil, err
		}
		list = append(list, algo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *AlgorithmRepository) GetByPublicID(ctx context.Context, publicId string) (*dto.AlgorithmDTO, error) {
	query := `
		SELECT public_id, slug, name, category, difficulty, content, created_at, updated_at
		FROM algorithms
		WHERE public_id = $1 AND status = 'approved'
	`

	var algo dto.AlgorithmDTO
	err := r.db.QueryRow(ctx, query, publicId).Scan(
		&algo.PublicId, &algo.Slug, &algo.Name, &algo.Category,
		&algo.Difficulty, &algo.Content, &algo.CreatedAt, &algo.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrAlgorithmNotFound
		}
		return nil, err
	}

	return &algo, nil
}

func (r *AlgorithmRepository) GetAdminAlgorithmById(ctx context.Context, algoId, userId string) (*dto.AlgorithmDTO, error) {
	query := `
		SELECT public_id, slug, name, category, difficulty, content, author_id, status, created_at, updated_at
		FROM algorithms
		WHERE public_id = $1 AND author_id = $2;
	`

	var algo dto.AlgorithmDTO
	err := r.db.QueryRow(ctx, query, algoId, userId).Scan(
		&algo.PublicId, &algo.Slug, &algo.Name, &algo.Category,
		&algo.Difficulty, &algo.Content, &algo.AuthorId, &algo.Status, &algo.CreatedAt, &algo.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrAlgorithmNotFound
		}
		return nil, err
	}

	return &algo, nil
}

func (r *AlgorithmRepository) PostAlgorithm(ctx context.Context, data models.NewAlgorithm) (*dto.AlgorithmDTO, error) {
	query := `
		INSERT INTO algorithms (public_id, slug, name, category, difficulty, content, author_id) VALUES
		($1, $2, $3, $4, $5, $6, $7)
		RETURNING public_id, slug;
	`

	var algo dto.AlgorithmDTO
	err := r.db.QueryRow(ctx, query, data.PublicId, data.Slug,
		data.Name, data.Category, data.Difficulty, data.Content, data.AuthorId,
	).Scan(&algo.PublicId, &algo.Slug)

	if err != nil {
		return nil, err
	}

	return &algo, nil
}

func (r *AlgorithmRepository) DeleteAlgorithm(ctx context.Context, publicId, userId string) error {
	return r.setStatus(ctx, publicId, userId, "deleted")
}

func (r *AlgorithmRepository) RestoreAlgorithm(ctx context.Context, publicId, userId string) error {
	return r.setStatus(ctx, publicId, userId, "pending")
}

func (r *AlgorithmRepository) PutAlgorithm(ctx context.Context, data models.PutAlgorithm, userId string) (*dto.AlgorithmDTO, error) {
	query := `
		UPDATE algorithms 
		SET slug = $1, name = $2, category = $3, difficulty = $4, content = $5, status = 'pending'
		WHERE public_id = $6 AND author_id = $7
		RETURNING public_id, slug, name, category, difficulty, content, created_at, updated_at;
	`

	var algo dto.AlgorithmDTO
	err := r.db.QueryRow(ctx, query, data.Slug, data.Name,
		data.Category, data.Difficulty, data.Content, data.PublicId, userId,
	).Scan(
		&algo.PublicId, &algo.Slug, &algo.Name, &algo.Category,
		&algo.Difficulty, &algo.Content, &algo.CreatedAt, &algo.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrAlgorithmNotFound
		}
		return nil, err
	}

	return &algo, nil
}

func (r *AlgorithmRepository) SitemapAlgorithms(ctx context.Context) ([]dto.SitemapItem, error) {
	query := `
		SELECT 
			CONCAT(slug, '-', public_id) AS slug, 
			updated_at
		FROM algorithms
		WHERE status = 'approved'
		ORDER BY updated_at DESC
		LIMIT 49000;
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]dto.SitemapItem, 0)

	for rows.Next() {
		var item dto.SitemapItem

		if err := rows.Scan(&item.Slug, &item.UpdatedAt); err != nil {
			return nil, err
		}

		list = append(list, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *AlgorithmRepository) setStatus(ctx context.Context, publicId, userId, status string) error {
	query := `
		UPDATE algorithms
		SET status = $1
		WHERE public_id = $2 AND author_id = $3;
	`

	res, err := r.db.Exec(ctx, query, status, publicId, userId)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return models.ErrAlgorithmNotFound
	}

	return nil
}
