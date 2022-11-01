package database

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func (db *Database) MigrateDB() (*migrate.Migrate, error) {
	driver, err := mysql.WithInstance(db.Client.DB, &mysql.Config{})

	if err != nil {
		return nil, fmt.Errorf("could not create the mysql driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "mysql", driver)

	if err != nil {
		return nil, err
	}

	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return nil, fmt.Errorf("could not up the migration: %w", err)
		}
	}

	return m, nil
}
