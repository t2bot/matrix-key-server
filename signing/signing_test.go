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
	"encoding/json"
	"testing"

	"github.com/t2bot/matrix-key-server/db/models"
	"golang.org/x/crypto/ed25519"
)

var seed, _ = DecodeUnpaddedBase64String("YJDBA9Xnr2sVqXD9Vj7XVUnmFZcZrlw8Md7kMW+3XA1")
var privateKey = ed25519.NewKeyFromSeed(seed)

const keyId = models.KeyID("ed25519:1")
const domain = "domain"

func TestSignObject_Simple(t *testing.T) {
	signed, _ := SignObject(map[string]interface{}{}, domain, keyId, privateKey)
	verifySignature(signed, "K8280/U9SSy9IVtjBuVeLr+HpOB4BQFWbg+UZaADMtTdGYI7Geitb76LTrr5QV/7Xg4ahLwYGYZzuHGZKM5ZAQ", domain, keyId, t)
}

func TestSignObject_WithContent(t *testing.T) {
	signed, _ := SignObject(map[string]interface{}{
		"one": 1,
		"two": "Two",
	}, domain, keyId, privateKey)
	verifySignature(signed, "KqmLSbO39/Bzb0QIYE82zqLwsA+PDzYIpIRA2sRQ4sL53+sN6/fpNSoqE7BP7vBZhG6kYdD13EIMJpvhJI+6Bw", domain, keyId, t)
}

func TestSignEvent_MinimalEvent(t *testing.T) {
	signed, _ := SignEvent(map[string]interface{}{
		"room_id":          "!x:domain",
		"sender":           "@a:domain",
		"origin":           "domain",
		"origin_server_ts": 1000000,
		"signatures":       map[string]interface{}{},
		"hashes":           map[string]interface{}{},
		"type":             "X",
		"content":          map[string]interface{}{},
		"prev_events":      []interface{}{},
		"auth_events":      []interface{}{},
		"depth":            3,
		"unsigned": map[string]interface{}{
			"age_ts": 1000000,
		},
	}, domain, keyId, privateKey)
	verifySignature(signed, "KxwGjPSDEtvnFgU00fwFz+l6d2pJM6XBIaMEn81SXPTRl16AqLAYqfIReFGZlHi5KLjAWbOoMszkwsQma+lYAg", domain, keyId, t)
}

func TestSignEvent_RedactableEvent(t *testing.T) {
	signed, _ := SignEvent(map[string]interface{}{
		"content": map[string]interface{}{
			"body": "Here is the message content",
		},
		"event_id":         "$0:domain",
		"origin":           "domain",
		"origin_server_ts": 1000000,
		"type":             "m.room.message",
		"room_id":          "!r:domain",
		"sender":           "@u:domain",
		"signatures":       map[string]interface{}{},
		"unsigned": map[string]interface{}{
			"age_ts": 1000000,
		},
	}, domain, keyId, privateKey)
	verifySignature(signed, "Wm+VzmOUOz08Ds+0NTWb1d4CZrVsJSikkeRxh6aCcUwu6pNC78FunoD7KNWzqFn241eYHYMGCA5McEiVPdhzBA", domain, keyId, t)
}

func verifySignature(signed map[string]interface{}, signature string, domain string, keyId models.KeyID, t *testing.T) {
	if sigs, ok := signed["signatures"]; ok {
		v := sigs.(map[string]interface{})
		if keySig, ok := v[domain]; ok {
			v = keySig.(map[string]interface{})
			if actual, ok := v[string(keyId)]; ok {
				actualSig := actual.(string)
				if actualSig != signature {
					t.Error("Mismatch signatures")
					t.Errorf("Got: %s", actualSig)
					t.Errorf("Expected: %s", signature)

					j, _ := json.Marshal(signed)
					t.Errorf("Got JSON: %s", string(j))

					t.Fail()
				}
			} else {
				t.Error("Missing key ID")
				t.Fail()
			}
		} else {
			t.Error("Missing domain signature")
			t.Fail()
		}
	} else {
		t.Error("Missing signatures")
		t.Fail()
	}
}
