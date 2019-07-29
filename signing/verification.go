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
	"errors"
	"fmt"

	"github.com/t2bot/matrix-key-server/util"
	"golang.org/x/crypto/ed25519"
)

func VerifySignatures(obj interface{}, publicKeys map[string]map[string]ed25519.PublicKey) error {
	m, err := util.InterfaceToMap(obj)
	if err != nil {
		return err
	}

	delete(m, "unsigned")

	var signatures map[string]interface{}
	if s, ok := m["signatures"]; !ok {
		signatures = nil
	} else {
		signatures = s.(map[string]interface{})
	}
	delete(m, "signatures")

	if len(signatures) == 0 {
		return errors.New("no signatures found")
	}

	canonical, err := EncodeCanonicalJson(m)
	if err != nil {
		return err
	}

	// Verify all the signatures now
	for domain, sig := range signatures {
		var domainKeys map[string]ed25519.PublicKey
		var ok bool
		if domainKeys, ok = publicKeys[domain]; !ok {
			return errors.New("missing public keys for " + domain)
		}

		keySigs := sig.(map[string]interface{})
		for keyId, b64 := range keySigs {
			var publicKey ed25519.PublicKey
			if publicKey, ok = domainKeys[keyId]; !ok {
				return errors.New(fmt.Sprintf("missing public key for %s %s", domain, keyId))
			}

			signature, err := DecodeUnpaddedBase64String(b64.(string))
			if err != nil {
				return err
			}

			valid := ed25519.Verify(publicKey, canonical, signature)
			if !valid {
				return errors.New(fmt.Sprintf("signature verification failed for %s %s", domain, keyId))
			}
		}
	}

	return nil
}
