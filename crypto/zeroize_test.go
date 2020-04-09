package crypto

import (
	"testing"
)

func TestZeroize(t *testing.T) {
	buffer := []byte("random buffer")

	Zeroize(buffer)

	for i := range buffer {
		if buffer[i] != 0 {
			t.Errorf("buffer not zeroized (index %d)", i)
		}
	}
}
