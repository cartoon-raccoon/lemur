package parser

import (
	"fmt"

	"github.com/cartoon-raccoon/lemur/ast"
	"github.com/cartoon-raccoon/lemur/lexer"
)

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.current}

	if !p.nextTokenIs(lexer.IDENT) {
		p.errors = append(p.errors, Err{
			Msg: fmt.Sprintf("Expected identifier, got `%s`", p.next.Type),
			Con: p.next.Pos,
		})
		return nil
	}

	p.advance()

	stmt.Name = &ast.Identifier{Token: p.current, Value: p.current.Literal}

	if !p.nextTokenIs(lexer.ASSIGN) {
		p.errors = append(p.errors, Err{
			Msg: fmt.Sprintf("Expected assignment operator, got `%s`", p.next.Type),
			Con: p.next.Pos,
		})
		return nil
	}

	// fixme: a bit hacky
	p.advance() // p.current is now assign, p.next is expr start
	p.advance() // p.current is now expr start

	val := p.parseExpression(LOWEST)
	if val == nil {
		return nil
	}
	stmt.Value = val

	// note: this does not account for if the user forgets to put a semicolon
	// The parser will happily continue advancing until it hits a semicolon,
	// whenever that may be
	p.advance() //puts the next token in p.current

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.current}

	p.advance() // p.current is now expr start

	stmt.Value = p.parseExpression(LOWEST)

	if stmt.Value == nil {
		return nil
	}
	p.advance()

	return stmt
}

func (p *Parser) parseExprStatement() *ast.ExprStatement {
	stmt := &ast.ExprStatement{Token: p.current}

	stmt.Expression = p.parseExpression(LOWEST)

	if stmt.Expression == nil {
		return nil
	}

	if p.nextTokenIs(lexer.SEMICOL) {
		p.advance()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.current}
	block.Statements = []ast.Statement{}

	p.advance()

	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		node := p.parseNode()
		if node == nil {
			return nil
		}
		stmt, ok := node.(ast.Statement)
		if stmt != nil && ok {
			block.Statements = append(block.Statements, stmt)
		} else {
			p.errors = append(p.errors, Err{
				Msg: "Only statements can be declared in blocks",
				Con: p.current.Pos,
			})
		}
		p.advance()
	}

	return block
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	while := &ast.WhileStatement{Token: p.current}

	p.advance()

	if !p.curTokenIs(lexer.LPAREN) {
		p.errors = append(p.errors, &Err{
			Msg: fmt.Sprintf("Expected `(`, got %s", p.current.Literal),
			Con: p.current.Pos,
		})
	}

	while.Condition = p.parseExpression(LOWEST)

	if !p.nextTokenIs(lexer.LBRACE) {
		p.errors = append(p.errors, &Err{
			Msg: fmt.Sprintf("Expected `{`, got %s", p.next.Literal),
			Con: p.next.Pos,
		})
	}

	p.advance()

	while.Body = p.parseBlockStatement()
	if while.Body == nil {
		return nil
	}

	return while
}
