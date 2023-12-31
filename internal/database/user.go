package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/BolajiOlajide/pgx-poc-db-store/internal/database/basestore"
	"github.com/BolajiOlajide/pgx-poc-db-store/internal/database/dbutil"
	"github.com/BolajiOlajide/pgx-poc-db-store/internal/types"
	"github.com/jackc/pgx/v5"
	"github.com/keegancsmith/sqlf"
)

const defaultUserLimit = 10

type userNotFoundErr struct {
	ID    string
	Email string
}

func (e userNotFoundErr) Error() string {
	if e.ID != "" {
		return fmt.Sprintf("user with ID %s not found", e.ID)
	}
	return fmt.Sprintf("user with email %s not found", e.Email)
}

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
	GetByID(ctx context.Context, userID string) (*types.User, error)
	GetByEmail(ctx context.Context, email string) (*types.User, error)
	Create(ctx context.Context, email string, username string) (*types.User, error)
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
LIMIT %s
`

func (u *userStore) List(ctx context.Context, opts ListUserArgs) ([]*types.User, error) {
	if opts.Limit == 0 {
		opts.Limit = defaultUserLimit
	}

	query := sqlf.Sprintf(
		listUsersFmtStr,
		sqlf.Join(userColumns, ", "),
		opts.Limit,
	)

	rows, err := u.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*types.User{}

	scanUserFunc := func(rows pgx.Rows) error {
		user, err := scanUser(rows)
		if err != nil {
			return err
		}
		users = append(users, user)
		return nil
	}

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

func (u *userStore) get(ctx context.Context, whereClause *sqlf.Query) (*types.User, error) {
	q := sqlf.Sprintf(
		getUserFmtStr,
		sqlf.Join(userColumns, ", "),
		whereClause,
	)

	return scanUser(u.QueryRow(ctx, q))
}

func (u *userStore) GetByID(ctx context.Context, userID string) (*types.User, error) {
	if userID == "" {
		return nil, errors.New("no user id provided")
	}

	user, err := u.get(ctx, sqlf.Sprintf("id = %s", userID))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, userNotFoundErr{ID: userID}
		}
		return nil, err
	}
	return user, nil
}

func (u *userStore) GetByEmail(ctx context.Context, email string) (*types.User, error) {
	if email == "" {
		return nil, errors.New("no email provided")
	}

	user, err := u.get(ctx, sqlf.Sprintf("email = %s", email))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, userNotFoundErr{Email: email}
		}
		return nil, err
	}
	return user, nil
}

const userCreateQueryFmtStr = `
INSERT INTO
	users (%s)
	VALUES (
		%s,
		%s
	)
	RETURNING %s
`

func (u *userStore) Create(ctx context.Context, email string, username string) (*types.User, error) {
	if email == "" {
		return nil, errors.New("no email provided")
	}

	if username == "" {
		return nil, errors.New("no username provided")
	}

	q := sqlf.Sprintf(
		userCreateQueryFmtStr,
		sqlf.Join(userInsertColumns, ", "),
		username,
		email,
		sqlf.Join(userColumns, ", "),
	)

	return scanUser(u.QueryRow(ctx, q))
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

func IsUserNotFoundErr(err error) bool {
	_, ok := err.(userNotFoundErr)
	return ok
}
