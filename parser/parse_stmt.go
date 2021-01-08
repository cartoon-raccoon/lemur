package parser

import (
	"fmt"

	"github.com/cartoon-raccoon/monkey-jit/ast"
	"github.com/cartoon-raccoon/monkey-jit/lexer"
)

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

func (p *Parser) parseBlockStatement() (*ast.BlockStatement, error) {
	block := &ast.BlockStatement{Token: p.current}
	block.Statements = []ast.Statement{}

	p.advance()

	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		node, err := p.parseNode()
		if err != nil {
			return block, err
		}
		stmt, ok := node.(ast.Statement)
		if stmt != nil && ok {
			block.Statements = append(block.Statements, stmt)
		} else {
			//todo: return err
		}
		p.advance()
	}

	return block, nil
}
