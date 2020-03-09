package parser

import (
	"strings"

	"gitlab.com/bloom42/lily/encoding/sane/lexer"
)

// Tree is the result of the parsing of a SANE file.
type Tree struct {
	Values   map[string]interface{} // string -> *Value, *Tree, []*Tree, comments
	Position lexer.Position
}

// Value represents a SANE value
type Value struct {
	Val      interface{} // string, int64, uint64, float64, bool or [] of any
	Position lexer.Position
}

// NewTree initialize a new Tree
func NewTree() *Tree {
	return &Tree{
		Values:   make(map[string]interface{}),
		Position: lexer.Position{Line: 1, Col: 1},
	}
}

// Set an element in the tree.
// Key is a dot-separated path (e.g. a.b.c).
// Creates all necessary intermediate trees, if needed.
func (t *Tree) Set(key string, value interface{}) {
	t.SetPath(strings.Split(key, "."), value)
}

func (t *Tree) SetPath(keys []string, value interface{}) {
	subtree := t
	for _, intermediateKey := range keys[:len(keys)-1] {
		nextTree, exists := subtree.Values[intermediateKey]
		if !exists {
			nextTree = NewTree()
			subtree.Values[intermediateKey] = nextTree // add new element here
		}
		switch node := nextTree.(type) {
		case *Tree:
			subtree = node
		case []*Tree:
			// go to most recent element
			if len(node) == 0 {
				// create element if it does not exist
				subtree.Values[intermediateKey] = append(node, NewTree())
			}
			subtree = node[len(node)-1]
		}
	}

	var toInsert interface{}

	switch value.(type) {
	case *Tree:
		//tt := value.(*Tree)
		toInsert = value
	case []*Tree:
		toInsert = value
	case *Value:
		tt := value.(*Value)
		toInsert = tt
	default:
		toInsert = &Value{Val: value}
	}

	subtree.Values[keys[len(keys)-1]] = toInsert
}

// Has returns a boolean indicating if the given key exists.
func (t *Tree) Has(key string) bool {
	if key == "" {
		return false
	}
	return t.HasPath(strings.Split(key, "."))
}

// HasPath returns true if the given path of keys exists, false otherwise.
func (t *Tree) HasPath(keys []string) bool {
	return t.GetPath(keys) != nil
}

// Get the value at key in the Tree.
// Key is a dot-separated path (e.g. a.b.c) without single/double quoted strings.
// If you need to retrieve non-bare keys, use GetPath.
// Returns nil if the path does not exist in the tree.
// If keys is of length zero, the current tree is returned.
func (t *Tree) Get(key string) interface{} {
	if key == "" {
		return t
	}
	return t.GetPath(strings.Split(key, "."))
}

// GetPath returns the element in the tree indicated by 'keys'.
// If keys is of length zero, the current tree is returned.
func (t *Tree) GetPath(keys []string) interface{} {
	if len(keys) == 0 {
		return t
	}
	subtree := t
	for _, intermediateKey := range keys[:len(keys)-1] {
		value, exists := subtree.Values[intermediateKey]
		if !exists {
			return nil
		}
		switch node := value.(type) {
		case *Tree:
			subtree = node
		case []*Tree:
			// go to most recent element
			if len(node) == 0 {
				return nil
			}
			subtree = node[len(node)-1]
		default:
			return nil // cannot navigate through other node types
		}
	}
	// branch based on final node type
	switch node := subtree.Values[keys[len(keys)-1]].(type) {
	case *Value:
		return node.Val
	default:
		return node
	}
}

// Keys returns the keys of the toplevel tree (does not recurse).
func (t *Tree) Keys() []string {
	keys := make([]string, len(t.Values))
	i := 0
	for k := range t.Values {
		keys[i] = k
		i++
	}
	return keys
}

// SetWithOptions is the same as Set, but allows you to provide formatting
// instructions to the key, that will be used by Marshal().
func (t *Tree) SetWithOptions(key string, value interface{}, multiline bool) {
	t.SetPathWithOptions(strings.Split(key, "."), value, multiline)
}

// SetPathWithOptions is the same as SetPath, but allows you to provide
// formatting instructions to the key, that will be reused by Marshal().
func (t *Tree) SetPathWithOptions(keys []string, value interface{}, multiline bool) {
	subtree := t
	for _, intermediateKey := range keys[:len(keys)-1] {
		nextTree, exists := subtree.Values[intermediateKey]
		if !exists {
			nextTree = NewTree()
			subtree.Values[intermediateKey] = nextTree // add new element here
		}
		switch node := nextTree.(type) {
		case *Tree:
			subtree = node
		case []*Tree:
			// go to most recent element
			if len(node) == 0 {
				// create element if it does not exist
				subtree.Values[intermediateKey] = append(node, NewTree())
			}
			subtree = node[len(node)-1]
		}
	}

	var toInsert interface{}

	switch value.(type) {
	case *Tree:
		toInsert = value
	case []*Tree:
		toInsert = value
	case *Value:
		tt := value.(*Value)
		toInsert = tt
	default:
		toInsert = &Value{Val: value}
	}

	subtree.Values[keys[len(keys)-1]] = toInsert
}
