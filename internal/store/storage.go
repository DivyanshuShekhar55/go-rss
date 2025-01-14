package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource Not Found")
	QueryTimeoutDuration = time.Second * 5
	ErrConflict          = errors.New("resource already exists")
)

type Storage struct {
	Users interface {
		GetByID(context.Context, int64) (*User, error)
		Create(context.Context, *sql.Tx, *User) (*User, error)
		Delete(context.Context, int64) error
	}
}

// create a function that takes in a func 'fn' as an input (update, delete, create)
// tries to perform fn with the transaction mechanism - if fails all changes made are rolled-back
func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// apply the function ...
	if err := fn(tx); err != nil {

		//if err undo all changes
		_ = tx.Rollback()
		return err
	}

	// if no error actually change the db
	return tx.Commit()
}
