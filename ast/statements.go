package ast

import (
	"bytes"

	"github.com/cartoon-raccoon/lemur/lexer"
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

// Context implements Node for LetStatement
func (ls *LetStatement) Context() lexer.Context {
	return ls.Token.Pos
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
	return out.String()
}

// Context implements Node for ExprStatement
func (es *ExprStatement) Context() lexer.Context {
	return es.Token.Pos
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

// Context implements Node for ReturnStatement
func (rs *ReturnStatement) Context() lexer.Context {
	return rs.Token.Pos
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

	out.WriteString("{\n")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
		out.WriteByte('\n')
	}
	out.WriteString("\n}")

	return out.String()
}

// Context implements Node for BlockStatement
func (bs *BlockStatement) Context() lexer.Context {
	return bs.Token.Pos
}

// WhileStatement represents a while loop
type WhileStatement struct {
	Token     lexer.Token
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode() {}

// TokenLiteral implements Node for WhileStatement
func (ws *WhileStatement) TokenLiteral() string {
	return ws.Token.Literal
}

// String implements Node for WhileStatement
func (ws *WhileStatement) String() string {
	var out bytes.Buffer

	out.WriteString("while ")
	out.WriteString(ws.Condition.String())
	out.WriteString(ws.Body.String())

	return out.String()
}

// Context implements Node for WhileStatement
func (ws *WhileStatement) Context() lexer.Context {
	return ws.Token.Pos
}

// BreakStatement represents a break statement
type BreakStatement struct {
	Token lexer.Token
}

func (bs *BreakStatement) statementNode() {}

// TokenLiteral implements Node for BreakStatement
func (bs *BreakStatement) TokenLiteral() string {
	return bs.Token.Literal
}

// String implements Node for BreakStatement
func (bs *BreakStatement) String() string {
	return bs.Token.Literal
}

// Context implements Node for BreakStatement
func (bs *BreakStatement) Context() lexer.Context {
	return bs.Token.Pos
}
