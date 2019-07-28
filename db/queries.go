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

var queries = map[string]string{
	selectAllSelfKeys:      "SELECT key_id, public_key_b64, private_key_b64, expires_ts FROM self_keys;",
	selectActiveSelfKeyIds: "SELECT key_id FROM self_keys WHERE expires_ts = 0;",
	selectSelfKey:          "SELECT public_key_b64, private_key_b64, expires_ts FROM self_keys WHERE key_id = $1;",
	insertActiveSelfKey:    "INSERT INTO self_keys (key_id, public_key_b64, private_key_b64) VALUES ($1, $2, $3);",
}
