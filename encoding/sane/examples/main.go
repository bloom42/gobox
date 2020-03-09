package main

import (
	"fmt"

	"gitlab.com/bloom42/lily/encoding/sane"
	"gitlab.com/bloom42/lily/encoding/sane/lexer"
	"gitlab.com/bloom42/lily/encoding/sane/parser"
)

type D struct {
	A string
}

type C struct {
	A int64 `sane:"a"`
	D []D   `sane:"d"`
}

type S struct {
	A     string          `sane:"a"`
	B     []int64         `sane:"b"`
	C     C               `sane:"c"`
	Z     Z               `sane:"z"`
	Rules map[string]Rule `sane:"rules"`
}

type Rule struct {
	Severity  *string                `sane:"severity"`
	Arguments map[string]interface{} `sane:"arguments"`
}

type Z struct {
	A string `sane:"a/a"`
	B string `sane:"b"`
}

func main() {
	str1 := `
a = "a"
b = [1, 2]
c = { a = 1, d = [ { a = "3.3" }, { a = "xxx" } ] }
z = {
  "a/a" = "a",
  "b" = "b"
}

rules = {
  arguments_limits = {
    arguments = { max = 1 }
  }
  blacklist = { arguments = { list = ["a", "b"] } }
}
`

	tokens := lexer.Lex([]byte(str1))
	fmt.Printf("%#v\n\n", tokens)
	tree, err := parser.Parse([]byte(str1))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n\n", tree.Values["c"].(*parser.Value).Val)

	var s S

	err = sane.Unmarshal([]byte(str1), &s)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n\n", s)

	b, err := sane.Marshal(s)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(b))
}
