package ast

import (
	"bytes"

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

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.Value + " ")
	out.WriteString("= ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")

	return out.String()

}

// ExprStatement - represents a bare expression in Monkey
type ExprStatement struct {
	Token      lexer.Token
	Expression Expression
}

func (es *ExprStatement) statementNode() {}

// TokenLiteral - Implements Node for ExpressionStatement
func (es *ExprStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExprStatement) String() string {
	var out bytes.Buffer
	if es.Expression != nil {
		out.WriteString(es.Expression.String())
	}
	return ""
}

// ReturnStatement - represents a return statement in the AST
type ReturnStatement struct {
	Token lexer.Token
	Value Expression
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral - Implements Node for ReturnStatement
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")
	if rs.Value != nil {
		out.WriteString(rs.Value.String())
	}
	out.WriteString(";")

	return out.String()

}

// BlockStatement represents a block of statements surrounded by braces
type BlockStatement struct {
	Token      lexer.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

// TokenLiteral implements Node for BlockStatement
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

// String implements Node for BlockStatement
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	out.WriteString("{")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	out.WriteString("}")

	return out.String()
}
