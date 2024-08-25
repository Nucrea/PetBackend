package integrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
)

// TODO: wrapper, connection pool
func NewPostgresConn(ctx context.Context, connUrl string) (*sql.DB, error) {
	connConf, err := pgx.ParseConnectionString(connUrl)
	if err != nil {
		return nil, fmt.Errorf("failed parsing postgres connection string: %v", err)
	}

	sqlDb := stdlib.OpenDB(connConf)
	if err := sqlDb.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed pinging postgres db: %v", err)
	}

	return sqlDb, nil
}
