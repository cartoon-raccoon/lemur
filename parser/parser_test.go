package parser

import (
	"testing"

	"github.com/cartoon-raccoon/monkey-jit/ast"
	"github.com/cartoon-raccoon/monkey-jit/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `let x = 10;
	let y = 5_000;
	let nice = 420.69;`

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
		t.Errorf("statement not let statement: got %T", stmt)
	}
	if letstmt.Name.Value != ident {
		t.Errorf("wrong identifier: got %q", letstmt.Name.Value)
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
