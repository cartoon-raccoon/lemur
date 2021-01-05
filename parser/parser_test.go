package parser

import (
	"testing"

	"github.com/cartoon-raccoon/monkey-jit/ast"
	"github.com/cartoon-raccoon/monkey-jit/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `let x = 10;
	let y = 5_000;
	let nice = 69;`

	l := lexer.New(input)
	p, err := New(l)
	if err != nil {
		t.Fatalf("Error instantiating parser: malformed input")
	}
	program, err := p.Parse()
	if program == nil || err != nil {
		t.Fatalf(err.Error())
	}

	tests := []struct {
		expectedIdent string
	}{
		{"x"},
		{"y"},
		{"nice"},
	}

	for i, ident := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, ident.expectedIdent) {
			return
		}
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, ident string) bool {
	if stmt.TokenLiteral() != lexer.LET {
		t.Errorf("error: wrong token, expected \"let\" got %q", stmt.TokenLiteral())
	}
	letstmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Fatalf("statement not let statement: got %T", stmt)
	}
	if letstmt.Name.Value != ident {
		t.Errorf("wrong identifier: got %q", letstmt.Name.Value)
	}
	_, ok = letstmt.Value.(*ast.Int)
	if !ok {
		t.Errorf("value is not int")
	}

	return true
}

func TestReturnStatement(t *testing.T) {
	input := `return x;
	return 0.78;
	return "what";`

	l := lexer.New(input)
	p, err := New(l)
	if err != nil {
		t.Fatalf("Error instantiating parser: malformed input")
	}
	program, err := p.Parse()
	if program == nil {
		t.Fatalf(err.Error())
	}

	if len := len(program.Statements); len != 3 {
		t.Fatalf("Wrong number of statements: expected 3, got %d", len)
	}

	for _, stmt := range program.Statements {
		if stmt.TokenLiteral() != lexer.RETURN {
			t.Fatalf("Error: expected RETURN, got %q", stmt.TokenLiteral())
		}
		_, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("statement is not a return statement")
		}
	}
}

func TestIdentExpr(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p, err := New(l)
	if err != nil {
		t.Fatalf("Error: lexer encountered an error - %s", err)
	}

	prog, err := p.Parse()
	if err != nil {
		t.Fatalf("Error: parser encountered an error - %s", err)
	}

	if len := len(prog.Statements); len != 1 {
		t.Fatalf("Not enough statements, got: %d", len)
	}

	if stmt, ok := prog.Statements[0].(*ast.ExprStatement); !ok {
		t.Fatalf("Statement is not expression statement, got %q", stmt.String())
		if ident, ok := stmt.Expression.(*ast.Identifier); ok {
			if ident.Token.Type != lexer.IDENT || ident.Value != "foobar" {
				t.Fatalf("Incorrect identifier, got %s", ident.Value)
			}
		} else {
			t.Fatalf("Statement is not identifier")
		}
	}
}

func TestIntExpr(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p, err := New(l)
	if err != nil {
		t.Fatalf("Error: lexer encountered an error - %s", err)
	}

	prog, err := p.Parse()
	if err != nil {
		t.Fatalf("Error: parser encountered an error - %s", err)
	}

	if len := len(prog.Statements); len != 1 {
		t.Fatalf("Not enough statements, got: %d", len)
	}

	if stmt, ok := prog.Statements[0].(*ast.ExprStatement); ok {
		if intlit, ok := stmt.Expression.(*ast.Int); ok {
			if intlit.Inner != 5 {
				t.Fatalf("Integer is not 5, got %d", intlit.Inner)
			}
		} else {
			t.Fatalf("Expression is not an integer")
		}
	} else {
		t.Fatalf("statement is not an expression statment")
	}
}
