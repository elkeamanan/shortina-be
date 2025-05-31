package repository

import (
	"context"
	linkDomain "elkeamanan/shortina/internal/link/domain"
	"elkeamanan/shortina/storage/postgres"
	"elkeamanan/shortina/util"
)

type LinkRepository interface {
	postgres.CommonRepository
	StoreLink(ctx context.Context, link *linkDomain.Link) error
	GetLinkRedirection(ctx context.Context, key string) (string, error)
	GetLink(ctx context.Context, pred linkDomain.LinkPredicate) (*linkDomain.Link, error)
	GetLinks(ctx context.Context, pred linkDomain.LinkPredicate, paginationParam *util.PaginationParam) ([]*linkDomain.Link, error)
	CountLinks(ctx context.Context, pred linkDomain.LinkPredicate) (uint32, error)
	UpdateLink(ctx context.Context, link *linkDomain.Link, pred linkDomain.LinkPredicate) error
}
