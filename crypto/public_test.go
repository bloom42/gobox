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

func TestPrivateKeyEncryptAnonymousDecrypt(t *testing.T) {
	message := []byte("this is a simple message")

	toPublicKey, toPrivateKey, err := GenerateKeyPair(RandReader())
	if err != nil {
		t.Error(err)
	}

	ciphertext, ephemeralPublicKey, err := toPublicKey.EncryptAnonymous(message)
	if err != nil {
		t.Error(err)
	}

	plaintext, err := toPrivateKey.Decrypt(ciphertext, ephemeralPublicKey)
	if err != nil {
		t.Error(err)
	}

	if string(message) != string(plaintext) {
		t.Errorf("Message (%s) and plaintext (%s) don't match", string(message), string(plaintext))
	}
}
