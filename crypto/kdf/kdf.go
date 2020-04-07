package kdf

import (
	"errors"
	"hash"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/blake2b"
)

// DeriveFromPassword derives a key from a human provided password using the argon2id Key Derivation
// Function
func DeriveFromPassword(password, salt []byte, keyLen uint32) ([]byte, error) {
	var time uint32 = 3
	var memory uint32 = 32 * 1024
	var threads uint8 = 4

	key := argon2.IDKey(password, salt, time, memory, threads, keyLen)
	if key == nil {
		return nil, errors.New("Error deriving key from password")
	}
	return key, nil
}

// DeriveFromKey derives a key from a high entropy key using the blake2b function
func DeriveFromKey(key []byte, keyLen int) ([]byte, error) {
	var err error
	var blake2bHash hash.Hash

	if keyLen != blake2b.Size && keyLen != blake2b.Size384 && keyLen != blake2b.Size256 {
		return nil, errors.New("Invalid keyLen parameter. Must be 32, 48 or 64")
	}

	if keyLen == blake2b.Size256 {
		blake2bHash, err = blake2b.New256(nil)
	} else if keyLen == blake2b.Size384 {
		blake2bHash, err = blake2b.New384(nil)
	} else if keyLen == blake2b.Size384 {
		blake2bHash, err = blake2b.New512(nil)
	} else {
		err = errors.New("Invalid keyLen parameter. Must be 32, 48 or 64")
	}
	if err != nil {
		return nil, err
	}

	return blake2bHash.Sum(key), nil
}
