package parser

import (
	"fmt"
	"strconv"

	"github.com/cartoon-raccoon/monkey-jit/ast"
	"github.com/cartoon-raccoon/monkey-jit/lexer"
)

//Parser - represents the parser for Monkey
type Parser struct {
	lexer   *lexer.Lexer
	current lexer.Token
	next    lexer.Token

	prefixParseFns map[string]prefixParseFn
	infixParseFns  map[string]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	// LOWEST - The lowest precedence an expression can take
	LOWEST
	// EQUALS - ==
	EQUALS
	// COMPARE - < or >
	COMPARE
	// SUM - a + b
	SUM
	// PRODUCT - a * b
	PRODUCT
	// PREFIX - -x or !x
	PREFIX
	// CALL - a function call
	CALL
)

var precedences = map[string]int{
	lexer.EQ:  EQUALS,
	lexer.NE:  EQUALS,
	lexer.LT:  COMPARE,
	lexer.GT:  COMPARE,
	lexer.LE:  COMPARE,
	lexer.GE:  COMPARE,
	lexer.ADD: SUM,
	lexer.SUB: SUM,
	lexer.MUL: PRODUCT,
	lexer.DIV: PRODUCT,
}

func getPrecedence(tt string) int {
	if pre, ok := precedences[tt]; ok {
		return pre
	}
	return LOWEST
}

//New - returns a new Parser
func New(l *lexer.Lexer) (*Parser, error) {
	p := &Parser{lexer: l}

	err1 := p.advance()
	err2 := p.advance()

	if err1 != nil {
		return p, err1
	}
	if err2 != nil {
		return p, err2
	}

	// Registering prefix parse functions
	p.prefixParseFns = make(map[string]prefixParseFn)
	p.registerPrefixFn(lexer.IDENT, p.parseIdentifier)
	p.registerPrefixFn(lexer.INTLIT, p.parseIntLiteral)
	p.registerPrefixFn(lexer.FLTLIT, p.parseFltLiteral)
	p.registerPrefixFn(lexer.STRLIT, p.parseStrLiteral)
	p.registerPrefixFn(lexer.TRUE, p.parseBoolLiteral)
	p.registerPrefixFn(lexer.FALSE, p.parseBoolLiteral)
	p.registerPrefixFn(lexer.BANG, p.parsePrefixExpr)
	p.registerPrefixFn(lexer.SUB, p.parsePrefixExpr)
	p.registerPrefixFn(lexer.LPAREN, p.parseGroupedExpr)

	// Registering infix parse functions
	p.infixParseFns = make(map[string]infixParseFn)
	p.registerInfixFn(lexer.ADD, p.parseInfixExpr)
	p.registerInfixFn(lexer.SUB, p.parseInfixExpr)
	p.registerInfixFn(lexer.MUL, p.parseInfixExpr)
	p.registerInfixFn(lexer.DIV, p.parseInfixExpr)
	p.registerInfixFn(lexer.LE, p.parseInfixExpr)
	p.registerInfixFn(lexer.GE, p.parseInfixExpr)
	p.registerInfixFn(lexer.LT, p.parseInfixExpr)
	p.registerInfixFn(lexer.GT, p.parseInfixExpr)
	p.registerInfixFn(lexer.EQ, p.parseInfixExpr)
	p.registerInfixFn(lexer.NE, p.parseInfixExpr)

	return p, nil
}

// Parse - parses a stream of tokens
func (p *Parser) Parse() (*ast.Program, error) {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(lexer.EOF) {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		program.Statements = append(program.Statements, stmt)
		p.advance()
	}

	return program, nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.current.Type {
	case lexer.LET:
		return p.parseLetStatement()
	case lexer.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExprStatement()
	}
}

func (p *Parser) parseLetStatement() (*ast.LetStatement, error) {
	stmt := &ast.LetStatement{Token: p.current}

	if !p.nextTokenIs(lexer.IDENT) {
		return stmt, Err{
			Msg: fmt.Sprintf("Expected identifier, got `%s`", p.next.Type),
			Con: p.next.Pos,
		}
	}

	p.advance()

	stmt.Name = &ast.Identifier{Token: p.current, Value: p.current.Literal}

	if !p.nextTokenIs(lexer.ASSIGN) {
		return stmt, Err{
			Msg: fmt.Sprintf("Expected assignment operator, got `%s`", p.next.Type),
			Con: p.next.Pos,
		}
	}

	// fixme: a bit hacky
	p.advance() // p.current is now assign, p.next is expr start
	p.advance() // p.current is now expr start

	val := p.parseExpression(LOWEST)
	if val == nil {
		return stmt, Err{
			Msg: fmt.Sprintf("Could not parse expression"),
			Con: p.current.Pos,
		}
	}
	stmt.Value = val

	// note: this does not account for if the user forgets to put a semicolon
	// The parser will happily continue advancing until it hits a semicolon,
	// whenever that may be
	p.advance() //puts the next token in p.current

	return stmt, nil
}

func (p *Parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	stmt := &ast.ReturnStatement{Token: p.current}

	p.advance()

	for !p.curTokenIs(lexer.SEMICOL) {
		p.advance()
	}

	// if !p.nextTokenIs(lexer.SEMICOL) {
	// 	return nil, Err{
	// 		Msg: fmt.Sprintf("Expected semicolon, got %s", p.next.Type),
	// 		Con: p.next.Pos,
	// 	}
	// }

	return stmt, nil
}

func (p *Parser) parseExprStatement() (*ast.ExprStatement, error) {
	stmt := &ast.ExprStatement{Token: p.current}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.nextTokenIs(lexer.SEMICOL) {
		p.advance()
	}

	return stmt, nil
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.current.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()

	for !p.nextTokenIs(lexer.SEMICOL) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.next.Type]
		if infix == nil {
			return leftExp
		}

		p.advance()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.current, Value: p.current.Literal}
}

func (p *Parser) parseIntLiteral() ast.Expression {
	lit := &ast.Int{Token: p.current}
	value, err := strconv.ParseInt(p.current.Literal, 0, 64)
	if err != nil {
		//todo: return err
		return nil
	}
	lit.Inner = value
	return lit
}

func (p *Parser) parseFltLiteral() ast.Expression {
	lit := &ast.Flt{Token: p.current}
	value, err := strconv.ParseFloat(p.current.Literal, 64)
	if err != nil {
		return nil
	}
	lit.Inner = value
	return lit
}

func (p *Parser) parseStrLiteral() ast.Expression {
	lit := &ast.Str{Token: p.current}
	lit.Inner = p.current.Literal
	return lit
}

func (p *Parser) parseBoolLiteral() ast.Expression {
	lit := &ast.Bool{Token: p.current}
	value, err := strconv.ParseBool(p.current.Literal)
	if err != nil {
		return nil
	}
	lit.Inner = value
	return lit
}

func (p *Parser) parseGroupedExpr() ast.Expression {
	p.advance()

	expr := p.parseExpression(LOWEST)

	if !p.nextTokenIs(lexer.RPAREN) {
		return nil
	}

	return expr
}

func (p *Parser) parsePrefixExpr() ast.Expression {
	expression := &ast.PrefixExpr{
		Token:    p.current,
		Operator: p.current.Literal,
	}

	p.advance() // move the next expression into current

	expression.Right = p.parseExpression(PREFIX)

	if expression.Right == nil {
		// todo: return error
		return nil
	}
	return expression
}

func (p *Parser) parseInfixExpr(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpr{
		Token:    p.current,
		Operator: p.current.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.advance()
	expr.Right = p.parseExpression(precedence)
	if expr.Right == nil {
		//todo: return error
		return nil
	}
	return expr
}

func (p *Parser) registerPrefixFn(tt string, fn prefixParseFn) {
	p.prefixParseFns[tt] = fn
}

func (p *Parser) registerInfixFn(tt string, fn infixParseFn) {
	p.infixParseFns[tt] = fn
}

func (p *Parser) curPrecedence() int {
	return getPrecedence(p.current.Type)
}

func (p *Parser) peekPrecedence() int {
	return getPrecedence(p.next.Type)
}

func (p *Parser) curTokenIs(t string) bool {
	return p.current.Type == t
}

func (p *Parser) nextTokenIs(t string) bool {
	return p.next.Type == t
}

func (p *Parser) advance() error {
	p.current = p.next
	next, err := p.lexer.NextToken()
	if err != nil {
		return err
	}
	p.next = next
	return nil
}

//Err represents the error that can be thrown by the parser
type Err struct {
	Msg string
	Con lexer.Context
}

func (e Err) Error() string {
	return fmt.Sprintf(
		"%s: line %d, col %d",
		e.Msg,
		e.Con.Line,
		e.Con.Col,
	)
}
