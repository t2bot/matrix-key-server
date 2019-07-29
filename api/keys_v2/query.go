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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/t2bot/matrix-key-server/api/api_models"
	"github.com/t2bot/matrix-key-server/api/common"
	"github.com/t2bot/matrix-key-server/db"
	"github.com/t2bot/matrix-key-server/db/models"
	"github.com/t2bot/matrix-key-server/keys"
	"github.com/t2bot/matrix-key-server/signing"
	"github.com/t2bot/matrix-key-server/util"
	"golang.org/x/crypto/ed25519"
)

type BatchedServerKeys struct {
	Keys []map[string]interface{} `json:"server_keys"`
}

type LookupCriteria struct {
	MinValidTs int64 `json:"minimum_valid_until_ts"`
}

type BatchKeyLookup struct {
	Keys map[string]map[string]LookupCriteria `json:"server_keys"`
}

func QueryKeysSingle(r *http.Request, log *logrus.Entry) interface{} {
	var err error

	params := mux.Vars(r)

	serverName := params["serverName"]
	minValidTsRaw := r.URL.Query().Get("minimum_valid_until_ts")

	minValidTs := util.NowMillis()
	if minValidTsRaw != "" {
		minValidTs, err = strconv.ParseInt(minValidTsRaw, 10, 64)
		if err != nil {
			log.Error(err)
			return common.BadRequest("invalid minimum timestamp")
		}
	}

	expanded, errLike := findAndPrepareKeys(serverName, minValidTs, log)
	if errLike != nil {
		return errLike
	}

	return &BatchedServerKeys{Keys: []map[string]interface{}{expanded}}
}

func QueryKeysBatch(r *http.Request, log *logrus.Entry) interface{} {
	lookup := &BatchKeyLookup{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		return common.InternalServerError("Failed to read body")
	}

	err = json.Unmarshal(b, &lookup)
	if err != nil {
		log.Error(err)
		return common.BadRequest("Body not JSON")
	}

	finalResp := &BatchedServerKeys{Keys: make([]map[string]interface{}, 0)}

	for domain, keySearches := range lookup.Keys {
		maxMinValidTs := int64(0)
		for _, q := range keySearches {
			if q.MinValidTs > maxMinValidTs {
				maxMinValidTs = q.MinValidTs
			}
		}

		if maxMinValidTs <= 0 {
			maxMinValidTs = util.NowMillis()
		}

		expanded, errLike := findAndPrepareKeys(domain, maxMinValidTs, log)
		if errLike != nil {
			return errLike
		}

		finalResp.Keys = append(finalResp.Keys, expanded)
	}

	return finalResp
}

func findAndPrepareKeys(serverName string, minValidTs int64, log *logrus.Entry) (map[string]interface{}, interface{}) {
	remoteKeys, err := keys.QueryRemoteKeys(models.ServerName(serverName), models.Timestamp(minValidTs))
	if err != nil {
		log.Error(err)
		return nil, common.InternalServerError("Fatal error retrieving keys")
	}
	if len(remoteKeys.Keys) == 0 {
		log.Warn("Did not get any keys from remote server")
		return map[string]interface{}{}, nil
	}

	if string(remoteKeys.ServerName) != serverName {
		log.Error("Got response from unexpected server")
		return nil, common.InternalServerError("Unexpected server_name")
	}

	publicKeys := map[string]map[string]ed25519.PublicKey{
		keys.SelfDomainName:           make(map[string]ed25519.PublicKey),
		string(remoteKeys.ServerName): make(map[string]ed25519.PublicKey),
	}

	unsignedResp := &api_models.ServerKeyResultUnsigned{
		ServerName:    string(remoteKeys.ServerName),
		ValidUntilTs:  int64(remoteKeys.ValidUntilTs),
		VerifyKeys:    make(map[models.KeyID]api_models.VerifyKey),
		OldVerifyKeys: make(map[models.KeyID]api_models.OldVerifyKey),
	}

	for _, k := range remoteKeys.Keys {
		if k.ExpiresTs > 0 {
			unsignedResp.OldVerifyKeys[k.ID] = api_models.OldVerifyKey{
				ExpiredTs: k.ExpiresTs,
				Key:       k.PublicKey,
			}
		} else {
			unsignedResp.VerifyKeys[k.ID] = api_models.VerifyKey{
				Key: k.PublicKey,
			}
			b, err := signing.DecodeUnpaddedBase64String(string(k.PublicKey))
			if err != nil {
				log.Error(err)
				return nil, common.InternalServerError("Failed to read remote public keys")
			}
			publicKeys[string(k.ServerName)][string(k.ID)] = ed25519.PublicKey(b)
		}
	}

	resp := &api_models.ServerKeyResult{
		ServerKeyResultUnsigned: unsignedResp,
		Signatures: api_models.Signatures{
			keys.SelfDomainName:           make(map[string]string),
			string(remoteKeys.ServerName): make(map[string]string),
		},
	}

	expanded, err := util.InterfaceToMap(resp)
	if err != nil {
		log.Error(err)
		return nil, common.InternalServerError("Failed to convert response")
	}
	if remoteKeys.NonStandardJSON != nil {
		for k, v := range remoteKeys.NonStandardJSON {
			expanded[k] = v
		}
	}

	ownKeys, err := db.GetAllOwnKeys()
	if err != nil {
		log.Error(err)
		return nil, common.InternalServerError("Failed to get own keys")
	}

	for _, key := range ownKeys {
		loaded, err := keys.LoadKey(key)
		if err != nil {
			log.Error(err)
			return nil, common.InternalServerError("Failed to load private key")
		}

		signature, err := signing.GetSignatureOfObject(expanded, loaded.Priv)
		if err != nil {
			log.Error(err)
			return nil, common.InternalServerError("Failed to sign response")
		}

		resp.Signatures[keys.SelfDomainName][string(loaded.ID)] = signature
		publicKeys[keys.SelfDomainName][string(loaded.ID)] = loaded.Pub
	}

	// Append the signatures for the remote server
	for _, sig := range remoteKeys.Signatures {
		resp.Signatures[string(sig.ServerName)][string(sig.KeyID)] = string(sig.Signature)
	}

	expanded["signatures"] = resp.Signatures

	// Do last minute verifications on signatures
	err = signing.VerifySignatures(expanded, publicKeys)
	if err != nil {
		log.Error(err)
		return nil, common.InternalServerError("Failed last-minute signature verifications")
	}

	return expanded, nil
}
