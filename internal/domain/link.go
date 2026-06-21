package domain

import "time"

type Link struct {
	ID          int64
	ShortCode   string
	OriginalURL string
	Clicks      int64
	CreatedAt   time.Time
	DeletedAt   *time.Time
}
