package parser

import (
	"fmt"

	"github.com/cartoon-raccoon/lemur/ast"
	"github.com/cartoon-raccoon/lemur/lexer"
)

func (p *Parser) parseFuncDecl() ast.Declaration {
	fndecl := &ast.FunctionDecl{Token: p.current}
	p.advance()
	if !p.curTokenIs(lexer.IDENT) {
		p.errors = append(p.errors, &Err{
			Msg: fmt.Sprintf("Expected identifier, got %s", p.current.Literal),
			Con: p.current.Pos,
		})
		return nil
	}
	fndecl.Name = p.parseIdentifier().(*ast.Identifier)

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
