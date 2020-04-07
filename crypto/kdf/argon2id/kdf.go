package argon2id

import (
	"errors"

	"golang.org/x/crypto/argon2"
)

// DeriveFromPassword derives a key from a human provided password
func DeriveFromPassword(password, salt []byte, time uint32, memory uint32, threads uint8, keyLen uint32) ([]byte, error) {
	key := argon2.Key(password, salt, time, memory, threads, keyLen)
	if key == nil {
		return errors.New("Error deriving key from password")
	}
}
