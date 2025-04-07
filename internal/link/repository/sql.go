package repository

import (
	"context"
	"database/sql"
	"elkeamanan/shortina/cmd/storage"
	"elkeamanan/shortina/internal/link/domain"

	utilSQL "elkeamanan/shortina/util/sql"

	sq "github.com/Masterminds/squirrel"
)

type linkRepository struct {
	storage.CommonRepository
}

func NewLinkRepository(st storage.CommonRepository) LinkRepository {
	return &linkRepository{CommonRepository: st}
}

func (r *linkRepository) StoreLink(ctx context.Context, msg *domain.Link) error {
	_, err := r.CreateBuilder(nil).
		Insert(domain.TableLink).
		Columns(domain.GetLinkColumns(utilSQL.InsertOperation)...).
		Values(domain.GetInsertLinkValues(msg)...).
		ExecContext(ctx)
	return err
}

func (r *linkRepository) GetLinkRedirection(ctx context.Context, key string) (string, error) {
	var redirection string
	err := r.CreateBuilder(nil).
		Select(domain.ColumnLinkRedirection).
		From(domain.TableLink).
		Where(sq.Eq{domain.ColumnLinkKey: key}).
		OrderBy("created_at DESC").ScanContext(ctx, &redirection)

	if err != nil && err == sql.ErrNoRows {
		return "", nil
	}

	return redirection, err
}
