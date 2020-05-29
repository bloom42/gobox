package crypto

import (
	"testing"
)

func TestGenerateKeypair(t *testing.T) {
	zeroPublicKey := make([]byte, PublicKeySize)
	zeroPrivateKey := make([]byte, PrivateKeySize)

	publicKey, privateKey, err := GenerateKeyPair(RandReader())
	if err != nil {
		t.Error(err)
	}

	if ConstantTimeCompare(zeroPrivateKey, privateKey) {
		t.Error("Generated private key is empty")
	}

	if ConstantTimeCompare(zeroPublicKey, publicKey) {
		t.Error("Generated public key is empty")
	}
}

func TestPrivateKeyEncryptDecrypt(t *testing.T) {
	message := []byte("this is a simple message")
	nonce, err := RandBytes(AEADNonceSize)
	if err != nil {
		t.Error(err)
	}

	toPublicKey, toPrivateKey, err := GenerateKeyPair(RandReader())
	if err != nil {
		t.Error(err)
	}

	fromPublicKey, fromPrivateKey, err := GenerateKeyPair(RandReader())
	if err != nil {
		t.Error(err)
	}

	ciphertext, err := toPublicKey.Encrypt(fromPrivateKey, nonce, message)
	if err != nil {
		t.Error(err)
	}

	plaintext, err := toPrivateKey.Decrypt(fromPublicKey, nonce, ciphertext)
	if err != nil {
		t.Error(err)
	}

	if !ConstantTimeCompare(message, plaintext) {
		t.Errorf("Message (%s) and plaintext (%s) don't match", string(message), string(plaintext))
	}
}

func TestPrivateKeyEncryptDecryptAnonymous(t *testing.T) {
	message := []byte("this is a simple message")

	toPublicKey, toPrivateKey, err := GenerateKeyPair(RandReader())
	if err != nil {
		t.Error(err)
	}

	ciphertext, ephemeralPublicKey, err := toPublicKey.EncryptAnonymous(message)
	if err != nil {
		t.Error(err)
	}

	plaintext, err := toPrivateKey.DecryptAnonymous(ephemeralPublicKey, ciphertext)
	if err != nil {
		t.Error(err)
	}

	if !ConstantTimeCompare(message, plaintext) {
		t.Errorf("Message (%s) and plaintext (%s) don't match", string(message), string(plaintext))
	}
}
