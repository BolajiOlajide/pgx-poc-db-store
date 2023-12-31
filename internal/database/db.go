package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/BolajiOlajide/pgx-poc-db-store/internal/database/basestore"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type db struct {
	*basestore.Store
	pool   *pgxpool.Pool
	logger *log.Logger
}

func (d *db) acquire(ctx context.Context) (*pgxpool.Conn, func(), error) {
	conn, err := d.pool.Acquire(ctx)
	if err != nil {
		return nil, nil, err
	}

	return conn, func() {
		conn.Release()
	}, nil
}

func (d *db) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return d.Handle().Query(ctx, query, args...)
}

func (d *db) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return d.Handle().Exec(ctx, query, args...)
}

func (d *db) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return d.Handle().QueryRow(ctx, query, args...)
}

func (d *db) GetSQLDB() *sql.DB {
	return stdlib.OpenDBFromPool(d.pool)
}

func (d *db) Close() {
	d.pool.Close()
}

func (d *db) WithTransact(ctx context.Context, f func(tx DB) error) error {
	return d.Store.WithTransact(ctx, func(tx *basestore.Store) error {
		return f(&db{logger: d.logger, Store: tx})
	})
}

func (d *db) Users() UserStore {
	return UsersWith(d.Store)
}

func (d *db) People() PeopleStore {
	return PeopleWith(d.Store)
}
