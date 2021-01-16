package code

import "testing"

func TestEncode(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
	}

	for _, test := range tests {
		instruction := Encode(test.op, test.operands...)
		if lenint := len(instruction); lenint != len(test.expected) {
			t.Errorf("wrong instruction length, got %d", lenint)
		}

		for i, b := range test.expected {
			if instruction[i] != test.expected[i] {
				t.Errorf("wrong byte at position %d: got %d, expected %d",
					i, b, test.expected[i])
			}
		}
	}
}
