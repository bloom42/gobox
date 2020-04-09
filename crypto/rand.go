package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"math/big"
)

// RandBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func RandBytes(n uint64) ([]byte, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// RandBytesHex returns securely generated random bytes encoded as a hex string.
func RandBytesHex(n uint64) (string, error) {
	data, err := RandBytes(n)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(data), nil
}

// RandBytesBase64 returns securely generated random bytes encoded as a standard base64 string.
func RandBytesBase64(n uint64) (string, error) {
	data, err := RandBytes(n)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

// RandInt64 returns a uniform random value in [min, max).
func RandInt64(min, max int64) (int64, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(max-min))
	if err != nil {
		return max, err
	}
	return n.Int64(), nil
}

// RandAlphabet returns a buffer a size n filled with random values taken from alphabet
func RandAlphabet(alphabet []byte, n uint64) ([]byte, error) {
	buffer := make([]byte, n)
	alphabetLen := int64(len(alphabet))

	for i := range buffer {
		n, err := RandInt64(0, alphabetLen)
		if err != nil {
			return nil, err
		}
		buffer[i] = alphabet[n]
	}
	return buffer, nil
}

// RandReader returns a cryptographically secure source of entropy which implements the `io.Reader`
// interface.
func RandReader() io.Reader {
	return rand.Reader
}
