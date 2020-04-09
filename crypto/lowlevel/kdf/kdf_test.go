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
	context := []byte("com.bloom42.lily")

	_, err = DeriveFromKey(key, context, 128)
	if err == nil {
		t.Error("Accept invalid keyLen")
	}

	_, err = DeriveFromKey(key, context, 65)
	if err == nil {
		t.Error("Accept invalid keyLen")
	}

	_, err = DeriveFromKey(key, context, 0)
	if err == nil {
		t.Error("Accept invalid keyLen")
	}

	_, err = DeriveFromKey(key, context, 1)
	if err != nil {
		t.Error("Reject valid keyLen")
	}

	_, err = DeriveFromKey(key, context, 64)
	if err != nil {
		t.Error("Reject valid keyLen")
	}
}
