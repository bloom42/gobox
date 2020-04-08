package curve25519

import (
	"gitlab.com/bloom42/lily/crypto/rand"
	"golang.org/x/crypto/curve25519"
)

const (
	// KeySize is the size of both private and public keys, in bytes.
	KeySize = curve25519.ScalarSize
)

// NewKeyPair genrates a new private and public key pair using a secure random source
// both keys are of `KeySize` size
func NewKeyPair() (publicKey, privateKey *[KeySize]byte, err error) {
	_, err = rand.Reader().Read(privateKey[:])
	if err != nil {
		return
	}
	curve25519.ScalarBaseMult(publicKey, privateKey)

	return
}
