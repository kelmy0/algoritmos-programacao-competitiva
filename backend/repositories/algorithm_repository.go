package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelmy0/algoritmos-programacao-competitiva/backend/models"
)

type AlgorithmRepository struct {
	db *pgxpool.Pool
}

func NewAlgorithmRepository(db *pgxpool.Pool) *AlgorithmRepository {
	return &AlgorithmRepository{db: db}
}

func (r *AlgorithmRepository) List(ctx context.Context, limit, offset int) ([]models.Algorithm, error) {
	query := `
		SELECT id, public_id, slug, name, category, difficulty
		FROM algorithms
		ORDER BY name ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Algorithm
	for rows.Next() {
		var algo models.Algorithm

		err := rows.Scan(
			&algo.Id, &algo.PublicId, &algo.Slug, &algo.Name, &algo.Category,
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

func (r *AlgorithmRepository) GetByPublicID(ctx context.Context, publicId string) (*models.Algorithm, error) {
	query := `
		SELECT id, public_id, slug, name, category, difficulty, content, created_at, updated_at
		FROM algorithms
		WHERE public_id = $1
	`

	var algo models.Algorithm
	err := r.db.QueryRow(ctx, query, publicId).Scan(
		&algo.Id, &algo.PublicId, &algo.Slug, &algo.Name, &algo.Category,
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

func (r *AlgorithmRepository) PostAlgorithm(ctx context.Context, data models.NewAlgorithm) (*models.Algorithm, error) {
	query := `
		INSERT INTO algorithms (public_id, slug, name, category, difficulty, content) VALUES
		($1, $2, $3, $4, $5, $6)
		RETURNING id, public_id, slug, name, category, difficulty, content, created_at, updated_at;
	`

	var algo models.Algorithm
	err := r.db.QueryRow(ctx, query, data.PublicId, data.Slug,
		data.Name, data.Category, data.Difficulty, data.Content,
	).Scan(
		&algo.Id, &algo.PublicId, &algo.Slug, &algo.Name, &algo.Category,
		&algo.Difficulty, &algo.Content, &algo.CreatedAt, &algo.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &algo, nil
}

func (r *AlgorithmRepository) DeleteAlgorithm(ctx context.Context, publicId string) (*models.Algorithm, error) {
	query := `
		DELETE FROM algorithms 
		WHERE public_id = $1
		RETURNING id, public_id, slug, name, category, difficulty, content, created_at, updated_at;
	`

	var algo models.Algorithm
	err := r.db.QueryRow(ctx, query, publicId).Scan(
		&algo.Id, &algo.PublicId, &algo.Slug, &algo.Name, &algo.Category,
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

func (r *AlgorithmRepository) PutAlgorithm(ctx context.Context, data models.PutAlgorithm) (*models.Algorithm, error) {
	query := `
		UPDATE algorithms 
		SET slug = $1, name = $2, category = $3, difficulty = $4, content = $5
		WHERE public_id = $6
		RETURNING id, public_id, slug, name, category, difficulty, content, created_at, updated_at;
	`

	var algo models.Algorithm
	err := r.db.QueryRow(ctx, query, data.Slug, data.Name,
		data.Category, data.Difficulty, data.Content, data.PublicId,
	).Scan(
		&algo.Id, &algo.PublicId, &algo.Slug, &algo.Name, &algo.Category,
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
