package store

import (
	"context"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// NOTE : the following is not a SQL structure, it is a Go struct for encoding/decoding json data
type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
}

// text is un-hashed password
type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.text = &text
	p.hash = hash
	return nil
}

func (p *password) Compare(text string) error {
	return bcrypt.CompareHashAndPassword(p.
		hash, []byte(text))
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `INSERT INTO users (username, password, email) VALUES ($1, $2, $3) RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password.hash,
		user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	if err != nil {
		// switch {
		// // check for errors from pgx driver and do it
		// default:
		// 	return err
		// }
		return err
	}
	return nil
}

func (s *UserStore) GetByID(ctx context.Context, userID int64) (*User, error) {
	query := `
		SELECT users.id, username, email, password, created_at
		FROM users
		WHERE users.id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	user := &User{}

	err := s.db.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		switch err {
		// look for pgx errors
		default:
			return nil, err
		}
	}
	return user, nil
}

// TODO : think of a better name when you understand what the func is actually doing
func (s *UserStore) CreateWithTx(ctx context.Context, user *User) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}
		return nil
	})
}
