package repositories

import (
	"context"
	"fmt"

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
		SELECT id, name, category, difficulty, content, created_at, updated_at
		FROM algorithms
		ORDER BY name ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query algorithms: %w", err)
	}
	defer rows.Close()

	var list []models.Algorithm
	for rows.Next() {
		var algo models.Algorithm

		err := rows.Scan(
			&algo.Id,
			&algo.Name,
			&algo.Category,
			&algo.Difficulty,
			&algo.Content,
			&algo.CreatedAt,
			&algo.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}
		list = append(list, algo)
	}

	return list, nil
}

func (r *AlgorithmRepository) GetById(ctx context.Context, id string) (*models.Algorithm, error) {
	query := `
		SELECT id, name, category, difficulty, content, created_at, updated_at
		FROM algorithms
		WHERE id = $1
	`

	var algo models.Algorithm
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&algo.Id,
		&algo.Name,
		&algo.Category,
		&algo.Difficulty,
		&algo.Content,
		&algo.CreatedAt,
		&algo.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &algo, nil
}
