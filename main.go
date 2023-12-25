package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"clevergo.tech/jsend"
	"github.com/BolajiOlajide/pgx-poc-db-store/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	migrate "github.com/rubenv/sql-migrate"
)

func main() {
	logger := log.New(os.Stdout, "pgx-testrunner", log.LstdFlags)
	runMigrations(logger)
	fmt.Println("hey hey hey!")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", rootHandler)
	r.Route("/people", func(ir chi.Router) {
		ir.Get("/", getPeople)
		ir.Post("/", createPeople)
	})
	r.Route("/user", func(ir chi.Router) {
		ir.Get("/", getUsers)
		ir.Get("/{id}", getUser)
		ir.Post("/", createUser)
		ir.Delete("/", deleteUser)
	})

	http.ListenAndServe(":3000", r)
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

func rootHandler(w http.ResponseWriter, r *http.Request) {
	jsend.Success(w, "hello world", http.StatusOK)
}

func getPeople(w http.ResponseWriter, r *http.Request) {
	jsend.Success(w, "hello people", http.StatusOK)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	jsend.Success(w, "hello users", http.StatusOK)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	jsend.Success(w, "hello user", http.StatusOK)
}

func createPeople(w http.ResponseWriter, r *http.Request) {
	jsend.Success(w, "hello create people", http.StatusCreated)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	jsend.Success(w, "hello create user", http.StatusCreated)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	jsend.Success(w, nil, http.StatusNoContent)
}
