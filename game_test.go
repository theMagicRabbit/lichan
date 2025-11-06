package main

import (
	"testing"
)

func TestMoveTokenizer(t *testing.T) {
	tests := []struct{
		Input struct{
			String string
			AtEnd bool
		}
		Expected struct{
			Advance int
			Token string
			Err error
		}
	}{
		{
			Input: struct{String string; AtEnd bool}{
				String: "O-O-O",
				AtEnd: false,
			},
			Expected: struct{Advance int; Token string; Err error}{
				Advance: 5,
				Token: "O-O-O",
				Err: nil,
			},
		},
	}
	for _, test := range tests {
		advance, token, err := tokenizerMoveString([]byte(test.Input.String), test.Input.AtEnd)
		if err != test.Expected.Err {
			t.Errorf("Error %v does not match expected: %v\n", err, test.Expected.Err)
			break
		}
		if advance != test.Expected.Advance {
			t.Errorf("Advance bytes: %v does not match expected: %v\n", advance, test.Expected.Advance)
			break
		}
		if string(token) != test.Expected.Token {
			t.Errorf("Result: %v does not match expected: %v\n", string(token), test.Expected.Token)
			break
		}
	}
}

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
