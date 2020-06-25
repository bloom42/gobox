package crypto

import (
	"golang.org/x/crypto/curve25519"
)

// Encrypt convert the PublicKey to a `curve25519` public key using `ToCurve25519PublicKey`,
// the privateKey to a `curve25519` prviate key using `ToCurve25519PrivateKey`,
// then perform a x25519 key exchange, and finally encrypt the message using `XChaCha20-Poly1305` with
// the shared secret as key and nonce as nonce.
func (publicKey PublicKey) Encrypt(fromPrivateKey PrivateKey, nonce []byte, message []byte) (ciphertext []byte, err error) {
	curve25519PublicKey := publicKey.ToCurve25519PublicKey()
	curve25519FromPrivateKey := fromPrivateKey.ToCurve25519PrivateKey()
	defer Zeroize(curve25519FromPrivateKey)

	sharedSecret, err := curve25519.X25519(curve25519FromPrivateKey, curve25519PublicKey)
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

// EncryptEphemeral generates an ephemeral keyPair and `Encrypt` message using the public key,
// the ephemeral privateKey and `blake2b(size=AEADNonceSize, message=ephemeralPublicKey || publicKey)` as nonce
func (publicKey PublicKey) EncryptEphemeral(message []byte) (ciphertext []byte, ephemeralPublicKey PublicKey, err error) {
	ephemeralPublicKey, ephemeralPrivateKey, err := GenerateKeyPair(RandReader())
	defer Zeroize(ephemeralPrivateKey)
	if err != nil {
		return
	}

	// generate nonce
	var nonceMessage []byte
	nonceMessage = append(nonceMessage, []byte(ephemeralPublicKey)...)
	nonceMessage = append(nonceMessage, []byte(publicKey)...)
	hash, err := NewHash(AEADNonceSize, nil)
	if err != nil {
		return
	}
	hash.Write(nonceMessage)
	nonce := hash.Sum(nil)

	ciphertext, err = publicKey.Encrypt(ephemeralPrivateKey, nonce, message)
	return
}

// Decrypt convert the privateKey to a `curve25519` prviate key using `ToCurve25519PrivateKey`,
// the fromPublicKey to a `curve25519` public keys using `ToCurve25519PublicKey`,
// then perform a x25519 key exchange, and finally decrypt the ciphertext using `XChaCha20-Poly1305` with
// the shared secret as key and nonce as nonce.
func (privateKey PrivateKey) Decrypt(fromPublicKey PublicKey, nonce []byte, ciphertext []byte) (plaintext []byte, err error) {
	curve25519FromPublicKey := fromPublicKey.ToCurve25519PublicKey()
	curve25519PrivateKey := privateKey.ToCurve25519PrivateKey()
	defer Zeroize(curve25519PrivateKey)

	sharedSecret, err := curve25519.X25519(curve25519PrivateKey, curve25519FromPublicKey)
	defer Zeroize(sharedSecret)
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

// DecryptEphemeral generates a noce with `blake2b(size=AEADNonceSize, message=ephemeralPublicKey || privateKey.PublicKey())`
// and decrypt the `ciphertext` using `Decrypt`
func (privateKey PrivateKey) DecryptEphemeral(ephemeralPublicKey PublicKey, ciphertext []byte) (plaintext []byte, err error) {
	var nonceMessage []byte
	myPublicKey := privateKey.Public()

	// generate nonce
	nonceMessage = append(nonceMessage, []byte(ephemeralPublicKey)...)
	nonceMessage = append(nonceMessage, []byte(myPublicKey)...)
	hash, err := NewHash(AEADNonceSize, nil)
	if err != nil {
		return
	}
	hash.Write(nonceMessage)
	nonce := hash.Sum(nil)

	return privateKey.Decrypt(ephemeralPublicKey, nonce, ciphertext)
}
