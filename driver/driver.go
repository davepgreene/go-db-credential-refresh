package driver

import (
	"database/sql/driver"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/lib/pq"
)

type factory func() (driver.Driver, Formatter, AuthError)

var (
	driverMu         sync.RWMutex
	driverFactories  = make(map[string]factory)
	availableDrivers = map[string]factory{
		"pgx":   pgxDriver,
		"mysql": mysqlDriver,
		"pq":    pqDriver,
	}
)

func init() {
	registerAllDrivers()
}

// Register registers a DB driver
// Note: Register behaves similarly to database/sql.Register except that it doesn't
// panic on duplicate registrations, it just ignores them and continues.
// The reason we Register drivers separately from database/sql is because
// 	a) most DB drivers already call database/sql.Register in an init() func
// 	b) we need to carry a lot more information along with the driver to ensure our
// 	   connector logic works correctly.
func Register(name string, f factory) error {
	driverMu.Lock()
	defer driverMu.Unlock()
	if f == nil {
		panic(fmt.Sprintf("attempted to register driver %s with a nil factory", name))
	}

	_, registered := driverFactories[name]
	if registered {
		return fmt.Errorf("driver factory %s already registered, ignoring", name)
	}
	driverFactories[name] = f

	return nil
}

func registerAllDrivers() {
	for k, v := range availableDrivers {
		if err := Register(k, v); err != nil {
			// We should never, EVER hit this condition. If this happens it means something
			// has fundamentally broken in pgx, pq, or go-mysql.
			panic(err)
		}
	}
}

func drivers() []string {
	driverMu.Lock()
	defer driverMu.Unlock()

	drivers := make([]string, 0, len(driverFactories))
	for k := range driverFactories {
		drivers = append(drivers, k)
	}

	sort.Strings(drivers)
	return drivers
}

// CreateDriver creates a Driver
func CreateDriver(name string) (driver.Driver, Formatter, AuthError, error) {
	driverMu.Lock()

	driverFactory, ok := driverFactories[name]
	if !ok {
		// Factory has not been registered.
		driverMu.Unlock()
		return nil, nil, nil, fmt.Errorf("invalid Driver name. Must be one of: %s", strings.Join(drivers(), ", "))
	}
	defer driverMu.Unlock()

	// Run the factory
	d, f, authError := driverFactory()
	return d, f, authError, nil
}

func mysqlDriver() (driver.Driver, Formatter, AuthError) {
	return &mysql.MySQLDriver{}, MysqlFormatter, MySQLAuthError
}

func pgxDriver() (driver.Driver, Formatter, AuthError) {
	return &stdlib.Driver{}, PgFormatter, PostgreSQLAuthError
}

func pqDriver() (driver.Driver, Formatter, AuthError) {
	return &pq.Driver{}, PgFormatter, PostgreSQLAuthError
}
