package main

import "testing"

func TestGameState(T *testing.T) {
	tests := []struct{
		Input *GameState
		Result *GameState
		Move string
		ResultString string
	} {
		{
			Input: &GameState{
				PlayerTurn: Black,
				Pieces: map[string]piece{
					"e8": {
						PieceType: King,
						PlayerColor: Black,
						Square: "e8",
					},
					"g8": {
						PieceType: Knight,
						PlayerColor: Black,
						Square: "g8",
					},
					"c6": {
						PieceType: Knight,
						PlayerColor: Black,
						Square: "c6",
					},
					"b5": {
						PieceType: Bishop,
						PlayerColor: White,
						Square: "b5",
					},
					"e1": {
						PieceType: King,
						PlayerColor: White,
						Square: "e1",
					},
				},
			},
			Result: &GameState{
				PlayerTurn: White,
				Pieces: map[string]piece{
					"e8": {
						PieceType: King,
						PlayerColor: Black,
						Square: "e8",
					},
					"e7": {
						PieceType: Knight,
						PlayerColor: Black,
						Square: "e7",
					},
					"c6": {
						PieceType: Knight,
						PlayerColor: Black,
						Square: "c6",
					},
					"b5": {
						PieceType: Bishop,
						PlayerColor: White,
						Square: "b5",
					},
					"e1": {
						PieceType: King,
						PlayerColor: White,
						Square: "e1",
					},
				},
			},
			Move: "Ne7",
			ResultString: "g8e7",
		},
		
	}

	for _, test := range tests {
		result, ems, err := test.Input.ApplyAndTranslateMove(test.Move, test.Input.PlayerTurn)
		if err != nil {
			T.Errorf("Unexpected error: %s\n", err.Error())
			continue
		}
		if ems != test.ResultString {
			T.Errorf("Result string %s does not match expected: %s\n", ems, test.ResultString)
		}

		if result != test.Result {
			T.Errorf("Result state %v does not match expected: %v\n", result, test.Result)
		}
	}
}
