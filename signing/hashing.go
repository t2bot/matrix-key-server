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
	"crypto/sha256"

	"github.com/t2bot/matrix-key-server/util"
)

func CalculateSha256ContentHash(ev interface{}) (string, error) {
	m, err := util.InterfaceToMap(ev)
	if err != nil {
		return "", err
	}

	delete(m, "unsigned")
	delete(m, "signatures")
	delete(m, "hashes")

	b, err := EncodeCanonicalJson(m)
	if err != nil {
		return "", err
	}

	h := sha256.Sum256(b)
	var hbytes = make([]byte, 0)
	for _, v := range h {
		hbytes = append(hbytes, v)
	}
	return EncodeUnpaddedBase64ToString(hbytes), nil
}
