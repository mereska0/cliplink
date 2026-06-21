package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mereska0/cliplink/internal/domain"
)

type PostgresLinkRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresLinkRepository(pool *pgxpool.Pool) *PostgresLinkRepository {
	return &PostgresLinkRepository{pool: pool}
}

func (pr *PostgresLinkRepository) Create(ctx context.Context, link *domain.Link) error {
	query := `
		INSERT INTO links (original_url)
		VALUES ($1)
		RETURNING id, created_at
	`

	err := pr.pool.QueryRow(ctx, query, link.OriginalURL).
		Scan(&link.ID, &link.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
func (pr *PostgresLinkRepository) SetShortCode(ctx context.Context, id int64, code string) error {
	query := `
		UPDATE links
		SET short_code = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	tag, err := pr.pool.Exec(ctx, query, code, id)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrAliasTaken
		}

		return err
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrLinkNotFound
	}
	return nil
}
func (pr *PostgresLinkRepository) GetByCode(ctx context.Context, code string) (*domain.Link, error) {
	query := `
		SELECT id, short_code, original_url, clicks, created_at, deleted_at
		FROM links
		WHERE short_code = $1
	`

	var link domain.Link

	err := pr.pool.QueryRow(ctx, query, code).Scan(
		&link.ID,
		&link.ShortCode,
		&link.OriginalURL,
		&link.Clicks,
		&link.CreatedAt,
		&link.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLinkNotFound
		}

		return nil, err
	}

	if link.DeletedAt != nil {
		return nil, domain.ErrDeletedLink
	}

	return &link, nil
}

func (pr *PostgresLinkRepository) List(ctx context.Context) ([]domain.Link, error) {
	query := `
		SELECT id, short_code, original_url, clicks, created_at, deleted_at
		FROM links
		WHERE deleted_at IS NULL
		ORDER BY id DESC
	`

	rows, err := pr.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []domain.Link

	for rows.Next() {
		var link domain.Link

		err := rows.Scan(
			&link.ID,
			&link.ShortCode,
			&link.OriginalURL,
			&link.Clicks,
			&link.CreatedAt,
			&link.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return links, nil
}

func (pr *PostgresLinkRepository) DeleteByCode(ctx context.Context, code string) error {
	query := `
		UPDATE links
		SET deleted_at = NOW()
		WHERE short_code = $1 AND deleted_at IS NULL
	`

	tag, err := pr.pool.Exec(ctx, query, code)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrLinkNotFound
	}

	return nil
}

func (pr *PostgresLinkRepository) IncrementClicks(ctx context.Context, code string) error {
	query := `
		UPDATE links
		SET clicks = clicks + 1
		WHERE short_code = $1 AND deleted_at IS NULL
	`

	tag, err := pr.pool.Exec(ctx, query, code)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrLinkNotFound
	}

	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}
