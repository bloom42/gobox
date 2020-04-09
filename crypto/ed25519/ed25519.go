package ed25519

import (
	"crypto"
	"crypto/ed25519"
	"crypto/sha512"
	"errors"
	"io"

	"golang.org/x/crypto/curve25519"
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

// Verify reports whether sig is a valid signature of message by publicKey.
// returns true if signature is valid. false otherwise.
func (public PublicKey) Verify(message, signature []byte) (bool, error) {
	if len(public) != PublicKeySize {
		return false, errors.New("ed25519: Invalid public key size")
	}

	return ed25519.Verify(ed25519.PublicKey(public), message, signature), nil
}

// PrivateKey is the type of Ed25519 private keys. It implements crypto.Signer.
type PrivateKey ed25519.PrivateKey

// Sign signs the given message with priv.
// Ed25519 performs two passes over messages to be signed and therefore cannot
// handle pre-hashed messages. Thus opts.HashFunc() must return zero to
// indicate the message hasn't been hashed. This can be achieved by passing
// crypto.Hash(0) as the value for opts.
func (priv PrivateKey) Sign(rand io.Reader, message []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	if len(priv) != PrivateKeySize {
		return nil, errors.New("ed25519: Invalid private key size")
	}

	return ed25519.Sign(ed25519.PrivateKey(priv), message), nil
}

// ToCurve25519PrivateKey returns a corresponding Curve25519 private key
// see here for more details: https://blog.filippo.io/using-ed25519-keys-for-encryption/
func (priv PrivateKey) ToCurve25519PrivateKey() []byte {
	// took from https://github.com/FiloSottile/age/blob/bbab440e198a4d67ba78591176c7853e62d29e04/internal/age/ssh.go#L289
	h := sha512.New()
	h.Write(ed25519.PrivateKey(priv).Seed())
	out := h.Sum(nil)
	return out[:curve25519.ScalarSize]
}

// Public returns the PublicKey corresponding to priv.
func (priv PrivateKey) Public() crypto.PublicKey {
	return ed25519.PrivateKey(priv).Public()
}

// Seed returns the private key seed corresponding to priv. It is provided for interoperability
// with RFC 8032. RFC 8032's private keys correspond to seeds in this package.
func (priv PrivateKey) Seed() []byte {
	return ed25519.PrivateKey(priv).Seed()
}

// GenerateKeyPair generates a public/private key pair using entropy from rand.
// If rand is nil, crypto/rand.Reader will be used.
func GenerateKeyPair(rand io.Reader) (PublicKey, PrivateKey, error) {
	public, private, err := ed25519.GenerateKey(rand)
	return PublicKey(public), PrivateKey(private), err
}

// NewPrivateKeyFromSeed calculates a private key from a seed. It will panic if
// len(seed) is not SeedSize. This function is provided for interoperability
// with RFC 8032. RFC 8032's private keys correspond to seeds in this
// package.
func NewPrivateKeyFromSeed(seed []byte) (PrivateKey, error) {
	if len(seed) != SeedSize {
		return nil, errors.New("ed25519: Invalid seed size")
	}

	private := ed25519.NewKeyFromSeed(seed)
	return PrivateKey(private), nil
}
