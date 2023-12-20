package database

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// local database so no cause for alarm
var dsn = "postgres://sourcegraph:sourcegraph@localhost/pgx-test"

type DB interface {
	QueryContext(ctx context.Context, q string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row

	Users() UserStore
	People() PeopleStore
}

func NewSQLConn(logger *log.Logger) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Println("error connecting to database:", err)
		return nil, err
	}
	return db, nil
}
