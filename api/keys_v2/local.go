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
	"github.com/t2bot/matrix-key-server/api/api_models"
	"github.com/t2bot/matrix-key-server/api/common"
	"github.com/t2bot/matrix-key-server/db"
	"github.com/t2bot/matrix-key-server/db/models"
	"github.com/t2bot/matrix-key-server/keys"
	"github.com/t2bot/matrix-key-server/signing"
	"github.com/t2bot/matrix-key-server/util"
)

func GetLocalKeys(r *http.Request, log *logrus.Entry) interface{} {
	ownKeys, err := db.GetAllOwnKeys()
	if err != nil {
		log.Error(err)
		return common.InternalServerError("Failed to get keys")
	}

	unsignedResp := &api_models.ServerKeyResultUnsigned{
		ServerName:    keys.SelfDomainName,
		ValidUntilTs:  util.NowMillis() + 86400000, // 24 hours
		VerifyKeys:    make(map[models.KeyID]api_models.VerifyKey),
		OldVerifyKeys: make(map[models.KeyID]api_models.OldVerifyKey),
	}

	for _, k := range ownKeys {
		if k.ExpiresTs > 0 {
			unsignedResp.OldVerifyKeys[k.ID] = api_models.OldVerifyKey{
				ExpiredTs: k.ExpiresTs,
				Key:       k.PublicKey,
			}
		} else {
			unsignedResp.VerifyKeys[k.ID] = api_models.VerifyKey{
				Key: k.PublicKey,
			}
		}
	}

	resp := &api_models.ServerKeyResult{
		ServerKeyResultUnsigned: unsignedResp,
		Signatures: api_models.Signatures{
			keys.SelfDomainName: map[string]string{},
		},
	}

	for _, key := range ownKeys {
		loaded, err := keys.LoadKey(key)
		if err != nil {
			log.Error(err)
			return common.InternalServerError("Failed to load private key")
		}

		signature, err := signing.GetSignatureOfObject(unsignedResp, loaded.Priv)
		if err != nil {
			log.Error(err)
			return common.InternalServerError("Failed to sign response")
		}

		resp.Signatures[keys.SelfDomainName][string(loaded.ID)] = signature
	}

	return resp
}
