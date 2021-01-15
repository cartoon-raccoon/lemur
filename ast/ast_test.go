package ast

import (
	"testing"

	"github.com/cartoon-raccoon/lemur/lexer"
)

func TestString(t *testing.T) {
	prog := Program{
		Statements: []Statement{
			&LetStatement{
				Token: lexer.Token{
					Type:    lexer.LET,
					Literal: lexer.LET,
				},
				Name: &Identifier{
					Token: lexer.Token{
						Type:    lexer.IDENT,
						Literal: "my_var",
					},
					Value: "my_var",
				},
				Value: &Identifier{
					Token: lexer.Token{
						Type:    lexer.IDENT,
						Literal: "another_var",
					},
					Value: "another_var",
				},
			},
		},
	}
	progstring := prog.String()
	if progstring != "let my_var = another_var;" {
		t.Fatalf("Incorrect p.String(): got %q", progstring)
	}
}
