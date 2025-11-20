package main

import (
	"fmt"
	"slices"
	"strings"
)

func (gs *GameState) ApplyAndTranslateMove(ms string, turn PlayerColor) (*GameState, string, error) {
	move, err := ParseMoveString(strings.TrimSpace(ms))
	if err != nil {
		return nil, "", err
	}

	if move.IsLongCastle && turn == Black {
		move.Target = "c8"
	}
	if move.IsLongCastle && turn == White {
		move.Target = "c1"
	}
	if move.IsShortCastle && turn == Black {
		move.Target = "g8"
	}
	if move.IsShortCastle && turn == White {
		move.Target = "g1"
	}

	var sourceSquare string
	if len(move.Discriminator) == 2 {
		sourceSquare = move.Discriminator
	} else {
		for _, boardPiece := range gs.Pieces {
			if !strings.Contains(boardPiece.Square, move.Discriminator) {
				continue
			}
			if boardPiece.PlayerColor != turn {
				continue
			}
			if boardPiece.PieceType != move.PieceType {
				continue
			}
			if isValid, err := gs.isValidMove(move, boardPiece); err != nil {
				return nil, "", err
			} else if !isValid {
				continue
			}
			sourceSquare = boardPiece.Square
			break
		}
	}

	var nextTurn PlayerColor
	if gs.PlayerTurn == White {
		nextTurn = Black
	} else {
		nextTurn = White
	}

	nextState := GameState{
		Pieces:     gs.Pieces,
		PlayerTurn: nextTurn,
	}
	movedPiece := nextState.Pieces[sourceSquare]
	movedPiece.Square = move.Target

	var promoteTo string
	if move.PromoteTo != "" {
		movedPiece.PieceType = move.PromoteTo
		switch move.PromoteTo {
		case Queen:
			promoteTo = "q"
		case Rook:
			promoteTo = "r"
		case Bishop:
			promoteTo = "b"
		case Knight:
			promoteTo = "n"
		}
	}

	if move.IsCapture {
		if _, targetExists := nextState.Pieces[move.Target]; targetExists {
			delete(nextState.Pieces, move.Target)
		} else if movedPiece.PieceType != Pawn {
			return nil, "", fmt.Errorf("No piece found on target square: %v\n", move.Target)
		} else {
			targetRank := string(sourceSquare[1])
			if !(turn == Black && "4" == targetRank) && !(turn == White && "5" == targetRank) {
				return nil, "", fmt.Errorf("Invalid capture to %v attempted\n", move.Target)
			}

			targetFile := move.Target[0]
			enPassantSquare := string(targetFile) + string(targetRank)
			if targetPawn, exists := nextState.Pieces[enPassantSquare]; exists &&
			targetPawn.PieceType == Pawn &&
			targetPawn.PlayerColor != movedPiece.PlayerColor {
				delete(nextState.Pieces, enPassantSquare)
			} else {
				return nil, "", fmt.Errorf("Invalid capture to %v attempted\n", move.Target)
			}
		}
	}

	if move.IsLongCastle {
		var rookSource string
		var rookDest string
		if turn == Black {
			rookSource = "a8"
			rookDest = "d8"
		} else {
			rookSource = "a1"
			rookDest = "d1"
		}
		qRook := nextState.Pieces[rookSource]
		qRook.Square = rookDest
		delete(nextState.Pieces, rookSource)
		nextState.Pieces[rookDest] = qRook
	}

	if move.IsShortCastle {
		var rookSource string
		var rookDest string
		if turn == Black {
			rookSource = "h8"
			rookDest = "f8"
		} else {
			rookSource = "h1"
			rookDest = "f1"
		}
		kRook := nextState.Pieces[rookSource]
		kRook.Square = rookDest
		delete(nextState.Pieces, rookSource)
		nextState.Pieces[rookDest] = kRook
	}

	nextState.Pieces[move.Target] = movedPiece
	delete(nextState.Pieces, sourceSquare)

	extendedMove := sourceSquare + move.Target + promoteTo
	if moveLen := len(extendedMove); !(moveLen == 4 || moveLen == 5) {
		return nil, extendedMove, fmt.Errorf("Stockfish move is wrong length. source: %s; dest: %s\n", sourceSquare, move.Target)
	}

	return &nextState, extendedMove, nil
}

func (gs *GameState) isValidMove(move *Move, p piece) (isValid bool, err error) {
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
	var calcFunc func(rune, rune, piece) []string
	switch p.PieceType {
	case King:
		calcFunc = gs.calcKingMoves
	case Queen:
		calcFunc = gs.calcQueenMoves
	case Rook:
		calcFunc = gs.calcRookMoves
	case Bishop:
		calcFunc = gs.calcBishopMoves
	case Knight:
		calcFunc = gs.calcKnightMoves
	case Pawn:
		calcFunc = gs.calcPawnMoves
	}
	squares = calcFunc(rank, file, p)
	return
}

func (gs *GameState) calcKingMoves(rank, file rune, p piece) (squares []string) {
	upRank := rank + 1
	downRank := rank - 1
	leftFile := file - 1
	rightFile := file + 1
	canMoveUp := upRank <= '8'
	canMoveDown := downRank >= '1'
	canMoveLeft := leftFile >= 'a'
	canMoveRight := rightFile <= 'h'

	if canMoveUp {
		upSquare, isValid := gs.checkGameSquare(upRank, file, p)
		if isValid {
			squares = append(squares, upSquare)
		}
		if canMoveLeft {
			upLeftSquare, isValid := gs.checkGameSquare(upRank, leftFile, p)
			if isValid {
				squares = append(squares, upLeftSquare)
			}
		}
		if canMoveRight {
			upRightSquare, isValid := gs.checkGameSquare(upRank, rightFile, p)
			if isValid {
				squares = append(squares, upRightSquare)
			}
		}
	}
	if canMoveDown {
		downSquare, isValid := gs.checkGameSquare(downRank, file, p)
		if isValid {
			squares = append(squares, downSquare)
		}
		if canMoveLeft {
			leftSquare, isValid := gs.checkGameSquare(downRank, leftFile, p)
			if isValid {
				squares = append(squares, leftSquare)
			}
		}
		if canMoveRight {
			rightSquare, isValid := gs.checkGameSquare(downRank, rightFile, p)
			if isValid {
				squares = append(squares, rightSquare)
			}
		}
	}
	if canMoveLeft {
		leftSquare, isValid := gs.checkGameSquare(rank, leftFile, p)
		if isValid {
			squares = append(squares, leftSquare)
		}
	}
	if canMoveRight {
		rightSquare, isValid := gs.checkGameSquare(rank, rightFile, p)
		if isValid {
			squares = append(squares, rightSquare)
		}
	}

	var startingSquare string
	if p.PlayerColor == Black {
		startingSquare = "e8"
	} else {
		startingSquare = "e1"
	}
	if p.Square == startingSquare {
		kingRookSquare := string(file+3) + string(rank)
		if otherPiece, ok := gs.Pieces[kingRookSquare]; ok &&
			otherPiece.PieceType == Rook &&
			otherPiece.PlayerColor == p.PlayerColor {
			kingBishopSquare := string(file+1) + string(rank)
			kingKnightSquare := string(file+2) + string(rank)
			_, knightSquareOccupied := gs.Pieces[kingKnightSquare]
			_, bishopSquareOccupied := gs.Pieces[kingBishopSquare]
			if !bishopSquareOccupied && !knightSquareOccupied {
				squares = append(squares, kingKnightSquare)
			}

		}
		queenRookSquare := string(file-4) + string(rank)
		if otherPiece, ok := gs.Pieces[queenRookSquare]; ok &&
			otherPiece.PieceType == Rook &&
			otherPiece.PlayerColor == p.PlayerColor {
			queenBishopSquare := string(file-2) + string(rank)
			queenKnightSquare := string(file-3) + string(rank)
			queenSquare := string(file-1) + string(rank)
			_, knightSquareOccupied := gs.Pieces[queenKnightSquare]
			_, bishopSquareOccupied := gs.Pieces[queenBishopSquare]
			_, queenSquareOccupied := gs.Pieces[queenSquare]
			if !bishopSquareOccupied && !knightSquareOccupied && !queenSquareOccupied {
				squares = append(squares, queenBishopSquare)
			}
		}
	}
	return
}

func (gs *GameState) calcQueenMoves(rank, file rune, p piece) (squares []string) {
	squares = append(squares, gs.calcRookMoves(rank, file, p)...)
	squares = append(squares, gs.calcBishopMoves(rank, file, p)...)
	return
}

func (gs *GameState) calcRookMoves(rank, file rune, p piece) (squares []string) {
	for r, f := rank+1, file; r <= '8'; r++ {
		canidateSquare, isValid := gs.checkGameSquare(r, f, p)
		if !isValid {
			break
		}
		squares = append(squares, canidateSquare)
	}
	for r, f := rank-1, file; r >= '1'; r-- {
		canidateSquare, isValid := gs.checkGameSquare(r, f, p)
		if !isValid {
			break
		}
		squares = append(squares, canidateSquare)
	}
	for r, f := rank, file-1; f >= 'a'; f-- {
		canidateSquare, isValid := gs.checkGameSquare(r, f, p)
		if !isValid {
			break
		}
		squares = append(squares, canidateSquare)
	}
	for r, f := rank, file+1; f <= 'h'; f++ {
		canidateSquare, isValid := gs.checkGameSquare(r, f, p)
		if !isValid {
			break
		}
		squares = append(squares, canidateSquare)
	}
	return
}

func (gs *GameState) calcBishopMoves(rank, file rune, p piece) (squares []string) {
	// up right
	for r, f := rank+1, file+1; r <= '8' && f <= 'h'; r, f = r+1, f+1 {
		canidateSquare, isValid := gs.checkGameSquare(r, f, p)
		if isValid {
			squares = append(squares, canidateSquare)
		}
	}
	// down right
	for r, f := rank-1, file+1; r >= '1' && f <= 'h'; r, f = r-1, f+1 {
		canidateSquare, isValid := gs.checkGameSquare(r, f, p)
		if isValid {
			squares = append(squares, canidateSquare)
		}
	}
	// down left
	for r, f := rank-1, file-1; r >= '1' && f >= 'a'; r, f = r-1, f-1 {
		canidateSquare, isValid := gs.checkGameSquare(r, f, p)
		if isValid {
			squares = append(squares, canidateSquare)
		}
	}
	// up left
	for r, f := rank+1, file-1; r <= '8' && f >= 'a'; r, f = r+1, f-1 {
		canidateSquare, isValid := gs.checkGameSquare(r, f, p)
		if isValid {
			squares = append(squares, canidateSquare)
		}
	}
	return
}

func (gs *GameState) calcKnightMoves(rank, file rune, p piece) (squares []string) {
	if upTwo := rank + 2; upTwo <= '8' {
		if leftOne := file - 1; leftOne >= 'a' {
			canidateSquare, isValid := gs.checkGameSquare(upTwo, leftOne, p)
			if isValid {
				squares = append(squares, canidateSquare)
			}
		}
		if rightOne := file + 1; rightOne <= 'h' {
			canidateSquare, isValid := gs.checkGameSquare(upTwo, rightOne, p)
			if isValid {
				squares = append(squares, canidateSquare)
			}
		}
	}
	if rightTwo := file + 2; rightTwo <= 'h' {
		if upOne := rank + 1; upOne <= '8' {
			canidateSquare, isValid := gs.checkGameSquare(upOne, rightTwo, p)
			if isValid {
				squares = append(squares, canidateSquare)
			}
		}
		if downOne := rank - 1; downOne >= '1' {
			canidateSquare, isValid := gs.checkGameSquare(downOne, rightTwo, p)
			if isValid {
				squares = append(squares, canidateSquare)
			}
		}
	}
	if downTwo := rank - 2; downTwo >= '1' {
		if leftOne := file - 1; leftOne >= 'a' {
			canidateSquare, isValid := gs.checkGameSquare(downTwo, leftOne, p)
			if isValid {
				squares = append(squares, canidateSquare)
			}
		}
		if rightOne := file + 1; rightOne <= 'h' {
			canidateSquare, isValid := gs.checkGameSquare(downTwo, rightOne, p)
			if isValid {
				squares = append(squares, canidateSquare)
			}
		}
	}
	if leftTwo := file - 2; leftTwo >= 'a' {
		if upOne := rank + 1; upOne <= '8' {
			canidateSquare, isValid := gs.checkGameSquare(upOne, leftTwo, p)
			if isValid {
				squares = append(squares, canidateSquare)
			}
		}
		if downOne := rank - 1; downOne >= '1' {
			canidateSquare, isValid := gs.checkGameSquare(downOne, leftTwo, p)
			if isValid {
				squares = append(squares, canidateSquare)
			}
		}
	}
	return
}

func (gs *GameState) calcPawnMoves(rank, file rune, p piece) (squares []string) {
	if p.PlayerColor == Black {
		startRank := '7'
		nextRank := rank - 1
		enPassentRank := '4'
		if canidateSquare, valid := gs.checkPawnMove(nextRank, file); valid {
			squares = append(squares, canidateSquare)
			if startRank == rank {
				if canidateSquare, valid = gs.checkPawnMove(rank-2, file); valid {
					squares = append(squares, canidateSquare)
				}
			}
		}
		if leftFile := file - 1; leftFile >= 'a' {
			captureSquare, validCapture := gs.canCapture(nextRank, leftFile, p)
			if validCapture {
				squares = append(squares, captureSquare)
			} else if rank == enPassentRank {
				enPassantSquare, validCapture := gs.canCapture(rank, leftFile, p)
				if validCapture {
					if otherPiece, _ := gs.Pieces[enPassantSquare]; otherPiece.PieceType == Pawn {
						if captureSquare, isEmpty := gs.checkPawnMove(nextRank, leftFile); isEmpty {
							squares = append(squares, captureSquare)
						}
					}
				}
			}
		}
		if rightFile := file + 1; rightFile <= 'h' {
			captureSquare, validCapture := gs.canCapture(nextRank, rightFile, p)
			if validCapture {
				squares = append(squares, captureSquare)
			} else if rank == enPassentRank {
				enPassantSquare, validCapture := gs.canCapture(rank, rightFile, p)
				if validCapture {
					if otherPiece, _ := gs.Pieces[enPassantSquare]; otherPiece.PieceType == Pawn {
						if captureSquare, isEmpty := gs.checkPawnMove(nextRank, rightFile); isEmpty {
							squares = append(squares, captureSquare)
						}
					}
				}
			}
		}
	} else {
		startRank := '2'
		nextRank := rank + 1
		enPassentRank := '5'
		if canidateSquare, valid := gs.checkPawnMove(nextRank, file); valid {
			squares = append(squares, canidateSquare)
			if startRank == rank {
				if canidateSquare, valid = gs.checkPawnMove(rank+2, file); valid {
					squares = append(squares, canidateSquare)
				}
			}
		}
		if leftFile := file - 1; leftFile >= 'a' {
			captureSquare, validCapture := gs.canCapture(nextRank, leftFile, p)
			if validCapture {
				squares = append(squares, captureSquare)
			} else if rank == enPassentRank {
				enPassantSquare, validCapture := gs.canCapture(rank, leftFile, p)
				if validCapture {
					if otherPiece, _ := gs.Pieces[enPassantSquare]; otherPiece.PieceType == Pawn {
						if captureSquare, isEmpty := gs.checkPawnMove(nextRank, leftFile); isEmpty {
							squares = append(squares, captureSquare)
						}
					}
				}
			}
		}
		if rightFile := file + 1; rightFile <= 'h' {
			captureSquare, validCapture := gs.canCapture(nextRank, rightFile, p)
			if validCapture {
				squares = append(squares, captureSquare)
			} else if rank == enPassentRank {
				enPassantSquare, validCapture := gs.canCapture(rank, rightFile, p)
				if validCapture {
					if otherPiece, _ := gs.Pieces[enPassantSquare]; otherPiece.PieceType == Pawn {
						if captureSquare, isEmpty := gs.checkPawnMove(nextRank, rightFile); isEmpty {
							squares = append(squares, captureSquare)
						}
					}
				}
			}
		}
	}
	return
}

func (gs *GameState) canCapture(r, f rune, p piece) (captureSquare string, validCapture bool) {
	captureSquare = string(f) + string(r)
	if otherPiece, occupiedSquare := gs.Pieces[captureSquare]; occupiedSquare {
		if otherPiece.PlayerColor != p.PlayerColor {
			validCapture = true
			return
		}
	}
	validCapture = false
	return
}

func (gs *GameState) checkPawnMove(r, f rune) (canidateSquare string, valid bool) {
	canidateSquare = string(f) + string(r)
	valid = false
	if _, occupiedSquare := gs.Pieces[canidateSquare]; !occupiedSquare {
		valid = true
	}
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
