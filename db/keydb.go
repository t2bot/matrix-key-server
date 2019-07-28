/*
 * Copyright 2019 Travis Ralston <travis@t2bot.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
