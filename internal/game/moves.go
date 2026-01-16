package game

func (g *Game) GenerateLegalMoves() []Move {
	pseudoMoves := g.GeneratePseudoLegalMoves()
	legalMoves := []Move{}

	for _, m := range pseudoMoves {
		if g.isMoveLegal(m) {
			legalMoves = append(legalMoves, m)
		}
	}

	return legalMoves
}

// isMoveLegal checks if a move is legal by simulating it and checking if king is safe
func (g *Game) isMoveLegal(m Move) bool {
	// Create a temporary board copy
	tempBoard := g.Board

	// Store original en passant target for restoration
	originalEP := g.EnPassantTarget

	// Execute move on temp board
	tempBoard[m.To] = tempBoard[m.From]
	tempBoard[m.From] = Piece{Type: Empty}

	if m.Promotion != Empty {
		tempBoard[m.To] = Piece{Type: m.Promotion, Color: g.Turn}
	}

	// Handle En Passant Capture
	if m.MoveType == MoveEnPassant {
		var captureSq int
		if g.Turn == White {
			captureSq = m.To - 8
		} else {
			captureSq = m.To + 8
		}
		tempBoard[captureSq] = Piece{Type: Empty}
	}

	// Handle Castling Rook Move
	if m.MoveType == MoveCastling {
		switch m.To {
		case 6: // White Short (g1)
			tempBoard[5] = tempBoard[7] // Rook: h1 -> f1
			tempBoard[7] = Piece{Type: Empty}
		case 2: // White Long (c1)
			tempBoard[3] = tempBoard[0] // Rook: a1 -> d1
			tempBoard[0] = Piece{Type: Empty}
		case 62: // Black Short (g8)
			tempBoard[61] = tempBoard[63] // Rook: h8 -> f8
			tempBoard[63] = Piece{Type: Empty}
		case 58: // Black Long (c8)
			tempBoard[59] = tempBoard[56] // Rook: a8 -> d8
			tempBoard[56] = Piece{Type: Empty}
		}
	}

	// Verify King Safety
	isLegal := !tempBoard.InCheck(g.Turn)

	// Restore en passant target
	g.EnPassantTarget = originalEP

	return isLegal
}

func (g *Game) GeneratePseudoLegalMoves() []Move {
	moves := []Move{}
	turn := g.Turn
	b := g.Board

	for sq := 0; sq < 64; sq++ {
		piece := b[sq]

		if piece.Type == Empty || piece.Color != turn {
			continue
		}

		switch piece.Type {
		case Pawn:
			moves = append(moves, g.getPawnMoves(sq)...)
		case Knight:
			moves = append(moves, b.getKnightMoves(sq)...)
		case Bishop, Rook, Queen:
			moves = append(moves, b.getSlidingMoves(sq)...)
		case King:
			moves = append(moves, g.getKingMoves(sq)...)
		}
	}
	return moves
}

// --- Stepping Pieces (Knight, King) ---

func (b *Board) getKnightMoves(sq int) []Move {
	moves := []Move{}
	offsets := [][2]int{
		{1, 2}, {1, -2}, {-1, 2}, {-1, -2},
		{2, 1}, {2, -1}, {-2, 1}, {-2, -1},
	}
	b.addSteppingMoves(sq, offsets, &moves)
	return moves
}

func (g *Game) getKingMoves(sq int) []Move {
	moves := []Move{}
	b := g.Board

	// 1. Normal King Moves
	offsets := [][2]int{
		{0, 1}, {0, -1}, {1, 0}, {-1, 0},
		{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
	}
	b.addSteppingMoves(sq, offsets, &moves)

	// 2. Castling - only add if not in check
	if b.InCheck(g.Turn) {
		return moves
	}

	if g.Turn == White {
		// White Short Castling (e1 -> g1)
		if g.Castling.WhiteKingSide && sq == 4 { // King must be on e1
			if b[5].Type == Empty && b[6].Type == Empty && // f1 and g1 empty
				b[7].Type == Rook && b[7].Color == White { // Rook on h1
				// Check that f1 and g1 are not attacked
				if !b.IsSquareAttacked(5, Black) && !b.IsSquareAttacked(6, Black) {
					moves = append(moves, Move{From: sq, To: 6, Piece: King, MoveType: MoveCastling})
				}
			}
		}

		// White Long Castling (e1 -> c1)
		if g.Castling.WhiteQueenSide && sq == 4 { // King must be on e1
			if b[3].Type == Empty && b[2].Type == Empty && b[1].Type == Empty && // d1, c1, b1 empty
				b[0].Type == Rook && b[0].Color == White { // Rook on a1
				// Check that d1 and c1 are not attacked (b1 doesn't matter)
				if !b.IsSquareAttacked(3, Black) && !b.IsSquareAttacked(2, Black) {
					moves = append(moves, Move{From: sq, To: 2, Piece: King, MoveType: MoveCastling})
				}
			}
		}
	} else {
		// Black Short Castling (e8 -> g8)
		if g.Castling.BlackKingSide && sq == 60 { // King must be on e8
			if b[61].Type == Empty && b[62].Type == Empty && // f8 and g8 empty
				b[63].Type == Rook && b[63].Color == Black { // Rook on h8
				// Check that f8 and g8 are not attacked
				if !b.IsSquareAttacked(61, White) && !b.IsSquareAttacked(62, White) {
					moves = append(moves, Move{From: sq, To: 62, Piece: King, MoveType: MoveCastling})
				}
			}
		}

		// Black Long Castling (e8 -> c8)
		if g.Castling.BlackQueenSide && sq == 60 { // King must be on e8
			if b[59].Type == Empty && b[58].Type == Empty && b[57].Type == Empty && // d8, c8, b8 empty
				b[56].Type == Rook && b[56].Color == Black { // Rook on a8
				// Check that d8 and c8 are not attacked (b8 doesn't matter)
				if !b.IsSquareAttacked(59, White) && !b.IsSquareAttacked(58, White) {
					moves = append(moves, Move{From: sq, To: 58, Piece: King, MoveType: MoveCastling})
				}
			}
		}
	}

	return moves
}

func (b *Board) addSteppingMoves(sq int, offsets [][2]int, moves *[]Move) {
	startRank := sq / 8
	startFile := sq % 8
	movingPiece := b[sq]

	for _, off := range offsets {
		dFile, dRank := off[0], off[1]
		targetRank := startRank + dRank
		targetFile := startFile + dFile

		if targetRank >= 0 && targetRank < 8 && targetFile >= 0 && targetFile < 8 {
			targetSq := targetRank*8 + targetFile
			targetPiece := b[targetSq]

			// Can only move to empty square or capture opponent's piece (NOT own piece)
			if targetPiece.Type == Empty || targetPiece.Color != movingPiece.Color {
				*moves = append(*moves, Move{From: sq, To: targetSq, Piece: movingPiece.Type})
			}
		}
	}
}

// --- Sliding Pieces ---

func (b *Board) getSlidingMoves(sq int) []Move {
	moves := []Move{}
	piece := b[sq]
	var directions [][2]int

	if piece.Type == Bishop || piece.Type == Queen {
		directions = append(directions, [][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}...)
	}
	if piece.Type == Rook || piece.Type == Queen {
		directions = append(directions, [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}...)
	}

	startRank := sq / 8
	startFile := sq % 8

	for _, dir := range directions {
		dFile, dRank := dir[0], dir[1]
		for i := 1; i < 8; i++ {
			targetRank := startRank + (dRank * i)
			targetFile := startFile + (dFile * i)

			if targetRank < 0 || targetRank > 7 || targetFile < 0 || targetFile > 7 {
				break
			}

			targetSq := targetRank*8 + targetFile
			targetPiece := b[targetSq]

			if targetPiece.Type == Empty {
				moves = append(moves, Move{From: sq, To: targetSq, Piece: piece.Type})
			} else {
				// Can capture opponent's piece, but not own piece
				if targetPiece.Color != piece.Color {
					moves = append(moves, Move{From: sq, To: targetSq, Piece: piece.Type})
				}
				break // Blocked by any piece
			}
		}
	}
	return moves
}

// --- Pawns ---

func (g *Game) getPawnMoves(sq int) []Move {
	moves := []Move{}
	b := g.Board
	piece := b[sq]

	rank := sq / 8
	file := sq % 8

	direction := 1
	startRank := 1
	promotionRank := 7

	if piece.Color == Black {
		direction = -1
		startRank = 6
		promotionRank = 0
	}

	// 1. Forward Move
	targetRank := rank + direction
	if targetRank >= 0 && targetRank < 8 {
		targetSq := targetRank*8 + file
		if b[targetSq].Type == Empty {
			// Check for promotion
			if targetRank == promotionRank {
				promotions := []PieceType{Queen, Rook, Bishop, Knight}
				for _, p := range promotions {
					moves = append(moves, Move{From: sq, To: targetSq, Piece: Pawn, Promotion: p})
				}
			} else {
				moves = append(moves, Move{From: sq, To: targetSq, Piece: Pawn})
			}

			// 2. Double Move from starting position
			if rank == startRank {
				doubleRank := rank + (2 * direction)
				doubleSq := doubleRank*8 + file
				if b[doubleSq].Type == Empty {
					moves = append(moves, Move{From: sq, To: doubleSq, Piece: Pawn})
				}
			}
		}
	}

	// 3. Captures (Normal + En Passant)
	captureOffsets := []int{-1, 1}
	for _, off := range captureOffsets {
		captureFile := file + off
		if captureFile >= 0 && captureFile < 8 {
			targetRank := rank + direction
			if targetRank >= 0 && targetRank < 8 {
				targetSq := targetRank*8 + captureFile
				targetPiece := b[targetSq]

				// Normal Capture - must be opponent's piece
				if targetPiece.Type != Empty && targetPiece.Color != piece.Color {
					if targetRank == promotionRank {
						promotions := []PieceType{Queen, Rook, Bishop, Knight}
						for _, p := range promotions {
							moves = append(moves, Move{From: sq, To: targetSq, Piece: Pawn, Promotion: p})
						}
					} else {
						moves = append(moves, Move{From: sq, To: targetSq, Piece: Pawn})
					}
				}

				// En Passant Capture
				if targetPiece.Type == Empty && targetSq == g.EnPassantTarget {
					moves = append(moves, Move{From: sq, To: targetSq, Piece: Pawn, MoveType: MoveEnPassant})
				}
			}
		}
	}

	return moves
}
