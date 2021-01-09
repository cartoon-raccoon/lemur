package eval

import (
	"testing"

	"github.com/cartoon-raccoon/monkey-jit/lexer"
	"github.com/cartoon-raccoon/monkey-jit/object"
	"github.com/cartoon-raccoon/monkey-jit/parser"
)

func TestExprStmtEval(t *testing.T) {
	tests := []struct {
		Input    string
		Expected object.Object
		Inspect  string
	}{
		{
			"1 + 1",
			&object.Integer{Value: 2},
			"2",
		},
		{
			"hello",
			&object.String{Value: "hello"},
			"hello",
		},
		{
			"420.69 + 7.4",
			&object.Float{Value: 428.09},
			"428.09",
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
		res := eval.Evaluate(prog)
		if res == nil {
			t.Fatalf("Test %d: Could not evaluate", i)
		}
		if res.Inspect() != test.Inspect {
			t.Errorf("Test %d: Expected %s, got %s", i, test.Inspect, res.Inspect())
		}
	}

}
