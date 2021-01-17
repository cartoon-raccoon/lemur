package compiler

import (
	"github.com/cartoon-raccoon/lemur/ast"
	"github.com/cartoon-raccoon/lemur/code"
	"github.com/cartoon-raccoon/lemur/object"
)

// Compiler walks the AST and compiles it into bytecode
type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}

// New returns a new compiler struct
func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

// Compile is the main compiler function and does all the heavy lifting
func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.ExprStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
	case *ast.InfixExpr:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Right)
		if err != nil {
			return err
		}
	case *ast.Int:
		integer := &object.Integer{Value: node.Inner}
		c.emit(code.OpConstant, c.addConstant(integer))
	case *ast.Flt:
		float := &object.Float{Value: node.Inner}
		c.emit(code.OpConstant, c.addConstant(float))
	case *ast.Str:
		str := &object.String{Value: node.Inner}
		c.emit(code.OpConstant, c.addConstant(str))
	case *ast.Bool:
		boolean := &object.Boolean{Value: node.Inner}
		c.emit(code.OpConstant, c.addConstant(boolean))
	}
	return nil
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Encode(op, operands...)
	pos := c.addInstruction(ins)
	return pos
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInst := len(c.instructions)
	c.instructions = append(c.instructions, ins...)

	return posNewInst
}

// Bytecode returns the compiled bytecode from the compiler
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

// Bytecode stores the compiled bytecode emitted by the compiler
type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
