package xchacha20poly1305

import (
	"crypto/cipher"

	"golang.org/x/crypto/chacha20poly1305"
)

const (
	// KeySize is the size of the key used by this AEAD, in bytes.
	KeySize = chacha20poly1305.KeySize

	// NonceSize is the size of the nonce used with the XChaCha20-Poly1305
	// variant of this AEAD, in bytes.
	NonceSize = chacha20poly1305.NonceSizeX
)

//  NewX returns a XChaCha20-Poly1305 AEAD that uses the given 256-bit key.
//
// XChaCha20-Poly1305 is a ChaCha20-Poly1305 variant that takes a longer nonce, suitable to be
// generated randomly without risk of collisions. It should be preferred when nonce uniqueness cannot
// be trivially ensured, or whenever nonces are randomly generated.
func New(key []byte) (cipher.AEAD, error) {
	return chacha20poly1305.NewX(key)
}
