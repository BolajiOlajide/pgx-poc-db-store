// SOURCE: https://sourcegraph.com/github.com/sourcegraph/sourcegraph/-/blob/internal/database/basestore/store.go
package basestore

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/keegancsmith/sqlf"
)

type Store struct {
	handle TransactableHandle
}

// ShareableStore is implemented by stores to explicitly allow distinct store instances
// to reference the store's underlying handle. This is used to share transactions between
// multiple stores. See `Store.With` for additional details.
type ShareableStore interface {
	// Handle returns the underlying transactable database handle.
	Handle() TransactableHandle
}

var _ ShareableStore = &Store{}

// NewWithHandle returns a new base store using the given database handle.
func NewWithHandle(handle TransactableHandle) *Store {
	return &Store{handle: handle}
}

// Handle returns the underlying transactable database handle.
func (s *Store) Handle() TransactableHandle {
	return s.handle
}

// With creates a new store with the underlying database handle from the given store.
// This method should be used when two distinct store instances need to perform an
// operation within the same shared transaction.
//
//	txn1 := store1.Transact(ctx) // Creates a transaction
//	txn2 := store2.With(txn1)    // References the same transaction
//
//	txn1.A(ctx) // Occurs within shared transaction
//	txn2.B(ctx) // Occurs within shared transaction
//	txn1.Done() // closes shared transaction
//
// Note that once a handle is shared between two stores, committing or rolling back
// a transaction will affect the handle of both stores. Most notably, two stores that
// share the same handle are unable to begin independent transactions.
func (s *Store) With(other ShareableStore) *Store {
	return &Store{handle: other.Handle()}
}

// Query performs QueryContext on the underlying connection.
func (s *Store) Query(ctx context.Context, query *sqlf.Query) (pgx.Rows, error) {
	return s.handle.Query(ctx, query.Query(sqlf.PostgresBindVar), query.Args()...)
}

// QueryRow performs QueryRowContext on the underlying connection.
func (s *Store) QueryRow(ctx context.Context, query *sqlf.Query) (pgx.Row, error) {
	return s.handle.QueryRow(ctx, query.Query(sqlf.PostgresBindVar), query.Args()...)
}

// Exec performs a query without returning any rows.
func (s *Store) Exec(ctx context.Context, query *sqlf.Query) error {
	_, err := s.ExecResult(ctx, query)
	return err
}

// ExecResult performs a query without returning any rows, but includes the
// result of the execution.
func (s *Store) ExecResult(ctx context.Context, query *sqlf.Query) (pgconn.CommandTag, error) {
	return s.handle.Exec(ctx, query.Query(sqlf.PostgresBindVar), query.Args()...)
}

// SetLocal performs the `SET LOCAL` query and returns a function to clear (aka to empty string) the setting.
// Calling this method only makes sense within a transaction, as the setting is unset after the transaction
// is either rolled back or committed. This does not perform argument parameterization.
func (s *Store) SetLocal(ctx context.Context, key, value string) (func(context.Context) error, error) {
	if !s.InTransaction() {
		return func(ctx context.Context) error { return nil }, ErrNotInTransaction
	}

	return func(ctx context.Context) error {
		return s.Exec(ctx, sqlf.Sprintf(fmt.Sprintf(`SET LOCAL "%s" TO ''`, key)))
	}, s.Exec(ctx, sqlf.Sprintf(fmt.Sprintf(`SET LOCAL "%s" TO "%s"`, key, value)))
}

// InTransaction returns true if the underlying database handle is in a transaction.
func (s *Store) InTransaction() bool {
	return s.handle.InTransaction()
}

// Transact returns a new store whose methods operate within the context of a new transaction
// or a new savepoint. This method will return an error if the underlying connection cannot be
// interface upgraded to a TxBeginner.
func (s *Store) Transact(ctx context.Context) (*Store, error) {
	handle, err := s.handle.Transact(ctx)
	if err != nil {
		return nil, err
	}

	return &Store{handle: handle}, nil
}

// Done performs a commit or rollback of the underlying transaction/savepoint depending
// on the value of the error parameter. The resulting error value is a multierror containing
// the error parameter along with any error that occurs during commit or rollback of the
// transaction/savepoint. If the store does not wrap a transaction the original error value
// is returned unchanged.
func (s *Store) Done(ctx context.Context, err error) error {
	return s.handle.Done(ctx, err)
}

var ErrPanicDuringTransaction = errors.New("encountered panic during transaction")

// WithTransact executes the callback using a transaction on the store. If the callback
// returns an error or panics, the transaction will be rolled back.
func (s *Store) WithTransact(ctx context.Context, f func(tx *Store) error) error {
	return InTransaction[*Store](ctx, s, f)
}

// InTransaction executes the callback using a transaction on the given transactable store. If
// the callback returns an error or panics, the transaction will be rolled back.
func InTransaction[T Transactable[T]](ctx context.Context, t Transactable[T], f func(tx T) error) (err error) {
	tx, err := t.Transact(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			// If we're panicking, roll back the transaction
			// even when err is nil.
			err = tx.Done(ctx, ErrPanicDuringTransaction)
			// Re-throw the panic after rolling back the transaction
			panic(r)
		} else {
			// If we're not panicking, roll back the transaction if the
			// operation on the transaction failed for whatever reason.
			err = tx.Done(ctx, err)
		}
	}()

	return f(tx)
}
