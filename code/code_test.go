package code

import "testing"

func TestEncode(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
	}

	for idx, test := range tests {
		instruction := Encode(test.op, test.operands...)
		if lenint := len(instruction); lenint != len(test.expected) {
			t.Errorf("wrong instruction length, got %d", lenint)
		}

		for i := range test.expected {
			if instruction[i] != test.expected[i] {
				t.Errorf("Test %d: wrong byte at position %d: got %d, expected %d",
					idx, i, instruction[i], test.expected[i])
			}
		}
	}
}

func TestInstructionsString(t *testing.T) {
	instructions := []Instructions{
		Encode(OpAdd),
		Encode(OpConstant, 2),
		Encode(OpConstant, 65535),
	}
	expected := `0000 OpAdd
0001 OpConstant 2
0004 OpConstant 65535
`

	concatted := Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}

	if str := concatted.String(); str != expected {
		t.Errorf("Wrong instructions, got: \n%s, expected: \n%s", str, expected)
	}
}
