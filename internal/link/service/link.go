package service

import (
	"context"
	linkDomain "elkeamanan/shortina/internal/link/domain"
	"elkeamanan/shortina/internal/link/repository"

	"github.com/google/uuid"
)

type linkService struct {
	linkRepository repository.LinkRepository
}

func NewLinkService(linkRepository repository.LinkRepository) LinkService {
	return &linkService{linkRepository: linkRepository}
}

func (s *linkService) StoreLink(ctx context.Context, request *linkDomain.StoreLinkRequest, actor *uuid.UUID) error {
	if err := request.Validate(); err != nil {
		return err
	}

	return s.linkRepository.StoreLink(ctx, request.ToNewLink(actor))
}

func (s *linkService) GetLinkRedirection(ctx context.Context, request *linkDomain.GetLinkRedirectionRequest) (*linkDomain.GetLinkRedirectionResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}
	redirection, err := s.linkRepository.GetLinkRedirection(ctx, request.Key)
	if err != nil {
		return nil, err
	}

	if redirection == "" {
		return nil, nil
	}

	return &linkDomain.GetLinkRedirectionResponse{
		Redirection: redirection,
	}, nil
}
