package parser

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"gitlab.com/bloom42/lily/encoding/sane/lexer"
)

func init() {
	numberUnderscoreInvalidRegexp = regexp.MustCompile(`([^\d]_|_[^\d])|_$|^_`)
	hexNumberUnderscoreInvalidRegexp = regexp.MustCompile(`(^0x_)|([^\da-f]_|_[^\da-f])|_$|^_`)
}

var numberUnderscoreInvalidRegexp *regexp.Regexp
var hexNumberUnderscoreInvalidRegexp *regexp.Regexp

// StateFn represents the state of the parser as a function that returns the next state.
type StateFn func() StateFn

type Parser struct {
	tokens      <-chan lexer.Token
	tree        *Tree
	err         error
	currentMap  []string // ordered depth keys of current map
	seenMapKeys []string
}

func Parse(input []byte) (*Tree, error) {
	lexer := lexer.NewLexer(input)
	return ParseTokens(lexer.Next())
}

func ParseTokens(tokens <-chan lexer.Token) (*Tree, error) {
	p := NewParser(tokens)

	p.run()
	return p.tree, p.err
}

func NewParser(tokens <-chan lexer.Token) *Parser {
	p := &Parser{
		tokens: tokens,
		tree:   NewTree(),
	}

	return p
}

func (p *Parser) run() {
	for state := p.parseStart; state != nil; {
		state = state()
	}
}

func (p *Parser) parseStart() StateFn {
	token := p.nextToken()

	// end of stream, parsing is finished
	if token == nil {
		return nil
	}

	switch token.Type {
	case lexer.TokenComment:
		return p.parseComment(*token, p.parseStart)
	case lexer.TokenKey, lexer.TokenString:
		return p.parseAssign(*token)
	case lexer.TokenEOF:
		return nil
	default:
		p.stateTokenError(*token, "unexpected token")
	}
	return nil
}

func (p *Parser) parseAssign(key lexer.Token) StateFn {
	p.expectNext(lexer.TokenEqual)
	if p.err != nil {
		return nil
	}

	token := p.nextToken()
	position, value := p.parseRvalue(token)
	if p.err != nil {
		return nil
	}

	if _, exist := p.tree.Values[key.Value]; exist {
		return p.stateTokenError(key, fmt.Sprintf("key '%s' is already defined", key.Value))
	}

	p.tree.Values[key.Value] = &Value{value, position}

	return p.parseStart
}

func (p *Parser) parseRvalue(tok *lexer.Token) (lexer.Position, interface{}) {
	if p.err != nil {
		return lexer.Position{}, nil
	}
	if tok == nil {
		p.stateError("value expected but got nil.")
		return lexer.Position{}, nil
	}
	if tok.Type == lexer.TokenEOF {
		p.stateTokenError(*tok, fmt.Sprintf("value expected but got '%s'.", *tok))
		return lexer.Position{}, nil
	}

	switch tok.Type {
	case lexer.TokenString:
		return tok.Position, tok.Value
	case lexer.TokenBoolean:
		return tok.Position, p.parseBoolean(*tok)
	case lexer.TokenInteger:
		return tok.Position, p.parseInteger(*tok)
	case lexer.TokenFloat:
		return tok.Position, p.parseFloat(*tok)
	case lexer.TokenLeftBrace:
		return tok.Position, p.parseMap(tok.Position)
	case lexer.TokenLeftBracket:
		return tok.Position, p.parseList()
	default:
		p.stateTokenError(*tok, "unexpected token")
		return tok.Position, nil
	}
}

func (p *Parser) parseComment(comment lexer.Token, previousState StateFn) StateFn {
	val := comment.Value
	value := Value{
		Val:      val,
		Position: comment.Position,
	}

	if p.tree.Values[""] == nil {
		p.tree.Values[""] = []interface{}{}
	}

	p.tree.Values[""] = append(p.tree.Values[""].([]interface{}), &value)
	return previousState
}

func (p *Parser) parseBoolean(token lexer.Token) bool {
	if token.Value == "true" {
		return true
	} else if token.Value == "false" {
		return false
	}
	p.stateTokenError(token, "invalid Boolean Token")
	return false
}

func (p *Parser) parseInteger(token lexer.Token) int64 {
	cleanedVal := cleanupNumberToken(token.Value)
	var err error
	var val int64

	if len(cleanedVal) >= 3 && cleanedVal[0] == '0' {
		switch cleanedVal[1] {
		case 'x':
			err = hexNumberContainsInvalidUnderscore(token.Value)
			if err != nil {
				p.stateTokenError(token, err.Error())
				return val
			}
			val, err = strconv.ParseInt(cleanedVal[2:], 16, 64)
		case 'o':
			err = numberContainsInvalidUnderscore(token.Value)
			if err != nil {
				p.stateTokenError(token, err.Error())
				return val
			}
			val, err = strconv.ParseInt(cleanedVal[2:], 8, 64)
		case 'b':
			err = numberContainsInvalidUnderscore(token.Value)
			if err != nil {
				p.stateTokenError(token, err.Error())
				return val
			}
			val, err = strconv.ParseInt(cleanedVal[2:], 2, 64)
		default:
			p.stateTokenError(token, "invalid base")
			return val
		}
	} else {
		err = numberContainsInvalidUnderscore(token.Value)
		if err != nil {
			p.stateTokenError(token, err.Error())
			return val
		}
		val, err = strconv.ParseInt(cleanedVal, 10, 64)
	}

	if err != nil {
		p.stateTokenError(token, err.Error())
	}
	return val
}

func (p *Parser) parseFloat(token lexer.Token) float64 {
	v := token.Value
	if v == "inf" || v == "+inf" {
		return math.Inf(-1)
	} else if v == "-inf" {
		return math.Inf(1)
	} else if v == "nan" || v == "+nan" || v == "-nan" {
		return math.NaN()
	}

	err := numberContainsInvalidUnderscore(token.Value)
	if err != nil {
		p.stateTokenError(token, err.Error())
		return 0
	}

	cleanedVal := cleanupNumberToken(token.Value)
	val, err := strconv.ParseFloat(cleanedVal, 64)
	if err != nil {
		p.stateTokenError(token, err.Error())
		return 0
	}
	return val
}

func (p *Parser) parseMap(position lexer.Position) *Tree {
	tree := NewTree()
	tree.Position = position
	var previous *lexer.Token

Loop:
	for {
		token := p.nextToken()
		if token == nil {
			p.stateError("unterminated map")
			return nil
		}
		if token.Type == lexer.TokenEOF {
			p.stateTokenError(*token, "unterminated map")
			return nil
		}
		switch token.Type {
		case lexer.TokenComment:
			continue
		case lexer.TokenRightBrace:
			break Loop
		case lexer.TokenKey, lexer.TokenString:
			if previous != nil && tokenIsComma(previous) == false && previous.Line == token.Line {
				p.stateTokenError(*token, "missing comma in inline map")
				return nil
			}
			p.expectNext(lexer.TokenEqual)
			if p.err != nil {
				return nil
			}
			key := token
			token = p.nextToken()
			position, value := p.parseRvalue(token)
			if p.err != nil {
				return nil
			}
			_, _ = position, value
			tree.SetPath([]string{key.Value}, value)
		case lexer.TokenComma:
			if previous == nil {
				p.stateTokenError(*token, "map cannot start with a comma")
				return nil
			}
			if tokenIsComma(previous) {
				p.stateTokenError(*previous, "trailing comma at the end of map")
				return nil
			}
		default:
			p.stateTokenError(*token, fmt.Sprintf("unexpected token in map"))
			return nil
		}
		previous = token
	}
	if tokenIsComma(previous) {
		p.stateTokenError(*previous, "trailing comma at the end of map")
	}
	return tree
}

func (p *Parser) parseList() interface{} {
	var list []interface{}
	listType := reflect.TypeOf(nil)
	for {
		token := p.nextToken()
		if token == nil {
			p.stateError("unterminated list")
			return nil
		}
		if token.Type == lexer.TokenEOF {
			p.stateTokenError(*token, "unterminated list")
			return nil
		}
		if token.Type == lexer.TokenComment {
			continue
		}
		position, val := p.parseRvalue(token)
		_ = position
		if listType == nil {
			listType = reflect.TypeOf(val)
		}
		if reflect.TypeOf(val) != listType {
			p.stateTokenError(*token, "mixed tpyes in list")
			return nil
		}
		list = append(list, val)
		token = p.nextToken()
		if token == nil || token.Type == lexer.TokenEOF {
			p.stateTokenError(*token, "unterminated list")
			return nil
		}
		if token.Type != lexer.TokenRightBracket && token.Type != lexer.TokenComma {
			p.stateTokenError(*token, "missing comma in list")
			return nil
		}
		if token.Type == lexer.TokenRightBracket {
			break
		}
	}
	// A list of Trees is actually an list of inline
	// map, which is a shorthand for a table list. If the
	// list was not converted from []interface{} to []*Tree,
	// the two notations would not be equivalent.
	if listType == reflect.TypeOf(NewTree()) {
		sanList := make([]*Tree, len(list))
		for i, v := range list {
			sanList[i] = v.(*Tree)
		}
		return sanList
	}
	return list
}

func (p *Parser) nextToken() *lexer.Token {
	token, ok := <-p.tokens
	if ok != true {
		return nil
	}

	if token.Type == lexer.TokenError {
		p.stateTokenError(token, token.Value)
		return nil
	}

	return &token
}

func (p *Parser) expectNext(typ lexer.TokenType) {
	it := p.nextToken()
	p.assertEqual(typ, it)
}

func (p *Parser) assertEqual(expected lexer.TokenType, got *lexer.Token) {
	if got == nil {
		if p.err == nil {
			p.stateError(fmt.Sprintf("expected '%s' but got nil.", expected))
		}
	}

	if expected != got.Type {
		p.stateTokenError(*got, fmt.Sprintf("expected '%s' but got '%s'.", expected, got.Type))
	}
}

func (p *Parser) stateTokenError(token lexer.Token, message string) StateFn {
	return p.stateError(fmt.Sprintf("parser: (%d:%d) %s: %s", token.Line, token.Col, token.Type, message))
}

func (p *Parser) stateError(message string) StateFn {
	p.err = errors.New(message)
	return nil
}

func numberContainsInvalidUnderscore(value string) error {
	if numberUnderscoreInvalidRegexp.MatchString(value) {
		return errors.New("invalid use of _ in number")
	}
	return nil
}

func hexNumberContainsInvalidUnderscore(value string) error {
	if hexNumberUnderscoreInvalidRegexp.MatchString(value) {
		return errors.New("invalid use of _ in hex number")
	}
	return nil
}

func cleanupNumberToken(value string) string {
	cleanedVal := strings.Replace(value, "_", "", -1)
	return cleanedVal
}

func tokenIsComma(t *lexer.Token) bool {
	return t != nil && t.Type == lexer.TokenComma
}
