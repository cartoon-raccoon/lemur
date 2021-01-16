package code

import (
	"encoding/binary"
	"fmt"
)

// Instructions = the program
type Instructions []byte

// Opcode represents a single instruction carried out by the VM
type Opcode byte

// OpConstant represents the opcodes
const (
	OpConstant Opcode = iota
)

// Definition defines a single instruction - opcode and operand widths
type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
}

// Lookup gets the definition of an Opcode
func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}

// Encode produces a raw instruction
func Encode(op Opcode, operands ...int) []byte {
	// lookup instruction definition
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	// calculate space needed for instruction
	instLen := 1
	for _, w := range def.OperandWidths {
		instLen += w
	}

	instruction := make([]byte, instLen)

	// putting the actual values
	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}

	return instruction
}
