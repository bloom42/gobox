package rand

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"math/big"
)

func Int64(min, max int64) (int64, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(max-min))
	if err != nil {
		return max, err
	}
	return n.Int64(), nil
}

func StringAlph(alphabet string, n uint) (string, error) {
	b := make([]byte, n)
	alphabetLen := len(alphabet)

	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(alphabetLen)))
		if err != nil {
			return "", err
		}
		b[i] = alphabet[n.Int64()]
	}
	return string(b), nil
}

// Bytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func Bytes(n uint64) ([]byte, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func TokenBase64(n uint64) (string, error) {
	data, err := Bytes(n)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

func TokenHex(n uint64) (string, error) {
	data, err := Bytes(n)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(data), nil
}

func Reader() io.Reader {
	return rand.Reader
}
