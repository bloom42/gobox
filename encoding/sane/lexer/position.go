package lexer

import (
	"fmt"
)

// Position of a document element within a SANE document.
//
// Line and Col are both 1-indexed positions for the element's line number and
// column number, respectively.  Values of zero or less will cause Invalid(),
// to return true.
type Position struct {
	Line int64 // line within the document
	Col  int64 // column within the line
}

// String representation of the position.
// Displays 1-indexed line and column numbers.
func (p Position) String() string {
	return fmt.Sprintf("(%d, %d)", p.Line, p.Col)
}

// Valid returns whether or not the position is valid (i.e. with negative or null values)
func (p Position) Valid() bool {
	return p.Line > 0 && p.Col > 0
}
