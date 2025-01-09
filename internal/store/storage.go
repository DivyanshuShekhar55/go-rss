package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Users interface {
		GetByID(context.Context,int64) (*User, error)
		Create(context.Context, *sql.Tx, *User) (*User, error)
		Delete(context.Context, int64) error
	}
}