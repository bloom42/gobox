package crypto

import (
	"hash"

	"golang.org/x/crypto/blake2b"
)

// HashSize is the size of a hash, in bytes.
type HashSize uint32

const (
	// HashSize256 is the size in bytes of a 256 bits hash
	HashSize256 HashSize = 32
	// HashSize384 is the size in bytes of a 384 bits hash
	HashSize384 HashSize = 48
	// HashSize512 is the size in bytes of a 512 bits hash
	HashSize512 HashSize = 64
)

// NewHash returns a new hash.Hash computing the BLAKE2b checksum with a custom length.
// size can be a value between 1 and 64.
// It is highly recommended to use values equal or greater than 32.
func NewHash(size HashSize) (hash.Hash, error) {
	return blake2b.New(int(size), nil)
}
