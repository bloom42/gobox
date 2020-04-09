package xchacha20poly1305

import (
	"crypto/cipher"
	"io"

	"gitlab.com/bloom42/lily/crypto/lowlevel/rand"
	"golang.org/x/crypto/chacha20poly1305"
)

const (
	// KeySize is the size of the key used by this AEAD, in bytes.
	KeySize = chacha20poly1305.KeySize

	// NonceSize is the size of the nonce used with the XChaCha20-Poly1305
	// variant of this AEAD, in bytes.
	NonceSize = chacha20poly1305.NonceSizeX
)

// NewKey generates a new random secret key.
func NewKey() (*[KeySize]byte, error) {
	key := new([KeySize]byte)
	_, err := io.ReadFull(rand.Reader(), key[:])
	if err != nil {
		return nil, err
	}

	return key, nil
}

// NewNonce generates a new random nonce.
func NewNonce() (*[NonceSize]byte, error) {
	nonce := new([NonceSize]byte)
	_, err := io.ReadFull(rand.Reader(), nonce[:])
	if err != nil {
		return nil, err
	}

	return nonce, nil
}

// New returns a XChaCha20-Poly1305 AEAD that uses the given 256-bit key.
//
// XChaCha20-Poly1305 is a ChaCha20-Poly1305 variant that takes a longer nonce, suitable to be
// generated randomly without risk of collisions. It should be preferred when nonce uniqueness cannot
// be trivially ensured, or whenever nonces are randomly generated.
func New(key []byte) (cipher.AEAD, error) {
	return chacha20poly1305.NewX(key)
}
