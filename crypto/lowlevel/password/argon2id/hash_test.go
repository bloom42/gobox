package argon2id

import (
	"regexp"
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	hashRX, err := regexp.Compile(`^\$argon2id\$v=19\$m=65536,t=3,p=2\$[A-Za-z0-9+/]{22}\$[A-Za-z0-9+/]{43}$`)
	if err != nil {
		t.Fatal(err)
	}

	hash1, err := HashPassword([]byte("pa$$word"), DefaultHashPasswordParams)
	if err != nil {
		t.Fatal(err)
	}

	if !hashRX.MatchString(hash1) {
		t.Errorf("hash %q not in correct format", hash1)
	}

	hash2, err := HashPassword([]byte("pa$$word"), DefaultHashPasswordParams)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Compare(hash1, hash2) == 0 {
		t.Error("hashes must be unique")
	}
}

func TestVerifyPassword(t *testing.T) {
	hash, err := HashPassword([]byte("pa$$word"), DefaultHashPasswordParams)
	if err != nil {
		t.Fatal(err)
	}

	match := VerifyPassword([]byte("pa$$word"), hash)

	if !match {
		t.Error("expected password and hash to match")
	}

	match = VerifyPassword([]byte("otherPa$$word"), hash)

	if match {
		t.Error("expected password and hash to not match")
	}
}
