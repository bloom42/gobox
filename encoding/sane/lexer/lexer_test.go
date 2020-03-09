package lexer

import (
	"reflect"
	"testing"
)

func testLex(t *testing.T, input string, expectedTokens []Token) {
	Tokens := Lex([]byte(input))

	if len(expectedTokens) != len(Tokens) {
		t.Fatalf(
			"Not enought Token lexed. Expected: %d\nGot: %d (%#v)\n",
			len(expectedTokens),
			len(Tokens),
			Tokens,
		)
	}

	for i := range expectedTokens {

		if reflect.DeepEqual(Tokens[i], expectedTokens[i]) != true {
			t.Errorf("Different Token.\nExpected: %s\nGot     : %s\n", expectedTokens[i], Tokens[i])
		}
	}
}

func TestSimpleComment(t *testing.T) {
	testLex(t, "# blahblah", []Token{
		{Position{1, 1}, TokenComment, "# blahblah"},
		{Position{1, 11}, TokenEOF, ""},
	})
}

func TestSimpleString(t *testing.T) {
	testLex(t, `str = "myString"`, []Token{
		{Position{1, 1}, TokenKey, "str"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenString, "myString"},
		{Position{1, 17}, TokenEOF, ""},
	})
}

func TestSimpleStringWithComment(t *testing.T) {
	testLex(t, `str = "myString" #comment`, []Token{
		{Position{1, 1}, TokenKey, "str"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenString, "myString"},
		{Position{1, 18}, TokenComment, "#comment"},
		{Position{1, 26}, TokenEOF, ""},
	})
}
func TestDigitalKey(t *testing.T) {
	testLex(t, `123 = "myString"`, []Token{
		{Position{1, 1}, TokenKey, "123"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenString, "myString"},
		{Position{1, 17}, TokenEOF, ""},
	})
}

func TestFloatWithExponent1(t *testing.T) {
	testLex(t, "a = 5e+22", []Token{
		{Position{1, 1}, TokenKey, "a"},
		{Position{1, 3}, TokenEqual, "="},
		{Position{1, 5}, TokenFloat, "5e+22"},
		{Position{1, 10}, TokenEOF, ""},
	})
}

func TestFloatWithExponent2(t *testing.T) {
	testLex(t, "a = 6.69e-22", []Token{
		{Position{1, 1}, TokenKey, "a"},
		{Position{1, 3}, TokenEqual, "="},
		{Position{1, 5}, TokenFloat, "6.69e-22"},
		{Position{1, 13}, TokenEOF, ""},
	})
}

func TestSimpleBoolean(t *testing.T) {
	testLex(t, "foo = false", []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenBoolean, "false"},
		{Position{1, 12}, TokenEOF, ""},
	})
}

func TestNestedStringList(t *testing.T) {
	testLex(t, `a = [ ["hello", "world"] ]`, []Token{
		{Position{1, 1}, TokenKey, "a"},
		{Position{1, 3}, TokenEqual, "="},
		{Position{1, 5}, TokenLeftBracket, "["},
		{Position{1, 7}, TokenLeftBracket, "["},
		{Position{1, 8}, TokenString, "hello"},
		{Position{1, 15}, TokenComma, ","},
		{Position{1, 17}, TokenString, "world"},
		{Position{1, 24}, TokenRightBracket, "]"},
		{Position{1, 26}, TokenRightBracket, "]"},
		{Position{1, 27}, TokenEOF, ""},
	})
}

func TestNestedIntegerList(t *testing.T) {
	testLex(t, "a = [ [42, 21], [10] ]", []Token{
		{Position{1, 1}, TokenKey, "a"},
		{Position{1, 3}, TokenEqual, "="},
		{Position{1, 5}, TokenLeftBracket, "["},
		{Position{1, 7}, TokenLeftBracket, "["},
		{Position{1, 8}, TokenInteger, "42"},
		{Position{1, 10}, TokenComma, ","},
		{Position{1, 12}, TokenInteger, "21"},
		{Position{1, 14}, TokenRightBracket, "]"},
		{Position{1, 15}, TokenComma, ","},
		{Position{1, 17}, TokenLeftBracket, "["},
		{Position{1, 18}, TokenInteger, "10"},
		{Position{1, 20}, TokenRightBracket, "]"},
		{Position{1, 22}, TokenRightBracket, "]"},
		{Position{1, 23}, TokenEOF, ""},
	})
}

func TestIntegerList(t *testing.T) {
	testLex(t, "a = [ 42, 21, 10, ]", []Token{
		{Position{1, 1}, TokenKey, "a"},
		{Position{1, 3}, TokenEqual, "="},
		{Position{1, 5}, TokenLeftBracket, "["},
		{Position{1, 7}, TokenInteger, "42"},
		{Position{1, 9}, TokenComma, ","},
		{Position{1, 11}, TokenInteger, "21"},
		{Position{1, 13}, TokenComma, ","},
		{Position{1, 15}, TokenInteger, "10"},
		{Position{1, 17}, TokenComma, ","},
		{Position{1, 19}, TokenRightBracket, "]"},
		{Position{1, 20}, TokenEOF, ""},
	})
}

func TestBooleanList(t *testing.T) {
	testLex(t, "foo = [true, false, true]", []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenLeftBracket, "["},
		{Position{1, 8}, TokenBoolean, "true"},
		{Position{1, 12}, TokenComma, ","},
		{Position{1, 14}, TokenBoolean, "false"},
		{Position{1, 19}, TokenComma, ","},
		{Position{1, 21}, TokenBoolean, "true"},
		{Position{1, 25}, TokenRightBracket, "]"},
		{Position{1, 26}, TokenEOF, ""},
	})
}

func TestBadList(t *testing.T) {
	testLex(t, "foo = [true", []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenLeftBracket, "["},
		{Position{1, 8}, TokenBoolean, "true"},
		{Position{1, 12}, TokenEOF, ""},
	})
}

func TestMultiString(t *testing.T) {
	str := `a = "tesla"
b = "spacex"
`
	testLex(t, str, []Token{
		{Position{1, 1}, TokenKey, "a"},
		{Position{1, 3}, TokenEqual, "="},
		{Position{1, 5}, TokenString, "tesla"},
		{Position{2, 1}, TokenKey, "b"},
		{Position{2, 3}, TokenEqual, "="},
		{Position{2, 5}, TokenString, "spacex"},
		{Position{3, 1}, TokenEOF, ""},
	})
}

func TestMultiInteger(t *testing.T) {
	testLex(t, "foo = 42\nbar=21", []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenInteger, "42"},
		{Position{2, 1}, TokenKey, "bar"},
		{Position{2, 4}, TokenEqual, "="},
		{Position{2, 5}, TokenInteger, "21"},
		{Position{2, 7}, TokenEOF, ""},
	})
}

func TestNestedLists(t *testing.T) {
	testLex(t, "foo = [[[]]]", []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenLeftBracket, "["},
		{Position{1, 8}, TokenLeftBracket, "["},
		{Position{1, 9}, TokenLeftBracket, "["},
		{Position{1, 10}, TokenRightBracket, "]"},
		{Position{1, 11}, TokenRightBracket, "]"},
		{Position{1, 12}, TokenRightBracket, "]"},
		{Position{1, 13}, TokenEOF, ""},
	})
}

func TestMultilineLists(t *testing.T) {
	testLex(t, `ignore = [".git/**"]
run = [
  'go run testdata/main.go'
]
`, []Token{
		{Position{1, 1}, TokenKey, "ignore"},
		{Position{1, 8}, TokenEqual, "="},
		{Position{1, 10}, TokenLeftBracket, "["},
		{Position{1, 11}, TokenString, ".git/**"},
		{Position{1, 20}, TokenRightBracket, "]"},
		{Position{2, 1}, TokenKey, "run"},
		{Position{2, 5}, TokenEqual, "="},
		{Position{2, 7}, TokenLeftBracket, "["},
		{Position{3, 3}, TokenString, "go run testdata/main.go"},
		{Position{4, 1}, TokenRightBracket, "]"},
		{Position{5, 1}, TokenEOF, ""},
	})
}
func TestSimpleMultilineString(t *testing.T) {
	testLex(t, `foo = """hello "literal" world"""`, []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenString, `hello "literal" world`},
		{Position{1, 34}, TokenEOF, ""},
	})
}

func TestKeyEqualStringNoEscape(t *testing.T) {
	testLex(t, "foo = \"hello \u0002\"", []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenError, "lexer: unescaped control character U+0002"},
	})
	testLex(t, "foo = \"hello \u001F\"", []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenError, "lexer: unescaped control character U+001F"},
	})
}

func TestKeyEqualStringUnicodeEscape(t *testing.T) {
	testLex(t, `foo = "hello \u2665"`, []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenString, "hello ♥"},
		{Position{1, 21}, TokenEOF, ""},
	})
	testLex(t, `foo = "hello \U000003B4"`, []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenString, "hello δ"},
		{Position{1, 25}, TokenEOF, ""},
	})
	testLex(t, `foo = "\uabcd"`, []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenString, "\uabcd"},
		{Position{1, 15}, TokenEOF, ""},
	})
	testLex(t, `foo = "\uABCD"`, []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenString, "\uABCD"},
		{Position{1, 15}, TokenEOF, ""},
	})
	testLex(t, `foo = "\U000bcdef"`, []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenString, "\U000bcdef"},
		{Position{1, 19}, TokenEOF, ""},
	})
	testLex(t, `foo = "\U000BCDEF"`, []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenString, "\U000BCDEF"},
		{Position{1, 19}, TokenEOF, ""},
	})
	testLex(t, `foo = "\u2"`, []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenError, "lexer: unfinished unicode escape"},
	})
	testLex(t, `foo = "\U2"`, []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenError, "lexer: unfinished unicode escape"},
	})
}

func TestLiteralString(t *testing.T) {
	testLex(t, `foo = 'C:\Users\nodejs\templates'`, []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenString, `C:\Users\nodejs\templates`},
		{Position{1, 34}, TokenEOF, ""},
	})
	testLex(t, `foo = '\\ServerX\admin$\system32\'`, []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenString, `\\ServerX\admin$\system32\`},
		{Position{1, 35}, TokenEOF, ""},
	})
	testLex(t, `foo = 'Tom "Dubs" Preston-Werner'`, []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenString, `Tom "Dubs" Preston-Werner`},
		{Position{1, 34}, TokenEOF, ""},
	})
	testLex(t, `foo = '<\i\c*\s*>'`, []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenString, `<\i\c*\s*>`},
		{Position{1, 19}, TokenEOF, ""},
	})
	testLex(t, `foo = 'C:\Users\nodejs\unfinis`, []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenError, "lexer: unclosed string"},
	})
}

func TestMultilineString(t *testing.T) {
	testLex(t, `foo = """hello "literal" world"""`, []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenString, `hello "literal" world`},
		{Position{1, 34}, TokenEOF, ""},
	})

	testLex(t, "foo = \"\"\"\nhello\\\n\"literal\"\\\nworld\"\"\"", []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenString, "hello\"literal\"world"},
		{Position{4, 9}, TokenEOF, ""},
	})

	testLex(t, "foo = \"\"\"\\\n    \\\n    \\\n    hello\\\nmultiline\\\nworld\"\"\"", []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenString, "hellomultilineworld"},
		{Position{6, 9}, TokenEOF, ""},
	})

	testLex(t, "key2 = \"\"\"\nThe quick brown \\\n\n\n  fox jumps over \\\n    the lazy dog.\"\"\"", []Token{
		{Position{1, 1}, TokenKey, "key2"},
		{Position{1, 6}, TokenEqual, "="},
		{Position{1, 8}, TokenString, "The quick brown fox jumps over the lazy dog."},
		{Position{6, 21}, TokenEOF, ""},
	})

	testLex(t, "key2 = \"\"\"\\\n       The quick brown \\\n       fox jumps over \\\n       the lazy dog.\\\n       \"\"\"", []Token{
		{Position{1, 1}, TokenKey, "key2"},
		{Position{1, 6}, TokenEqual, "="},
		{Position{1, 8}, TokenString, "The quick brown fox jumps over the lazy dog."},
		{Position{5, 11}, TokenEOF, ""},
	})

	testLex(t, `key2 = "Roses are red\nViolets are blue"`, []Token{
		{Position{1, 1}, TokenKey, "key2"},
		{Position{1, 6}, TokenEqual, "="},
		{Position{1, 8}, TokenString, "Roses are red\nViolets are blue"},
		{Position{1, 41}, TokenEOF, ""},
	})

	testLex(t, "key2 = \"\"\"\nRoses are red\nViolets are blue\"\"\"", []Token{
		{Position{1, 1}, TokenKey, "key2"},
		{Position{1, 6}, TokenEqual, "="},
		{Position{1, 8}, TokenString, "Roses are red\nViolets are blue"},
		{Position{3, 20}, TokenEOF, ""},
	})
}

func TestMultilineLiteralString(t *testing.T) {
	testLex(t, `foo = '''hello 'literal' world'''`, []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenString, `hello 'literal' world`},
		{Position{1, 34}, TokenEOF, ""},
	})

	testLex(t, "foo = '''\nhello\n'literal'\nworld'''", []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenString, "hello\n'literal'\nworld"},
		{Position{4, 9}, TokenEOF, ""},
	})
	testLex(t, "foo = '''\nhello\n'literal'\nworld'''", []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenString, "hello\n'literal'\nworld"},
		{Position{4, 9}, TokenEOF, ""},
	})
}

func TestQuotedKey(t *testing.T) {
	testLex(t, "\"a b\" = 42", []Token{
		{Position{1, 1}, TokenKey, "a b"},
		{Position{1, 7}, TokenEqual, "="},
		{Position{1, 9}, TokenInteger, "42"},
		{Position{1, 11}, TokenEOF, ""},
	})
}

func TestSimpleBrace(t *testing.T) {
	testLex(t, "map = { x = 1 }", []Token{
		{Position{1, 1}, TokenKey, "map"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenLeftBrace, "{"},
		{Position{1, 9}, TokenKey, "x"},
		{Position{1, 11}, TokenEqual, "="},
		{Position{1, 13}, TokenInteger, "1"},
		{Position{1, 15}, TokenRightBrace, "}"},
		{Position{1, 16}, TokenEOF, ""},
	})
}

//func TestDigitalKeyMap(t *testing.T) {
//	testLex(t, `123 = { 5 = 6 }`, []Token{
//		{Position{1, 1}, TokenKey, "123"},
//		{Position{1, 5}, TokenEqual, "="},
//		{Position{1, 7}, TokenLeftBrace, "{"},
//		{Position{1, 9}, TokenKey, "5"},
//		{Position{1, 11}, TokenEqual, "="},
//		{Position{1, 13}, TokenInteger, "6"},
//		{Position{1, 15}, TokenRightBrace, "}"},
//		{Position{1, 16}, TokenEOF, ""},
//	})
//}

func TestStringKeyMap(t *testing.T) {
	testLex(t, "a = {\n  \"b/b\" = \"c\"\n}", []Token{
		{Position{1, 1}, TokenKey, "a"},
		{Position{1, 3}, TokenEqual, "="},
		{Position{1, 5}, TokenLeftBrace, "{"},
		{Position{2, 3}, TokenString, "b/b"},
		{Position{2, 9}, TokenEqual, "="},
		{Position{2, 11}, TokenString, "c"},
		{Position{3, 1}, TokenRightBrace, "}"},
		{Position{3, 2}, TokenEOF, ""},
	})
}

func TestLexUnknownValue(t *testing.T) {
	testLex(t, `a = !b`, []Token{
		{Position{1, 1}, TokenKey, "a"},
		{Position{1, 3}, TokenEqual, "="},
		{Position{1, 5}, TokenError, "lexer: no value can start with !"},
	})

	testLex(t, `a = \b`, []Token{
		{Position{1, 1}, TokenKey, "a"},
		{Position{1, 3}, TokenEqual, "="},
		{Position{1, 5}, TokenError, "lexer: no value can start with \\"},
	})
}

func TestUnicodeString(t *testing.T) {
	testLex(t, `foo = "hello ♥ world"`, []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenString, "hello ♥ world"},
		{Position{1, 22}, TokenEOF, ""},
	})
}

func TestEscapeInString(t *testing.T) {
	testLex(t, `foo = "\b\f\/"`, []Token{
		{Position{1, 1}, TokenKey, "foo"},
		{Position{1, 5}, TokenEqual, "="},
		{Position{1, 7}, TokenString, "\b\f/"},
		{Position{1, 15}, TokenEOF, ""},
	})
}

func TestListInMap(t *testing.T) {
	str := `docker = {
  login = false
  images = [
    "astrocorp/signal_api:latest"
  ]
}
`
	testLex(t, str, []Token{
		{Position{1, 1}, TokenKey, "docker"},
		{Position{1, 8}, TokenEqual, "="},
		{Position{1, 10}, TokenLeftBrace, "{"},
		{Position{2, 3}, TokenKey, "login"},
		{Position{2, 9}, TokenEqual, "="},
		{Position{2, 11}, TokenBoolean, "false"},
		{Position{3, 3}, TokenKey, "images"},
		{Position{3, 10}, TokenEqual, "="},
		{Position{3, 12}, TokenLeftBracket, "["},
		{Position{4, 5}, TokenString, "astrocorp/signal_api:latest"},
		{Position{5, 3}, TokenRightBracket, "]"},
		{Position{6, 1}, TokenRightBrace, "}"},
		{Position{7, 1}, TokenEOF, ""},
	})
}
