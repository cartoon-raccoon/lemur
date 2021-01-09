package object

import (
	"fmt"
	"strings"
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

	//IDENT - Identifier
	IDENT = "IDENT_OBJ"

	//PRES - Program result
	PRES = "PROG_RES"
)

// Object represents any object returnable from an expression
type Object interface {
	Type() string
	Inspect() string
	Display()
}

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

// ProgramResult is the results returned by a program
type ProgramResult struct {
	Results []Object
}

// Type implements Object for ProgramResult
func (pr *ProgramResult) Type() string { return PRES }

// Inspect implements Object for ProgramResult
func (pr *ProgramResult) Inspect() string {
	pres := []string{}
	for _, res := range pr.Results {
		if !isNull(res) {
			pres = append(pres, res.Inspect())
		}
	}
	return strings.Join(pres, "\n")
}

// Display implements Object for ProgramResult
func (pr *ProgramResult) Display() {
	inspect := pr.Inspect()
	if len(strings.TrimSpace(inspect)) != 0 {
		fmt.Printf("%s\n", inspect)
	}
}

func isNull(o Object) bool {
	switch o.(type) {
	case *Null:
		return true
	default:
		return false
	}
}
