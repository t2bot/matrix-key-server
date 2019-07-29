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
)

var statements = map[string]*sql.Stmt{}

const selectAllSelfKeys = "selectAllSelfKeys"
const selectActiveSelfKeyIds = "selectActiveSelfKeyIds"
const selectSelfKey = "selectSelfKey"
const insertActiveSelfKey = "insertActiveSelfKey"
const selectRemoteServer = "selectRemoteServer"
const selectRemoteKeys = "selectRemoteKeys"
const selectRemoteSignatures = "selectRemoteSignatures"
const deleteRemoteKeys = "deleteRemoteKeys"
const deleteRemoteSignatures = "deleteRemoteSignatures"
const upsertRemoteServer = "upsertRemoteServer"
const insertRemoteKey = "insertRemoteKey"
const insertRemoteSignature = "insertRemoteSignature"

var queries = map[string]string{
	selectAllSelfKeys:      "SELECT key_id, public_key_b64, private_key_b64, expires_ts FROM self_keys;",
	selectActiveSelfKeyIds: "SELECT key_id FROM self_keys WHERE expires_ts = 0;",
	selectSelfKey:          "SELECT public_key_b64, private_key_b64, expires_ts FROM self_keys WHERE key_id = $1;",
	insertActiveSelfKey:    "INSERT INTO self_keys (key_id, public_key_b64, private_key_b64) VALUES ($1, $2, $3);",
	selectRemoteServer:     "SELECT updated_ts, valid_until_ts, nonstandard_json FROM remote_servers WHERE server_name = $1",
	selectRemoteKeys:       "SELECT key_id, public_key_b64, expires_ts FROM remote_keys WHERE server_name = $1",
	selectRemoteSignatures: "SELECT key_id, signature_b64 FROM remote_signatures WHERE server_name = $1",
	deleteRemoteKeys:       "DELETE FROM remote_keys WHERE server_name = $1;",
	deleteRemoteSignatures: "DELETE FROM remote_signatures WHERE server_name = $1;",
	upsertRemoteServer:     "INSERT INTO remote_servers (server_name, updated_ts, valid_until_ts, nonstandard_json) VALUES ($1, $2, $3, $4) ON CONFLICT (server_name) DO UPDATE SET updated_ts = $2, valid_until_ts = $3, nonstandard_json = $4;",
	insertRemoteKey:        "INSERT INTO remote_keys (server_name, key_id, public_key_b64, expires_ts) VALUES ($1, $2, $3, $4);",
	insertRemoteSignature:  "INSERT INTO remote_signatures (server_name, key_id, signature_b64) VALUES ($1, $2, $3);",
}
