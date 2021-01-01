package lexer

import "testing"

func TestNextToken(t *testing.T) {

	input := `let five = 5;
let ten = 10;
fn hello() {
	five += 10;
	let sum = five + 10;
}
let add = fn(a, b) {
	return a + b
}
let sum2 = add(15, 35);`

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
	}

	for i, tt := range tests {
		tok := l.nextToken()

		if tok.Type != tt.expectedToken {
			t.Fatalf(
				"token %d: wrong token, expected %q, got %q",
				i,
				tt.expectedToken,
				tok.Type,
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
