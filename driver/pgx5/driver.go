package pgx5

import (
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/davepgreene/go-db-credential-refresh/driver"
)

const pgxName = "pgx"

// WithPgxDriver registers the pgx v5 driver
func WithPgxDriver() error {
	return driver.Register(pgxName, pgxDriver)
}

func pgxDriver() *driver.Driver {
	return &driver.Driver{
		Driver:    &stdlib.Driver{},
		Formatter: driver.PgFormatter,
		AuthError: driver.PostgreSQLAuthError,
	}
}
