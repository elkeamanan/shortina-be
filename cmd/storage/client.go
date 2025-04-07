package storage

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type CommonRepository interface {
	Ping(ctx context.Context) error
	CreateBuilder(tx *sql.Tx) sq.StatementBuilderType
	RunInSQLTransaction(ctx context.Context, isolationLevel sql.IsolationLevel, fn func(tx *sql.Tx) error) error
}

type commonRepository struct {
	client *sql.DB
	dbName string
}

func (c *commonRepository) Ping(ctx context.Context) error {
	return c.client.PingContext(ctx)
}

func (c *commonRepository) CreateBuilder(tx *sql.Tx) sq.StatementBuilderType {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(c.client)
	if tx != nil {
		builder = builder.RunWith(tx)
	}
	return builder
}

func (c *commonRepository) RunInSQLTransaction(ctx context.Context, isolationLevel sql.IsolationLevel, fn func(tx *sql.Tx) error) error {
	tx, err := c.client.BeginTx(ctx, &sql.TxOptions{
		Isolation: isolationLevel,
	})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}
