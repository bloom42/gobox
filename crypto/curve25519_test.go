package crypto

import (
	"testing"
)

func TestCurve25519EncryptDecrypt(t *testing.T) {
	message := []byte("this is a simple message")
	nonce, err := RandBytes(AEADNonceSize)
	if err != nil {
		t.Error(err)
	}

	toPublicKey, toPrivateKey, err := GenerateCurve25519KeyPair()
	if err != nil {
		t.Error(err)
	}

	fromPublicKey, fromPrivateKey, err := GenerateCurve25519KeyPair()
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

func TestCurve25519EncryptDecryptEphemeral(t *testing.T) {
	message := []byte("this is a simple message")

	toPublicKey, toPrivateKey, err := GenerateCurve25519KeyPair()
	if err != nil {
		t.Error(err)
	}

	ciphertext, ephemeralPublicKey, err := toPublicKey.EncryptEphemeral(message)
	if err != nil {
		t.Error(err)
	}

	plaintext, err := toPrivateKey.DecryptEphemeral(ephemeralPublicKey, ciphertext)
	if err != nil {
		t.Error(err)
	}

	if !ConstantTimeCompare(message, plaintext) {
		t.Errorf("Message (%s) and plaintext (%s) don't match", string(message), string(plaintext))
	}
}
