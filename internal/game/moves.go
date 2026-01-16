package game

// GeneratePseudoLegalMoves returns all possible moves for the active color
func (b *Board) GeneratePseudoLegalMoves(turn Color) []Move {
	moves := []Move{}

	for sq := 0; sq < 64; sq++ {
		piece := b[sq]

		// Skip empty squares or enemy pieces
		if piece.Type == Empty || piece.Color != turn {
			continue
		}

		switch piece.Type {
		case Pawn:
			moves = append(moves, b.getPawnMoves(sq)...)
		case Knight:
			moves = append(moves, b.getKnightMoves(sq)...)
		case Bishop, Rook, Queen:
			moves = append(moves, b.getSlidingMoves(sq)...)
		case King:
			moves = append(moves, b.getKingMoves(sq)...)
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

func (b *Board) getKingMoves(sq int) []Move {
	moves := []Move{}
	offsets := [][2]int{
		{0, 1}, {0, -1}, {1, 0}, {-1, 0}, // Orthogonal
		{1, 1}, {1, -1}, {-1, 1}, {-1, -1}, // Diagonal
	}
	b.addSteppingMoves(sq, offsets, &moves)
	return moves
}

// Helper for Knight/King to reduce code duplication
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

// --- Sliding Pieces (Bishop, Rook, Queen) ---

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

		// Loop to slide in that direction
		for i := 1; i < 8; i++ {
			targetRank := startRank + (dRank * i)
			targetFile := startFile + (dFile * i)

			// 1. Check boundaries
			if targetRank < 0 || targetRank > 7 || targetFile < 0 || targetFile > 7 {
				break // Hit edge of board
			}

			targetSq := targetRank*8 + targetFile
			targetPiece := b[targetSq]

			// 2. Check collisions
			if targetPiece.Type == Empty {
				// Empty square, valid move, keep sliding
				moves = append(moves, Move{From: sq, To: targetSq, Piece: piece.Type})
			} else {
				// Hit a piece
				if targetPiece.Color != piece.Color {
					// Enemy piece: capture it, then stop
					moves = append(moves, Move{From: sq, To: targetSq, Piece: piece.Type})
				}
				break // Blocked by piece (friend or foe), stop sliding
			}
		}
	}
	return moves
}

// --- Pawns ---

func (b *Board) getPawnMoves(sq int) []Move {
	moves := []Move{}
	piece := b[sq]

	rank := sq / 8
	file := sq % 8

	// Define direction and start rank based on color
	direction := 1 // White moves up
	startRank := 1 // White starts on rank 2 (index 1)

	if piece.Color == Black {
		direction = -1
		startRank = 6 // Black starts on rank 7 (index 6)
	}

	// 1. Forward Move (1 step)
	targetRank := rank + direction
	if targetRank >= 0 && targetRank < 8 {
		targetSq := targetRank*8 + file
		if b[targetSq].Type == Empty {
			// Check for promotion (Rank 8 or Rank 1)
			if targetRank == 7 || targetRank == 0 {
				// Add promotion moves (Queen, Rook, Bishop, Knight)
				promotions := []PieceType{Queen, Rook, Bishop, Knight}
				for _, p := range promotions {
					moves = append(moves, Move{From: sq, To: targetSq, Piece: Pawn, Promotion: p})
				}
			} else {
				moves = append(moves, Move{From: sq, To: targetSq, Piece: Pawn})
			}

			// 2. Double Move (2 steps) - Only if 1 step was valid
			if rank == startRank {
				doubleRank := rank + (2 * direction)
				doubleSq := doubleRank*8 + file
				if b[doubleSq].Type == Empty {
					moves = append(moves, Move{From: sq, To: doubleSq, Piece: Pawn})
				}
			}
		}
	}

	// 3. Captures (Diagonal)
	captureOffsets := []int{-1, 1} // Left and Right file
	for _, off := range captureOffsets {
		captureFile := file + off
		if captureFile >= 0 && captureFile < 8 {
			targetRank := rank + direction
			if targetRank >= 0 && targetRank < 8 {
				targetSq := targetRank*8 + captureFile
				targetPiece := b[targetSq]

				// Must contain enemy to capture
				if targetPiece.Type != Empty && targetPiece.Color != piece.Color {
					// Check promotion on capture
					if targetRank == 7 || targetRank == 0 {
						promotions := []PieceType{Queen, Rook, Bishop, Knight}
						for _, p := range promotions {
							moves = append(moves, Move{From: sq, To: targetSq, Piece: Pawn, Promotion: p})
						}
					} else {
						moves = append(moves, Move{From: sq, To: targetSq, Piece: Pawn})
					}
				}
			}
		}
	}

	return moves
}
