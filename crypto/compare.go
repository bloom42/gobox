package crypto

import (
	"crypto/subtle"
)

func ConstantTimeCompare(x, y []byte) bool {
	res := subtle.ConstantTimeCompare(x, y)
	return res == 1
}
