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

package keys

import (
	"encoding/json"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/t2bot/matrix-key-server/api/api_models"
	"github.com/t2bot/matrix-key-server/db"
	"github.com/t2bot/matrix-key-server/db/models"
	"github.com/t2bot/matrix-key-server/federation"
	"github.com/t2bot/matrix-key-server/signing"
	"github.com/t2bot/matrix-key-server/util"
	"golang.org/x/crypto/ed25519"
)

func QueryRemoteKeys(serverName models.ServerName, minValidUntilTs models.Timestamp) (*models.CachedRemoteKeys, error) {
	s, err := db.GetRemoteServerMetadata(serverName)
	if err != nil {
		return nil, err
	}
	if s != nil {
		now := util.NowMillis()
		sevenDaysFromUpdate := int64(s.UpdatedTs) + int64(604800000)

		isBeyondLifespan := now > sevenDaysFromUpdate
		isHalfwayDead := int64(s.UpdatedTs+s.ValidUntilTs)/2 < now
		isMinimallyAccepted := s.ValidUntilTs >= minValidUntilTs

		if isMinimallyAccepted && !isHalfwayDead && !isBeyondLifespan {
			return packageCachedKeysFor(s)
		}
	}

	// Cache miss: fetch new keys
	// TODO: Rate limit: https://github.com/turt2live/matrix-key-server/issues/2
	url, hostname, err := federation.GetServerApiUrl(string(serverName))
	if err != nil {
		logrus.Error(err)

		if s != nil {
			// Continue to serve the last known response from the dead server
			return packageCachedKeysFor(s)
		} else {
			// else if the server is dead and we have no keys then return nothing back
			return &models.CachedRemoteKeys{
				Keys:       make([]*models.RemoteKey, 0),
				Signatures: make([]*models.RemoteSignature, 0),
				RemoteServer: &models.RemoteServer{
					ServerName:   serverName,
					ValidUntilTs: models.Timestamp(util.NowMillis()),
					UpdatedTs:    models.Timestamp(util.NowMillis()),
				},
			}, nil
		}
	}

	keysUrl := url + "/_matrix/key/v2/server"
	keysResponse, err := federation.FederatedGet(keysUrl, hostname)
	if err != nil {
		return nil, err
	}

	c, err := ioutil.ReadAll(keysResponse.Body)
	if err != nil {
		return nil, err
	}

	keyInfo := api_models.ServerKeyResult{}
	err = json.Unmarshal(c, &keyInfo)
	if err != nil {
		return nil, err
	}

	publicKeys, err := grabPublicKeys(keyInfo)
	if err != nil {
		return nil, err
	}

	additionalFields := models.AdditionalJSON{}
	fullyUnmarshalled := make(map[string]interface{})
	err = json.Unmarshal(c, &fullyUnmarshalled)
	if err != nil {
		return nil, err
	}
	m, err := util.InterfaceToMap(keyInfo)
	if err != nil {
		return nil, err
	}
	for k, v := range fullyUnmarshalled {
		if _, ok := m[k]; !ok {
			additionalFields[k] = v
			m[k] = v
		}
	}

	err = signing.VerifySignatures(m, publicKeys)
	if err != nil {
		return nil, err
	}

	return storeRemoteKeys(keyInfo, additionalFields)
}

func packageCachedKeysFor(server *models.RemoteServer) (*models.CachedRemoteKeys, error) {
	keys, err := db.GetAllRemoteServerKeys(server.ServerName)
	if err != nil {
		return nil, err
	}

	sigs, err := db.GetAllRemoteServerSignatures(server.ServerName)
	if err != nil {
		return nil, err
	}

	return &models.CachedRemoteKeys{
		RemoteServer: server,
		Keys:         keys,
		Signatures:   sigs,
	}, nil
}

func storeRemoteKeys(keyInfo api_models.ServerKeyResult, additionalJson models.AdditionalJSON) (*models.CachedRemoteKeys, error) {
	res := &models.CachedRemoteKeys{
		RemoteServer: &models.RemoteServer{
			ServerName:      models.ServerName(keyInfo.ServerName),
			UpdatedTs:       models.Timestamp(util.NowMillis()),
			ValidUntilTs:    models.Timestamp(keyInfo.ValidUntilTs),
			NonStandardJSON: additionalJson,
		},
		Keys:       make([]*models.RemoteKey, 0),
		Signatures: make([]*models.RemoteSignature, 0),
	}

	err := db.UpsertRemoteServer(res.ServerName, res.UpdatedTs, res.ValidUntilTs, additionalJson)
	if err != nil {
		return nil, err
	}

	err = db.DeleteRemoteServerKeys(res.ServerName)
	if err != nil {
		return nil, err
	}

	err = db.DeleteRemoteServerSignatures(res.ServerName)
	if err != nil {
		return nil, err
	}

	for keyId, key := range keyInfo.VerifyKeys {
		cachedKey := &models.RemoteKey{
			ServerName: res.ServerName,
			ID:         models.KeyID(keyId),
			PublicKey:  models.UnpaddedBase64EncodedData(key.Key),
			ExpiresTs:  models.Timestamp(0),
		}
		err = db.AddRemoteServerKey(cachedKey.ServerName, cachedKey.ID, cachedKey.PublicKey, cachedKey.ExpiresTs)
		if err != nil {
			return nil, err
		}
		res.Keys = append(res.Keys, cachedKey)
	}

	for keyId, key := range keyInfo.OldVerifyKeys {
		cachedKey := &models.RemoteKey{
			ServerName: res.ServerName,
			ID:         models.KeyID(keyId),
			PublicKey:  models.UnpaddedBase64EncodedData(key.Key),
			ExpiresTs:  models.Timestamp(key.ExpiredTs),
		}
		err = db.AddRemoteServerKey(cachedKey.ServerName, cachedKey.ID, cachedKey.PublicKey, cachedKey.ExpiresTs)
		if err != nil {
			return nil, err
		}
		res.Keys = append(res.Keys, cachedKey)
	}

	for _, sig := range keyInfo.Signatures {
		for keyId, signature := range sig {
			cachedSignature := &models.RemoteSignature{
				ServerName: res.ServerName,
				KeyID:      models.KeyID(keyId),
				Signature:  models.UnpaddedBase64EncodedData(signature),
			}
			err = db.AddRemoteServerSignature(cachedSignature.ServerName, cachedSignature.KeyID, cachedSignature.Signature)
			if err != nil {
				return nil, err
			}
			res.Signatures = append(res.Signatures, cachedSignature)
		}
	}

	return res, nil
}

func grabPublicKeys(keyInfo api_models.ServerKeyResult) (map[string]map[string]ed25519.PublicKey, error) {
	keys := make(map[string]map[string]ed25519.PublicKey)
	keys[keyInfo.ServerName] = make(map[string]ed25519.PublicKey)
	for keyId, encodedKey := range keyInfo.VerifyKeys {
		b, err := signing.DecodeUnpaddedBase64String(string(encodedKey.Key))
		if err != nil {
			return nil, err
		}

		keys[keyInfo.ServerName][string(keyId)] = ed25519.PublicKey(b)
	}

	return keys, nil
}
