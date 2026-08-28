package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/lib/pq"
)

func Open(dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("database DSN is required")
	}
	return sql.Open("postgres", dsn)
}

func Migration() ([]byte, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("cannot resolve migration path")
	}
	path := filepath.Join(filepath.Dir(file), "migrations/001_init.sql")
	if _, err := os.Stat(path); err != nil {
		path = "/migrations/001_init.sql"
	}
	return os.ReadFile(path)
}

func ApplyMigrations(db *sql.DB) error {
	migration, err := Migration()
	if err != nil {
		return err
	}
	_, err = db.Exec(string(migration))
	return err
}
