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

	"github.com/t2bot/matrix-key-server/db/models"
)

func GetAllOwnKeys() ([]*models.OwnKey, error) {
	r, err := statements[selectAllSelfKeys].Query()
	if err == sql.ErrNoRows {
		return make([]*models.OwnKey, 0), nil
	}
	if err != nil {
		return nil, err
	}

	var results []*models.OwnKey
	for r.Next() {
		v := &models.OwnKey{}
		err = r.Scan(&v.ID, &v.PublicKey, &v.PrivateKey, &v.ExpiresTs)
		if err != nil {
			return nil, err
		}
		results = append(results, v)
	}

	return results, nil
}

func GetOwnActiveKeyIds() ([]models.KeyID, error) {
	r, err := statements[selectActiveSelfKeyIds].Query()
	if err == sql.ErrNoRows {
		return make([]models.KeyID, 0), nil
	}
	if err != nil {
		return nil, err
	}

	var results []models.KeyID
	for r.Next() {
		var v models.KeyID
		err = r.Scan(&v)
		if err != nil {
			return nil, err
		}
		results = append(results, v)
	}

	return results, nil
}

func GetOwnKey(id models.KeyID) (*models.OwnKey, error) {
	r := statements[selectSelfKey].QueryRow(id)

	var key = &models.OwnKey{ID: id}

	err := r.Scan(&key.PublicKey, &key.PrivateKey, &key.ExpiresTs)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return key, nil
}

func AddOwnActiveKey(id models.KeyID, publicKey models.Base64EncodedKeyData, privateKey models.Base64EncodedKeyData) error {
	_, err := statements[insertActiveSelfKey].Exec(id, publicKey, privateKey)
	if err != nil {
		return err
	}
	return nil
}
