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

	"github.com/t2bot/matrix-key-server/util"
)

type RedactionAlgorithm string

const RedactionV1 = RedactionAlgorithm("v1")

func RedactObject(obj interface{}, algorithm RedactionAlgorithm) (map[string]interface{}, map[string]interface{}, error) {
	if algorithm == RedactionV1 {
		return redactObjectV1(obj)
	} else {
		return nil, nil, errors.New("unknown redaction algorithm")
	}
}

func redactObjectV1(obj interface{}) (map[string]interface{}, map[string]interface{}, error) {
	keepKeys := map[string]bool{
		"event_id":         true,
		"type":             true,
		"room_id":          true,
		"sender":           true,
		"state_key":        true,
		"content":          true,
		"hashes":           true,
		"signatures":       true,
		"depth":            true,
		"prev_events":      true,
		"prev_state":       true,
		"auth_events":      true,
		"origin":           true,
		"origin_server_ts": true,
		"membership":       true,
	}
	keepContentKeysIfType := map[string][]string{
		"m.room.member":             {"membership"},
		"m.room.create":             {"creator"},
		"m.room.join_rules":         {"join_rule"},
		"m.room.power_levels":       {"ban", "events", "events_default", "kick", "redact", "state_default", "users", "users_default"},
		"m.room.aliases":            {"aliases"},
		"m.room.history_visibility": {"history_visibility"},
	}

	m, err := util.InterfaceToMap(obj)
	if err != nil {
		return nil, nil, err
	}

	var unsigned map[string]interface{}
	if u, ok := m["unsigned"]; !ok {
		unsigned = nil
	} else {
		unsigned = u.(map[string]interface{})
	}

	p, err := pruneObject(m, keepKeys, keepContentKeysIfType)
	return p, unsigned, err
}

func pruneObject(obj map[string]interface{}, keepKeys map[string]bool, keepContentKeysIfType map[string][]string) (map[string]interface{}, error) {
	var m = make(map[string]interface{})

	for k, v := range obj {
		if val, ok := keepKeys[k]; ok && val {
			m[k] = v
		}
	}

	if t, ok := m["type"]; ok {
		eventType := t.(string)
		if c, ok := m["content"]; ok {
			eventContent := c.(map[string]interface{})
			newContent := make(map[string]interface{})
			if keysToKeep, ok := keepContentKeysIfType[eventType]; ok {
				for _, k := range keysToKeep {
					if v, ok := eventContent[k]; ok {
						newContent[k] = v
					}
				}
			}
			m["content"] = newContent
		}
	} else if _, ok := m["content"]; ok {
		m["content"] = make(map[string]interface{})
	}

	return m, nil
}
