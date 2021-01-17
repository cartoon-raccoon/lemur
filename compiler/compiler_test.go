package compiler

import (
	"fmt"
	"testing"

	"github.com/cartoon-raccoon/lemur/code"
	"github.com/cartoon-raccoon/lemur/lexer"
	"github.com/cartoon-raccoon/lemur/object"
	"github.com/cartoon-raccoon/lemur/parser"
)

type compilerTestCase struct {
	input             string
	expectedConstants []interface{}
	expectedInsts     []code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1 + 2",
			expectedConstants: []interface{}{1, 2},
			expectedInsts: []code.Instructions{
				code.Encode(code.OpConstant, 0),
				code.Encode(code.OpConstant, 1),
			},
		},
	}

	runCompilerTests(t, tests)
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, test := range tests {
		l := lexer.New(test.input)
		p, err := parser.New(l)
		if err != nil {
			t.Errorf("Error in parsing: %s", err)
			continue
		}
		prog := p.Parse()
		if errors := p.CheckErrors(); errors != nil {
			t.Errorf("Errors during parsing:")
			for _, e := range errors {
				t.Logf("%s", e)
			}
			continue
		}

		compiler := New()
		err = compiler.Compile(prog)
		if err != nil {
			t.Fatalf("Error while compiling: %s", err)
		}

		bytecode := compiler.Bytecode()

		err = testInstructions(test.expectedInsts, bytecode.Instructions)
		if err != nil {
			t.Errorf("Error in instructions: %s", err)
		}
	}
}

func testInstructions(expected []code.Instructions, actual code.Instructions) error {
	concatted := concatInstructions(expected)

	if len(actual) != len(concatted) {
		return fmt.Errorf("wrong instructions length, got %q, expected %q",
			actual,
			concatted,
		)
	}

	for i, ins := range concatted {
		if actual[i] != ins {
			return fmt.Errorf("wrong instruction at %d, expected %q, got %q",
				i, concatted, actual)
		}
	}
	return nil
}

func testConstants(
	t *testing.T,
	expected []interface{},
	actual []object.Object,
) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants, got %d, expected %d",
			len(actual), len(expected))
	}

	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s", i, err)
			}
		}
	}

	return nil
}

func concatInstructions(s []code.Instructions) code.Instructions {
	out := code.Instructions{}

	for _, ins := range s {
		out = append(out, ins...)
	}

	return out
}

func testIntegerObject(constant int64, obj object.Object) error {
	result, ok := obj.(*object.Integer)
	if !ok {
		return fmt.Errorf("Expected integer, got %T", obj)
	}

	if result.Value != constant {
		return fmt.Errorf("Values do not equate: got %d, expected %d",
			result.Value, constant)
	}

	return nil
}
