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
}
