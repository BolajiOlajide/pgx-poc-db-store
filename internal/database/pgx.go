package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func createPgxPoolConfig(logger *log.Logger) *pgxpool.Config {
	const defaultMaxConns = int32(4)
	const defaultMinConns = int32(1)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	dbConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		logger.Println("Before acquiring the connection pool to the database!!")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		logger.Println("After releasing the connection pool to the database!!")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		logger.Println("Closed the connection pool to the database!!")
	}

	dbConfig.BeforeConnect = func(ctx context.Context, cc *pgx.ConnConfig) error {
		logger.Println("Before connecting to database")
		return nil
	}

	dbConfig.AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
		logger.Println("after connecting to database!!")
		return nil
	}

	return dbConfig
}
