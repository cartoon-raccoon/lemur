package ast

import (
	"bytes"
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
