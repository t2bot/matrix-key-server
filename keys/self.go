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
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/t2bot/matrix-key-server/db"
	"github.com/t2bot/matrix-key-server/db/models"
	"github.com/t2bot/matrix-key-server/signing"
	"github.com/t2bot/matrix-key-server/util"
	"golang.org/x/crypto/ed25519"
)

var SelfDomainName string

type SelfKey struct {
	RawKey *models.OwnKey
	ID     models.KeyID
	Pub    ed25519.PublicKey
	Priv   ed25519.PrivateKey
}

var ownKey *SelfKey

func GetSelfKey() (*SelfKey, error) {
	if ownKey == nil {
		activeKeyIds, err := db.GetOwnActiveKeyIds()
		if err != nil {
			return nil, err
		}

		if len(activeKeyIds) == 0 {
			logrus.Warn("No active key for this server: generating one now")

			pub, priv, err := ed25519.GenerateKey(rand.Reader)
			if err != nil {
				return nil, err
			}

			keyId, err := util.GenerateRandomString(8)
			if err != nil {
				return nil, err
			}
			keyId = fmt.Sprintf("ed25519:%s", keyId)

			pubEncoded := signing.EncodeUnpaddedBase64ToString(pub)
			privEncoded := signing.EncodeUnpaddedBase64ToString(priv)

			dbKey := &models.OwnKey{
				ID:         models.KeyID(keyId),
				PublicKey:  models.UnpaddedBase64EncodedData(pubEncoded),
				PrivateKey: models.UnpaddedBase64EncodedData(privEncoded),
				ExpiresTs:  0,
			}
			err = db.AddOwnActiveKey(dbKey.ID, dbKey.PublicKey, dbKey.PrivateKey)
			if err != nil {
				return nil, err
			}

			ownKey = &SelfKey{
				RawKey: dbKey,
				ID:     dbKey.ID,
				Pub:    pub,
				Priv:   priv,
			}
		} else {
			logrus.Infof("There are %d active keys for this server", len(activeKeyIds))

			k, err := db.GetOwnKey(activeKeyIds[0])
			if err != nil {
				return nil, err
			}
			if k == nil {
				return nil, errors.New("failed to fetch active key")
			}

			loadedKey, err := LoadKey(k)
			if err != nil {
				return nil, err
			}
			ownKey = loadedKey
		}
	}

	return ownKey, nil
}

func LoadKey(key *models.OwnKey) (*SelfKey, error) {
	pubDecoded, err := signing.DecodeUnpaddedBase64String(string(key.PublicKey))
	if err != nil {
		return nil, err
	}

	privDecoded, err := signing.DecodeUnpaddedBase64String(string(key.PrivateKey))
	if err != nil {
		return nil, err
	}

	return &SelfKey{
		RawKey: key,
		ID:     key.ID,
		Pub:    ed25519.PublicKey(pubDecoded),
		Priv:   ed25519.PrivateKey(privDecoded),
	}, nil
}
