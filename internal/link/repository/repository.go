package repository

import (
	"context"
	"elkeamanan/shortina/cmd/storage"
	linkDomain "elkeamanan/shortina/internal/link/domain"
)

type LinkRepository interface {
	storage.CommonRepository
	StoreLink(ctx context.Context, link *linkDomain.Link) error
	GetLinkRedirection(ctx context.Context, key string) (string, error)
}
