package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type Database struct {
	Client *sqlx.DB
}

// NewDatabase - returns a pointer to a database object
func NewDatabase() (*Database, error) {
	log.Info("Setting up new database connection")

	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DB"),
	)

	db, err := sqlx.Open("mysql", connectionString)
	if err != nil {
		return &Database{}, fmt.Errorf("could not connect to database: %w", err)
	}

	log.Info("Setting up new database connection successfully")

	return &Database{
		Client: db,
	}, nil
}

func (d *Database) Ping(ctx context.Context) error {
	return d.Client.DB.PingContext(ctx)
}
