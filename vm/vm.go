package vm

import (
	"fmt"

	"github.com/cartoon-raccoon/lemur/code"
	"github.com/cartoon-raccoon/lemur/compiler"
	"github.com/cartoon-raccoon/lemur/eval"
	"github.com/cartoon-raccoon/lemur/lexer"
	"github.com/cartoon-raccoon/lemur/object"
)

// StackSize is the maximum size the stack can take
const StackSize = 2048

// True - an invariant true object
var True = &object.Boolean{Value: true}

// False - an invariant false object
var False = &object.Boolean{Value: false}

// VM represents the virtual machine
type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int // stack pointer. Top of stack is stack[sp - 1]

	ip int //instruction pointer
}

// New returns a new VM
func New() *VM {
	return &VM{
		stack: make([]object.Object, StackSize),
		sp:    0,
	}
}

// LastPopped returns the item just popped from the stack
func (vm *VM) LastPopped() object.Object {
	return vm.stack[vm.sp]
}

// Run executes the code that is given to it via the VM
func (vm *VM) Run(bc *compiler.Bytecode) error {
	vm.instructions = bc.Instructions
	vm.constants = bc.Constants
	vm.sp = 0
	for vm.ip = 0; vm.ip < len(vm.instructions); vm.ip++ {
		op := code.Opcode(vm.instructions[vm.ip])

		switch op {
		case code.OpPush:
			constIndex := code.ReadUint16(vm.instructions[vm.ip+1:])
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
			vm.ip += 2

		case code.OpPop:
			vm.pop()

		case code.OpAdd:
			right, err := vm.pop()
			if err != nil {
				return err
			}
			left, err := vm.pop()
			if err != nil {
				return err
			}

			result := eval.EvaluateSides(left, right, "+", lexer.Context{})
			vm.push(result)

		case code.OpSub:
			right, err := vm.pop()
			if err != nil {
				return err
			}
			left, err := vm.pop()
			if err != nil {
				return err
			}

			result := eval.EvaluateSides(left, right, "-", lexer.Context{})
			vm.push(result)

		case code.OpMul:
			right, err := vm.pop()
			if err != nil {
				return err
			}
			left, err := vm.pop()
			if err != nil {
				return err
			}

			result := eval.EvaluateSides(left, right, "*", lexer.Context{})
			vm.push(result)

		case code.OpDiv:
			right, err := vm.pop()
			if err != nil {
				return err
			}
			left, err := vm.pop()
			if err != nil {
				return err
			}

			result := eval.EvaluateSides(left, right, "/", lexer.Context{})
			vm.push(result)

		case code.OpBWAnd:
			right, err := vm.pop()
			if err != nil {
				return err
			}
			left, err := vm.pop()
			if err != nil {
				return err
			}

			result := eval.EvaluateSides(left, right, "&", lexer.Context{})
			vm.push(result)

		case code.OpBWOr:
			right, err := vm.pop()
			if err != nil {
				return err
			}
			left, err := vm.pop()
			if err != nil {
				return err
			}

			result := eval.EvaluateSides(left, right, "|", lexer.Context{})
			vm.push(result)
		case code.OpBWXOR:
			right, err := vm.pop()
			if err != nil {
				return err
			}
			left, err := vm.pop()
			if err != nil {
				return err
			}

			result := eval.EvaluateSides(left, right, "^", lexer.Context{})
			vm.push(result)
		case code.OpEq:
			right, err := vm.pop()
			if err != nil {
				return err
			}
			left, err := vm.pop()
			if err != nil {
				return err
			}

			result := eval.EvaluateComp(left, right, "==", lexer.Context{})
			vm.push(result)
		case code.OpNE:
			right, err := vm.pop()
			if err != nil {
				return err
			}
			left, err := vm.pop()
			if err != nil {
				return err
			}

			result := eval.EvaluateComp(left, right, "!=", lexer.Context{})
			vm.push(result)
		case code.OpGT:
			right, err := vm.pop()
			if err != nil {
				return err
			}
			left, err := vm.pop()
			if err != nil {
				return err
			}

			result := eval.EvaluateComp(left, right, ">", lexer.Context{})
			vm.push(result)
		case code.OpGE:
			right, err := vm.pop()
			if err != nil {
				return err
			}
			left, err := vm.pop()
			if err != nil {
				return err
			}

			result := eval.EvaluateComp(left, right, ">=", lexer.Context{})
			vm.push(result)
		case code.OpTrue:
			vm.push(eval.TRUE)
		case code.OpFalse:
			vm.push(eval.FALSE)
		case code.OpBang:
			op, err := vm.pop()
			if err != nil {
				return err
			}
			vm.push(vm.evalPrefixBang(op))
		case code.OpMinus:
			op, err := vm.pop()
			if err != nil {
				return err
			}
			vm.push(vm.evalPrefixMinus(op))
		case code.OpBWNOT:
			op, err := vm.pop()
			if err != nil {
				return err
			}
			obj, ok := op.(*object.Integer)
			if !ok {
				vm.push(&object.Exception{
					Msg: "Cannot use bitwise NOT on non-integer\n",
					Con: lexer.Context{},
				})
			} else {
				vm.push(&object.Integer{Value: ^obj.Value})
			}
		}
	}
	return nil
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) pop() (object.Object, error) {
	if vm.sp == 0 {
		return nil, fmt.Errorf("stack underflow")
	}

	o := vm.stack[vm.sp-1]
	vm.sp--
	return o, nil
}

func (vm *VM) evalPrefixBang(op object.Object) object.Object {
	if op != eval.TRUE || op != eval.FALSE {
		return &object.Exception{
			Msg: "Cannot use bang on operand\n",
			Con: lexer.Context{},
		}
	}
	if op == eval.TRUE {
		return eval.FALSE
	}
	return eval.TRUE
}
func (vm *VM) evalPrefixMinus(op object.Object) object.Object {
	switch obj := op.(type) {
	case *object.Integer:
		return &object.Integer{Value: -obj.Value}
	case *object.Float:
		return &object.Float{Value: -obj.Value}
	default:
		return &object.Exception{
			Msg: "Cannot use - on operand\n",
			Con: lexer.Context{},
		}
	}
}
