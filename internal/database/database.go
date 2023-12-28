package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

// local database so no cause for alarm
var dsn = "postgres://sourcegraph:sourcegraph@localhost/pgx-test"

type DB interface {
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, query string, args ...any) (pgx.Row, error)

	GetSQLDB() *sql.DB

	// Users() UserStore
	// People() PeopleStore
}

type db struct {
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
	conn, release, err := d.acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer release()

	return conn.Query(ctx, query, args...)
}

func (d *db) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	conn, release, err := d.acquire(ctx)
	if err != nil {
		return pgconn.CommandTag{}, err
	}
	defer release()

	return conn.Exec(ctx, query, args...)
}

func (d *db) QueryRow(ctx context.Context, query string, args ...any) (pgx.Row, error) {
	conn, release, err := d.acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer release()

	return conn.QueryRow(ctx, query, args...), err
}

func NewSQLConn(logger *log.Logger) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Println("error connecting to database:", err)
		return nil, err
	}
	return db, nil
}

func New(ctx context.Context, logger *log.Logger) DB {
	// Create database connection
	connPool, err := pgxpool.NewWithConfig(ctx, createPgxPoolConfig(logger))
	if err != nil {
		log.Fatal("Error while creating connection to the database!!")
	}

	conn, err := connPool.Acquire(ctx)
	if err != nil {
		log.Fatal("Error while acquiring connection to the database!!")
	}
	defer conn.Release()

	err = conn.Ping(ctx)
	if err != nil {
		log.Fatal("Error while pinging the database!!")
	}

	return &db{
		pool:   connPool,
		logger: logger,
	}
}

func (d *db) GetSQLDB() *sql.DB {
	return stdlib.OpenDBFromPool(d.pool)
}

func createPgxPoolConfig(logger *log.Logger) *pgxpool.Config {
	const defaultMaxConns = int32(4)
	const defaultMinConns = int32(1)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	dbConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		log.Println("Before acquiring the connection pool to the database!!")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		log.Println("After releasing the connection pool to the database!!")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		log.Println("Closed the connection pool to the database!!")
	}

	return dbConfig
}
