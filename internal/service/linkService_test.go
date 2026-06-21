package service

import (
	"context"
	"errors"
	"testing"

	"github.com/mereska0/cliplink/internal/domain"
	"github.com/mereska0/cliplink/internal/encoder"
	"github.com/mereska0/cliplink/internal/repository/memory"
)

func newTestService() *LinkService {
	repo := memory.NewInMemoryRepository()
	encoder := encoder.NewBase62Encoder()

	return NewLinkService(repo, encoder)
}

func TestLinkService_CreateLink(t *testing.T) {
	ctx := context.Background()
	linkService := newTestService()

	link, err := linkService.CreateLink(ctx, CreateLinkInput{
		OriginalURL: "https://example.com",
	})

	if err != nil {
		t.Fatalf("CreateLink returned error: %v", err)
	}

	if link.ID == 0 {
		t.Fatal("expected link ID to be set")
	}

	if link.ShortCode == "" {
		t.Fatal("expected short code to be set")
	}

	if link.OriginalURL != "https://example.com" {
		t.Fatalf("OriginalURL = %s, want %s", link.OriginalURL, "https://example.com")
	}

	if link.Clicks != 0 {
		t.Fatalf("Clicks = %d, want 0", link.Clicks)
	}
}

func TestLinkService_CreateLink_InvalidURL(t *testing.T) {
	ctx := context.Background()
	linkService := newTestService()

	_, err := linkService.CreateLink(ctx, CreateLinkInput{
		OriginalURL: "not-a-url",
	})

	if !errors.Is(err, domain.ErrInvalidURL) {
		t.Fatalf("expected ErrInvalidURL, got %v", err)
	}
}

func TestLinkService_CreateLink_CustomAlias(t *testing.T) {
	ctx := context.Background()
	linkService := newTestService()

	link, err := linkService.CreateLink(ctx, CreateLinkInput{
		OriginalURL: "https://example.com",
		CustomAlias: "my-link",
	})

	if err != nil {
		t.Fatalf("CreateLink returned error: %v", err)
	}

	if link.ShortCode != "my-link" {
		t.Fatalf("ShortCode = %s, want my-link", link.ShortCode)
	}
}

func TestLinkService_CreateLink_DuplicateCustomAlias(t *testing.T) {
	ctx := context.Background()
	linkService := newTestService()

	_, err := linkService.CreateLink(ctx, CreateLinkInput{
		OriginalURL: "https://example.com",
		CustomAlias: "same",
	})
	if err != nil {
		t.Fatalf("first CreateLink returned error: %v", err)
	}

	_, err = linkService.CreateLink(ctx, CreateLinkInput{
		OriginalURL: "https://google.com",
		CustomAlias: "same",
	})

	if !errors.Is(err, domain.ErrAliasTaken) {
		t.Fatalf("expected ErrAliasTaken, got %v", err)
	}
}

func TestLinkService_RegisterClick(t *testing.T) {
	ctx := context.Background()
	linkService := newTestService()

	link, err := linkService.CreateLink(ctx, CreateLinkInput{
		OriginalURL: "https://example.com",
	})
	if err != nil {
		t.Fatalf("CreateLink returned error: %v", err)
	}

	openedLink, err := linkService.RegisterClick(ctx, link.ShortCode)
	if err != nil {
		t.Fatalf("RegisterClick returned error: %v", err)
	}

	if openedLink.Clicks != 1 {
		t.Fatalf("Clicks = %d, want 1", openedLink.Clicks)
	}
}

func TestLinkService_RegisterClick_UnknownCode(t *testing.T) {
	ctx := context.Background()
	linkService := newTestService()

	_, err := linkService.RegisterClick(ctx, "unknown")

	if !errors.Is(err, domain.ErrLinkNotFound) {
		t.Fatalf("expected ErrLinkNotFound, got %v", err)
	}
}

func TestLinkService_ListLinks(t *testing.T) {
	ctx := context.Background()
	linkService := newTestService()

	_, err := linkService.CreateLink(ctx, CreateLinkInput{
		OriginalURL: "https://example.com",
	})
	if err != nil {
		t.Fatalf("CreateLink returned error: %v", err)
	}

	_, err = linkService.CreateLink(ctx, CreateLinkInput{
		OriginalURL: "https://google.com",
	})
	if err != nil {
		t.Fatalf("CreateLink returned error: %v", err)
	}

	links, err := linkService.ListLinks(ctx)
	if err != nil {
		t.Fatalf("ListLinks returned error: %v", err)
	}

	if len(links) != 2 {
		t.Fatalf("len(links) = %d, want 2", len(links))
	}
}

func TestLinkService_DeleteLink(t *testing.T) {
	ctx := context.Background()
	linkService := newTestService()

	link, err := linkService.CreateLink(ctx, CreateLinkInput{
		OriginalURL: "https://example.com",
	})
	if err != nil {
		t.Fatalf("CreateLink returned error: %v", err)
	}

	err = linkService.DeleteLink(ctx, link.ShortCode)
	if err != nil {
		t.Fatalf("DeleteLink returned error: %v", err)
	}

	_, err = linkService.GetLink(ctx, link.ShortCode)

	if !errors.Is(err, domain.ErrDeletedLink) {
		t.Fatalf("expected ErrDeletedLink, got %v", err)
	}
}
