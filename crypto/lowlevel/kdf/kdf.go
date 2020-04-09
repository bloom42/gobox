package kdf

import (
	"errors"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/blake2b"
)

const (
	KeySize256 = 32
	KeySize384 = 48
	KeySize512 = 64
)

// DeriveFromPassword derives a key from a human provided password using the argon2id Key Derivation
// Function
func DeriveFromPassword(password, salt []byte, keyLen uint32) ([]byte, error) {
	var time uint32 = 2
	var memory uint32 = 32 * 1024
	var threads uint8 = 4

	key := argon2.IDKey(password, salt, time, memory, threads, keyLen)
	if key == nil {
		return nil, errors.New("Error deriving key from password")
	}
	return key, nil
}

// DeriveFromKey derives a key from a high entropy key using the blake2b function
func DeriveFromKey(key, context []byte, keyLen uint8) ([]byte, error) {
	if keyLen < 1 || keyLen > 64 {
		return nil, errors.New("keyLen must be between 1 and 64")
	}

	blake2bHash, err := blake2b.New(int(keyLen), nil)
	if err != nil {
		return nil, err
	}

	blake2bHash.Write(key)
	return blake2bHash.Sum(context), nil
}