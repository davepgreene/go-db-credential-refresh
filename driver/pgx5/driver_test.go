package pgx5

import (
	"context"
	"slices"
	"testing"

	"github.com/davepgreene/go-db-credential-refresh/driver"
)

func TestCanRegisterPgxV5Driver(t *testing.T) {
	driver.UnregisterAll()

	if err := WithPgxDriver(); err != nil {
		t.Fatal(err)
	}

	drivers := driver.Drivers()
	if !slices.Contains(drivers, "pgx") {
		t.Fatal("`pgx` was not found in registered driver list")
	}
}

func TestFailureIfDriverAlreadyRegistered(t *testing.T) {
	driver.UnregisterAll()

	if err := WithPgxDriver(); err != nil {
		t.Fatal(err)
	}
	if err := WithPgxDriver(); err == nil {
		t.Fatal("expected an error from trying to register an already registered driver but didn't get one")
	}
}

func TestFailureIfDriverNotRegistered(t *testing.T) {
	driver.UnregisterAll()

	s := &noopStore{}

	_, err := driver.NewConnector(s, "pgx", &driver.Config{
		Opts: nil,
	})
	if err == nil {
		t.Fatal("expected an error because the driver wasn't registered but didn't get one")
	}
}

type noopStore struct{}

func (s *noopStore) Get(ctx context.Context) (driver.Credentials, error) {
	return nil, nil
}

func (s *noopStore) Refresh(ctx context.Context) (driver.Credentials, error) {
	return nil, nil
}
