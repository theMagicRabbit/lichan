package main

import (
	"fmt"
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

func isValidMove(move Move, p piece, ) (isValid bool, err error) {

	return
}

func calculatePossibleMoves(p piece, gs *GameState) (squares []string, err error) {
	square := p.Square
	if !squareRE.MatchString(square) {
		err = fmt.Errorf("Not a valid square: %v\n", square)
		return
	}
	if _, ok := IsValidPieceType[p.PieceType]; !ok {
		err = fmt.Errorf("Not a valid piece: %v\n", p.PieceType)
		return
	}
	var calcFunc func(string)([]string)
	switch p.PieceType {
	case King:
		calcFunc = calcKingMoves
	case Queen:
		calcFunc = calcQueenMoves
	case Rook:
		calcFunc = calcRookMoves
	case Bishop:
		calcFunc = calcBishopMoves
	case Knight:
		calcFunc = calcKnightMoves
	case Pawn:
		calcFunc = calcPawnMoves
	}
	squares = calcFunc(square)
	return
}

func calcKingMoves(square string) (squares []string) {
	return
}

func calcQueenMoves(square string) (squares []string) {
	return
}

func calcRookMoves(square string) (squares []string) {
	file := string(square[0])
	rank := string(square[1])
	for _, r := range rankTokens {
		squareCandiate := file + r
		if squareCandiate == square {
			continue
		}
		squares = append(squares, squareCandiate)
	}
	for _, f := range fileTokens {
		squareCandiate := f + rank
		if squareCandiate == square {
			continue
		}
		squares = append(squares, squareCandiate)
	}
	return
}

func calcBishopMoves(square string) (squares []string) {
	return
}

func calcKnightMoves(square string) (squares []string) {
	return
}

func calcPawnMoves(square string) (squares []string) {
	return
}

