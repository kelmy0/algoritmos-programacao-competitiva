package repositories

import (
	"context"
	"errors"
	"strings"
	"time"

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

func (r *AlgorithmRepository) List(ctx context.Context, limit, offset int) (list []dto.ListAlgorithmDTO, err error) {
	limit = max(1, min(100, limit))
	offset = max(0, offset)

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

	list = make([]dto.ListAlgorithmDTO, 0, limit)

	for rows.Next() {
		var algo dto.ListAlgorithmDTO

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

func (r *AlgorithmRepository) ListAdmin(ctx context.Context, limit, offset int, userId, status string) (list []dto.ListAlgorithmDTO, err error) {
	limit = max(1, min(100, limit))
	offset = max(0, offset)

	args := make([]any, 0, 4)
	args = append(args, userId)

	query := `
		SELECT public_id, slug, name, category, difficulty, status
		FROM algorithms
		WHERE author_id = $1
	`

	if status != "" {
		args = append(args, status)
		query += " AND status = $2"
	}

	args = append(args, limit, offset)

	if status != "" {
		query += " ORDER BY updated_at DESC LIMIT $3 OFFSET $4"
	} else {
		query += " ORDER BY updated_at DESC LIMIT $2 OFFSET $3"
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list = make([]dto.ListAlgorithmDTO, 0, limit)
	for rows.Next() {
		var algo dto.ListAlgorithmDTO

		err = rows.Scan(
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

func (r *AlgorithmRepository) ListModeration(ctx context.Context, limit, offset int, status string) (list []dto.ListAlgorithmDTO, err error) {
	limit = max(1, min(100, limit))
	offset = max(0, offset)

	args := make([]any, 0, 4)
	args = append(args, "deleted")

	var query strings.Builder
	query.Grow(200)
	query.WriteString(`
		SELECT public_id, slug, name, category, difficulty, status
		FROM algorithms
		WHERE status != $1
	`)

	if status != "" && status != "deleted" {
		args = append(args, status)
		query.WriteString(" AND status = $2")
	}

	args = append(args, limit, offset)

	if len(args) == 4 {
		query.WriteString(" ORDER BY updated_at DESC LIMIT $3 OFFSET $4;")
	} else {
		query.WriteString(" ORDER BY updated_at DESC LIMIT $2 OFFSET $3;")
	}

	rows, err := r.db.Query(ctx, query.String(), args...)
	if err != nil {
		return nil, models.ErrFailQueryingAlgorithm
	}
	defer rows.Close()

	list = make([]dto.ListAlgorithmDTO, 0, limit)

	for rows.Next() {
		var algo dto.ListAlgorithmDTO

		err = rows.Scan(
			&algo.PublicId,
			&algo.Slug,
			&algo.Name,
			&algo.Category,
			&algo.Difficulty,
			&algo.Status,
		)
		if err != nil {
			return nil, models.ErrFailQueryingAlgorithm
		}
		list = append(list, algo)
	}

	if err = rows.Err(); err != nil {
		return nil, models.ErrFailQueryingAlgorithm
	}

	return list, nil
}

func (r *AlgorithmRepository) GetByPublicID(ctx context.Context, publicId string) (algo *models.Algorithm, err error) {
	query := `
		SELECT public_id, slug, name, category, difficulty, content, created_at, updated_at
		FROM algorithms
		WHERE public_id = $1 AND status = 'approved'
	`

	var temp models.Algorithm
	err = r.db.QueryRow(ctx, query, publicId).Scan(
		&temp.PublicId, &temp.Slug, &temp.Name, &temp.Category,
		&temp.Difficulty, &temp.Content, &temp.CreatedAt, &temp.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrAlgorithmNotFound
		}
		return nil, err
	}

	algo = &temp
	return algo, nil
}

func (r *AlgorithmRepository) GetAdminAlgorithmById(ctx context.Context, algoId, userId string) (algo *models.Algorithm, err error) {
	query := `
		SELECT public_id, slug, name, category, difficulty, content, author_id, status, created_at, updated_at
		FROM algorithms
		WHERE public_id = $1 AND author_id = $2;
	`

	var temp models.Algorithm
	err = r.db.QueryRow(ctx, query, algoId, userId).Scan(
		&temp.PublicId, &temp.Slug, &temp.Name, &temp.Category,
		&temp.Difficulty, &temp.Content, &temp.AuthorId, &temp.Status, &temp.CreatedAt, &temp.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrAlgorithmNotFound
		}
		return nil, err
	}

	algo = &temp
	return algo, nil
}

func (r *AlgorithmRepository) PostAlgorithm(ctx context.Context, data models.PostAlgorithm) (algo dto.PostAlgorithmResponse, err error) {
	query := `
		INSERT INTO algorithms (public_id, slug, name, category, difficulty, content, author_id) VALUES
		($1, $2, $3, $4, $5, $6, $7)
		RETURNING public_id, slug;
	`

	err = r.db.QueryRow(ctx, query, data.PublicId, data.Slug,
		data.Name, data.Category, data.Difficulty, data.Content, data.AuthorId,
	).Scan(&algo.PublicId, &algo.Slug)

	if err != nil {
		return algo, err
	}

	return algo, nil
}

func (r *AlgorithmRepository) DeleteAlgorithm(ctx context.Context, publicId, userId string) error {
	return r.setStatus(ctx, publicId, userId, "deleted")
}

func (r *AlgorithmRepository) RestoreAlgorithm(ctx context.Context, publicId, userId string) error {
	return r.setStatus(ctx, publicId, userId, "pending")
}

func (r *AlgorithmRepository) PutAlgorithm(ctx context.Context, data models.PutAlgorithm, userId string) (algo dto.PutAlgorithmResponse, err error) {
	query := `
		UPDATE algorithms 
		SET slug = $1, name = $2, category = $3, difficulty = $4, content = $5, status = 'pending'
		WHERE public_id = $6 AND author_id = $7
		RETURNING public_id, slug;
	`

	err = r.db.QueryRow(ctx, query, data.Slug, data.Name,
		data.Category, data.Difficulty, data.Content, data.PublicId, userId,
	).Scan(&algo.PublicId, &algo.Slug)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return algo, models.ErrAlgorithmNotFound
		}
		return algo, err
	}

	return algo, nil
}

func (r *AlgorithmRepository) SitemapAlgorithms(ctx context.Context) (data []dto.SitemapItem, err error) {
	query := `
		SELECT slug, public_id, updated_at
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

	data = make([]dto.SitemapItem, 0, 1000)

	var slug, publicId string
	var updatedAt time.Time

	for rows.Next() {
		if err := rows.Scan(&slug, &publicId, &updatedAt); err != nil {
			return nil, err
		}

		data = append(data, dto.SitemapItem{
			Slug:      slug + "-" + publicId,
			UpdatedAt: updatedAt,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func (r *AlgorithmRepository) setStatus(ctx context.Context, publicId, userId, status string) error {
	query := `
		UPDATE algorithms
		SET status = $1
		WHERE public_id = $2 AND author_id = $3 AND status != $1;
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
