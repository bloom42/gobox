package lexer

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// StateFn represents the state of the lexer as a function that returns the next state.
type StateFn func() StateFn

type validRuneFn func(r rune) bool

// Lexer is a SANE lexer
type Lexer struct {
	input             []rune
	pos               int64
	currentTokenStart int64
	currentTokenStop  int64
	tokens            chan Token
	line              int64
	col               int64
	depth             int64
	endbufferLine     int64
	endbufferCol      int64
}

// Lex returns all the lexed tokens for the given input
func Lex(input []byte) []Token {
	ret := []Token{}
	lexer := NewLexer(input)

	for token := range lexer.Next() {
		ret = append(ret, token)
	}

	return ret
}

// NewLexer returns a new Lexer initialized with input
func NewLexer(input []byte) *Lexer {
	lx := &Lexer{
		input:         bytes.Runes(input),
		line:          1,
		col:           1,
		endbufferLine: 1,
		endbufferCol:  1,
		tokens:        make(chan Token),
	}

	go lx.run()
	return lx
}

// Next returns the next lexed token
func (lx *Lexer) Next() <-chan Token {
	return lx.tokens
}

// real lexer entry point
func (lx *Lexer) run() {
	for state := lx.lexText; state != nil; {
		state = state()
	}
	close(lx.tokens)
}

// State functions
func (lx *Lexer) lexText() StateFn {
	for {
		next := lx.peek()
		switch next {
		case Hash:
			return lx.lexComment(lx.lexRvalue)
		case Equal:
			return lx.lexEqual
		case '\n':
			lx.skip()
			continue
		}

		if isWhitespace(next) {
			lx.skip()
		}

		if lx.depth > 0 {
			return lx.lexRvalue
		}

		if isKeyStartChar(next) {
			return lx.lexKey
		}

		if next == EOF {
			lx.next()
			break
		}
	}

	lx.emit(TokenEOF)
	return nil
}

func (lx *Lexer) emitWithValue(t TokenType, value string) {
	lx.tokens <- Token{
		Position: Position{lx.line, lx.col},
		Type:     t,
		Value:    value,
	}
	lx.ignore()
}

func (lx *Lexer) emit(t TokenType) {
	lx.emitWithValue(t, string(lx.input[lx.currentTokenStart:lx.currentTokenStop]))
}

func (lx *Lexer) ignore() {
	lx.currentTokenStart = lx.currentTokenStop
	lx.line = lx.endbufferLine
	lx.col = lx.endbufferCol
}

func (lx *Lexer) skip() {
	lx.next()
	lx.ignore()
}

func (lx *Lexer) next() rune {
	r := lx.read()

	if r != EOF {
		lx.currentTokenStop += 1
	}
	return r
}

func (lx *Lexer) read() rune {
	r := lx.peek()
	if r == '\n' {
		lx.endbufferLine += 1
		lx.endbufferCol = 1
	} else {
		lx.endbufferCol++
	}

	lx.pos += 1
	return r
}

func (lx *Lexer) peek() rune {
	if lx.pos >= int64(len(lx.input)) {
		return EOF
	}
	return lx.input[lx.pos]
}

func (lx *Lexer) lexComment(previousState StateFn) StateFn {
	var commentValue bytes.Buffer

	for next := lx.peek(); next != '\n' && next != EOF; next = lx.peek() {
		lx.next()
		commentValue.WriteRune(next)
	}
	strCommentValue := commentValue.String()
	lx.emitWithValue(TokenComment, strCommentValue)

	lx.ignore()
	return previousState
}

func (lx *Lexer) follow(next string) bool {
	return next == lx.peekString(len(next))
}

func (lx *Lexer) peekString(size int) string {
	maxIdx := int64(len(lx.input))
	upperIdx := lx.pos + int64(size)
	if upperIdx > maxIdx {
		upperIdx = maxIdx
	}
	return string(lx.input[lx.pos:upperIdx])
}

func (lx *Lexer) lexEqual() StateFn {
	lx.next()
	lx.emit(TokenEqual)
	return lx.lexRvalue
}

func (lx *Lexer) lexRvalue() StateFn {
	for {
		next := lx.peek()
		switch next {
		case Dot:
			return lx.errorf("lexer: cannot start Float with a dot")
		case Equal:
			return lx.lexEqual
		case LeftBrace:
			lx.depth += 1
			return lx.lexLeftBrace
		case RightBrace:
			lx.depth -= 1
			return lx.lexRightBrace
		case LeftBracket:
			lx.depth += 1
			return lx.lexLeftBracket
		case RightBracket:
			lx.depth -= 1
			return lx.lexRightBracket
		case Hash:
			return lx.lexComment(lx.lexRvalue)
		case DoubleQuote:
			return lx.lexString
		case SingleQuote:
			return lx.lexLiteralString
		case Comma:
			return lx.lexComma
		// case '\r':
		// 	fallthrough
		case '\n':
			lx.skip()
			if lx.depth == 0 {
				return lx.lexText
			}
			return lx.lexRvalue
		case '_':
			return lx.errorf("lexer: cannot start number with underscore")
		}

		if lx.follow("true") || lx.follow("false") {
			return lx.lexBoolean
		}

		if lx.follow("inf") {
			return lx.lexInf
		}

		if lx.follow("nan") {
			return lx.lexNan
		}

		if isWhitespace(next) {
			lx.skip()
			continue
		}

		if next == EOF {
			lx.next()
			break
		}

		//possibleDate := lx.peekString(35)
		//dateMatch := dateRegexp.FindString(possibleDate)
		//if dateMatch != "" {
		//	lx.fastForward(len(dateMatch))
		//	return lx.lexDate
		//}

		if next == '+' || next == '-' || isDigit(next) {
			return lx.lexNumber
		}

		if isAlphanumeric(next) {
			return lx.lexKey
		}

		return lx.errorf("lexer: no value can start with %c", next)
	}

	lx.emit(TokenEOF)
	return nil
}

func (lx *Lexer) lexComma() StateFn {
	lx.next()
	lx.emit(TokenComma)
	return lx.lexRvalue
}

func (lx *Lexer) lexBoolean() StateFn {
	if lx.follow("true") {
		lx.fastForward(4)
		lx.emitWithValue(TokenBoolean, "true")
	} else {
		lx.fastForward(5)
		lx.emitWithValue(TokenBoolean, "false")
	}
	return lx.lexRvalue
}

func (lx *Lexer) fastForward(n int) {
	for i := 0; i < n; i++ {
		lx.next()
	}
}

func (lx *Lexer) lexLeftBracket() StateFn {
	lx.next()
	lx.emit(TokenLeftBracket)
	return lx.lexRvalue
}

func (lx *Lexer) lexRightBracket() StateFn {
	lx.next()
	lx.emit(TokenRightBracket)
	return lx.lexRvalue
}

func (lx *Lexer) lexLeftBrace() StateFn {
	lx.next()
	lx.emit(TokenLeftBrace)
	return lx.lexRvalue
}

func (lx *Lexer) lexRightBrace() StateFn {
	lx.next()
	lx.emit(TokenRightBrace)
	return lx.lexRvalue
}

func isValidOctalRune(r rune) bool {
	return r >= '0' && r <= '7' || r == '_'
}

func isValidBinaryRune(r rune) bool {
	return r == '0' || r == '1' || r == '_'
}

func isValidHexRune(r rune) bool {
	return r >= 'a' && r <= 'f' ||
		r >= 'A' && r <= 'F' ||
		r >= '0' && r <= '9' ||
		r == '_'
}

func (lx *Lexer) lexNumber() StateFn {
	r := lx.peek()

	if r == '0' {
		follow := lx.peekString(2)
		if len(follow) == 2 {
			var isValidRune validRuneFn
			switch follow[1] {
			case 'x':
				isValidRune = isValidHexRune
			case 'o':
				isValidRune = isValidOctalRune
			case 'b':
				isValidRune = isValidBinaryRune
			default:
				if follow[1] >= 'a' && follow[1] <= 'z' || follow[1] >= 'A' && follow[1] <= 'Z' {
					return lx.errorf("lexer unknown number base: %s. possible options are x (hex) o (octal) b (binary)", string(follow[1]))
				}
			}

			if isValidRune != nil {
				lx.next()
				lx.next()
				digitSeen := false
				for {
					next := lx.peek()
					if !isValidRune(next) {
						break
					}
					digitSeen = true
					lx.next()
				}

				if !digitSeen {
					return lx.errorf("lexer: number needs at least one digit")
				}

				lx.emit(TokenInteger)

				return lx.lexRvalue
			}
		}
	}

	if r == '+' || r == '-' {
		lx.next()
		if lx.follow("inf") {
			return lx.lexInf
		}
		if lx.follow("nan") {
			return lx.lexNan
		}
	}

	pointSeen := false
	expSeen := false
	digitSeen := false
	for {
		next := lx.peek()
		if next == '.' {
			if pointSeen {
				return lx.errorf("lexer: cannot have two dots in one float")
			}
			lx.next()
			if !isDigit(lx.peek()) {
				return lx.errorf("lexer: float cannot end with a dot")
			}
			pointSeen = true
		} else if next == 'e' || next == 'E' {
			expSeen = true
			lx.next()
			r := lx.peek()
			if r == '+' || r == '-' {
				lx.next()
			}
		} else if isDigit(next) {
			digitSeen = true
			lx.next()
		} else if next == '_' {
			lx.next()
		} else {
			break
		}
		if pointSeen && !digitSeen {
			return lx.errorf("lexer: cannot start float with a dot")
		}
	}

	if !digitSeen {
		return lx.errorf("lexer: no digit in that number")
	}
	if pointSeen || expSeen {
		lx.emit(TokenFloat)
	} else {
		lx.emit(TokenInteger)
	}
	return lx.lexRvalue
}

func (lx *Lexer) lexInf() StateFn {
	lx.fastForward(3)
	lx.emitWithValue(TokenFloat, "inf")
	return lx.lexRvalue
}

func (lx *Lexer) lexNan() StateFn {
	lx.fastForward(3)
	lx.emitWithValue(TokenFloat, "nan")
	return lx.lexRvalue
}

func (lx *Lexer) lexString() StateFn {
	lineOrig := lx.line
	colOrig := lx.col
	lx.skip()

	// handle special case for triple-quote
	terminator := `"`
	discardLeadingNewLine := false
	acceptNewLines := false
	if lx.follow(`""`) {
		lx.skip()
		lx.skip()
		terminator = `"""`
		discardLeadingNewLine = true
		acceptNewLines = true
	}

	str, err := lx.lexStringAsString(terminator, discardLeadingNewLine, acceptNewLines)

	if err != nil {
		lx.setLineCol(lineOrig, colOrig)
		return lx.errorf(err.Error())
	}

	// ugly, but to put the pos the starting quote
	colAfter := lx.col
	lineAfter := lx.line
	lx.setLineCol(lineOrig, colOrig)
	lx.emitWithValue(TokenString, str)
	lx.setLineCol(lineAfter, colAfter)

	lx.fastForward(len(terminator))
	lx.ignore()
	return lx.lexRvalue
}

// Parse the key and emits its value without escape sequences.
// bare keys, basic string keys and literal string keys are supported.
func (lx *Lexer) lexKey() StateFn {
	growingString := ""

	for r := lx.peek(); isKeyChar(r) || r == '\n'; r = lx.peek() {
		if r == '"' {
			lx.next()
			str, err := lx.lexStringAsString(`"`, false, true)
			if err != nil {
				return lx.errorf(fmt.Sprintf("lexer: %s", err))
			}
			growingString += str
			lx.next()
			continue
		} else if r == '\'' {
			lx.next()
			str, err := lx.lexLiteralStringAsString(`'`, false)
			if err != nil {
				return lx.errorf(err.Error())
			}
			growingString += str
			lx.next()
			continue
		} else if r == '\n' {
			return lx.errorf("lexer: keys cannot contain new lines")
		} else if isWhitespace(r) {
			break
		} else if !isValidBareChar(r) {
			return lx.errorf("lexer: keys cannot contain %c character", r)
		}
		growingString += string(r)
		lx.next()
	}

	if growingString == "" {
		return lx.errorf("lexer: keys cannot be emtpy")
	}

	lx.emitWithValue(TokenKey, growingString)
	return lx.lexText
}

// Lex a string and return the results as a string.
// Terminator is the substring indicating the end of the token.
// The resulting string does not include the terminator.
func (lx *Lexer) lexStringAsString(terminator string, discardLeadingNewLine, acceptNewLines bool) (string, error) {
	growingString := ""

	if discardLeadingNewLine {
		if lx.peek() == '\n' {
			lx.skip()
		}
	}

	for {
		if lx.follow(terminator) {
			return growingString, nil
		}

		if lx.follow("\\") {
			lx.next()
			switch lx.peek() {
			// case '\r':
			// 	fallthrough
			case '\n':
				fallthrough
			case '\t':
				fallthrough
			case ' ':
				// skip all whitespace chars following backslash
				for strings.ContainsRune("\n ", lx.peek()) {
					lx.next()
				}
			case '"':
				growingString += "\""
				lx.next()
			case 'n':
				growingString += "\n"
				lx.next()
			case 'b':
				growingString += "\b"
				lx.next()
			case 'f':
				growingString += "\f"
				lx.next()
			case '/':
				growingString += "/"
				lx.next()
			case 't':
				growingString += "\t"
				lx.next()
			case 'r':
				growingString += "\r"
				lx.next()
			case '\\':
				growingString += "\\"
				lx.next()
			case 'u':
				lx.next()
				code := ""
				for i := 0; i < 4; i++ {
					c := lx.peek()
					if !isHexDigit(c) {
						return "", errors.New("lexer: unfinished unicode escape")
					}
					lx.next()
					code = code + string(c)
				}
				intcode, err := strconv.ParseInt(code, 16, 32)
				if err != nil {
					return "", errors.New("lexer: invalid unicode escape: \\u" + code)
				}
				growingString += string(rune(intcode))
			case 'U':
				lx.next()
				code := ""
				for i := 0; i < 8; i++ {
					c := lx.peek()
					if !isHexDigit(c) {
						return "", errors.New("lexer: unfinished unicode escape")
					}
					lx.next()
					code = code + string(c)
				}
				intcode, err := strconv.ParseInt(code, 16, 64)
				if err != nil {
					return "", errors.New("lexer: invalid unicode escape: \\U" + code)
				}
				growingString += string(rune(intcode))
			default:
				return "", errors.New("lexer: invalid escape sequence: \\" + string(lx.peek()))
			}
		} else {
			r := lx.peek()

			if 0x00 <= r && r <= 0x1F && !(acceptNewLines && (r == '\n')) {
				return "", fmt.Errorf("lexer: unescaped control character %U", r)
			}
			lx.next()
			growingString += string(r)
		}

		if lx.peek() == EOF {
			break
		}
	}

	return "", errors.New("lexer: unclosed string")
}

func (lx *Lexer) lexLiteralStringAsString(terminator string, discardLeadingNewLine bool) (string, error) {
	growingString := ""

	if discardLeadingNewLine {
		if lx.peek() == '\n' {
			lx.skip()
		}
	}

	// find end of string
	for {
		if lx.follow(terminator) {
			return growingString, nil
		}

		next := lx.peek()
		if next == EOF {
			break
		}
		growingString += string(lx.next())
	}

	return "", errors.New("lexer: unclosed string")
}

func (lx *Lexer) lexLiteralString() StateFn {
	colOrig := lx.col
	lineOrig := lx.line
	lx.skip()

	// handle special case for triple-quote
	terminator := "'"
	discardLeadingNewLine := false
	if lx.follow("''") {
		lx.skip()
		lx.skip()
		terminator = "'''"
		discardLeadingNewLine = true
	}

	str, err := lx.lexLiteralStringAsString(terminator, discardLeadingNewLine)
	if err != nil {
		lx.setLineCol(lineOrig, colOrig)
		return lx.errorf(err.Error())
	}

	// ugly, but to put the pos the starting quote
	colAfter := lx.col
	lineAfter := lx.line
	lx.setLineCol(lineOrig, colOrig)
	lx.emitWithValue(TokenString, str)
	lx.setLineCol(lineAfter, colAfter)

	lx.fastForward(len(terminator))
	lx.ignore()
	return lx.lexRvalue
}

// quick function to set good offset for strings
func (lx *Lexer) setLineCol(line, col int64) {
	lx.line = line
	lx.col = col
}

// errorf stops all lexing by emitting an error and returning `nil`.
// Note that any value that is a character is escaped if it's a special
// character (newlines, tabs, etc.).
func (lx *Lexer) errorf(format string, values ...interface{}) StateFn {
	lx.tokens <- Token{
		Position{Line: lx.line, Col: lx.col},
		TokenError,
		fmt.Sprintf(format, values...),
	}
	return nil
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isHexDigit(r rune) bool {
	return isDigit(r) ||
		(r >= 'a' && r <= 'f') ||
		(r >= 'A' && r <= 'F')
}

func isHexadecimal(r rune) bool {
	return (r >= '0' && r <= '9') ||
		(r >= 'a' && r <= 'f') ||
		(r >= 'A' && r <= 'F')
}

func isBareKeyChar(r rune) bool {
	return (r >= 'A' && r <= 'Z') ||
		(r >= 'a' && r <= 'z') ||
		(r >= '0' && r <= '9') ||
		r == '_' ||
		r == '-'
}

// isWhitespace returns true if `r` is a whitespace character according to the spec.
func isWhitespace(r rune) bool {
	return r == ' '
}

// isWhitespace returns true if `r` is a newline character according to the spec.
func isNL(r rune) bool {
	return r == '\n'
}

func isAlphanumeric(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isKeyChar(r rune) bool {
	// Keys start with the first character that isn't whitespace or [ and end
	// with the last non-whitespace character before the Equal sign. Keys
	// cannot contain a # character."
	return !(r == '\n' || r == EOF || r == '=')
}

func isKeyStartChar(r rune) bool {
	return !(isWhitespace(r) || r == '\n' || r == EOF || r == '[')
}

func isValidBareChar(r rune) bool {
	return isAlphanumeric(r) || r == '-' || unicode.IsNumber(r)
}
