package kdf

import (
	"errors"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/blake2b"
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
func DeriveFromKey(key []byte, keyLen uint8) ([]byte, error) {
	if keyLen > 64 {
		return nil, errors.New("Invalid keylen parameter. Must be inferior or equal to 64")
	}

	blake2bHash, err := blake2b.New(int(keyLen), nil)
	if err != nil {
		return nil, err
	}

	return blake2bHash.Sum(key), nil
}
