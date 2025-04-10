package domain

import (
	"errors"
	"net/url"

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

	_, err := url.ParseRequestURI(req.Redirection)
	if err != nil {
		return errors.New("redirection data is not a valid")
	}

	return nil
}

func (req *StoreLinkRequest) ToLink() *Link {
	return &Link{
		ID:          uuid.New(),
		Key:         req.Key,
		Redirection: req.Redirection,
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
