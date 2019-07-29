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
	"testing"

	"golang.org/x/crypto/ed25519"
)

var publicKeys = map[string]map[string]ed25519.PublicKey{
	domain: {
		string(keyId): makePublicKey(privateKey),
	},
}

func makePublicKey(priv ed25519.PrivateKey) ed25519.PublicKey {
	publicKey := make([]byte, ed25519.PublicKeySize)
	copy(publicKey, priv[32:])
	return ed25519.PublicKey(publicKey)
}

func TestVerifyObject_Simple(t *testing.T) {
	// We should be able to sign something and get the same verification back
	signed, _ := SignObject(map[string]interface{}{}, domain, keyId, privateKey)
	err := VerifySignatures(signed, publicKeys)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestVerifyObject_WithContent(t *testing.T) {
	// We should be able to sign something and get the same verification back
	signed, _ := SignObject(map[string]interface{}{
		"one": 1,
		"two": "Two",
	}, domain, keyId, privateKey)
	err := VerifySignatures(signed, publicKeys)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
