package crypto

import (
	"errors"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/blake2b"
)

const (
	// KeySize256 is the size in bytes of a 256 bits key
	KeySize256 = 32
	// KeySize384 is the size in bytes of a 384 bits key
	KeySize384 = 48
	// KeySize512 is the size in bytes of a 512 bits key
	KeySize512 = 64
	// KeySize1024 is the size in bytes of a 1024 bits key
	KeySize1024 = 128
	// KeySize2048 is the size in bytes of a 2048 bits key
	KeySize2048 = 256
	// KeySize4096 is the size in bytes of a 4096 bits key
	KeySize4096 = 512
)

// DeriveKeyFromPassword derives a key from a human provided password using the argon2id Key Derivation
// Function
func DeriveKeyFromPassword(password, salt []byte, keySize uint32) ([]byte, error) {
	var time uint32 = 2
	var memory uint32 = 32 * 1024
	var threads uint8 = 4

	key := argon2.IDKey(password, salt, time, memory, threads, keySize)
	if key == nil {
		return nil, errors.New("crypto: Deriving key from password")
	}
	return key, nil
}

// DeriveKeyFromKey derives a key from a high entropy key using the blake2b function
func DeriveKeyFromKey(key, info []byte, keySize uint8) ([]byte, error) {
	if keySize < 1 || keySize > 64 {
		return nil, errors.New("crypto: keySize must be between 1 and 64")
	}

	blake2bHash, err := blake2b.New(int(keySize), key)
	if err != nil {
		return nil, err
	}

	blake2bHash.Write(info)
	return blake2bHash.Sum(nil), nil
}
