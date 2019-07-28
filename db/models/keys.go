package models

type KeyID string
type Base64EncodedKeyData string
type Timestamp int64

type OwnKey struct {
	ID         KeyID
	PublicKey  Base64EncodedKeyData
	PrivateKey Base64EncodedKeyData
	ExpiresTs  Timestamp
}
