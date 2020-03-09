package sane

import (
	"reflect"
	"testing"
)

type test1 struct {
	A string
	B int64 `sane:"c"`
}

// AssertEqual checks if values are equal
func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	// debug.PrintStack()
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}

func TestSimpleStruct(t *testing.T) {
	str := `
a = "test"
c = 1
`
	var s test1

	err := Unmarshal([]byte(str), &s)
	if err != nil {
		t.Error(err)
	}

	AssertEqual(t, "test", s.A)
	AssertEqual(t, int64(1), s.B)
}
