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
	return nil
}

// Bytecode returns the compiled bytecode from the compiler
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: code.Instructions{},
		Constants:    []object.Object{},
	}
}

// Bytecode stores the compiled bytecode emitted by the compiler
type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
