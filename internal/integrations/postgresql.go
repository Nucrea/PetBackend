package integrations

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type SqlDB interface {
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Close() error
	Conn(ctx context.Context) (*sql.Conn, error)
	Driver() driver.Driver
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Ping() error
	PingContext(ctx context.Context) error
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	SetConnMaxIdleTime(d time.Duration)
	SetConnMaxLifetime(d time.Duration)
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
	Stats() sql.DBStats
}

type Postgres struct {
	*sql.DB
	pool *pgxpool.Pool
}

func NewPostgresConn(ctx context.Context, postgresUrl string) (SqlDB, error) {
	connUrl, err := url.Parse(postgresUrl)
	if err != nil {
		return nil, err
	}

	query, _ := url.ParseQuery(connUrl.RawQuery)
	query.Set("pool_max_conns", "16")
	connUrl.RawQuery = query.Encode()

	pool, err := pgxpool.New(ctx, connUrl.String())
	if err != nil {
		return nil, err
	}

	db := stdlib.OpenDBFromPool(pool)
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed pinging postgres db: %v", err)
	}

	return &Postgres{
		DB:   db,
		pool: pool,
	}, nil
}
