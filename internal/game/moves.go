package game

// GeneratePseudoLegalMoves returns all possible moves for the active color
// It ignores whether the King is in check (we handle that in Phase 5)
func (b *Board) GeneratePseudoLegalMoves(turn Color) []Move {
	moves := []Move{}

	for sq := 0; sq < 64; sq++ {
		piece := b[sq]

		// Skip empty squares or enemy pieces
		if piece.Type == Empty || piece.Color != turn {
			continue
		}

		switch piece.Type {
		case Knight:
			moves = append(moves, b.getKnightMoves(sq)...)
		case King:
			moves = append(moves, b.getKingMoves(sq)...)
			// We will add Sliders and Pawns in the next step
		}
	}
	return moves
}

// --- Stepping Piece Logic ---

func (b *Board) getKnightMoves(sq int) []Move {
	moves := []Move{}
	// A Knight has 8 potential jumps.
	// We represent them as (change in file, change in rank)
	offsets := [][2]int{
		{1, 2}, {1, -2}, {-1, 2}, {-1, -2},
		{2, 1}, {2, -1}, {-2, 1}, {-2, -1},
	}

	startRank := sq / 8
	startFile := sq % 8

	for _, off := range offsets {
		dFile, dRank := off[0], off[1]
		targetRank := startRank + dRank
		targetFile := startFile + dFile

		// Check if it's on the board
		if targetRank >= 0 && targetRank < 8 && targetFile >= 0 && targetFile < 8 {
			targetSq := targetRank*8 + targetFile
			targetPiece := b[targetSq]

			// Valid if empty OR contains enemy
			if targetPiece.Type == Empty || targetPiece.Color != b[sq].Color {
				moves = append(moves, Move{From: sq, To: targetSq, Piece: Knight})
			}
		}
	}
	return moves
}

func (b *Board) getKingMoves(sq int) []Move {
	moves := []Move{}
	// King moves 1 step in any direction
	offsets := [][2]int{
		{0, 1}, {0, -1}, {1, 0}, {-1, 0}, // Up, Down, Right, Left
		{1, 1}, {1, -1}, {-1, 1}, {-1, -1}, // Diagonals
	}

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
				moves = append(moves, Move{From: sq, To: targetSq, Piece: King})
			}
		}
	}
	return moves
}
