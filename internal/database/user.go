package database

import (
	"context"
	"errors"

	"github.com/BolajiOlajide/pgx-poc-db-store/internal/database/basestore"
	"github.com/BolajiOlajide/pgx-poc-db-store/internal/database/dbutil"
	"github.com/BolajiOlajide/pgx-poc-db-store/internal/types"
	"github.com/jackc/pgx/v5"
	"github.com/keegancsmith/sqlf"
)

var userColumns = []*sqlf.Query{
	sqlf.Sprintf("users.id"),
	sqlf.Sprintf("users.username"),
	sqlf.Sprintf("users.email"),
}

var userInsertColumns = []*sqlf.Query{
	sqlf.Sprintf("username"),
	sqlf.Sprintf("email"),
}

type UserStore interface {
	basestore.ShareableStore

	List(ctx context.Context, opts ListUserArgs) ([]*types.User, error)
	Get(ctx context.Context, userID string) (*types.User, error)
	// Create(ctx context.Context, user *types.User) error
}

type ListUserArgs struct {
	Limit int
}

func UsersWith(other basestore.ShareableStore) UserStore {
	return &userStore{Store: basestore.NewWithHandle(other.Handle())}
}

type userStore struct {
	*basestore.Store
}

var _ UserStore = &userStore{}

const listUsersFmtStr = `
SELECT %s FROM users
`

func (u *userStore) List(ctx context.Context, opts ListUserArgs) ([]*types.User, error) {
	var users []*types.User

	scanUserFunc := func(rows pgx.Rows) error {
		user, err := scanUser(rows)
		if err != nil {
			return err
		}
		users = append(users, user)
		return nil
	}

	rows, err := u.Query(ctx, sqlf.Sprintf(listUsersFmtStr, sqlf.Join(userColumns, ", ")))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := scanUserFunc(rows); err != nil {
			return nil, err
		}
	}

	return users, rows.Err()
}

const getUserFmtStr = `
SELECT %s FROM users
WHERE %s
LIMIT 1`

func (u *userStore) Get(ctx context.Context, userID string) (*types.User, error) {
	if userID == "" {
		return nil, errors.New("no user id provided")
	}

	q := sqlf.Sprintf(
		getUserFmtStr,
		sqlf.Join(userColumns, ", "),
		sqlf.Sprintf("id = %s", userID),
	)

	row, err := u.QueryRow(ctx, q)
	if err != nil {
		return nil, err
	}

	return scanUser(row)
}

func scanUser(sc dbutil.Scanner) (*types.User, error) {
	var user types.User
	if err := sc.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
	); err != nil {
		return nil, err
	}

	return &user, nil
}
