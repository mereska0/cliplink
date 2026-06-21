package memory

import (
	"context"
	"sync"
	"time"

	"github.com/mereska0/cliplink/internal/domain"
)

type InMemoryRepository struct {
	mu          sync.RWMutex
	linksByCode map[string]*domain.Link
	linksByID   map[int64]*domain.Link
	nextID      int64
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		linksByCode: make(map[string]*domain.Link),
		linksByID:   make(map[int64]*domain.Link),
		nextID:      1,
	}
}

func (imr *InMemoryRepository) Create(ctx context.Context, link *domain.Link) error {
	imr.mu.Lock()
	defer imr.mu.Unlock()
	link.ID = imr.nextID
	link.CreatedAt = time.Now()
	imr.nextID++
	imr.linksByID[link.ID] = link
	return nil
}

func (imr *InMemoryRepository) SetShortCode(ctx context.Context, id int64, code string) error {
	imr.mu.Lock()
	defer imr.mu.Unlock()
	if _, exists := imr.linksByCode[code]; exists {
		return domain.ErrAliasTaken
	}
	link, exists := imr.linksByID[id]
	if !exists {
		return domain.ErrLinkNotFound
	}
	link.ShortCode = code
	imr.linksByCode[code] = link
	return nil
}

func (imr *InMemoryRepository) GetByCode(ctx context.Context, code string) (*domain.Link, error) {
	imr.mu.RLock()
	defer imr.mu.RUnlock()
	link, exists := imr.linksByCode[code]
	if !exists {
		return nil, domain.ErrLinkNotFound
	}
	if link.DeletedAt != nil {
		return nil, domain.ErrDeletedLink
	}
	copyLink := *link
	return &copyLink, nil
}

func (imr *InMemoryRepository) List(ctx context.Context) ([]domain.Link, error) {
	imr.mu.RLock()
	defer imr.mu.RUnlock()
	result := make([]domain.Link, 0, len(imr.linksByCode))
	for _, link := range imr.linksByCode {
		if link.DeletedAt != nil {
			continue
		}
		result = append(result, *link)
	}
	return result, nil
}

func (imr *InMemoryRepository) DeleteByCode(ctx context.Context, code string) error {
	imr.mu.Lock()
	defer imr.mu.Unlock()
	link, exists := imr.linksByCode[code]
	if !exists {
		return domain.ErrLinkNotFound
	}
	if link.DeletedAt != nil {
		return domain.ErrDeletedLink
	}
	now := time.Now()
	link.DeletedAt = &now
	return nil
}

func (imr *InMemoryRepository) IncrementClicks(ctx context.Context, code string) error {
	imr.mu.Lock()
	defer imr.mu.Unlock()
	link, exists := imr.linksByCode[code]
	if !exists {
		return domain.ErrLinkNotFound
	}
	if link.DeletedAt != nil {
		return domain.ErrDeletedLink
	}
	link.Clicks++
	return nil
}
