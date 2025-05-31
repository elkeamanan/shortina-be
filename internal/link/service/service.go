package service

import (
	"context"
	linkDomain "elkeamanan/shortina/internal/link/domain"

	"github.com/google/uuid"
)

type LinkService interface {
	StoreLink(ctx context.Context, request *linkDomain.StoreLinkRequest, actor *uuid.UUID) error
	GetLinkRedirection(ctx context.Context, request *linkDomain.GetLinkRedirectionRequest) (*linkDomain.GetLinkRedirectionResponse, error)
}
