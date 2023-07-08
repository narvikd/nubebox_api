package dbconfig

import (
	"api/internal/cfg"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // imports postgres driver for side effects
	"github.com/narvikd/errorskit"
	"time"
)

// InitDB returns a DB pool pointer (*sql.DB), it needs to be called on init or main
func InitDB(config *cfg.Config) (*sql.DB, error) {
	conn, errOpen := sql.Open("postgres", getDbURL(config))
	if errOpen != nil {
		return nil, errorskit.Wrap(errOpen, "there was a problem connecting with the DB on InitDB")
	}

	// See "Important settings" section (SQLDB). More info in: https://www.alexedwards.net/blog/configuring-sqldb
	conn.SetConnMaxLifetime(time.Minute * 3)
	conn.SetMaxOpenConns(25) // Original was 10
	conn.SetMaxIdleConns(25) // Original was 10

	if errPing := Ping(conn); errPing != nil {
		return nil, errPing
	}
	return conn, nil
}

// Ping sends a ping to the DB pool to make sure the connection was successful / is alive.
func Ping(db *sql.DB) error {
	if errPing := db.Ping(); errPing != nil {
		return errorskit.Wrap(errPing, "could not ping the DB. Maybe the DB connection is down")
	}
	return nil
}

// CloseDB only logs its error. Reason: https://stackoverflow.com/questions/50787804/does-db-close-need-to-be-called
func CloseDB(db *sql.DB) error {
	if err := db.Close(); err != nil {
		return errorskit.Wrap(err, "there was a problem closing the DB connection")
	}
	return nil
}

func getDbURL(config *cfg.Config) string {
	// TODO: Refactor for prod
	const (
		password = ""
		sslMode  = "disable"
	)
	c := config.Database
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User,
		password,
		c.Ip,
		c.Port,
		c.DBName,
		sslMode,
	)
}
