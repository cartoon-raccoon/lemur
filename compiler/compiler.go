package compiler

import (
	"fmt"

	"github.com/cartoon-raccoon/lemur/ast"
	"github.com/cartoon-raccoon/lemur/code"
	"github.com/cartoon-raccoon/lemur/lexer"
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
		c.emit(code.OpPop)
	case *ast.PrefixExpr:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}
		switch node.Operator {
		case lexer.BANG:
			c.emit(code.OpBang)
		case lexer.SUB:
			c.emit(code.OpMinus)
		case lexer.BWNOT:
			c.emit(code.OpBWNOT)
		default:
			return fmt.Errorf("Unknown operator %s", node.Operator)
		}
	case *ast.InfixExpr:
		if node.Operator == lexer.LT || node.Operator == lexer.LE {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}
			err = c.Compile(node.Left)
			if err != nil {
				return err
			}
			if node.Operator == lexer.LT {
				c.emit(code.OpGT)
			} else {
				c.emit(code.OpGE)
			}
			return nil
		}
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Right)
		if err != nil {
			return err
		}
		switch node.Operator {
		case lexer.ADD:
			c.emit(code.OpAdd)
		case lexer.SUB:
			c.emit(code.OpSub)
		case lexer.MUL:
			c.emit(code.OpMul)
		case lexer.DIV:
			c.emit(code.OpDiv)
		case lexer.BWAND:
			c.emit(code.OpBWAnd)
		case lexer.BWOR:
			c.emit(code.OpBWOr)
		case lexer.BWNOT:
			c.emit(code.OpBWXOR)
		case lexer.EQ:
			c.emit(code.OpEq)
		case lexer.NE:
			c.emit(code.OpNE)
		case lexer.GT:
			c.emit(code.OpGT)
		case lexer.GE:
			c.emit(code.OpGE)
		default:
			return fmt.Errorf("unknown operator: %s", node.Operator)
		}
	case *ast.Int:
		integer := &object.Integer{Value: node.Inner}
		c.emit(code.OpPush, c.addConstant(integer))
	case *ast.Flt:
		float := &object.Float{Value: node.Inner}
		c.emit(code.OpPush, c.addConstant(float))
	case *ast.Str:
		str := &object.String{Value: node.Inner}
		c.emit(code.OpPush, c.addConstant(str))
	case *ast.Bool:
		if node.Inner {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
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
