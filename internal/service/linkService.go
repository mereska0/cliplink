package service

import (
	"context"
	"net/url"

	"github.com/mereska0/cliplink/internal/domain"
	"github.com/mereska0/cliplink/internal/encoder"
	"github.com/mereska0/cliplink/internal/repository"
)

const codeOffset = 1024

type CreateLinkInput struct {
	OriginalURL string
	CustomAlias string
}

type LinkService struct {
	repo    repository.LinkRepository
	encoder encoder.CodeEncoder
}

func NewLinkService(
	repo repository.LinkRepository,
	encoder encoder.CodeEncoder,
) *LinkService {
	return &LinkService{
		repo:    repo,
		encoder: encoder,
	}
}

func (s *LinkService) CreateLink(ctx context.Context, input CreateLinkInput) (*domain.Link, error) {
	if !isValidURL(input.OriginalURL) {
		return nil, domain.ErrInvalidURL
	}
	link := &domain.Link{
		OriginalURL: input.OriginalURL,
	}
	if err := s.repo.Create(ctx, link); err != nil {
		return nil, err
	}
	var code string
	if input.CustomAlias != "" {
		code = input.CustomAlias
	} else {
		code = s.encoder.Encode(link.ID + codeOffset)
	}
	if err := s.repo.SetShortCode(ctx, link.ID, code); err != nil {
		return nil, err
	}
	link.ShortCode = code
	return link, nil
}

func (s *LinkService) GetLink(ctx context.Context, code string) (*domain.Link, error) {
	return s.repo.GetByCode(ctx, code)
}

func (s *LinkService) ListLinks(ctx context.Context) ([]domain.Link, error) {
	return s.repo.List(ctx)
}

func (s *LinkService) DeleteLink(ctx context.Context, code string) error {
	return s.repo.DeleteByCode(ctx, code)
}

func (s *LinkService) RegisterClick(ctx context.Context, code string) (*domain.Link, error) {
	link, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if err := s.repo.IncrementClicks(ctx, code); err != nil {
		return nil, err
	}
	link.Clicks++
	return link, nil
}

func isValidURL(rawURL string) bool {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	if parsedURL.Host == "" {
		return false
	}

	return true
}
