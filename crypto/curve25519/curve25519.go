package curve25519

import (
	"fmt"

	"gitlab.com/bloom42/lily/crypto/rand"
	"golang.org/x/crypto/curve25519"
)

const (
	// KeySize is the size of both private and public keys, in bytes
	KeySize = curve25519.ScalarSize
)

// NewKeyPair genrates a new private and public key pair using a secure random source
// both keys are of `KeySize` size
func NewKeyPair() (publicKey, privateKey []byte, err error) {
	privateKey, err = rand.Bytes(KeySize)
	if err != nil {
		return nil, nil, fmt.Errorf("internal error: %v", err)
	}
	publicKey, err = curve25519.X25519(privateKey, curve25519.Basepoint)
	if err != nil {
		return nil, nil, fmt.Errorf("internal error: %v", err)
	}

	return
}
