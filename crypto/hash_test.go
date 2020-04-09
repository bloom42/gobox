package crypto

import "testing"

func TestHashSizes(t *testing.T) {
	if HashSize256 != 32 {
		t.Error("HashSize256 != 32")
	}
	if HashSize384 != 48 {
		t.Error("HashSize384 != 48")
	}
	if HashSize512 != 64 {
		t.Error("HashSize512 != 64")
	}
}
