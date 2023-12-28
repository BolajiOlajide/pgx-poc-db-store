package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BolajiOlajide/pgx-poc-db-store/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	migrate "github.com/rubenv/sql-migrate"
)

func main() {
	logger := log.New(os.Stdout, "pgx-testrunner", log.LstdFlags)
	ctx := context.Background()

	db := database.New(ctx, logger)

	runMigrations(db, logger)
	fmt.Println("hey hey hey!")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	s := newServer(db, r)
	s.setupRoutes()

	http.ListenAndServe(":3000", r)
}

func runMigrations(db database.DB, logger *log.Logger) {
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}
	n, err := migrate.Exec(db.GetSQLDB(), "postgres", migrations, migrate.Up)
	if err != nil {
		logger.Println("error running migrations:", err)
		os.Exit(1)
	}
	log.Printf("Applied %d migrations!\n", n)
}
