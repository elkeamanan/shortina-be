package service

import (
	"context"
	linkDomain "elkeamanan/shortina/internal/link/domain"
)

type LinkService interface {
	StoreLink(ctx context.Context, request *linkDomain.StoreLinkRequest) error
	GetLinkRedirection(ctx context.Context, request *linkDomain.GetLinkRedirectionRequest) (*linkDomain.GetLinkRedirectionResponse, error)
}
