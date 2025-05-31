package service

import (
	"context"
	linkDomain "elkeamanan/shortina/internal/link/domain"
	"elkeamanan/shortina/util"

	"github.com/google/uuid"
)

type LinkService interface {
	StoreLink(ctx context.Context, request *linkDomain.StoreLinkRequest, actor *uuid.UUID) error
	GetLinkRedirection(ctx context.Context, key string) (string, error)
	GetLink(ctx context.Context, pred linkDomain.LinkPredicate) (*linkDomain.Link, error)
	GetPaginatedLinks(ctx context.Context, pred linkDomain.LinkPredicate, paginationParam util.PaginationParam) ([]*linkDomain.Link, *util.PaginationResponse, error)
	UpdateLink(ctx context.Context, link *linkDomain.Link, pred linkDomain.LinkPredicate) error
}
