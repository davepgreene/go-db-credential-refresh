package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// Ping the database to verify DSN provided by the user is valid and the
// server accessible. If the ping fails exit the program with an error.
func Ping(ctx context.Context, db *sql.DB) error {
	log.Println("PING")
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	log.Println("PONG")

	return nil
}

// QueryUsers gets a list of all users in the DB
func QueryUsers(ctx context.Context, db *sql.DB, exclude map[string]bool) ([]string, error) {
	if exclude == nil {
		exclude = make(map[string]bool)
	}

	rows, err := db.QueryContext(ctx, `SELECT u.usename FROM pg_catalog.pg_user u;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	users := make([]string, 0)

	for rows.Next() {
		var user string
		if err := rows.Scan(&user); err != nil {
			return nil, err
		}
		// Exclude the user that Vault is using to create new roles
		if ok := exclude[user]; !ok {
			users = append(users, user)
		}
	}

	return users, nil
}
