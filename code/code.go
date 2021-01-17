package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Instructions = the program
type Instructions []byte

// Opcode represents a single instruction carried out by the VM
type Opcode byte

const (
	// OpPush - retrieves a constant via its operand and pushes it
	OpPush Opcode = iota
	// OpPop - Pops the topmost element off the stack
	OpPop
	// OpAdd - Pops the top two values off the stack, adds them and pushes the result
	OpAdd
	// OpSub - Does the same thing as OpAdd, but performs subtraction
	OpSub
)

// Definition defines a single instruction - opcode and operand widths
type Definition struct {
	// The name of the opcode
	Name string
	// The length of the entire instruction
	TotalLength int
	// Lengths of all the operands
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpPush: {"OpPush", 3, []int{2}},
	OpPop:  {"OpPop", 1, []int{}},
	OpAdd:  {"OpAdd", 1, []int{}},
	OpSub:  {"OpSub", 1, []int{}},
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
		return nil
	}

	// calculate space needed for instruction
	instLen := 1
	for _, w := range def.OperandWidths {
		instLen += w
	}

	instruction := make([]byte, instLen)
	instruction[0] = byte(op)

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

func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0

	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			return out.String()
		}

		operands, read := ReadOperands(def, ins[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))

		i += 1 + read
	}

	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("Lengths do not match defined %d\n", operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	default:
		return fmt.Sprintf("Error: unhandled operand count for %s\n", def.Name)
	}

}

// ReadOperands disassembles a single instruction
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {

	offset := 0
	operands := make([]int, len(def.OperandWidths))

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		}

		offset += width
	}

	return operands, offset
}

// ReadUint16 reads the next two bytes in the instructions as a 16 bit number
func ReadUint16(ins []byte) uint16 {
	return binary.BigEndian.Uint16(ins)
}
