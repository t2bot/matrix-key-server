package db

import (
	"database/sql"

	_ "github.com/lib/pq" // postgres driver
	"github.com/t2bot/matrix-key-server/db/migrations"
)

type Database struct {
	db *sql.DB
}

var dbInstance *Database

func Setup(dbUrl string) error {
	dbInstance := &Database{}
	var err error

	if dbInstance.db, err = sql.Open("postgres", dbUrl); err != nil {
		return err
	}

	dbInstance.db.SetMaxOpenConns(15)
	dbInstance.db.SetMaxIdleConns(5)

	fnCalls := make([]func() error, 0)
	fnCalls = append(fnCalls, func() error { return prepareMigrations(dbInstance.db) })
	fnCalls = append(fnCalls, func() error { return applyMigration(dbInstance.db, migrations.Up20190727160045AddKeyTables) })
	fnCalls = append(fnCalls, func() error { return prepareStatements(dbInstance.db) })

	for _, fn := range fnCalls {
		err = fn()
		if err != nil {
			return err
		}
	}

	return nil
}

func prepareStatements(db *sql.DB) error {
	for k, v := range queries {
		s, err := db.Prepare(v)
		if err != nil {
			return err
		}
		statements[k] = s
	}
	return nil
}
