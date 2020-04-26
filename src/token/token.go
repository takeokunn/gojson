package token

import "fmt"

type Type string

const (
	Illegal      Type = "ILLEGAL"
	EOF          Type = "EOF"
	String       Type = "STRING"
	Number       Type = "NUMBER"
	LeftBrace    Type = "{"
	RightBrace   Type = "}"
	LeftBracket  Type = "["
	RightBracket Type = "]"
	Comma        Type = ","
	Colon        Type = ":"
	True         Type = "TRUE"
	False        Type = "FALSE"
	Null         Type = "NULL"
)

type Token struct {
	Type    Type
	Literal string
	Line    int
	Start   int
	End     int
}

var validJsonIdentifiers = map[string]Type{
	"true":  True,
	"false": False,
	"null":  Null,
}

func LookupIdentfier(identifier string) (Type, error) {
	if token, ok := validJsonIdentifiers[identifier]; ok {
		return token, nil
	}
	return "", fmt.Errorf("Expected a valid JSON identifier. Found: %s", identifier)
}

var escapes = map[rune]int{
	'"':  0,
	'\\': 1,
	'/':  2,
	'b':  3,
	'f':  4,
	'n':  5,
	'r':  6,
	't':  7,
	'u':  8,
}

var escapeChars = map[string]string{
	"b": "\b",
	"f": "\f",
	"n": "\n",
	"r": "\r",
	"t": "\t",
}
