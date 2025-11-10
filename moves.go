package main

import "strings"


func (GS *GameState) ApplyMove(ms string, turn PlayerColor) (*GameState, error) {
	move, err := ParseMoveString(strings.TrimSpace(ms))
	if err != nil {
		return nil, err
	}

	for _, boardPiece := range GS.Pieces {
		if boardPiece.PlayerColor != turn {
			continue
		}
		if boardPiece.PieceType != move.PieceType {
			continue
		}
	}
	return nil, nil
}
