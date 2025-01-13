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
