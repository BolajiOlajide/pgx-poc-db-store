package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/BolajiOlajide/pgx-poc-db-store/internal/database/basestore"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// local database so no cause for alarm
var dsn = "postgres://sourcegraph:sourcegraph@localhost/pgx-test"

type DB interface {
	basestore.ShareableStore

	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row

	Users() UserStore
	// People() PeopleStore

	WithTransact(context.Context, func(tx DB) error) error
	GetSQLDB() *sql.DB
	Close()
}

var _ DB = (*db)(nil)

func New(ctx context.Context, logger *log.Logger) DB {
	connPool, err := pgxpool.NewWithConfig(ctx, createPgxPoolConfig(logger))
	if err != nil {
		logger.Fatal("Error while creating connection to the database!!")
	}

	conn, err := connPool.Acquire(ctx)
	if err != nil {
		logger.Fatal("Error while acquiring connection to the database!!")
	}
	defer conn.Release()

	err = conn.Ping(ctx)
	if err != nil {
		logger.Fatal("Error while pinging the database!!")
	}

	return &db{
		pool:   connPool,
		logger: logger,
		Store:  basestore.NewWithHandle(basestore.NewHandleWithDB(logger, connPool, pgx.TxOptions{})),
	}
}
