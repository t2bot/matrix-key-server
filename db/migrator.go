package db

import (
	"database/sql"
	"reflect"
	"runtime"

	"github.com/sirupsen/logrus"
)

func prepareMigrations(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS migrations(func_name TEXT NOT NULL PRIMARY KEY);")
	if err != nil {
		return err
	}

	return nil
}

func applyMigration(db *sql.DB, up func(db *sql.DB) error) error {
	fnName := runtime.FuncForPC(reflect.ValueOf(up).Pointer()).Name()

	logrus.Info("Applying migration: ", fnName)

	r := db.QueryRow("SELECT func_name FROM migrations WHERE func_name=$1 LIMIT 1", fnName)
	var dbFnName string
	err := r.Scan(&dbFnName)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == nil {
		// Already executed
		return nil
	}

	logrus.Info("Running migration: ", fnName)
	err = up(db)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO migrations (func_name) VALUES ($1)", fnName)
	if err != nil {
		return err
	}

	return nil
}
