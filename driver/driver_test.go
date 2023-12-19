package driver

import (
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
			t.Fatal("expected a panic when attempting to register a nil factory")
		}
	}()

	if err := Register("a driver", nil); err != nil {
		t.Fatal(err)
	}
}

func TestCanRegisterAValidFactory(t *testing.T) {
	unregisterAllDrivers()
	defer func() {
		if r := recover(); r != nil {
			t.Fatal("did not expect a panic when attempting to register a valid factory")
		}
	}()

	var fn factory = func() *Driver {
		return &Driver{}
	}

	Register("a driver", fn) //nolint:errcheck
	d := drivers()
	if len(d) != 1 {
		t.Fatalf("expected one driver to be registered but got %d", len(d))
	}
}

func TestCantRegisterMultipleFactoriesWithTheSameName(t *testing.T) {
	unregisterAllDrivers()
	var fn factory = func() *Driver {
		return &Driver{}
	}

	if err := Register(driverName, fn); err != nil {
		t.Fatal(err)
	}

	if err := Register(driverName, fn); err == nil {
		t.Fatal("expected an error registering a duplicate driver but didn't get one")
	}
	d := drivers()
	if len(d) != 1 || d[0] != driverName {
		t.Fatalf("expected one driver to be registered but got %d", len(d))
	}
}

func TestCanCreateADriverInstance(t *testing.T) {
	unregisterAllDrivers()

	if err := Register("a driver", func() *Driver {
		return &Driver{
			Driver:    nil,
			Formatter: MysqlFormatter,
			AuthError: func(e error) bool { return true },
		}
	}); err != nil {
		t.Fatal(err)
	}

	ds := drivers()
	if len(ds) != 1 || ds[0] != "a driver" {
		t.Fatalf("expected one driver to be registered but got %d", len(ds))
	}

	d, err := CreateDriver("a driver")
	if err != nil {
		t.Fatal(err)
	}

	if d.Driver != nil {
		t.Fatalf("expected a nil driver but got %v", d)
	}

	// test formatter
	expectedFormatter := MysqlFormatter("user", "pass", "host", 0, "", nil)
	if d.Formatter("user", "pass", "host", 0, "", nil) != expectedFormatter {
		t.Fatal("Formatter should be mysqlFormatter but wasn't")
	}

	if !d.AuthError(errors.New("foo")) {
		t.Fatal("AuthError should be true but wasn't")
	}
}

func TestCantCreateMissingDriver(t *testing.T) {
	unregisterAllDrivers()

	_, err := CreateDriver("a driver") //nolint:dogsled
	if err == nil {
		t.Fatal("expected an error but didn't get one")
	}

	if !strings.Contains(err.Error(), "invalid Driver name") {
		t.Fatalf("expected 'invalid Driver name' in error but got %s", err)
	}
}
