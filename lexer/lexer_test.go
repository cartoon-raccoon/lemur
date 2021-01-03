package lexer

import "testing"

func TestNextToken(t *testing.T) {

	input := `let five = 5;
let ten = 10;
fn hello() {
	five += 10;
	let sum = five + 10;
}
let str = "Hello i am\n cool\n";
let add = fn(a, b) {
	return a + b
}
let sum2 = add(15, 35);
let thou = 1_000.57;
if thou == 1_000.57 {
	thou += 4;
} else {
	hello();
}`

	l := New(input)

	tests := []struct {
		expectedToken   string
		expectedLiteral string
	}{
		{LET, "LET"},
		{IDENT, "five"},
		{ASSIGN, "="},
		{IDENT, "5"},
		{SEMICOL, ";"},
		{LET, "LET"},
		{IDENT, "ten"},
		{ASSIGN, "="},
		{IDENT, "10"},
		{SEMICOL, ";"},
		{FUNCTION, "FUNCTION"},
		{IDENT, "hello"},
		{LPAREN, "("},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{IDENT, "five"},
		{ADDASSIGN, "+="},
		{IDENT, "10"},
		{SEMICOL, ";"},
		{LET, "LET"},
		{IDENT, "sum"},
		{ASSIGN, "="},
		{IDENT, "five"},
		{ADD, "+"},
		{IDENT, "10"},
		{SEMICOL, ";"},
		{RBRACE, "}"},
		{LET, "LET"},
		{IDENT, "str"},
		{ASSIGN, "="},
		{STRLIT, "Hello i am\\n cool\\n"},
		{SEMICOL, ";"},
		{LET, "LET"},
		{IDENT, "add"},
		{ASSIGN, "="},
		{FUNCTION, "FUNCTION"},
		{LPAREN, "("},
		{IDENT, "a"},
		{COMMA, ","},
		{IDENT, "b"},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{RETURN, "RETURN"},
		{IDENT, "a"},
		{ADD, "+"},
		{IDENT, "b"},
	}

	for i, tt := range tests {
		tok, err := l.nextToken()
		if err != nil {
			t.Fatalf("Error on token %d, expected %q", i, tt.expectedToken)
		}
		if tok.Type != tt.expectedToken {
			t.Fatalf(
				"token %d: wrong token, expected %q, got %q with literal '%s'",
				i,
				tt.expectedToken,
				tok.Type,
				tok.Literal,
			)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf(
				"token %d, wrong literal, expected %q, got %q",
				i,
				tt.expectedLiteral,
				tok.Literal,
			)
		}
	}
}
