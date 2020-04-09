package crypto

import (
	"crypto/cipher"
	"crypto/rand"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
)

const (
	// AEADKeySize is the size of the key used by this AEAD, in bytes.
	AEADKeySize = chacha20poly1305.KeySize

	// AEADNonceSize is the size of the nonce used with the XChaCha20-Poly1305
	// variant of this AEAD, in bytes.
	AEADNonceSize = chacha20poly1305.NonceSizeX
)

// NewAEADKey generates a new random secret key.
func NewAEADKey() (*[AEADKeySize]byte, error) {
	key := new([AEADKeySize]byte)
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		return nil, err
	}

	return key, nil
}

// NewAEADNonce generates a new random nonce.
func NewAEADNonce() (*[AEADNonceSize]byte, error) {
	nonce := new([AEADNonceSize]byte)
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return nil, err
	}

	return nonce, nil
}

// NewAEAD returns a XChaCha20-Poly1305 AEAD that uses the given 256-bit key.
//
// XChaCha20-Poly1305 is a ChaCha20-Poly1305 variant that takes a longer nonce, suitable to be
// generated randomly without risk of collisions. It should be preferred when nonce uniqueness cannot
// be trivially ensured, or whenever nonces are randomly generated.
func NewAEAD(key []byte) (cipher.AEAD, error) {
	return chacha20poly1305.NewX(key)
}
