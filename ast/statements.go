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
