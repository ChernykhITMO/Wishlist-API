package postgres

import (
	"context"
	"fmt"
	"github.com/ChernykhITMO/Wishlist-API/internal/config"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg config.DBConfig) (*repository, error) {
	const op = "internal.platform.postgres.Open"

	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("%s: parse config: %w", op, err)
	}

	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("%s: new pool: %w", op, err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("%s: ping: %w", op, err)
	}

	return &repository{pool: pool}, nil
}

func (r *repository) Register(ctx context.Context, id uuid.UUID, email, password string) error {
	const query = `
	INSERT INTO users (id, email, password_hash)
	VALUES ($1, $2, $3)
	`

	if _, err := r.pool.Exec(ctx, query, id, email, password); err != nil {
		return err
	}

	return nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (uuid.UUID, string, error) {
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
