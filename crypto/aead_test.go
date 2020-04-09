package crypto

import "testing"

func TestAEADKeySize(t *testing.T) {
	if AEADKeySize != 32 {
		t.Error("AEADKeySize != 32")
	}
}

func TestAEADNonceSize(t *testing.T) {
	if AEADNonceSize != 24 {
		t.Error("AEADNonceSize != 24")
	}
}

func TestNewAEADKey(t *testing.T) {
	key, err := NewAEADKey()
	if err != nil {
		t.Error(err)
	}

	if len(key) != AEADKeySize {
		t.Errorf("generated key have bad size (%d)", len(key))
	}
}

func TestNewAEADNonce(t *testing.T) {
	nonce, err := NewAEADNonce()
	if err != nil {
		t.Error(err)
	}

	if len(nonce) != AEADNonceSize {
		t.Errorf("generated nonce have bad size (%d)", len(nonce))
	}
}
