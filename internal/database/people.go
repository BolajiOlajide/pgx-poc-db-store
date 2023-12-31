package database

import (
	"context"
	"errors"

	"github.com/BolajiOlajide/pgx-poc-db-store/internal/database/basestore"
	"github.com/BolajiOlajide/pgx-poc-db-store/internal/database/dbutil"
	"github.com/BolajiOlajide/pgx-poc-db-store/internal/types"
	"github.com/keegancsmith/sqlf"
)

var peopleColumns = []*sqlf.Query{
	sqlf.Sprintf("id"),
	sqlf.Sprintf("user_id"),
}

var peopleInsertColumns = []*sqlf.Query{
	sqlf.Sprintf("user_id"),
}

type PeopleStore interface {
	Create(ctx context.Context, userID string) (*types.People, error)
}

type peopleStore struct {
	*basestore.Store
}

var _ PeopleStore = &peopleStore{}

func PeopleWith(other basestore.ShareableStore) PeopleStore {
	return &peopleStore{Store: basestore.NewWithHandle(other.Handle())}
}

const peopleCreateQueryFmtStr = `
INSERT INTO
	people (%s)
	VALUES (%s)
	RETURNING %s
`

func (p *peopleStore) Create(ctx context.Context, userID string) (*types.People, error) {
	if userID == "" {
		return nil, errors.New("no user id provided")
	}

	q := sqlf.Sprintf(
		peopleCreateQueryFmtStr,
		sqlf.Join(peopleInsertColumns, ", "),
		userID,
		sqlf.Join(peopleColumns, ", "),
	)
	return scanPeople(p.QueryRow(ctx, q))
}

func scanPeople(sc dbutil.Scanner) (*types.People, error) {
	var person types.People
	if err := sc.Scan(
		&person.ID,
		&person.UserID,
	); err != nil {
		return nil, err
	}

	return &person, nil
}
