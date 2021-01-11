package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/cartoon-raccoon/monkey-jit/ast"
	"github.com/cartoon-raccoon/monkey-jit/lexer"
)

const (
	//INTEGER - Integer
	INTEGER = "INT_OBJ"
	//FLOAT - Float
	FLOAT = "FLT_OBJ"
	//STRING - String
	STRING = "STR_OBJ"
	//BOOLEAN - Boolean
	BOOLEAN = "BOOL_OBJ"
	//NULL - Null value
	NULL = "NULL_OBJ"
	//RETURN - Return value of a block
	RETURN = "RETURN_OBJ"
	//FUNCTION - Function object
	FUNCTION = "FUNC_OBJ"
	//BUILTIN - Builtin function
	BUILTIN = "BUILTIN_OBJ"

	//ERROR - Error object
	ERROR = "ERROR_OBJ"

	//IDENT - Identifier
	IDENT = "IDENT_OBJ"

	//PRES - Program result
	PRES = "PROG_RES"

	//ENVIRONMENT - Environment
	ENVIRONMENT = "ENV_OBJ"
)

// Object represents any object returnable from an expression
type Object interface {
	Type() string
	Inspect() string
	Display()
}

// Environment represents the execution environment
type Environment struct {
	Data  map[string]Object
	Outer *Environment
}

// NewEnv - Returns a new fresh environment
func NewEnv() *Environment {
	env := &Environment{}
	env.Data = make(map[string]Object)

	return env
}

// NewEnclosedEnv nests the provided environment inside a new one
// It mimics a new stack frame
func NewEnclosedEnv(outer *Environment) *Environment {
	env := NewEnv()
	env.Outer = outer

	return env
}

// Get recursively gets a variable from the environment and all its outer envs
func (env *Environment) Get(ident string) (Object, bool) {
	obj, ok := env.Data[ident]
	if !ok && env.Outer != nil {
		obj, ok = env.Outer.Get(ident)
	}
	return obj, ok
}

// Set adds a variable to the environment
func (env *Environment) Set(ident string, val Object) Object {
	env.Data[ident] = val
	return val
}

// Type implements Object for Environment
func (env *Environment) Type() string { return ENVIRONMENT }

// Inspect implements Object for Environment
func (env *Environment) Inspect() string {
	//todo
	return "Environment"
}

// Display implements Object for Environment
func (env *Environment) Display() {}

// Integer represents an integer
type Integer struct {
	Value int64
}

// Type implements Object for Int
func (i *Integer) Type() string { return INTEGER }

// Inspect implements Object for Int
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// Display implements Object for Int
func (i *Integer) Display() {
	fmt.Printf("%d\n", i.Value)
}

// Float represents a Float
type Float struct {
	Value float64
}

// Type implements Object for Float
func (f *Float) Type() string { return FLOAT }

// Inspect implements Object for Float
func (f *Float) Inspect() string {
	return fmt.Sprintf("%f", f.Value)
}

// Display implements Object for Float
func (f *Float) Display() {
	fmt.Printf("%f\n", f.Value)
}

// String represents a string
type String struct {
	Value string
}

// Type implements Object for String
func (s *String) Type() string { return STRING }

// Inspect implements Object for String
func (s *String) Inspect() string {
	return s.Value
}

// Display implements Object for Float
func (s *String) Display() {
	fmt.Printf("%s\n", s.Value)
}

// Boolean represents a bool
type Boolean struct {
	Value bool
}

// Type implements Object for Boolean
func (b *Boolean) Type() string { return BOOLEAN }

// Inspect implements Object for Boolean
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

// Display implements Object for Boolean
func (b *Boolean) Display() {
	fmt.Printf("%t\n", b.Value)
}

// Null represents a null value
type Null struct{}

// Type implements Object for Identifier
func (n *Null) Type() string { return NULL }

// Inspect implements Object for Identifier
func (n *Null) Inspect() string {
	return "Null"
}

// Display implements Object for Null
func (n *Null) Display() {}

// Return - A wrapper type for a value returned by a return statement
type Return struct {
	Inner Object
}

// Type implements Object for Return
func (r *Return) Type() string { return RETURN }

// Inspect implements Object for Return
func (r *Return) Inspect() string {
	return fmt.Sprintf("%s", r.Inner.Inspect())
}

// Display implements Object for Return
func (r *Return) Display() {
	r.Inner.Display()
}

// Function represents a function in the environment
type Function struct {
	Params []*ast.Identifier
	Body   *ast.BlockStatement
	Env    *Environment
}

// Type implements Object for Function
func (f *Function) Type() string { return FUNCTION }

// Inspect implements Object for Function
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}

	for _, param := range f.Params {
		params = append(params, param.String())
	}

	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(f.Body.String())

	return out.String()

}

// Display implements Object for Function
func (f *Function) Display() {
	fmt.Printf("%s\n", f.Inspect())
}

// BuiltinFn is a function that is implemented within the interpreter itself
type BuiltinFn func(ctxt lexer.Context, args ...Object) Object

type Builtin struct {
	Fn BuiltinFn
}

func (b *Builtin) Type() string { return BUILTIN }
func (b *Builtin) Inspect() string {
	return "builtin function"
}

func (b *Builtin) Display() {}

// StmtResults is the results returned by a program
type StmtResults struct {
	Results []Object
}

// Type implements Object for StmtResults
func (pr *StmtResults) Type() string { return PRES }

// Inspect implements Object for StmtResults
func (pr *StmtResults) Inspect() string {
	pres := []string{}
	for _, res := range pr.Results {
		if !IsNull(res) {
			pres = append(pres, res.Inspect())
		}
	}
	return strings.Join(pres, "\n")
}

// Display implements Object for StmtResults
func (pr *StmtResults) Display() {
	inspect := pr.Inspect()
	if len(strings.TrimSpace(inspect)) != 0 {
		fmt.Printf("%s\n", inspect)
	}
}

// Exception - an error type to return
type Exception struct {
	Msg string
	Con lexer.Context
}

// Type implements Object for Exception
func (ex *Exception) Type() string { return ERROR }

// Inspect implements Object for Exception
func (ex *Exception) Inspect() string {
	return fmt.Sprintf("%s - Line %d, Col %d", ex.Msg, ex.Con.Line, ex.Con.Col)
}

// Display implements Object for Exception
func (ex *Exception) Display() {
	fmt.Printf(ex.Inspect())
}

// IsNull checks whether a result is Null
func IsNull(o Object) bool {
	switch o.(type) {
	case *Null:
		return true
	default:
		return false
	}
}

// IsErr checks whether a result is Exception
func IsErr(o Object) bool {
	switch o.(type) {
	case *Exception:
		return true
	default:
		return false
	}
}
