package lexer

import (
	"testing"
)

func TestPositionString(t *testing.T) {
	p := Position{123, 456}
	expected := "(123, 456)"
	value := p.String()

	if value != expected {
		t.Errorf("Expected %v, got %v instead", expected, value)
	}
}

func TestValid(t *testing.T) {
	for i, v := range []Position{
		{0, 1234},
		{1234, 0},
		{0, 0},
	} {
		if v.Valid() {
			t.Errorf("Position at %v is invalid: %v", i, v)
		}
	}
}
