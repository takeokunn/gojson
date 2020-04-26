package parser

import (
	"errors"
	"strconv"
	"strings"

	"fmt"

	"github.com/takeokunn/gojson/src/ast"
	"github.com/takeokunn/gojson/src/lexer"
	"github.com/takeokunn/gojson/src/token"
)

type Parser struct {
	lexer        *lexer.Lexer
	errors       []string
	currentToken token.Token
	peekToken    token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) ParseJson() (ast.RootNode, error) {
	var rootNode ast.RootNode
	if p.currentTokenTypeIs(token.LeftBracket) {
		rootNode.Type = ast.ArrayRoot
	}

	val := p.parserValue()
	if val == nil {
		p.parseError(fmt.Sprintf(
			"Error parsing JSON expected a value, got: %v:",
			p.currentToken.Literal,
		))
		return ast.RootNode{}, errors.New(p.Errors())
	}

	rootNode.RootValue = &val
	return rootNode, nil
}

func (p *Parser) currentTokenTypeIs(t token.Type) bool {
	return p.currentToken.Type == t
}

func (p *Parser) parserValue() ast.Value {
	switch p.currentToken.Type {
	case token.LeftBrace:
		return p.parseJsonObject()
	case token.LeftBracket:
		return p.parseArrayObject()
	default:
		return p.parseJsonLiteral()
	}
}

func (p *Parser) parseJsonObject() ast.Value {
	obj := ast.Object{Type: "Object"}
	objState := ast.ObjStart

	for !p.currentTokenTypeIs(token.EOF) {
		switch objState {
		case ast.ObjStart:
			if p.currentTokenTypeIs(token.LeftBrace) {
				objState = ast.ObjOpen
				obj.Start = p.currentToken.Start
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf(
					"Error parsing JSON object Expected `{` token, got: %s",
					p.currentToken.Literal,
				))
				return nil
			}
		case ast.ObjOpen:
			if p.currentTokenTypeIs(token.RightBrace) {
				p.nextToken()
				obj.End = p.currentToken.End
				return obj
			}
			prop := p.parseProperty()
			obj.Children = append(obj.Children, prop)
			objState = ast.ObjProperty
		case ast.ObjProperty:
			if p.currentTokenTypeIs(token.RightBrace) {
				p.nextToken()
				obj.End = p.currentToken.Start
				return obj
			} else if p.currentTokenTypeIs(token.Comma) {
				objState = ast.ObjComma
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf(
					"Error parsing property. Expected RightBrace or Comma token, got: %s",
					p.currentToken.Literal,
				))
				return nil
			}
		case ast.ObjComma:
			prop := p.parseProperty()
			if prop.Value != nil {
				obj.Children = append(obj.Children, prop)
				objState = ast.ObjProperty
			}
		}
	}

	obj.End = p.currentToken.Start
	return obj
}

func (p *Parser) parseArrayObject() ast.Value {
	array := ast.Array{Type: "Array"}
	arrayState := ast.ArrayStart

	for !p.currentTokenTypeIs(token.EOF) {
		switch arrayState {
		case ast.ArrayStart:
			if p.currentTokenTypeIs(token.LeftBracket) {
				array.Start = p.currentToken.Start
				arrayState = ast.ArrayOpen
				p.nextToken()
			}
		case ast.ArrayOpen:
			if p.currentTokenTypeIs(token.RightBracket) {
				array.End = p.currentToken.End
				p.nextToken()
				return array
			}
			val := p.parserValue()
			array.Children = append(array.Children, val)
			arrayState = ast.ArrayValue
			if p.peekTokenTypeIs(token.RightBracket) {
				p.nextToken()
			}
		case ast.ArrayValue:
			if p.currentTokenTypeIs(token.RightBracket) {
				array.End = p.currentToken.End
				p.nextToken()
				return array
			} else if p.currentTokenTypeIs(token.Comma) {
				arrayState = ast.ArrayComma
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf(
					"Error parsing property. Expected RightBrace or Comma token, got: %s",
					p.currentToken.Literal,
				))
			}
		case ast.ArrayComma:
			val := p.parserValue()
			array.Children = append(array.Children, val)
			arrayState = ast.ArrayValue
		}
	}

	array.End = p.currentToken.Start
	return array
}

func (p *Parser) parseJsonLiteral() ast.Value {
	val := ast.Literal{Type: "Literal"}
	defer p.nextToken()

	switch p.currentToken.Type {
	case token.String:
		val.Value = p.parseString()
	case token.Number:
		ct := p.currentToken.Literal
		f, err := strconv.ParseFloat(ct, 64)
		if err != nil {
			p.parseError("error parsing JSON number, incorrect syntax")
			val.Value = ct
			return val
		}
		val.Value = f
	case token.True:
		val.Value = true
	case token.False:
		val.Value = false
	default:
		val.Value = "null"
	}
	return val
}

func (p *Parser) parseProperty() ast.Property {
	prop := ast.Property{Type: "Property"}
	propertyState := ast.PropertyStart

	for !p.currentTokenTypeIs(token.EOF) {
		switch propertyState {
		case ast.PropertyStart:
			if p.currentTokenTypeIs(token.String) {
				key := ast.Identifier{
					Type: "Identifier",
					Value: p.parseString(),
				}
				prop.Key = key
				propertyState = ast.PropertyKey
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf(
					"Error parsing property start. Expected String token, got: %s",
					p.currentToken.Literal,
				))
			}
		case ast.PropertyKey:
			if p.currentTokenTypeIs(token.Colon) {
				propertyState = ast.PropertyColon
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf(
					"Error parsing property. Expected Colon token, got: %s",
					p.currentToken.Literal,
				))
			}
		case ast.PropertyColon:
			val := p.parserValue()
			prop.Value = val
			return prop
		}
	}

	return prop
}

func (p *Parser) parseString() string {
	return p.currentToken.Literal
}

func (p *Parser) expectPeekType(t token.Type) bool {
	if p.peekTokenTypeIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekTokenTypeIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("Line: %d: Expected next token to be %s, got: %s instead", p.currentToken.Line, t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseError(msg string) {
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() string {
	return strings.Join(p.errors, ", ")
}
