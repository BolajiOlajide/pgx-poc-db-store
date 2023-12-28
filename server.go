package main

import (
	"net/http"

	"clevergo.tech/jsend"
	"github.com/BolajiOlajide/pgx-poc-db-store/internal/database"
	"github.com/go-chi/chi/v5"
)

type server struct {
	db     database.DB
	router *chi.Mux
}

func newServer(db database.DB, r *chi.Mux) *server {
	return &server{
		db:     db,
		router: r,
	}
}

func (s *server) start() {
	http.ListenAndServe(":3000", s.router)
}

func (s *server) gracefulShutdown() {
	s.db.Close()
}

func (s *server) setupRoutes() {
	s.router.Get("/", s.rootHandler)
	s.router.Route("/people", func(ir chi.Router) {
		ir.Get("/", s.getPeople)
		ir.Post("/", s.createPeople)
	})
	s.router.Route("/user", func(ir chi.Router) {
		ir.Get("/", s.getUsers)
		ir.Get("/{id}", s.getUser)
		ir.Post("/", s.createUser)
		ir.Delete("/", s.deleteUser)
	})
}

func (s *server) rootHandler(w http.ResponseWriter, r *http.Request) {
	jsend.Success(w, "hello world", http.StatusOK)
}

func (s *server) getPeople(w http.ResponseWriter, r *http.Request) {
	jsend.Success(w, "hello people", http.StatusOK)
}

func (s *server) getUsers(w http.ResponseWriter, r *http.Request) {
	jsend.Success(w, "hello users", http.StatusOK)
}

func (s *server) getUser(w http.ResponseWriter, r *http.Request) {
	jsend.Success(w, "hello user", http.StatusOK)
}

func (s *server) createPeople(w http.ResponseWriter, r *http.Request) {
	jsend.Success(w, "hello create people", http.StatusCreated)
}

func (s *server) createUser(w http.ResponseWriter, r *http.Request) {
	jsend.Success(w, "hello create user", http.StatusCreated)
}

func (s *server) deleteUser(w http.ResponseWriter, r *http.Request) {
	jsend.Success(w, nil, http.StatusNoContent)
}
