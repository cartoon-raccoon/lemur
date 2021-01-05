package ast

import (
	"bytes"
	"strconv"

	"github.com/cartoon-raccoon/monkey-jit/lexer"
)

// Node defines the general behaviour for a node in the AST
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement defines a statement in Monkey syntax
type Statement interface {
	Node
	statementNode()
}

// Expression defines an expression in Monkey Syntax
type Expression interface {
	Node
	expressionNode()
}

// Program represents the entire parsed program
type Program struct {
	Statements []Statement
}

// TokenLiteral implements Node for string
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()

}

// Identifier - represents a name bound to a function or a variable
type Identifier struct {
	Token lexer.Token
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral - implements Node for Identifier
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// String - implements Node for Identifier
func (i *Identifier) String() string {
	return i.Value
}

// Literal defines a literal in Monkey: string, int or float
type Literal interface {
	Literal()
	String() string
}

// Int represents an integer literal in the Monkey AST
type Int struct {
	Token lexer.Token
	Inner int64
}

// Literal implements Literal for Int
func (i *Int) Literal()        {}
func (i *Int) expressionNode() {}

// TokenLiteral implements Node for Int
func (i *Int) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Int) String() string {
	return i.Token.Literal
}

func intFromRaw(raw string) Literal {
	num, err := strconv.ParseInt(raw, 0, 64)
	if err != nil {
		return nil
	}
	return &Int{Inner: num}
}

// Flt represents a float literal in the Monkey AST
type Flt struct {
	Token lexer.Token
	Inner float64
}

// Literal implements Literal for Flt
func (f *Flt) Literal()        {}
func (f *Flt) expressionNode() {}

// TokenLiteral implements Node for Flt
func (f *Flt) TokenLiteral() string {
	return f.Token.Literal
}
func (f *Flt) String() string {
	return f.Token.Literal
}

func fltFromRaw(raw string) Literal {
	num, err := strconv.ParseFloat(raw, 0)
	if err != nil {
		return nil
	}
	return &Flt{Inner: num}
}

// Str represents a string literal in the Monkey AST
type Str struct {
	Token lexer.Token
	Inner string
}

// Literal implements Literal for Str
func (s *Str) Literal()        {}
func (s *Str) expressionNode() {}

// TokenLiteral implements Node for Str
func (s *Str) TokenLiteral() string {
	return s.Token.Literal
}
func (s *Str) String() string {
	return s.Inner
}

func strFromLit(raw string) Literal {
	// todo: add method to parse escaped chars
	return &Str{Inner: raw}
}
