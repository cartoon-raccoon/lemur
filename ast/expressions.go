package ast

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/cartoon-raccoon/monkey-jit/lexer"
)

//*----------| Identifier |----------*/

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

//*----------| PrefixExpr |----------*/

// PrefixExpr represents a prefixed expression, such as ! or -
type PrefixExpr struct {
	Token    lexer.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpr) expressionNode() {}

// TokenLiteral implements Node for PrefixExpr
func (pe *PrefixExpr) TokenLiteral() string {
	return pe.Token.Literal
}

// String implements Node for PrefixExpr
func (pe *PrefixExpr) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

//*----------| InfixExpr |----------*/

// InfixExpr represents an expression with an infixed operator
type InfixExpr struct {
	Token    lexer.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpr) expressionNode() {}

// TokenLiteral implements Node for InfixExpr
func (ie *InfixExpr) TokenLiteral() string {
	return ie.Token.Literal
}

// String implements Node InfixExpr
func (ie *InfixExpr) String() string {
	var out bytes.Buffer

	out.WriteString("(" + ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String() + ")")

	return out.String()
}

//*----------| IfExpression |----------*/

// IfExpression represents an if statement/expr in Monkey
type IfExpression struct {
	Token       lexer.Token
	Condition   Expression
	Result      *BlockStatement
	Alternative Node //! I don't like using a Node here, find a better alternative
	//* Alternative can only be *BlockStatement or *IfExpression
	//* Must perform a runtime check
}

func (ie *IfExpression) expressionNode() {}

// TokenLiteral implements Node for IfExpression
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// String implements Node for IfExpression
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if (")
	out.WriteString(ie.Condition.String())
	out.WriteString(") ")
	out.WriteString(ie.Result.String())
	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

//*----------| Fnliteral |----------*/

// FnLiteral represents a function declaration in Monkey
type FnLiteral struct {
	Token  lexer.Token
	Params []*Identifier
	Body   *BlockStatement
}

func (fl *FnLiteral) expressionNode() {}

// TokenLiteral implements Node for FnLiteral
func (fl *FnLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

// String implements Node for FnLiteral
func (fl *FnLiteral) String() string {
	var out bytes.Buffer

	params := []string{}

	for _, p := range fl.Params {
		params = append(params, p.String())
	}
	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

//*----------| FunctionCall |----------*/

// FunctionCall defines a function call in Monkey
type FunctionCall struct {
	Token  lexer.Token
	Ident  Expression
	Params []Expression
}

func (fc *FunctionCall) expressionNode() {}

// TokenLiteral implements Node for FunctionCall
func (fc *FunctionCall) TokenLiteral() string {
	return fc.Token.Literal
}

// String implements Node for FunctionCall
func (fc *FunctionCall) String() string {
	var out bytes.Buffer

	params := []string{}

	for _, p := range fc.Params {
		params = append(params, p.String())
	}
	out.WriteString(fc.Ident.String())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")

	return out.String()
}

//*----------| Literals |----------*/

// Literal defines a literal in Monkey: string, int or float
type Literal interface {
	Node
	Literal()
	String() string
}

//! Int

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

//! Flt

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

//! Str

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

//! Bool

// Bool represents a boolean literal in the Monkey AST
type Bool struct {
	Token lexer.Token
	Inner bool
}

// Literal implements Literal for Bool
func (b *Bool) Literal()        {}
func (b *Bool) expressionNode() {}

// TokenLiteral implements Node for Bool
func (b *Bool) TokenLiteral() string {
	return b.Token.Literal
}

// TokenLiteral implements Node for Bool
func (b *Bool) String() string {
	return b.Token.Literal
}
