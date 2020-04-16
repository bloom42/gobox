// Package base64 provides a high-level API to encode and decode unpadded alternate base64 encoding
// defined in RFC 4648
package base64

import (
	"encoding/base64"
)

// DecodeURLUnpaddedString returns the bytes represented by the base64 string str.
func DecodeURLUnpaddedString(str string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(str)
}

// EncodeToURLUnpaddedString returns the base64 encoding of data.
func EncodeToURLUnpaddedString(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}
