package repository

import (
	"context"
	"database/sql"
	"elkeamanan/shortina/internal/user/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, tx *sql.Tx, user *domain.User) error
	GetUser(ctx context.Context, pred domain.UserPredicate) (*domain.User, error)
	UpdateUser(ctx context.Context, tx *sql.Tx, user *domain.User, pred domain.UserPredicate) error
}
