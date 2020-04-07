package curve25519

import (
	"fmt"

	"gitlab.com/bloom42/lily/crypto/rand"
	"golang.org/x/crypto/curve25519"
)

const (
	// KeySize is the size of both private and public keys
	KeySize = curve25519.ScalarSize
)

func NewKeyPair() (publicKey, privateKey []byte, err error) {
	privateKey = make([]byte, curve25519.ScalarSize)

	if _, err = rand.Reader().Read(privateKey); err != nil {
		return nil, nil, fmt.Errorf("internal error: %v", err)
	}
	publicKey, err = curve25519.X25519(privateKey, curve25519.Basepoint)
	if err != nil {
		return nil, nil, fmt.Errorf("internal error: %v", err)
	}

	return
}
