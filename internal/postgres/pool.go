package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	if dsn == "" {
		return nil, fmt.Errorf("database url is required")
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse database config: %w", err)
	}

	config.MaxConns = 10
	config.MinConns = 0

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

func ClosePool(ctx context.Context, pool *pgxpool.Pool) error {
	if pool == nil {
		return nil
	}

	closed := make(chan struct{})
	go func() {
		pool.Close()
		close(closed)
	}()

	select {
	case <-closed:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
