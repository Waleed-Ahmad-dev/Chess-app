package game

func (g *Game) GenerateLegalMoves() []Move {
	// 1. Generate all physical moves (now includes Castling/EP candidates)
	pseudoMoves := g.GeneratePseudoLegalMoves()
	legalMoves := []Move{}

	// 2. Filter them
	for _, m := range pseudoMoves {
		// Create a temporary board copy
		tempBoard := g.Board

		// Execute move on temp board explicitly for check testing
		// Logic mirrors MakeMove but without updating History/State

		// Move the main piece
		tempBoard[m.To] = tempBoard[m.From]
		tempBoard[m.From] = Piece{Type: Empty}
		if m.Promotion != Empty {
			tempBoard[m.To] = Piece{Type: m.Promotion, Color: g.Turn}
		}

		// Handle En Passant Capture on Temp Board
		if m.MoveType == MoveEnPassant {
			// Remove the captured pawn
			var captureSq int
			if g.Turn == White {
				captureSq = m.To - 8 // Pawn is below the target
			} else {
				captureSq = m.To + 8 // Pawn is above the target
			}
			tempBoard[captureSq] = Piece{Type: Empty}
		}

		// Handle Castling Rook Move on Temp Board
		if m.MoveType == MoveCastling {
			// King has already moved in lines above. Move the rook.
			if m.To == 62 { // White Short (g1) -> Move Rook h1 to f1
				tempBoard[61] = tempBoard[63]
				tempBoard[63] = Piece{Type: Empty}
			} else if m.To == 58 { // White Long (c1) -> Move Rook a1 to d1
				tempBoard[59] = tempBoard[56]
				tempBoard[56] = Piece{Type: Empty}
			} else if m.To == 6 { // Black Short (g8) -> Move Rook h8 to f8
				tempBoard[5] = tempBoard[7]
				tempBoard[7] = Piece{Type: Empty}
			} else if m.To == 2 { // Black Long (c8) -> Move Rook a8 to d8
				tempBoard[3] = tempBoard[0]
				tempBoard[0] = Piece{Type: Empty}
			}
		}

		// 3. Verify King Safety
		if !tempBoard.InCheck(g.Turn) {
			legalMoves = append(legalMoves, m)
		}
	}

	return legalMoves
}

// GeneratePseudoLegalMoves returns all possible moves for the active color
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

	// 2. Castling
	// Requirements: King not moved, Rook not moved, Path empty, Not in Check, Path not Attacked

	if b.InCheck(g.Turn) {
		return moves // Cannot castle out of check
	}

	if g.Turn == White {
		// White Short (King e1 -> g1)
		if g.Castling.WhiteKingSide && b[61].Type == Empty && b[62].Type == Empty {
			// Check if f1(61) or g1(62) is attacked
			if !b.IsSquareAttacked(61, Black) && !b.IsSquareAttacked(62, Black) {
				moves = append(moves, Move{From: sq, To: 62, Piece: King, MoveType: MoveCastling})
			}
		}
		// White Long (King e1 -> c1)
		if g.Castling.WhiteQueenSide && b[59].Type == Empty && b[58].Type == Empty && b[57].Type == Empty {
			if !b.IsSquareAttacked(59, Black) && !b.IsSquareAttacked(58, Black) {
				moves = append(moves, Move{From: sq, To: 58, Piece: King, MoveType: MoveCastling})
			}
		}
	} else {
		// Black Short (King e8 -> g8)
		if g.Castling.BlackKingSide && b[5].Type == Empty && b[6].Type == Empty {
			if !b.IsSquareAttacked(5, White) && !b.IsSquareAttacked(6, White) {
				moves = append(moves, Move{From: sq, To: 6, Piece: King, MoveType: MoveCastling})
			}
		}
		// Black Long (King e8 -> c8)
		if g.Castling.BlackQueenSide && b[3].Type == Empty && b[2].Type == Empty && b[1].Type == Empty {
			if !b.IsSquareAttacked(3, White) && !b.IsSquareAttacked(2, White) {
				moves = append(moves, Move{From: sq, To: 2, Piece: King, MoveType: MoveCastling})
			}
		}
	}

	return moves
}

func (b *Board) addSteppingMoves(sq int, offsets [][2]int, moves *[]Move) {
	startRank := sq / 8
	startFile := sq % 8

	for _, off := range offsets {
		dFile, dRank := off[0], off[1]
		targetRank := startRank + dRank
		targetFile := startFile + dFile

		if targetRank >= 0 && targetRank < 8 && targetFile >= 0 && targetFile < 8 {
			targetSq := targetRank*8 + targetFile
			targetPiece := b[targetSq]

			if targetPiece.Type == Empty || targetPiece.Color != b[sq].Color {
				*moves = append(*moves, Move{From: sq, To: targetSq, Piece: b[sq].Type})
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
				if targetPiece.Color != piece.Color {
					moves = append(moves, Move{From: sq, To: targetSq, Piece: piece.Type})
				}
				break
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
	if piece.Color == Black {
		direction = -1
		startRank = 6
	}

	// 1. Forward Move
	targetRank := rank + direction
	if targetRank >= 0 && targetRank < 8 {
		targetSq := targetRank*8 + file
		if b[targetSq].Type == Empty {
			// Promotion?
			if targetRank == 7 || targetRank == 0 {
				promotions := []PieceType{Queen, Rook, Bishop, Knight}
				for _, p := range promotions {
					moves = append(moves, Move{From: sq, To: targetSq, Piece: Pawn, Promotion: p})
				}
			} else {
				moves = append(moves, Move{From: sq, To: targetSq, Piece: Pawn})
			}

			// 2. Double Move
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
			targetSq := targetRank*8 + captureFile

			// Normal Capture
			if targetRank >= 0 && targetRank < 8 {
				targetPiece := b[targetSq]
				if targetPiece.Type != Empty && targetPiece.Color != piece.Color {
					if targetRank == 7 || targetRank == 0 {
						promotions := []PieceType{Queen, Rook, Bishop, Knight}
						for _, p := range promotions {
							moves = append(moves, Move{From: sq, To: targetSq, Piece: Pawn, Promotion: p})
						}
					} else {
						moves = append(moves, Move{From: sq, To: targetSq, Piece: Pawn})
					}
				}

				// En Passant Capture
				// Condition: Target square is empty AND matches EnPassantTarget
				if targetPiece.Type == Empty && targetSq == g.EnPassantTarget {
					moves = append(moves, Move{From: sq, To: targetSq, Piece: Pawn, MoveType: MoveEnPassant})
				}
			}
		}
	}

	return moves
}
