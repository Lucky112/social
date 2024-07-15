package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	pool *pgxpool.Pool
}

func ViaPGX(ctx context.Context, cfg Config) (Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.connectionURL())
	if err != nil {
		return Pool{}, fmt.Errorf("creating new pgx pool: %v", err)
	}

	return Pool{pool}, nil
}

func (p Pool) Close() {
	p.pool.Close()
}
