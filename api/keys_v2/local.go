package keys_v2

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/t2bot/matrix-key-server/api/common"
	"github.com/t2bot/matrix-key-server/db"
	"github.com/t2bot/matrix-key-server/db/models"
	"github.com/t2bot/matrix-key-server/keys"
	"github.com/t2bot/matrix-key-server/util"
)

type VerifyKey struct {
	Key models.Base64EncodedKeyData `json:"key"`
}

type OldVerifyKey struct {
	Key       models.Base64EncodedKeyData `json:"key"`
	ExpiredTs models.Timestamp            `json:"expired_ts"`
}

type ServerKeyResult struct {
	ServerName    string                        `json:"server_name"`
	ValidUntilTs  int64                         `json:"valid_until_ts"`
	VerifyKeys    map[models.KeyID]VerifyKey    `json:"verify_keys"`
	OldVerifyKeys map[models.KeyID]OldVerifyKey `json:"old_verify_keys"`
}

func GetLocalKeys(r *http.Request, log *logrus.Entry) interface{} {
	ownKeys, err := db.GetAllOwnKeys()
	if err != nil {
		logrus.Error(err)
		return common.InternalServerError("Failed to get keys")
	}

	resp := &ServerKeyResult{
		ServerName:    keys.SelfDomainName,
		ValidUntilTs:  util.NowMillis() + 86400000, // 24 hours
		VerifyKeys:    make(map[models.KeyID]VerifyKey),
		OldVerifyKeys: make(map[models.KeyID]OldVerifyKey),
	}

	for _, k := range ownKeys {
		if k.ExpiresTs > 0 {
			resp.OldVerifyKeys[k.ID] = OldVerifyKey{
				ExpiredTs: k.ExpiresTs,
				Key:       k.PublicKey,
			}
		} else {
			resp.VerifyKeys[k.ID] = VerifyKey{
				Key: k.PublicKey,
			}
		}
	}

	// TODO: Sign object

	return resp
}
