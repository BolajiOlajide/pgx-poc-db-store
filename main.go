package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BolajiOlajide/pgx-poc-db-store/internal/database"
	migrate "github.com/rubenv/sql-migrate"
)

func main() {
	logger := log.New(os.Stdout, "pgx-testrunner", log.LstdFlags)
	runMigrations(logger)
	fmt.Println("hey hey hey!")
}

func runMigrations(logger *log.Logger) {
	db, err := database.NewSQLConn(logger)
	if err != nil {
		os.Exit(1)
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}
	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		logger.Println("error running migrations:", err)
		os.Exit(1)
	}
	log.Printf("Applied %d migrations!\n", n)
}
