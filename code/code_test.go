package code

import "testing"

func TestEncode(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpPush, []int{65534}, []byte{byte(OpPush), 255, 254}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
		{OpSub, []int{}, []byte{byte(OpSub)}},
		{OpPop, []int{}, []byte{byte(OpPop)}},
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
		Encode(OpSub),
		Encode(OpAdd),
		Encode(OpPop),
		Encode(OpPush, 2),
		Encode(OpPush, 65535),
	}
	expected := `0000 OpSub
0001 OpAdd
0002 OpPop
0003 OpConstant 2
0006 OpConstant 65535
`

	concatted := Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}

	if str := concatted.String(); str != expected {
		t.Errorf("Wrong instructions, got: \n%s, expected: \n%s", str, expected)
	}
}
