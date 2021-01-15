package eval

import (
	"testing"

	"github.com/cartoon-raccoon/lemur/lexer"
	"github.com/cartoon-raccoon/lemur/object"
	"github.com/cartoon-raccoon/lemur/parser"
)

func TestExprStmtEval(t *testing.T) {
	tests := []struct {
		Input    string
		Expected object.Object
		Inspect  string
	}{
		{
			"1 + (1 * 7) / 2",
			&object.Integer{Value: 4},
			"4",
		},
		{
			"(5 + 10 * 2 + 15 / 3) == 2 + -10",
			&object.Boolean{Value: false},
			"false",
		},
		{
			"\"hello\"",
			&object.String{Value: "hello"},
			"hello",
		},
		{
			"420.69 + 7.4",
			&object.Float{Value: 428.09},
			"428.090000",
		},
	}

	for i, test := range tests {
		l := lexer.New(test.Input)
		p, err := parser.New(l)
		if err != nil {
			t.Errorf("Test %d: Error while beginning lexing", i)
			continue
		}
		prog := p.Parse()
		if p.CheckErrors() != nil {
			t.Errorf("Test %d: Errors while parsing", i)
			for _, err := range p.CheckErrors() {
				t.Logf(err.Error())
			}
			t.Log("Aborting tests")
			t.FailNow()
		}
		if prog == nil {
			t.Fatalf("Test %d: Program is nil", i)
		}
		eval := Evaluator{}
		env := object.NewEnv()
		res := eval.Evaluate(prog, env)
		if res == nil {
			t.Fatalf("Test %d: Could not evaluate", i)
		}
		if res.Inspect() != test.Inspect {
			t.Errorf("Test %d: Expected %s, got %s", i, test.Inspect, res.Inspect())
		}
	}
}

func TestIfExprEval(t *testing.T) {
	input := "if (6 < 7) { return 5; } else { return \"hello\"; }"
	expected := "5"
	l := lexer.New(input)
	p, err := parser.New(l)
	if err != nil {
		t.Fatalf("Error while beginning lexing")
	}
	prog := p.Parse()
	if p.CheckErrors() != nil {
		t.Errorf("Errors while parsing")
		for _, err := range p.CheckErrors() {
			t.Logf(err.Error())
		}
		t.Log("Aborting tests")
		t.FailNow()
	}
	if prog == nil {
		t.Fatalf("Program is nil")
	}
	eval := Evaluator{}
	env := object.NewEnv()
	res := eval.Evaluate(prog, env)
	if res == nil {
		t.Fatalf("Could not evaluate program")
	}
	if res.Inspect() != expected {
		t.Fatalf("Error: expected %s got %s", expected, res.Inspect())
	}
}
