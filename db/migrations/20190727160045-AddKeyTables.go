package migrations

import (
	"database/sql"
)

func Up20190727160045AddKeyTables(db *sql.DB) error {
	var err error

	_, err = db.Exec("CREATE TABLE self_keys (key_id VARCHAR(255) NOT NULL PRIMARY KEY, public_key_b64 VARCHAR(255) NOT NULL, private_key_b64 VARCHAR(255) NOT NULL, expires_ts BIGINT NOT NULL DEFAULT 0);")
	if err != nil {
		return err
	}

	return nil
}
