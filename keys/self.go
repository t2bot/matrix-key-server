package keys

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/t2bot/matrix-key-server/db"
	"github.com/t2bot/matrix-key-server/db/models"
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

			pubEncoded := base64.RawStdEncoding.EncodeToString(pub)
			privEncoded := base64.RawStdEncoding.EncodeToString(priv)

			dbKey := &models.OwnKey{
				ID:         models.KeyID(keyId),
				PublicKey:  models.Base64EncodedKeyData(pubEncoded),
				PrivateKey: models.Base64EncodedKeyData(privEncoded),
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
			logrus.Info("There are %d active keys for this server", len(activeKeyIds))

			k, err := db.GetOwnKey(activeKeyIds[0])
			if err != nil {
				return nil, err
			}
			if k == nil {
				return nil, errors.New("failed to fetch active key")
			}

			pubDecoded, err := base64.RawStdEncoding.DecodeString(string(k.PublicKey))
			if err != nil {
				return nil, err
			}

			privDecoded, err := base64.RawStdEncoding.DecodeString(string(k.PrivateKey))
			if err != nil {
				return nil, err
			}

			ownKey = &SelfKey{
				RawKey: k,
				ID:     k.ID,
				Pub:    ed25519.PublicKey(pubDecoded),
				Priv:   ed25519.PrivateKey(privDecoded),
			}
		}
	}

	return ownKey, nil
}
