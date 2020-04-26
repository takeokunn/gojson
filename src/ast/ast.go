package ast

type RootNodeType int
type Value interface{}
type state int

type Identifier struct {
	Type  string
	Value string
}

type Literal struct {
	Type  string
	Value Value
}

type Property struct {
	Type  string
	Key   Identifier
	Value Value
}

type Object struct {
	Type     string
	Children []Property
	Start    int
	End      int
}

type Array struct {
	Type     string
	Children []Value
	Start    int
	End      int
}

type RootNode struct {
	RootValue *Value
	Type      RootNodeType
}

const (
	ObjectRoot RootNodeType = iota
	ArrayRoot
)

const (
	// object
	ObjStart state = iota
	ObjOpen
	ObjProperty
	ObjComma

	// property
	PropertyStart
	PropertyKey
	PropertyColon

	// array
	ArrayStart
	ArrayOpen
	ArrayValue
	ArrayComma

	// string
	StringStart
	StringQuoteOrChar
	Escape

	// number
	NumberStart
	NumberMinus
	NumberZero
	NumberDigit
	NumberPoint
	NumberDigitFraction
	NumberExp
	NumberExpDigitOrSign
)
