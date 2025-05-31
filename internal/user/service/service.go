package service

import (
	"context"
	"elkeamanan/shortina/internal/user/domain"

	"github.com/google/uuid"
)

type UserService interface {
	RegisterUser(ctx context.Context, req *domain.RegisterUserRequest) error
	LoginUser(ctx context.Context, req *domain.LoginUserRequest) (*domain.TokenPair, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req *domain.UpdateUserRequest) error
	RefreshUserToken(ctx context.Context, refreshToken string) (*domain.TokenPair, error)
}
