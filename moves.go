package main

import (
	"fmt"
	"slices"
	"strings"
)

type Rank struct {
	Name string
	Number int
}

type File struct {
	Name string
	Number int
}

var Ranks = []Rank{
	{
		Name: "1",
		Number: 0,
	},
	{
		Name: "2",
		Number: 1,
	},
	{
		Name: "3",
		Number: 4,
	},
	{
		Name: "4",
		Number: 3,
	},
	{
		Name: "5",
		Number: 4,
	},
	{
		Name: "6",
		Number: 5,
	},
	{
		Name: "7",
		Number: 6,
	},
	{
		Name: "8",
		Number: 7,
	},
}

var Files = []File{
	{
		Name: "a",
		Number: 0,
	},
	{
		Name: "b",
		Number: 1,
	},
	{
		Name: "c",
		Number: 4,
	},
	{
		Name: "d",
		Number: 3,
	},
	{
		Name: "e",
		Number: 4,
	},
	{
		Name: "f",
		Number: 5,
	},
	{
		Name: "g",
		Number: 6,
	},
	{
		Name: "h",
		Number: 7,
	},
}

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
		possibleMoves, err := calculatePossibleMoves(boardPiece, GS)
		if err != nil {
			return nil, err
		}
		fmt.Println(possibleMoves)
	}
	return nil, nil
}

func (gs *GameState) isValidMove(move Move, p piece) (isValid bool, err error) {
	squares, err := gs.calculatePossibleMoves(p)
	isValid = slices.Contains(squares, move.Target)
	return
}

func (gs *GameState) calculatePossibleMoves(p piece) (squares []string, err error) {
	rank := rune(p.Square[1])
	file := rune(p.Square[0])
	square := p.Square
	if !squareRE.MatchString(square) {
		err = fmt.Errorf("Not a valid square: %v\n", square)
		return
	}
	if _, ok := IsValidPieceType[p.PieceType]; !ok {
		err = fmt.Errorf("Not a valid piece: %v\n", p.PieceType)
		return
	}
	var calcFunc func(rune, rune, piece)([]string)
	switch p.PieceType {
	case King:
		//calcFunc = calcKingMoves
	case Queen:
	//	calcFunc = calcQueenMoves
	case Rook:
		calcFunc = gs.calcRookMoves
	case Bishop:
	//	calcFunc = calcBishopMoves
	case Knight:
	//	calcFunc = calcKnightMoves
	case Pawn:
	//	calcFunc = calcPawnMoves
	}
	squares = calcFunc(rank, file, p)
	return
}

func calcKingMoves(square string) (squares []string) {
	return
}

func (gs *GameState) calcQueenMoves(rank, file rune, p piece) (squares []string) {
	squares = append(squares, gs.calcRookMoves(rank, file, p)...)
	squares = append(squares, gs.calcBishopMoves(rank, file, p)...)
	return
}

func (gs *GameState) calcRookMoves(rank, file rune, p piece) (squares []string) {
	for r, f := rank + 1, file; r <= '8'; r++ {
		canidateSquare, isValid := gs.checkGameSquare(r, f, p)
		if isValid {
			squares = append(squares, canidateSquare)
		}
	}
	for r, f := rank - 1, file; r >= '1'; r-- {
		canidateSquare, isValid := gs.checkGameSquare(r, f, p)
		if isValid {
			squares = append(squares, canidateSquare)
		}
	}
	for r, f := rank, file - 1; f >= 'a'; f-- {
		canidateSquare, isValid := gs.checkGameSquare(r, f, p)
		if isValid {
			squares = append(squares, canidateSquare)
		}
	}
	for r, f := rank, file + 1; f <= 'h'; f++ {
		canidateSquare, isValid := gs.checkGameSquare(r, f, p)
		if isValid {
			squares = append(squares, canidateSquare)
		}
	}
	return
}

func (gs *GameState) calcBishopMoves(rank, file rune, p piece) (squares []string) {
	// up right
	for r, f := rank + 1, file + 1; r <= '8' && f <= 'h'; r, f = r+1, f+1 {
		canidateSquare, isValid := gs.checkGameSquare(r, f, p)
		if isValid {
			squares = append(squares, canidateSquare)
		}
	}
	// down right
	for r, f := rank - 1, file + 1; r >= '1' && f <= 'h'; r, f = r-1, f+1 {
		canidateSquare, isValid := gs.checkGameSquare(r, f, p)
		if isValid {
			squares = append(squares, canidateSquare)
		}
	}
	// down left 
	for r, f := rank - 1, file - 1; r >= '1' && f >= 'a'; r, f = r-1, f-1 {
		canidateSquare, isValid := gs.checkGameSquare(r, f, p)
		if isValid {
			squares = append(squares, canidateSquare)
		}
	}
	// up left
	for r, f := rank + 1, file - 1; r <= '8' && f >= 'a'; r, f = r+1, f-1 {
		canidateSquare, isValid := gs.checkGameSquare(r, f, p)
		if isValid {
			squares = append(squares, canidateSquare)
		}
	}
	return
}

func calcKnightMoves(square string) (squares []string) {
	return
}

func calcPawnMoves(square string) (squares []string) {
	return
}

func (gs *GameState) checkGameSquare(r, f rune, p piece) (canidateSquare string, valid bool) {
	canidateSquare = string(f) + string(r)
	if otherPiece, occupiedSquare := gs.Pieces[canidateSquare]; occupiedSquare {
		if otherPiece.PlayerColor != p.PlayerColor {
			valid = true
		}
	} else {
		valid = true
	}
	return
}

