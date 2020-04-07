package kdf

import (
	"testing"
)

func TestKeySize(t *testing.T) {
	if KeySize256 != 32 {
		t.Error("KeySize256 != 32")
	}
	if KeySize384 != 48 {
		t.Error("KeySize384 != 48")
	}
	if KeySize512 != 64 {
		t.Error("KeySize512 != 64")
	}
}

func TestDeriveFromKeyKeyLen(t *testing.T) {
	var err error
	key := []byte("some random data")

	_, err = DeriveFromKey(key, 128)
	if err == nil {
		t.Error("Accept invalid keyLen")
	}

	_, err = DeriveFromKey(key, 65)
	if err == nil {
		t.Error("Accept invalid keyLen")
	}

	_, err = DeriveFromKey(key, 64)
	if err != nil {
		t.Error("Reject valid keyLen")
	}
}
