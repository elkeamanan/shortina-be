package repository

import (
	"context"
	"database/sql"
	"elkeamanan/shortina/internal/link/domain"
	"elkeamanan/shortina/storage/postgres"
	"elkeamanan/shortina/util"

	sq "github.com/Masterminds/squirrel"
)

type linkRepository struct {
	postgres.CommonRepository
}

func NewLinkRepository(st postgres.CommonRepository) LinkRepository {
	return &linkRepository{CommonRepository: st}
}

func (r *linkRepository) StoreLink(ctx context.Context, link *domain.Link) error {
	_, err := r.CreateBuilder(nil).
		Insert(domain.TableLink).
		Columns(domain.GetLinkColumns(util.InsertOperation)...).
		Values(domain.GetInsertLinkValues(link)...).
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

func (r *linkRepository) selectBuilder(pred domain.LinkPredicate, paginationParam *util.PaginationParam) sq.SelectBuilder {
	selectBuilder := r.CreateBuilder(nil).
		Select(domain.GetLinkColumns(util.SelectOperation)...).
		From(domain.TableLink).
		Where(pred.ToWherePredicate())

	if paginationParam != nil {
		selectBuilder = selectBuilder.Offset(uint64(paginationParam.PageSize) - 1).Limit(uint64(paginationParam.PageSize))
	}

	return selectBuilder
}

func (r *linkRepository) GetLinks(ctx context.Context, pred domain.LinkPredicate, paginationParam *util.PaginationParam) ([]*domain.Link, error) {
	rows, err := r.selectBuilder(pred, paginationParam).QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanLinkRows(rows)
}

func (r *linkRepository) GetLink(ctx context.Context, pred domain.LinkPredicate) (*domain.Link, error) {
	links, err := r.GetLinks(ctx, pred, &util.PaginationParam{CurrentPage: 1, PageSize: 1})
	if err != nil {
		return nil, err
	}

	if len(links) == 0 {
		return nil, nil
	}

	return links[0], nil
}

func (r *linkRepository) CountLinks(ctx context.Context, pred domain.LinkPredicate) (result uint32, err error) {
	err = r.CreateBuilder(nil).
		Select("COUNT(*)").
		From(domain.TableLink).
		Where(pred.ToWherePredicate()).
		ScanContext(ctx, &result)

	return
}

func (r *linkRepository) UpdateLink(ctx context.Context, link *domain.Link, pred domain.LinkPredicate) error {
	_, err := r.CreateBuilder(nil).
		Update(domain.TableLink).
		SetMap(domain.GetUpdateLinkMap(link)).
		Where(pred.ToWherePredicate()).
		ExecContext(ctx)

	return err
}
