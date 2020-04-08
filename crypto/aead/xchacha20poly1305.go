package aead

import (
	"crypto/cipher"
	"io"

	"gitlab.com/bloom42/lily/crypto/rand"
	"golang.org/x/crypto/chacha20poly1305"
)

const (
	// XChaCha20Poly1305KeySize is the size of the key used by this AEAD, in bytes.
	XChaCha20Poly1305KeySize = chacha20poly1305.KeySize

	// XChaCha20Poly1305NonceSize is the size of the nonce used with the XChaCha20-Poly1305
	// variant of this AEAD, in bytes.
	XChaCha20Poly1305NonceSize = chacha20poly1305.NonceSizeX
)

// NewXChaCha20Poly1305Key generates a new random secret key.
func NewXChaCha20Poly1305Key() (*[XChaCha20Poly1305KeySize]byte, error) {
	key := new([XChaCha20Poly1305KeySize]byte)
	_, err := io.ReadFull(rand.Reader(), key[:])
	if err != nil {
		return nil, err
	}

	return key, nil
}

// NewXChaCha20Poly1305Nonce generates a new random nonce.
func NewXChaCha20Poly1305Nonce() (*[XChaCha20Poly1305NonceSize]byte, error) {
	nonce := new([XChaCha20Poly1305NonceSize]byte)
	_, err := io.ReadFull(rand.Reader(), nonce[:])
	if err != nil {
		return nil, err
	}

	return nonce, nil
}

// NewXChaCha20Poly1305 returns a XChaCha20-Poly1305 AEAD that uses the given 256-bit key.
//
// XChaCha20-Poly1305 is a ChaCha20-Poly1305 variant that takes a longer nonce, suitable to be
// generated randomly without risk of collisions. It should be preferred when nonce uniqueness cannot
// be trivially ensured, or whenever nonces are randomly generated.
func NewXChaCha20Poly1305(key []byte) (cipher.AEAD, error) {
	return chacha20poly1305.NewX(key)
}
