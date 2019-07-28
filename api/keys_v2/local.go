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

package keys_v2

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/t2bot/matrix-key-server/api/common"
	"github.com/t2bot/matrix-key-server/db"
	"github.com/t2bot/matrix-key-server/db/models"
	"github.com/t2bot/matrix-key-server/keys"
	"github.com/t2bot/matrix-key-server/util"
)

type VerifyKey struct {
	Key models.Base64EncodedKeyData `json:"key"`
}

type OldVerifyKey struct {
	Key       models.Base64EncodedKeyData `json:"key"`
	ExpiredTs models.Timestamp            `json:"expired_ts"`
}

type ServerKeyResult struct {
	ServerName    string                        `json:"server_name"`
	ValidUntilTs  int64                         `json:"valid_until_ts"`
	VerifyKeys    map[models.KeyID]VerifyKey    `json:"verify_keys"`
	OldVerifyKeys map[models.KeyID]OldVerifyKey `json:"old_verify_keys"`
}

func GetLocalKeys(r *http.Request, log *logrus.Entry) interface{} {
	ownKeys, err := db.GetAllOwnKeys()
	if err != nil {
		logrus.Error(err)
		return common.InternalServerError("Failed to get keys")
	}

	resp := &ServerKeyResult{
		ServerName:    keys.SelfDomainName,
		ValidUntilTs:  util.NowMillis() + 86400000, // 24 hours
		VerifyKeys:    make(map[models.KeyID]VerifyKey),
		OldVerifyKeys: make(map[models.KeyID]OldVerifyKey),
	}

	for _, k := range ownKeys {
		if k.ExpiresTs > 0 {
			resp.OldVerifyKeys[k.ID] = OldVerifyKey{
				ExpiredTs: k.ExpiresTs,
				Key:       k.PublicKey,
			}
		} else {
			resp.VerifyKeys[k.ID] = VerifyKey{
				Key: k.PublicKey,
			}
		}
	}

	// TODO: Sign object

	return resp
}
