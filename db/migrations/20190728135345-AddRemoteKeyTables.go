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

package migrations

import (
	"database/sql"
)

func Up20190728135345AddRemoteKeyTables(db *sql.DB) error {
	var err error

	_, err = db.Exec("CREATE TABLE remote_servers (server_name VARCHAR(255) NOT NULL PRIMARY KEY, updated_ts BIGINT NOT NULL, valid_until_ts BIGINT NOT NULL, nonstandard_json JSON NOT NULL);")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE remote_keys (server_name VARCHAR(255) NOT NULL REFERENCES remote_servers(server_name), key_id VARCHAR(255) NOT NULL, public_key_b64 VARCHAR(255) NOT NULL, expires_ts BIGINT NOT NULL DEFAULT 0, PRIMARY KEY (server_name, key_id));")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE remote_signatures (server_name VARCHAR(255) NOT NULL REFERENCES remote_servers(server_name), key_id VARCHAR(255) NOT NULL, signature_b64 VARCHAR(255) NOT NULL, PRIMARY KEY (server_name, key_id));")
	if err != nil {
		return err
	}

	return nil
}
