package crypto

import (
	"testing"
)

func TestGenerateEd25519Keypair(t *testing.T) {
	zeroPublicKey := make([]byte, Ed25519PublicKeySize)
	zeroPrivateKey := make([]byte, Ed25519PrivateKeySize)

	publicKey, privateKey, err := GenerateEd25519KeyPair()
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
