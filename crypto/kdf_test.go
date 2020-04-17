package crypto

import (
	"bytes"
	"testing"
)

func TestKeySizes(t *testing.T) {
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

func TestDeriveKeyFromKeyKeyLen(t *testing.T) {
	var err error
	key := []byte("some random data")
	context := []byte("com.bloom42.lily")

	_, err = DeriveKeyFromKey(key, context, 128)
	if err == nil {
		t.Error("Accept invalid keyLen")
	}

	_, err = DeriveKeyFromKey(key, context, 65)
	if err == nil {
		t.Error("Accept invalid keyLen")
	}

	_, err = DeriveKeyFromKey(key, context, 0)
	if err == nil {
		t.Error("Accept invalid keyLen")
	}

	_, err = DeriveKeyFromKey(key, context, 1)
	if err != nil {
		t.Error("Reject valid keyLen")
	}

	_, err = DeriveKeyFromKey(key, context, 64)
	if err != nil {
		t.Error("Reject valid keyLen")
	}
}

func TestDeriveKeyFromKeyContext(t *testing.T) {
	key := []byte("some random data")
	context1 := []byte("com.bloom42.lily1")
	context2 := []byte("com.bloom42.lily2")

	subKey1, err := DeriveKeyFromKey(key, context1, KeySize512)
	if err != nil {
		t.Error(err)
	}

	subKey2, err := DeriveKeyFromKey(key, context2, KeySize512)
	if err != nil {
		t.Error(err)
	}

	if bytes.Equal(subKey1, subKey2) {
		t.Error("subKey1 and subKey2 are equal")
	}

	subKey3, err := DeriveKeyFromKey(key, context1, KeySize512)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(subKey1, subKey3) {
		t.Error("subKey1 and subKey3 are different")
	}
}
