package crypto

import (
	"crypto"
	"crypto/ed25519"
	"crypto/sha512"
	"errors"
	"io"
	"math/big"

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

	// PrivateKeySignerOpts must be used for `PrivateKey.Sign`
	PrivateKeySignerOpts = crypto.Hash(0)
)

// PublicKey is the type of Ed25519 public keys.
type PublicKey ed25519.PublicKey

// Verify reports whether sig is a valid signature of message by publicKey.
// returns true if signature is valid. false otherwise.
func (public PublicKey) Verify(message, signature []byte) (bool, error) {
	if len(public) != PublicKeySize {
		return false, errors.New("crypto: Invalid public key size")
	}

	return ed25519.Verify(ed25519.PublicKey(public), message, signature), nil
}

var curve25519P, _ = new(big.Int).SetString("57896044618658097711785492504343953926634992332820282019728792003956564819949", 10)

// ToCurve25519PublicKey returns the corresponding Curve25519 public key.
//
// See here for more details: https://blog.filippo.io/using-ed25519-keys-for-encryption
func (public PublicKey) ToCurve25519PublicKey() []byte {
	// taken from https://github.com/FiloSottile/age/blob/master/internal/agessh/agessh.go#L179

	// ed25519.PublicKey is a little endian representation of the y-coordinate,
	// with the most significant bit set based on the sign of the x-coordinate.
	bigEndianY := make([]byte, PublicKeySize)
	for i, b := range public {
		bigEndianY[PublicKeySize-i-1] = b
	}
	bigEndianY[0] &= 0b0111_1111

	// The Montgomery u-coordinate is derived through the bilinear map
	//
	//     u = (1 + y) / (1 - y)
	//
	// See https://blog.filippo.io/using-ed25519-keys-for-encryption.
	y := new(big.Int).SetBytes(bigEndianY)
	denom := big.NewInt(1)
	denom.ModInverse(denom.Sub(denom, y), curve25519P) // 1 / (1 - y)
	u := y.Mul(y.Add(y, big.NewInt(1)), denom)
	u.Mod(u, curve25519P)

	out := make([]byte, curve25519.PointSize)
	uBytes := u.Bytes()
	for i, b := range uBytes {
		out[len(uBytes)-i-1] = b
	}

	return out
}

// Encrypt convert the PublicKey to a `curve25519` public key using `ToCurve25519PublicKey`,
// then perform a key exchange, then encrypt the message using XChaCha20-Poly1305
func (public PublicKey) Encrypt(message []byte, privateKey PrivateKey, nonce []byte) (ciphertext []byte, err error) {
	curve25519PublicKey := public.ToCurve25519PublicKey()
	curve25519PrivateKey := privateKey.ToCurve25519PrivateKey()
	defer Zeroize(curve25519PrivateKey)

	sharedSecret, err := curve25519.X25519(curve25519PrivateKey, curve25519PublicKey)
	defer Zeroize(sharedSecret)
	if err != nil {
		return
	}

	cipher, err := NewAEAD(sharedSecret)
	if err != nil {
		return
	}

	ciphertext = cipher.Seal(nil, nonce, message, nil)
	return
}

func (public PublicKey) EncryptAnonymous(message []byte) (ciphertext []byte, ephemeralPublicKey PublicKey, err error) {
	ephemeralPublicKey, ephemeralPrivateKey, err := GenerateKeyPair(RandReader())
	defer Zeroize(ephemeralPrivateKey)
	if err != nil {
		return
	}

	// generate nonce
	var nonceMessage []byte
	nonceMessage = append(nonceMessage, []byte(public)...)
	nonceMessage = append(nonceMessage, []byte(ephemeralPublicKey)...)
	hash, err := NewHash(AEADNonceSize, nil)
	if err != nil {
		return
	}
	hash.Write(nonceMessage)
	nonce := hash.Sum(nil)

	ciphertext, err = public.Encrypt(message, ephemeralPrivateKey, nonce)
	return
}

// PrivateKey is the type of Ed25519 private keys. It implements crypto.Signer.
type PrivateKey ed25519.PrivateKey

// Sign signs the given message with priv.
// Ed25519 performs two passes over messages to be signed and therefore cannot
// handle pre-hashed messages. Thus opts.HashFunc() must return zero to
// indicate the message hasn't been hashed. This can be achieved by passing
// crypto.PrivateKeySignerOpts as the value for opts.
func (priv PrivateKey) Sign(rand io.Reader, message []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	if len(priv) != PrivateKeySize {
		return nil, errors.New("crpyto: Invalid private key size")
	}

	return ed25519.Sign(ed25519.PrivateKey(priv), message), nil
}

// ToCurve25519PrivateKey returns a corresponding Curve25519 private key.
//
// See here for more details: https://blog.filippo.io/using-ed25519-keys-for-encryption
func (priv PrivateKey) ToCurve25519PrivateKey() []byte {
	// taken from https://github.com/FiloSottile/age/blob/292c3aaeea0695dbba356dfe18a70f10efb17d75/internal/agessh/agessh.go#L294
	h := sha512.New()
	h.Write(ed25519.PrivateKey(priv).Seed())
	out := h.Sum(nil)
	return out[:curve25519.ScalarSize]
}

// Public returns the PublicKey corresponding to priv.
func (priv PrivateKey) Public() PublicKey {
	return PublicKey(ed25519.PrivateKey(priv).Public().(ed25519.PublicKey))
}

// Seed returns the private key seed corresponding to priv. It is provided for interoperability
// with RFC 8032. RFC 8032's private keys correspond to seeds in this package.
func (priv PrivateKey) Seed() []byte {
	return ed25519.PrivateKey(priv).Seed()
}

func (priv PrivateKey) Decrypt(ciphertext []byte, publicKey PublicKey, nonce []byte) (plaintext []byte, err error) {
	curve25519PublicKey := publicKey.ToCurve25519PublicKey()
	curve25519PrivateKey := priv.ToCurve25519PrivateKey()
	defer Zeroize(curve25519PrivateKey)

	sharedSecret, err := curve25519.X25519(curve25519PrivateKey, curve25519PublicKey)
	if err != nil {
		return
	}

	cipher, err := NewAEAD(sharedSecret)
	if err != nil {
		return
	}

	plaintext, err = cipher.Open(nil, nonce, ciphertext, nil)
	return
}

func (priv PrivateKey) DecryptAnonymous(ciphertext []byte, publicKey PublicKey) (plaintext []byte, err error) {
	var nonceMessage []byte
	myPublicKey := priv.Public()

	// generate nonce
	nonceMessage = append(nonceMessage, []byte(myPublicKey)...)
	nonceMessage = append(nonceMessage, []byte(publicKey)...)
	hash, err := NewHash(AEADNonceSize, nil)
	if err != nil {
		return
	}
	hash.Write(nonceMessage)
	nonce := hash.Sum(nil)

	return priv.Decrypt(ciphertext, publicKey, nonce)
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
		return nil, errors.New("crypto: Invalid seed size")
	}

	private := ed25519.NewKeyFromSeed(seed)
	return PrivateKey(private), nil
}
