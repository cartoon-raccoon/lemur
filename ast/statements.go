package ast

import (
	"github.com/cartoon-raccoon/monkey-jit/lexer"
)

// LetStatement - represents a let statement in the AST
type LetStatement struct {
	Token lexer.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

// TokenLiteral - Implements Node for LetStatement
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// ReturnStatement - represents a return statement in the AST
type ReturnStatement struct {
	Token lexer.Token
	Value Expression
}

func (ls *ReturnStatement) statementNode() {}

// TokenLiteral - Implements Node for ReturnStatement
func (ls *ReturnStatement) TokenLiteral() string {
	return ls.Token.Literal
}
