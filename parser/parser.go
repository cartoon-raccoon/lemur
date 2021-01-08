package parser

import (
	"fmt"

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
	// BITWISE - <<, >>, &, |, ~
	BITWISE
	// LOGIC - &&, ||
	LOGIC
	// PREFIX - -x or !x
	PREFIX
	// CALL - a function call
	CALL
)

var precedences = map[string]int{
	lexer.EQ:    EQUALS,
	lexer.NE:    EQUALS,
	lexer.LT:    COMPARE,
	lexer.GT:    COMPARE,
	lexer.LE:    COMPARE,
	lexer.GE:    COMPARE,
	lexer.ADD:   SUM,
	lexer.SUB:   SUM,
	lexer.MUL:   PRODUCT,
	lexer.DIV:   PRODUCT,
	lexer.BWAND: BITWISE,
	lexer.BWOR:  BITWISE,
	lexer.BWNOT: BITWISE,
	lexer.BSR:   BITWISE,
	lexer.BSL:   BITWISE,
	lexer.LOR:   LOGIC,
	lexer.LAND:  LOGIC,
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
	p.registerPrefixFn(lexer.IF, p.parseIfExpression)
	p.registerPrefixFn(lexer.FUNCTION, p.parseFnLiteral)

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
	p.registerInfixFn(lexer.BSL, p.parseInfixExpr)
	p.registerInfixFn(lexer.BSR, p.parseInfixExpr)
	p.registerInfixFn(lexer.BWAND, p.parseInfixExpr)
	p.registerInfixFn(lexer.BWOR, p.parseInfixExpr)
	p.registerInfixFn(lexer.BWNOT, p.parseInfixExpr)
	p.registerInfixFn(lexer.LAND, p.parseInfixExpr)
	p.registerInfixFn(lexer.LOR, p.parseInfixExpr)

	return p, nil
}

// Parse - parses a stream of tokens
func (p *Parser) Parse() (*ast.Program, error) {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(lexer.EOF) {
		node, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		switch node.(type) {
		case ast.Statement:
			program.Statements = append(program.Statements, node.(ast.Statement))
		case ast.Declaration:
			switch node.(ast.Declaration).(type) {
			case *ast.FunctionDecl:
				program.Functions = append(program.Functions, *node.(ast.Declaration).(*ast.FunctionDecl))
			default:
				//todo: return err
				return nil, nil
			}
		default:
			//todo: return err
			return nil, nil
		}

		p.advance()
	}

	return program, nil
}

func (p *Parser) parseNode() (ast.Node, error) {
	switch p.current.Type {
	case lexer.LET:
		return p.parseLetStatement()
	case lexer.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExprStatement()
	}
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.current.Type]
	if prefix == nil {
		//todo: return error
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
