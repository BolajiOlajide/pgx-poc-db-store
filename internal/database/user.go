package database

import (
	"context"

	"github.com/BolajiOlajide/pgx-poc-db-store/internal/types"
)

type UserStore interface {
	List(ctx context.Context, opts ListUserArgs) ([]*types.User, error)
	Get(ctx context.Context, id string) (*types.User, error)
	Create(ctx context.Context, user *types.User) error
}

type ListUserArgs struct {
	Limit int
}
