package parser

import (
	"github.com/cartoon-raccoon/monkey-jit/ast"
	"github.com/cartoon-raccoon/monkey-jit/lexer"
)

type Parser struct {
	lexer   *lexer.Lexer
	current lexer.Token
	next    lexer.Token
}

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

func (p *Parser) Parse() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(lexer.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.advance()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.current.Type {
	case lexer.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.current}

	if !p.nextTokenIs(lexer.IDENT) {
		return nil
	}

	p.advance()

	stmt.Name = &ast.Identifier{Token: p.current, Value: p.current.Literal}

	if !p.nextTokenIs(lexer.ASSIGN) {
		return nil
	}

	for !p.curTokenIs(lexer.SEMICOL) {
		p.advance()
	}

	return stmt
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
