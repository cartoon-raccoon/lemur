package ast

import (
	"github.com/cartoon-raccoon/monkey-jit/lexer"
)

// Node defines the general behaviour for a node in the AST
type Node interface {
	TokenLiteral() string
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

type LetStatement struct {
	Token lexer.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

type Identifier struct {
	Token lexer.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
