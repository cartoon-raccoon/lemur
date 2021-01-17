package vm

import (
	"fmt"
	"testing"

	"github.com/cartoon-raccoon/lemur/compiler"
	"github.com/cartoon-raccoon/lemur/lexer"
	"github.com/cartoon-raccoon/lemur/object"
	"github.com/cartoon-raccoon/lemur/parser"
)

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

type vmTestCase struct {
	input    string
	expected interface{}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3}, //fixme
	}

	runVMTests(t, tests)
}

func runVMTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p, err := parser.New(l)
		if err != nil {

		}

		prog := p.Parse()

		comp := compiler.New()

		err = comp.Compile(prog)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New()

		err = vm.Run(comp.Bytecode())
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		stackElem := vm.LastPopped()

		testExpectedObject(t, tt.expected, stackElem)
	}
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testintobj: %s", err)
		}
	}
}
