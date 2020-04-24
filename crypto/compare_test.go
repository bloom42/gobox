package crypto

import "testing"

func TestConstantTimeCompare(t *testing.T) {
	a := []byte("helloWorld")
	b := []byte("helloWorld")
	c := []byte("helloWarld")
	d := []byte("helloWorl")

	if !ConstantTimeCompare(a, b) {
		t.Errorf("%s != %s", a, b)
	}

	if ConstantTimeCompare(a, c) {
		t.Errorf("%s == %s", a, c)
	}

	if ConstantTimeCompare(a, d) {
		t.Errorf("%s == %s", a, d)
	}
}
