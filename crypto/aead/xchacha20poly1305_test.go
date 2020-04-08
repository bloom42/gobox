package aead

import "testing"

func TestXChaCha20Poly1305KeySize(t *testing.T) {
	if XChaCha20Poly1305KeySize != 32 {
		t.Error("XChaCha20Poly1305KeySize != 32")
	}
}

func TestXChaCha20Poly1305NonceSize(t *testing.T) {
	if XChaCha20Poly1305NonceSize != 24 {
		t.Error("XChaCha20Poly1305NonceSize != 24")
	}
}

func TestNewXChaCha20Poly1305Key(t *testing.T) {
	key, err := NewXChaCha20Poly1305Key()
	if err != nil {
		t.Error(err)
	}

	if len(key) != XChaCha20Poly1305KeySize {
		t.Errorf("generated key have bad size (%d)", len(key))
	}
}

func TestNewXChaCha20Poly1305Nonce(t *testing.T) {
	nonce, err := NewXChaCha20Poly1305Nonce()
	if err != nil {
		t.Error(err)
	}

	if len(nonce) != XChaCha20Poly1305NonceSize {
		t.Errorf("generated nonce have bad size (%d)", len(nonce))
	}
}
