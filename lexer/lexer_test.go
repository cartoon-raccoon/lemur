package lexer

import "testing"

func TestNextToken(t *testing.T) {

	input := `let five = 5;
let ten = 10.0;
fn hello() {
	five = [5, 6.9, "hello", 10];
	let sum = five + 10_000;
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
		{LET, "let"},
		{IDENT, "five"},
		{ASSIGN, "="},
		{INTLIT, "5"},
		{SEMICOL, ";"},
		{LET, "let"},
		{IDENT, "ten"},
		{ASSIGN, "="},
		{FLTLIT, "10.0"},
		{SEMICOL, ";"},
		{FUNCTION, "fn"},
		{IDENT, "hello"},
		{LPAREN, "("},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{IDENT, "five"},
		{ASSIGN, "="},
		{LSBRKT, "["},
		{INTLIT, "5"},
		{COMMA, ","},
		{FLTLIT, "6.9"},
		{COMMA, ","},
		{STRLIT, "hello"},
		{COMMA, ","},
		{INTLIT, "10"},
		{RSBRKT, "]"},
		{SEMICOL, ";"},
		{LET, "let"},
		{IDENT, "sum"},
		{ASSIGN, "="},
		{IDENT, "five"},
		{ADD, "+"},
		{INTLIT, "10_000"},
		{SEMICOL, ";"},
		{RBRACE, "}"},
		{LET, "let"},
		{IDENT, "str"},
		{ASSIGN, "="},
		{STRLIT, "Hello i am\\n cool\\n"},
		{SEMICOL, ";"},
		{LET, "let"},
		{IDENT, "add"},
		{ASSIGN, "="},
		{FUNCTION, "fn"},
		{LPAREN, "("},
		{IDENT, "a"},
		{COMMA, ","},
		{IDENT, "b"},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{RETURN, "return"},
		{IDENT, "a"},
		{ADD, "+"},
		{IDENT, "b"},
	}

	for i, tt := range tests {
		tok, err := l.NextToken()
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
