package main

import (
	"fmt"
	"github.com/takeokunn/gojson/src/ast"
	"github.com/takeokunn/gojson/src/lexer"
	"github.com/takeokunn/gojson/src/parser"
)

type Client struct {
	input       []rune
	tree        *ast.RootNode
}

func NewFromString(jsonStr string) (*Client, error) {
	l := lexer.New(jsonStr)
	p := parser.New(l)
	tree, err := p.ParseJson()
	if err != nil {
		return nil, err
	}
	return &Client{tree: &tree, input: l.Input}, nil
}


func main() {
	var exampleJSON = `{ "string": "a neat string", "bool": true, "PI": 3.14159 }`

	c, err := NewFromString(exampleJSON)
	if err != nil {
		fmt.Printf("\nError creating client: %v\n", err)
	}
	fmt.Println(string(c.input))
}
