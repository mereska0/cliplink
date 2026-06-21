package repository

import (
	"context"

	"github.com/mereska0/cliplink/internal/domain"
)

type LinkRepository interface {
	Create(ctx context.Context, link *domain.Link) error
	SetShortCode(ctx context.Context, id int64, code string) error
	GetByCode(ctx context.Context, code string) (*domain.Link, error)
	List(ctx context.Context) ([]domain.Link, error)
	DeleteByCode(ctx context.Context, code string) error
	IncrementClicks(ctx context.Context, code string) error
}
