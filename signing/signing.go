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

package signing

import (
	"github.com/t2bot/matrix-key-server/db/models"
	"github.com/t2bot/matrix-key-server/util"
	"golang.org/x/crypto/ed25519"
)

func SignEvent(obj interface{}, domain string, keyId models.KeyID, key ed25519.PrivateKey) (map[string]interface{}, error) {
	redacted, unsigned, err := RedactObject(obj, RedactionV1)
	if err != nil {
		return nil, err
	}

	contentHash, err := CalculateSha256ContentHash(obj)
	if err != nil {
		return nil, err
	}
	redacted["hashes"] = map[string]interface{}{
		"sha256": contentHash,
	}

	// This is taken off again right away, but it is faster to pass it down like this
	if unsigned != nil {
		redacted["unsigned"] = unsigned
	}

	return SignObject(redacted, domain, keyId, key)
}

func SignObject(obj map[string]interface{}, domain string, keyId models.KeyID, key ed25519.PrivateKey) (map[string]interface{}, error) {
	var unsigned map[string]interface{}
	if u, ok := obj["unsigned"]; !ok {
		unsigned = nil
	} else {
		unsigned = u.(map[string]interface{})
	}
	delete(obj, "unsigned")

	var signatures map[string]interface{}
	if s, ok := obj["signatures"]; !ok {
		signatures = nil
	} else {
		signatures = s.(map[string]interface{})
	}
	delete(obj, "signatures")

	canonical, err := EncodeCanonicalJson(obj)
	if err != nil {
		return nil, err
	}

	signature := ed25519.Sign(key, canonical)

	m, err := util.InterfaceToMap(obj)
	if err != nil {
		return nil, err
	}

	if signatures != nil {
		m["signatures"] = signatures
	}

	var signaturesObj map[string]interface{}
	if sigs, ok := m["signatures"]; !ok {
		signaturesObj = make(map[string]interface{})
		m["signatures"] = signaturesObj
	} else {
		signaturesObj = sigs.(map[string]interface{})
	}

	var domainObj map[string]interface{}
	if domains, ok := signaturesObj[domain]; !ok {
		domainObj = make(map[string]interface{})
		signaturesObj[domain] = domainObj
	} else {
		domainObj = domains.(map[string]interface{})
	}

	domainObj[string(keyId)] = EncodeUnpaddedBase64ToString(signature)
	if unsigned != nil {
		m["unsigned"] = unsigned
	}

	return m, nil
}

func GetSignatureOfObject(obj interface{}, key ed25519.PrivateKey) (string, error) {
	m, err := util.InterfaceToMap(obj)
	if err != nil {
		return "", err
	}

	s, err := SignObject(m, "a", "a", key)
	if err != nil {
		return "", nil
	}

	sigs := s["signatures"].(map[string]interface{})
	domains := sigs["a"].(map[string]interface{})
	return domains["a"].(string), nil
}
