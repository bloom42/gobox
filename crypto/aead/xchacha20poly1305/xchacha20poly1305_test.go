package xchacha20poly1305

import "testing"

func TestKeySize(t *testing.T) {
	if KeySize != 32 {
		t.Error("KeySize != 32")
	}
}

func TestNonceSize(t *testing.T) {
	if NonceSize != 24 {
		t.Error("NonceSize != 24")
	}
}

func TestNewKey(t *testing.T) {
	key, err := NewKey()
	if err != nil {
		t.Error(err)
	}

	if len(key) != KeySize {
		t.Errorf("generated key have bad size (%d)", len(key))
	}
}

func TestNewNonce(t *testing.T) {
	nonce, err := NewNonce()
	if err != nil {
		t.Error(err)
	}

	if len(nonce) != NonceSize {
		t.Errorf("generated nonce have bad size (%d)", len(nonce))
	}
}
