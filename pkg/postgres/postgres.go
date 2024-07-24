package postgres

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

const (
	defaultMaxPoolSize  = 4
	defaultConnAttempts = 10
	defaultConnTimeout  = 2 * time.Second
)

type PgxPool interface {
	Close()
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration
	Builder      squirrel.StatementBuilderType // Генератор sql запросов, не является orm!
	Pool         PgxPool
}

func NewPG(url string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  defaultMaxPoolSize,
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}

	for _, option := range opts {
		option(pg)
	}

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	poolConfig.MaxConns = int32(pg.maxPoolSize)
	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}
		log.Printf("Postgres trying to connect, attemps left: %d", pg.connAttempts)
		time.Sleep(pg.connTimeout)
		pg.connAttempts--
	}
	if err != nil {
		return nil, fmt.Errorf("error connect to postgres, %w", err)
	}
	return pg, err
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
