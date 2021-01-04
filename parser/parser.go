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
	default:
		return nil, nil
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

	for !p.curTokenIs(lexer.SEMICOL) {

		p.advance()
	}

	return stmt, nil
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
