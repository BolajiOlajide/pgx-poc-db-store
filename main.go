package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BolajiOlajide/pgx-poc-db-store/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	migrate "github.com/rubenv/sql-migrate"
)

func main() {
	logger := log.New(os.Stdout, "pgx-testrunner: ", log.LstdFlags)
	ctx := context.Background()

	db := database.New(ctx, logger)

	runMigrations(db, logger)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	s := newServer(db, r)
	s.setupRoutes()

	// Create a channel to receive the interrupt signal
	interruptChan := make(chan os.Signal, 1)

	// Notify the channel for the interrupt signal
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	go s.start()

	<-interruptChan

	s.gracefulShutdown()
	logger.Println("Server stopped gracefully")
}

func runMigrations(db database.DB, logger *log.Logger) {
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	sdb := db.GetSQLDB()
	defer sdb.Close()

	n, err := migrate.Exec(sdb, "postgres", migrations, migrate.Up)
	if err != nil {
		logger.Println("error running migrations:", err)
		os.Exit(1)
	}
	logger.Printf("Applied %d migrations!\n", n)
}
