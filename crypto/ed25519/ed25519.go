package ed25519

import (
	"crypto/ed25519"
	"io"
)

const (
	// PublicKeySize is the size, in bytes, of public keys as used in this package.
	PublicKeySize = ed25519.PublicKeySize
	// PrivateKeySize is the size, in bytes, of private keys as used in this package.
	PrivateKeySize = ed25519.PrivateKeySize
	// SignatureSize is the size, in bytes, of signatures generated and verified by this package.
	SignatureSize = ed25519.SignatureSize
	// SeedSize is the size, in bytes, of private key seeds. These are the private key representations used by RFC 8032.
	SeedSize = ed25519.SeedSize
)

// PublicKey is the type of Ed25519 public keys.
type PublicKey ed25519.PublicKey

// PrivateKey is the type of Ed25519 private keys. It implements crypto.Signer.
type PrivateKey ed25519.PrivateKey

// GenerateKeyPair generates a public/private key pair using entropy from rand.
// If rand is nil, crypto/rand.Reader will be used.
func GenerateKeyPair(rand io.Reader) (PublicKey, PrivateKey, error) {
	return ed25519.Sign(rand)
}
