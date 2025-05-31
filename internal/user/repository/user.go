package repository

import (
	"context"
	"database/sql"
	"elkeamanan/shortina/cmd/storage"
	"elkeamanan/shortina/internal/user/domain"
	sqlUtil "elkeamanan/shortina/util/sql"
)

type userRepository struct {
	storage.CommonRepository
}

func NewUserRepository(st storage.CommonRepository) UserRepository {
	return &userRepository{CommonRepository: st}
}

func (r *userRepository) CreateUser(ctx context.Context, tx *sql.Tx, user *domain.User) error {
	_, err := r.CreateBuilder(tx).
		Insert(domain.TableUsers).
		Values(domain.GetInsertUsersValues(user)...).
		ExecContext(ctx)
	return err
}

func (r *userRepository) GetUsers(ctx context.Context, pred domain.UserPredicate) ([]*domain.User, error) {
	rows, err := r.CreateBuilder(nil).
		Select(domain.GetUsersColumns(sqlUtil.SelectOperation)...).
		From(domain.TableUsers).
		Where(pred.ToWherePredicate()).QueryContext(ctx)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanUserRows(rows)

}

func (r *userRepository) GetUser(ctx context.Context, pred domain.UserPredicate) (*domain.User, error) {
	users, err := r.GetUsers(ctx, pred)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, nil
	}

	return users[0], nil
}

func (r *userRepository) UpdateUser(ctx context.Context, tx *sql.Tx, user *domain.User, pred domain.UserPredicate) error {
	_, err := r.CreateBuilder(tx).
		Update(domain.TableUsers).
		SetMap(domain.GetUpdateUsersMap(user)).
		Where(pred.ToWherePredicate()).
		ExecContext(ctx)
	return err
}
