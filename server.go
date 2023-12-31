package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"clevergo.tech/jsend"
	"github.com/BolajiOlajide/pgx-poc-db-store/internal/database"
	"github.com/BolajiOlajide/pgx-poc-db-store/internal/types"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
		ir.Get("/{userID}", s.getUser)
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
	users, err := s.db.Users().List(r.Context(), database.ListUserArgs{})
	if err != nil {
		jsend.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsend.Success(w, users, http.StatusOK)
}

func (s *server) getUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		jsend.Error(w, "userID not provided", http.StatusBadRequest)
		return
	}

	_, err := uuid.Parse(userID)
	if err != nil {
		jsend.Error(w, "invalid uuid", http.StatusBadRequest)
		return
	}

	user, err := s.db.Users().GetByID(r.Context(), userID)
	if err != nil {
		jsend.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsend.Success(w, user, http.StatusOK)
}

func (s *server) createUser(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Username string `json:"username"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		jsend.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.Email == "" {
		jsend.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	if body.Username == "" {
		jsend.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	var newUser *types.User
	status := http.StatusBadRequest
	err = s.db.WithTransact(ctx, func(tx database.DB) error {
		user, err := tx.Users().GetByEmail(ctx, body.Email)
		if err != nil && !database.IsUserNotFoundErr(err) {
			return err
		}

		if user != nil {
			status = http.StatusConflict
			return errors.New("user with that email exists")
		}

		newUser, err = tx.Users().Create(ctx, body.Email, body.Username)
		if err != nil {
			return err
		}

		fmt.Println("created user, ===>>", newUser.ID)

		_, err = tx.People().Create(ctx, newUser.ID)
		if err != nil {
			return err
		}

		status = http.StatusOK
		return nil
	})
	if err != nil {
		jsend.Error(w, err.Error(), status)
		return
	}

	jsend.Success(w, newUser, status)
}

func (s *server) createPeople(w http.ResponseWriter, r *http.Request) {
	jsend.Success(w, "hello create people", http.StatusCreated)
}

func (s *server) deleteUser(w http.ResponseWriter, r *http.Request) {
	jsend.Success(w, nil, http.StatusNoContent)
}
