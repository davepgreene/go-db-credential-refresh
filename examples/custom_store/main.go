package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"time"

	"github.com/davepgreene/go-db-credential-refresh/driver"
	"github.com/davepgreene/go-db-credential-refresh/examples/db"
	"github.com/davepgreene/go-db-credential-refresh/store"
)

const (
	role     = "role"
	host     = "localhost"
	username = "postgres"
	password = "postgres"
	dbName   = "postgres"
	port     = 5432
)

// This example assumes a PostgreSQL database running on localhost:5432
// docker run --name pg -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres
func main() {
	err := Run()
	if err == nil {
		return
	}

	if err == context.Canceled {
		return
	}

	fmt.Println(err)
	os.Exit(1)
}

func Run() error {
	var err error

	url, close := setupTestServer()
	defer close()

	store, err := NewHTTPTestConnectingStore(url, "GET", nil, func(r *http.Response) (driver.Credentials, error) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		var v map[string]interface{}
		if err := json.Unmarshal([]byte(b), &v); err != nil {
			return nil, err
		}

		username, ok := v["username"].(string)
		if !ok {
			return nil, errors.New("missing username")
		}

		password, ok := v["password"].(string)
		if !ok {
			return nil, errors.New("missing password")
		}

		return &store.Credential{
			Username: username,
			Password: password,
		}, nil
	})
	if err != nil {
		return err
	}

	c, err := driver.NewConnector(store, "pgx", &driver.Config{
		Host: host,
		Port: port,
		DB:   dbName,
		Opts: map[string]string{
			"sslmode": "disable",
		},
	})
	if err != nil {
		return err
	}

	// Use the built in database/sql package to work with the connector
	database := sql.OpenDB(c)

	// In order to demonstrate the creation and revocation of roles we need to set the
	// connection lifetime very short. In a production environment, Vault role TTLs and
	// connection lifetime should be tuned based on database performance requirements.
	database.SetConnMaxLifetime(2 * time.Second)
	database.SetMaxIdleConns(2)
	database.SetMaxOpenConns(5)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appSignal := make(chan os.Signal)
	signal.Notify(appSignal, os.Interrupt)

	go func() {
		<-appSignal
		cancel()
	}()

	for {
		// First ping the DB to open a connection
		err = db.Ping(ctx, database)
		if err != nil {
			fmt.Println(err)
			break
		}

		// Sleep long enough that the creds should expire
		time.Sleep(3 * time.Second)

		// Now get users
		users, err := db.QueryUsers(ctx, database, nil)
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println(users)
	}

	return err
}

func setupTestServer() (string, func()) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password)
		json.NewEncoder(w).Encode(json.RawMessage(resp))
	}))

	return ts.URL, ts.Close
}
