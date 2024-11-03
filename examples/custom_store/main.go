package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
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

	maxOpenConns = 5
	maxIdleConns = 2
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

	log.Fatal(err)
}

func Run() error {
	var err error

	u, serverClose := setupTestServer()
	defer serverClose()

	s, err := NewHTTPTestConnectingStore(
		u,
		"GET",
		nil,
		func(r *http.Response) (driver.Credentials, error) {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				return nil, err
			}

			var v map[string]any
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
		},
	)
	if err != nil {
		return err
	}

	c, err := driver.NewConnector(s, "pgx", &driver.Config{
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
	database.SetMaxIdleConns(maxIdleConns)
	database.SetMaxOpenConns(maxOpenConns)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appSignal := make(chan os.Signal, 2)
	signal.Notify(appSignal, os.Interrupt)

	go func() {
		<-appSignal
		cancel()
	}()

	for {
		// First ping the DB to open a connection
		err = db.Ping(ctx, database)
		if err != nil {
			log.Println(err)

			break
		}

		// Sleep long enough that the creds should expire
		time.Sleep(3 * time.Second)

		// Now get users
		users, err := db.QueryUsers(ctx, database, nil)
		if err != nil {
			log.Println(err)

			break
		}

		log.Printf("Users: %v", users)
	}

	return err
}

func setupTestServer() (string, func()) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		log.Println("Server: retrieved credentials")
		w.Header().Set("Content-Type", "application/json")
		resp := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password)
		if err := json.NewEncoder(w).Encode(json.RawMessage(resp)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))

	return ts.URL, ts.Close
}
