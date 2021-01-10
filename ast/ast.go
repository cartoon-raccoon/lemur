package ast

import (
	"bytes"
	"strings"

	"github.com/cartoon-raccoon/monkey-jit/lexer"
)

/*
Node, Statement and Expression are generic interfaces that
define the grammar of the language. Every language construct
can be defined as either a statement or expression, and the
rest are defined as literals that are stored as separate data.

Node defines a single node on the AST. Both Statement and Expression
implement Node, so that they can be represented on the AST.
*/

// Node defines the general behaviour for a node in the AST
type Node interface {
	TokenLiteral() string
	Context() lexer.Context
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

// Declaration defines a top level declaration (e.g. class, fn)
type Declaration interface {
	Node
	declarationNode()
}

// Program represents the entire parsed program
type Program struct {
	Statements []Statement
	Functions  []FunctionDecl
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

// Context implements Node for Program
func (p *Program) Context() lexer.Context {
	return lexer.Context{}
}

// FunctionDecl represents a function declaration
type FunctionDecl struct {
	Token  lexer.Token
	Name   string
	Params []*Identifier
	Body   *BlockStatement
}

func (fd *FunctionDecl) declarationNode() {}

// TokenLiteral implements Node for FunctionDecl
func (fd *FunctionDecl) TokenLiteral() string {
	return fd.Token.Literal
}

// String implements Node for FunctionDecl
func (fd *FunctionDecl) String() string {
	var out bytes.Buffer

	params := []string{}

	for _, p := range fd.Params {
		params = append(params, p.String())
	}
	out.WriteString("fn ")
	out.WriteString(fd.Name)
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fd.Body.String())

	return out.String()
}

// Context implements Node for FunctionDecl
func (fd *FunctionDecl) Context() lexer.Context {
	return fd.Token.Pos
}
