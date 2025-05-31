package domain

import (
	"elkeamanan/shortina/util"
	"errors"

	"github.com/google/uuid"
)

type StoreLinkRequest struct {
	Key         string `json:"key"`
	Redirection string `json:"redirection"`
}

func (req *StoreLinkRequest) Validate() error {
	if req == nil {
		return errors.New("store link request is missing")
	}

	if req.Key == "" {
		return errors.New("cannot store link without shortcut")
	}

	if req.Redirection == "" {
		return errors.New("cannot store link without redirection")
	}

	return nil
}

func (req *StoreLinkRequest) ToNewLink(createdBy *uuid.UUID) *Link {
	return &Link{
		ID:          uuid.New(),
		Key:         req.Key,
		Redirection: req.Redirection,
		Status:      LinkStatusActive,
		CreatedBy:   createdBy,
	}
}

type GetLinkRedirectionRequest struct {
	Key string
}

func (req *GetLinkRedirectionRequest) Validate() error {
	if req == nil {
		return errors.New("get link request is missing")
	}

	if req.Key == "" {
		return errors.New("cannot get link without identifier")
	}

	return nil
}

type GetLinkRedirectionResponse struct {
	Redirection string `json:"redirection"`
}

type GetLinkResponse struct {
	Link *Link `json:"link"`
}

type GetLinksResponse struct {
	Pagination *util.PaginationResponse `json:"pagination"`
	Links      []*Link                  `json:"links"`
}

type UpdateLinkRequest struct {
	Key         string     `json:"key"`
	Redirection string     `json:"redirection"`
	Status      LinkStatus `json:"status"`
}
