package postgres

import (
	"context"
	"fmt"

	"github.com/ChernykhITMO/Wishlist-API/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, cfg config.DBConfig) (*pgxpool.Pool, error) {
	const op = "services.platform.postgres.NewPool"

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

	return pool, nil
}
