package pgx4

import (
	"github.com/jackc/pgx/v4/stdlib"

	"github.com/davepgreene/go-db-credential-refresh/driver"
)

const pgxName = "pgxv4"

// WithPgxDriver registers the pgx v4 driver
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
