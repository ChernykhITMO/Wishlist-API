package repository

import (
	"context"
	"errors"
	"github.com/ChernykhITMO/Wishlist-API/internal/auth/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const uniqueViolationCode = "23505"

type repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
}

func (r *repository) Create(ctx context.Context, id uuid.UUID, email, password string) error {
	const query = `
	INSERT INTO users (id, email, password_hash)
	VALUES ($1, $2, $3)
	`

	_, err := r.pool.Exec(ctx, query, id, email, password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == uniqueViolationCode {
				return domain.ErrEmailAlreadyExists
			}
		}
		return err
	}

	return nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (uuid.UUID, string, error) {
	const query = `SELECT id, password_hash FROM users WHERE email=$1`
	var (
		id       uuid.UUID
		password string
	)

	if err := r.pool.QueryRow(ctx, query, email).Scan(&id, &password); err != nil {
		return uuid.Nil, "", err
	}

	return id, password, nil
}
