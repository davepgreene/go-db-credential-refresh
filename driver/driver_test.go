package driver

import (
	"database/sql/driver"
	"errors"
	"strings"
	"testing"
)

const (
	driverName = "a driver"
)

func unregisterAllDrivers() {
	driverMu.Lock()
	defer driverMu.Unlock()
	// For tests.
	driverFactories = make(map[string]factory)
}

func TestAvailableDriversAreRegistered(t *testing.T) {
	registerAllDrivers()
	// This is a brittle test but afaik the only way to test init() behaviors
	for name := range availableDrivers {
		if _, ok := driverFactories[name]; !ok {
			t.Fatalf("driver %s was not registered", name)
		}
	}
}

func TestCantRegisterADriverWithoutAFactory(t *testing.T) {
	unregisterAllDrivers()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected a panic when attempting to register a nil factory")
		}
	}()

	if err := Register("a driver", nil); err != nil {
		t.Error(err)
	}
}

func TestCanRegisterAValidFactory(t *testing.T) {
	unregisterAllDrivers()
	defer func() {
		if r := recover(); r != nil {
			t.Error("did not expect a panic when attempting to register a valid factory")
		}
	}()

	var fn factory = func() (driver.Driver, Formatter, AuthError) {
		return nil, nil, nil
	}

	Register("a driver", fn) //nolint:errcheck
	d := drivers()
	if len(d) != 1 {
		t.Errorf("expected one driver to be registered but got %d", len(d))
	}
}

func TestCantRegisterMultipleFactoriesWithTheSameName(t *testing.T) {
	unregisterAllDrivers()
	var fn factory = func() (driver.Driver, Formatter, AuthError) {
		return nil, nil, nil
	}

	if err := Register(driverName, fn); err != nil {
		t.Error(err)
	}

	if err := Register(driverName, fn); err == nil {
		t.Error("expected an error registering a duplicate driver but didn't get one")
	}
	d := drivers()
	if len(d) != 1 || d[0] != driverName {
		t.Errorf("expected one driver to be registered but got %d", len(d))
	}
}

func TestCanCreateADriverInstance(t *testing.T) {
	unregisterAllDrivers()

	if err := Register("a driver", func() (driver.Driver, Formatter, AuthError) {
		return nil, MysqlFormatter, func(e error) bool { return true }
	}); err != nil {
		t.Error(err)
	}

	ds := drivers()
	if len(ds) != 1 || ds[0] != "a driver" {
		t.Errorf("expected one driver to be registered but got %d", len(ds))
	}

	d, f, a, err := CreateDriver("a driver")
	if err != nil {
		t.Error(err)
	}

	if d != nil {
		t.Errorf("expected a nil driver but got %v", d)
	}

	// test formatter
	if f("user", "pass", "host", 0, "", nil) != MysqlFormatter("user", "pass", "host", 0, "", nil) {
		t.Error("Formatter should be mysqlFormatter but wasn't")
	}

	if !a(errors.New("foo")) {
		t.Error("AuthError should be true but wasn't")
	}
}

func TestCantCreateMissingDriver(t *testing.T) {
	unregisterAllDrivers()

	_, _, _, err := CreateDriver("a driver") //nolint:dogsled
	if err == nil {
		t.Error("expected an error but didn't get one")
	}

	if !strings.Contains(err.Error(), "invalid Driver name") {
		t.Errorf("expected 'invalid Driver name' in error but got %s", err)
	}
}
