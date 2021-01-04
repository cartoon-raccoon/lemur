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
	program := p.Parse()
	if program == nil {
		t.Fatalf("Parser could not generate AST")
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
	if stmt.TokenLiteral() != "let" {
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
