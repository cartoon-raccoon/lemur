package parser

import (
	"fmt"

	"github.com/cartoon-raccoon/monkey-jit/ast"
	"github.com/cartoon-raccoon/monkey-jit/lexer"
)

func (p *Parser) parseFuncDecl() ast.Declaration {
	fndecl := &ast.FunctionDecl{Token: p.current}
	p.advance()
	expr := p.parseExpression(LOWEST)

	ident, ok := expr.(*ast.Identifier)
	if !ok {
		p.errors = append(p.errors, &Err{
			Msg: fmt.Sprintf("Expected identifier, got %s", p.current.Type),
			Con: p.current.Pos,
		})
		return nil
	}
	fndecl.Name = ident

	p.advance()
	if !p.curTokenIs(lexer.LPAREN) {
		p.errors = append(p.errors, &Err{
			Msg: fmt.Sprintf("Expected `(`, got %s", p.current.Literal),
			Con: p.current.Pos,
		})
	}
	fndecl.Params = p.parseFunctionParams()
	if fndecl.Params == nil {
		return nil
	}
	p.advance()
	// p.next should now be lbrace
	if p.nextTokenIs(lexer.LBRACE) {
		p.advance()
		body := p.parseBlockStatement()
		if body == nil {
			return nil
		}
		fndecl.Body = body
	} else {
		p.errors = append(p.errors, &Err{
			Msg: fmt.Sprintf("Expected block, got %s", p.next.Literal),
			Con: p.next.Pos,
		})
		return nil
	}

	return fndecl
}
