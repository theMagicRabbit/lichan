package main

var fenBoardOrder = [64]string{
	"a8", "b8", "c8", "d8", "e8", "f8", "g8", "h8",
	"a7", "b7", "c7", "d7", "e7", "f7", "g7", "h7",
	"a6", "b6", "c6", "d6", "e6", "f6", "g6", "h6",
	"a5", "b5", "c5", "d5", "e5", "f5", "g5", "h5",
	"a4", "b4", "c4", "d4", "e4", "f4", "g4", "h4",
	"a3", "b3", "c3", "d3", "e3", "f3", "g3", "h3",
	"a2", "b2", "c2", "d2", "e2", "f2", "g2", "h2",
	"a1", "b1", "c1", "d1", "e1", "f1", "g1", "h1",
}

func initalGameState() (*GameState) {
	gs := &GameState{
		PlayerTurn: White,
		Pieces: map[string]piece{
			"a1": {
				PieceType: Rook,
				Square: "a1",
				PlayerColor: White,
			},
			"b1": {
				PieceType: Knight,
				Square: "b1",
				PlayerColor: White,
			},
			"c1": {
				PieceType: Bishop,
				Square: "c1",
				PlayerColor: White,
			},
			"d1": {
				PieceType: Queen,
				Square: "d1",
				PlayerColor: White,
			},
			"e1": {
				PieceType: King,
				Square: "e1",
				PlayerColor: White,
			},
			"f1": {
				PieceType: Bishop,
				Square: "f1",
				PlayerColor: White,
			},
			"g1": {
				PieceType: Knight,
				Square: "g1",
				PlayerColor: White,
			},
			"h1": {
				PieceType: Rook,
				Square: "h1",
				PlayerColor: White,
			},
			"a2": {
				PieceType: Pawn,
				Square: "a2",
				PlayerColor: White,
			},
			"b2": {
				PieceType: Pawn,
				Square: "b2",
				PlayerColor: White,
				},
			"c2": {
				PieceType: Pawn,
				Square: "c2",
				PlayerColor: White,
			},
			"d2": {
				PieceType: Pawn,
				Square: "d2",
				PlayerColor: White,
			},
			"e2": {
				PieceType: Pawn,
				Square: "e2",
				PlayerColor: White,
			},
			"f2": {
				PieceType: Pawn,
				Square: "f2",
				PlayerColor: White,
			},
			"g2": {
				PieceType: Pawn,
				Square: "g2",
				PlayerColor: White,
			},
			"h2": {
				PieceType: Pawn,
				Square: "h2",
				PlayerColor: White,
			},
			"a7": {
				PieceType: Pawn,
				Square: "a7",
				PlayerColor: Black,
			},
			"b7": {
				PieceType: Pawn,
				Square: "b7",
				PlayerColor: Black,
			},
			"c7": {
				PieceType: Pawn,
				Square: "c7",
				PlayerColor: Black,
			},
			"d7": {
				PieceType: Pawn,
				Square: "d7",
				PlayerColor: Black,
			},
			"e7": {
				PieceType: Pawn,
				Square: "e7",
				PlayerColor: Black,
			},
			"f7": {
				PieceType: Pawn,
				Square: "f7",
				PlayerColor: Black,
			},
			"g7": {
				PieceType: Pawn,
				Square: "g7",
				PlayerColor: Black,
			},
			"h7": {
				PieceType: Pawn,
				Square: "h7",
				PlayerColor: Black,
			},
			"a8": {
				PieceType: Rook,
				Square: "a8",
				PlayerColor: Black,
			},
			"b8": {
				PieceType: Knight,
				Square: "b8",
				PlayerColor: Black,
			},
			"c8": {
				PieceType: Bishop,
				Square: "c8",
				PlayerColor: Black,
			},
			"d8": {
				PieceType: Queen,
				Square: "d8",
				PlayerColor: Black,
			},
			"e8": {
				PieceType: King,
				Square: "e8",
				PlayerColor: Black,
			},
			"f8": {
				PieceType: Bishop,
				Square: "f8",
				PlayerColor: Black,
			},
			"g8": {
				PieceType: Knight,
				Square: "g8",
				PlayerColor: Black,
			},
			"h8": {
				PieceType: Rook,
				Square: "h8",
				PlayerColor: Black,
			},
		},
	}
	return gs
}
