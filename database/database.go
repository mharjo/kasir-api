package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDBPool(conn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(conn)
	if err != nil {
		return nil, err
	}

	// ✅ BEST PRACTICE for Supabase Transaction Pooler (PgBouncer):
	// Disable prepared statements / statement cache by using Simple Protocol.
	cfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	// Pool tuning (safe defaults)
	cfg.MaxConns = 5
	cfg.MinConns = 1
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.HealthCheckPeriod = 30 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// ✅ Real DB check (not Ping)
	if err := DBCheck(ctx, pool); err != nil {
		pool.Close()
		return nil, fmt.Errorf("db check failed: %w", err)
	}

	log.Println("Database connected successfully (pgxpool + simple protocol)")
	return pool, nil
}

func DBCheck(ctx context.Context, pool *pgxpool.Pool) error {
	cctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return pool.QueryRow(cctx, "select 1").Scan(new(int))
}
