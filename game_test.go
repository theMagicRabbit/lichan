package main

import (
	"testing"
)

func TestParseMoveString(t *testing.T) {
	tests := []struct{
		Input string
		Expected Move
		E error
	}{
		{
			Input: "e4",
			Expected: Move{
				PieceType: Pawn,
				Target: "e4",
			},
			E: nil,
		},
		{
			Input: "O-O-O",
			Expected: Move{
				Target: "O-O-O",
				IsLongCastle: true,
			},
			E: nil,
		},
	}
	for _, test := range tests {
		result, err := ParseMoveString(test.Input)
		if err != test.E {
			t.Errorf("Error %v does not match expected: %v\n", err, test.E)
			break
		}
		if test.Expected != *result {
			t.Errorf("Result %v does not match expected: %v\n", *result, test.Expected)
		}
	}
}
