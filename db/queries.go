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
