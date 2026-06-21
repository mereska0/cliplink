package domain

import "errors"

var (
	ErrInvalidURL   = errors.New("invalid url")
	ErrLinkNotFound = errors.New("link not found")
	ErrAliasTaken   = errors.New("alias already taken")
	ErrDeletedLink  = errors.New("link is deleted")
)
