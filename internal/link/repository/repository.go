package repository

import (
	"context"
	linkDomain "elkeamanan/shortina/internal/link/domain"
	"elkeamanan/shortina/storage/postgres"
)

type LinkRepository interface {
	postgres.CommonRepository
	StoreLink(ctx context.Context, link *linkDomain.Link) error
	GetLinkRedirection(ctx context.Context, key string) (string, error)
}
