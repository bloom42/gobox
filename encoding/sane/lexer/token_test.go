package lexer

import (
	"testing"
)

type tokenStringTest struct {
	token TokenType
	str   string
}

func TestTokenString(t *testing.T) {
	toTest := []tokenStringTest{
		{TokenError, "Error"},
		{TokenEOF, "EOF"},
		{TokenComma, "Comma"},
		{TokenEqual, "Equal"},
		{TokenLeftBracket, "LeftBracket"},
		{TokenRightBracket, "RightBracket"},
		{TokenLeftBrace, "LeftBrace"},
		{TokenRightBrace, "RightBrace"},
		{TokenComment, "Comment"},
		{TokenKey, "Key"},
		{TokenBoolean, "Boolean"},
		{TokenInteger, "Integer"},
		{TokenFloat, "Float"},
		{TokenString, "String"},
	}

	for _, pair := range toTest {
		str := pair.token.String()
		if str != pair.str {
			t.Errorf("Expected %s, got %s instead", pair.str, str)
		}
	}
}
