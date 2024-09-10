package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func RunMigrations(dbURL string) error {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}
