package kdf

import (
	"testing"
)

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
