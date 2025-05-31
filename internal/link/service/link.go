package service

import (
	"context"
	linkDomain "elkeamanan/shortina/internal/link/domain"
	"elkeamanan/shortina/internal/link/repository"
	"elkeamanan/shortina/util"
	"errors"

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

func (s *linkService) GetLinkRedirection(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", errors.New("key is required")
	}

	return s.linkRepository.GetLinkRedirection(ctx, key)
}

func (s *linkService) GetLink(ctx context.Context, pred linkDomain.LinkPredicate) (*linkDomain.Link, error) {
	return s.linkRepository.GetLink(ctx, pred)
}

func (s *linkService) GetPaginatedLinks(ctx context.Context, pred linkDomain.LinkPredicate, paginationParam util.PaginationParam) ([]*linkDomain.Link, *util.PaginationResponse, error) {
	links, err := s.linkRepository.GetLinks(ctx, pred, &paginationParam)
	if err != nil {
		return nil, nil, err
	}

	count, err := s.linkRepository.CountLinks(ctx, pred)
	if err != nil {
		return nil, nil, err
	}

	paginationResponse := util.GeneratePaginationResponse(paginationParam.PageSize, paginationParam.CurrentPage, count)
	return links, &paginationResponse, nil
}

func (s *linkService) UpdateLink(ctx context.Context, link *linkDomain.Link, pred linkDomain.LinkPredicate) error {
	return s.linkRepository.UpdateLink(ctx, link, pred)
}
