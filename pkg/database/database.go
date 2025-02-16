package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresDB struct {
	pool *pgxpool.Pool
}

const (
	defaultMaxConns          = 10
	defaultMinConns          = 2
	defaultMaxConnLifetime   = time.Hour
	defaultMaxConnIdleTime   = time.Minute * 30
	defaultHealthCheckPeriod = time.Minute
)

type PoolConfig struct {
	MaxConns          int32
	MinConns          int32
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
}

type DB interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (int64, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
	Ping(ctx context.Context) error
	Close()
}

func New(ctx context.Context, connStr string, cfg *PoolConfig) (DB, error) {
	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: failed to parse connection string: %w", err)
	}

	if cfg != nil {
		applyConfig(poolConfig, cfg)
	} else {
		poolConfig.MaxConns = defaultMaxConns
		poolConfig.MinConns = defaultMinConns
		poolConfig.MaxConnLifetime = defaultMaxConnLifetime
		poolConfig.MaxConnIdleTime = defaultMaxConnIdleTime
		poolConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewWithConfig: failed to connect to database: %w", err)
	}

	return &postgresDB{pool: pool}, nil
}

func (db *postgresDB) Exec(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, fmt.Errorf("context error: %w", err)
	}

	tag, err := db.pool.Exec(ctx, sql, args...)
	if err != nil {
		return 0, fmt.Errorf("db.pool.Exec: failed to execute query: %w", err)
	}
	return tag.RowsAffected(), nil
}

func (db *postgresDB) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context error: %w", err)
	}

	rows, err := db.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("db.pool.Query: failed to execute query: %w", err)
	}
	return rows, nil
}

func (db *postgresDB) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	if err := ctx.Err(); err != nil {
		return nil
	}

	return db.pool.QueryRow(ctx, sql, args...)
}

func (db *postgresDB) BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("context error: %w", err)
	}

	tx, err := db.pool.BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("db.pool.BeginTx: failed to begin transaction: %w", err)
	}
	return tx, nil
}

func (db *postgresDB) Ping(ctx context.Context) error {
	if err := db.pool.Ping(ctx); err != nil {
		return fmt.Errorf("db.pool.Ping: database connection is not alive: %w", err)
	}
	return nil
}

func (db *postgresDB) Close() {
	db.pool.Close()
}

func applyConfig(poolConfig *pgxpool.Config, cfg *PoolConfig) {
	if cfg.MaxConns > 0 {
		poolConfig.MaxConns = cfg.MaxConns
	}
	if cfg.MinConns > 0 {
		poolConfig.MinConns = cfg.MinConns
	}
	if cfg.MaxConnLifetime > 0 {
		poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
	}
	if cfg.MaxConnIdleTime > 0 {
		poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	}
	if cfg.HealthCheckPeriod > 0 {
		poolConfig.HealthCheckPeriod = cfg.HealthCheckPeriod
	}
}
