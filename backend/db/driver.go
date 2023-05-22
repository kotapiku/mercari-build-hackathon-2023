package db

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

func PrepareDB(ctx context.Context) (*sql.DB, error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get current path: %w")
	}

	db, err := sql.Open("sqlite3", filepath.Join(path, "db", "mercari.sqlite3"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create DB: %w")
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to ping DB: %w")
	}

	f, err := os.ReadFile(filepath.Join(path, "sql", "01_schema.sql"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to open schema.sql %w")
	}

	if _, err = db.ExecContext(ctx, string(f)); err != nil {
		return nil, errors.Wrap(err, "failed to exec query: %w")
	}

	return db, nil
}
