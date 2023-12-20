package database

import (
	"context"

	"github.com/BolajiOlajide/pgx-poc-db-store/internal/types"
)

type PeopleStore interface {
	List(ctx context.Context, opts ListUserArgs) ([]*types.People, error)
	Get(ctx context.Context, id int64) (*types.People, error)
	Create(ctx context.Context, people *types.People) error
}

type ListPeopleArgs struct {
	Limit int
}
