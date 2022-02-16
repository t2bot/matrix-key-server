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

package custom

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/t2bot/matrix-key-server/api/common"
	"github.com/t2bot/matrix-key-server/db/models"
	"github.com/t2bot/matrix-key-server/keys"
	"github.com/t2bot/matrix-key-server/signing"
	"golang.org/x/crypto/ed25519"
)

func VerifyAuthHeader(r *http.Request, log *logrus.Entry) interface{} {
	method := r.Header.Get("X-Keys-Method")
	uri := r.Header.Get("X-Keys-URI")
	destination := r.Header.Get("X-Keys-Destination")

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		return common.InternalServerError("Failed to read body")
	}

	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "X-Matrix ") {
		return common.UnauthorizedError()
	}

	// Parse the header
	authParams := map[string]string{}
	paramCsv := strings.Split(strings.Split(auth, " ")[1], ",")
	for _, p := range paramCsv {
		vals := strings.Split(p, "=")
		if len(vals) != 2 {
			continue
		}

		k := vals[0]
		v := vals[1]
		if v[0] == '"' {
			v = v[1 : len(v)-1]
		}

		authParams[k] = v
	}

	origin := authParams["origin"]
	keyId := authParams["key"]
	signature := authParams["sig"]

	obj := map[string]interface{}{
		"method":      method,
		"uri":         uri,
		"origin":      origin,
		"destination": destination,
		"signatures": map[string]interface{}{
			origin: map[string]interface{}{
				keyId: signature,
			},
		},
	}
	if len(b) > 0 {
		parsed := make(map[string]interface{})
		err = json.Unmarshal(b, &parsed)
		if err != nil {
			log.Warn("Error parsing JSON body", err)
		} else {
			obj["content"] = parsed
		}
	}

	validKeys, err := keys.QueryRemoteKeys(models.ServerName(origin), 0)
	if err != nil {
		log.Error(err)
		return common.InternalServerError("Failed to get remote server keys")
	}

	publicKeys := map[string]map[string]ed25519.PublicKey{
		string(validKeys.ServerName): {},
	}
	for _, k := range validKeys.Keys {
		b, err := signing.DecodeUnpaddedBase64String(string(k.PublicKey))
		if err != nil {
			log.Error(err)
			return common.InternalServerError("Failed to parse remote server keys")
		}

		publicKeys[string(k.ServerName)][string(k.ID)] = ed25519.PublicKey(b)
	}

	err = signing.VerifySignatures(obj, publicKeys)
	if err != nil {
		log.Error(err)
		return common.UnauthorizedError()
	}

	return common.EmptyResponse{}
}
