package parser

import (
	"fmt"
	"strconv"

	"github.com/cartoon-raccoon/monkey-jit/ast"
	"github.com/cartoon-raccoon/monkey-jit/lexer"
)

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.current, Value: p.current.Literal}
}

func (p *Parser) parseIntLiteral() ast.Expression {
	lit := &ast.Int{Token: p.current}
	value, err := strconv.ParseInt(p.current.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, Err{
			Msg: fmt.Sprintf("Unable to parse %s as integer", p.current.Literal),
			Con: p.current.Pos,
		})
		return nil
	}
	lit.Inner = value
	return lit
}

func (p *Parser) parseFltLiteral() ast.Expression {
	lit := &ast.Flt{Token: p.current}
	value, err := strconv.ParseFloat(p.current.Literal, 64)
	if err != nil {
		p.errors = append(p.errors, Err{
			Msg: fmt.Sprintf("Unable to parse %s as float", p.current.Literal),
			Con: p.current.Pos,
		})
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
		p.errors = append(p.errors, Err{
			Msg: fmt.Sprintf("Unable to parse %s as bool", p.current.Literal),
			Con: p.current.Pos,
		})
		return nil
	}
	lit.Inner = value
	return lit
}

func (p *Parser) parseGroupedExpr() ast.Expression {
	p.advance()

	expr := p.parseExpression(LOWEST)

	if expr == nil {
		return nil
	}

	if !p.nextTokenIs(lexer.RPAREN) {
		p.errors = append(p.errors, Err{
			Msg: fmt.Sprintf("Expected `)`, got %s", p.next.Literal),
			Con: p.next.Pos,
		})
		return nil
	}

	p.advance()

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
		return nil
	}
	return expr
}

func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{Token: p.current}

	if !p.nextTokenIs(lexer.LPAREN) {
		p.errors = append(p.errors, Err{
			Msg: fmt.Sprintf("Expected `(` after 'if', got %s", p.next.Literal),
			Con: p.next.Pos,
		})
		return nil
	}
	p.advance()
	// p.current is now LPAREN

	expr.Condition = p.parseExpression(LOWEST)

	if !p.nextTokenIs(lexer.LBRACE) {
		p.errors = append(p.errors, Err{
			Msg: fmt.Sprintf("Expected start of block, got %s", p.next.Literal),
			Con: p.next.Pos,
		})
		return nil
	}
	p.advance()

	res := p.parseBlockStatement()
	if res == nil {
		return nil
	}
	// p.current is now RBRACE

	expr.Result = res

	if p.nextTokenIs(lexer.ELSE) {
		p.advance()

		if p.nextTokenIs(lexer.LBRACE) {
			p.advance()

			alt := p.parseBlockStatement()
			if alt == nil {
				return nil
			}
			expr.Alternative = alt
		} else if p.nextTokenIs(lexer.IF) {
			p.advance()
			expr.Alternative = p.parseIfExpression()
		} else {
			p.errors = append(p.errors, Err{
				Msg: fmt.Sprintf("Expected 'if' or `{`, got %s", p.next.Literal),
				Con: p.next.Pos,
			})
			return nil
		}

	}

	return expr
}

func (p *Parser) parseFnLiteral() ast.Expression {
	lit := &ast.FnLiteral{Token: p.current}

	if p.current.Type != lexer.FUNCTION {
		panic(fmt.Sprintf("Wrong token, got %s", p.current.Type))
	}

	if !p.nextTokenIs(lexer.LPAREN) {
		p.errors = append(p.errors, Err{
			Msg: fmt.Sprintf("Expected `(`, got %s", p.next.Literal),
			Con: p.next.Pos,
		})
		return nil
	}

	p.advance()
	// p.current is now lparen

	lit.Params = p.parseFunctionParams()
	// p.next is now rparen

	p.advance()
	// p.next should now be lbrace
	if p.nextTokenIs(lexer.LBRACE) {
		p.advance()
		body := p.parseBlockStatement()
		if body == nil {
			return nil
		}
		lit.Body = body
	}

	return lit
}

func (p *Parser) parseFunctionParams() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	//p.current is lparen
	if p.nextTokenIs(lexer.RPAREN) {
		return identifiers
	}

	p.advance()
	//p.current should now be IDENT

	ident := &ast.Identifier{Token: p.current, Value: p.current.Literal}
	identifiers = append(identifiers, ident)

	// p.next should now be comma
	for p.nextTokenIs(lexer.COMMA) {
		p.advance() //p.current == comma
		p.advance() //p.current == ident
		ident = &ast.Identifier{Token: p.current, Value: p.current.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.nextTokenIs(lexer.RPAREN) {
		p.errors = append(p.errors, Err{
			Msg: fmt.Sprintf("Expected `)`, got %s", p.next.Literal),
			Con: p.next.Pos,
		})
		return nil
	}

	return identifiers
}

func (p *Parser) parseFunctionCall(fn ast.Expression) ast.Expression {
	exp := &ast.FunctionCall{Token: p.current, Ident: fn}
	exp.Params = p.parseCallArguments()

	if exp.Params == nil {
		return nil
	}

	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.nextTokenIs(lexer.RPAREN) {
		p.advance()
		return args
	}

	p.advance()
	args = append(args, p.parseExpression(LOWEST))

	for p.nextTokenIs(lexer.COMMA) {
		p.advance()
		p.advance()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.nextTokenIs(lexer.RPAREN) {
		p.errors = append(p.errors, Err{
			Msg: fmt.Sprintf("ParseCallArgs: Expected `)`, got %s", p.next.Literal),
			Con: p.next.Pos,
		})
		return nil
	}

	p.advance()

	return args
}
