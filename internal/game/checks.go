package game

// InCheck returns true if the King of the given color is under attack
func (b *Board) InCheck(color Color) bool {
	// 1. Find the King
	kingPos := -1
	for i := 0; i < 64; i++ {
		piece := b[i]
		if piece.Type == King && piece.Color == color {
			kingPos = i
			break
		}
	}

	// Should not happen, but safe guard
	if kingPos == -1 {
		return false
	}

	// 2. Check if that square is attacked by the opponent
	enemyColor := Black
	if color == Black {
		enemyColor = White
	}

	return b.IsSquareAttacked(kingPos, enemyColor)
}

// IsSquareAttacked checks if 'sq' is attacked by pieces of 'attackerColor'
func (b *Board) IsSquareAttacked(sq int, attackerColor Color) bool {
	// Strategy: "Look out" from the square as if we were a piece.
	// If we hit an enemy piece that can move that way, we are attacked.

	// 1. Check for Pawn attacks
	// (Note: We look in the REVERSE direction of pawn movement)
	pawnDir := 1
	if attackerColor == White {
		pawnDir = -1 // White pawns come from below, so we look down
	} else {
		pawnDir = 1 // Black pawns come from above, so we look up
	}

	// Check diagonals explicitly
	rank := sq / 8
	file := sq % 8

	// Look for attacker pawns
	checkRank := rank - pawnDir // The rank an attacking pawn would be on
	if checkRank >= 0 && checkRank < 8 {
		for _, offset := range []int{-1, 1} {
			checkFile := file + offset
			if checkFile >= 0 && checkFile < 8 {
				idx := checkRank*8 + checkFile
				target := b[idx]
				if target.Type == Pawn && target.Color == attackerColor {
					return true
				}
			}
		}
	}

	// 2. Check for Knight attacks
	knightOffsets := []int{
		-17, -15, -10, -6, 6, 10, 15, 17,
	}
	for _, off := range knightOffsets {
		targetSq := sq + off
		if IsOnBoard(targetSq) {
			// Validate knight jump logic (prevent wrapping)
			// A valid jump changes rank by 1 or 2
			tRank := targetSq / 8
			if abs(tRank-rank) == 1 || abs(tRank-rank) == 2 {
				target := b[targetSq]
				if target.Type == Knight && target.Color == attackerColor {
					return true
				}
			}
		}
	}

	// 3. Check for Sliding attacks (Rook/Queen) - Orthogonal
	dirs := []int{-8, 8, -1, 1} // Up, Down, Left, Right
	for _, d := range dirs {
		if checkSlider(b, sq, d, attackerColor, Rook) {
			return true
		}
	}

	// 4. Check for Sliding attacks (Bishop/Queen) - Diagonal
	dirs = []int{-9, -7, 7, 9}
	for _, d := range dirs {
		if checkSlider(b, sq, d, attackerColor, Bishop) {
			return true
		}
	}

	// 5. Check for King attacks (standard 1 square radius)
	kingOffsets := []int{-9, -8, -7, -1, 1, 7, 8, 9}
	for _, off := range kingOffsets {
		targetSq := sq + off
		if IsOnBoard(targetSq) {
			tRank := targetSq / 8
			if abs(tRank-rank) <= 1 { // Ensure strict adjacency
				target := b[targetSq]
				if target.Type == King && target.Color == attackerColor {
					return true
				}
			}
		}
	}

	return false
}

// Helper to check sliding attacks
func checkSlider(b *Board, start int, step int, color Color, pType PieceType) bool {
	cursor := start
	startRank := start / 8
	startFile := start % 8

	for {
		cursor += step
		if !IsOnBoard(cursor) {
			break
		}

		// Prevent wrapping
		currRank := cursor / 8
		currFile := cursor % 8

		// If moving horizontally, rank must stay same
		if abs(step) == 1 && currRank != startRank {
			break
		}

		// If moving diagonally (step is not 1 or 8), rank diff must match file diff
		if abs(step) != 1 && abs(step) != 8 {
			if abs(currRank-startRank) != abs(currFile-startFile) {
				break
			}
		}

		target := b[cursor]
		if target.Type == Empty {
			continue
		}

		// Found a piece
		if target.Color == color {
			if target.Type == pType || target.Type == Queen {
				return true
			}
		}
		break // Blocked by any piece
	}
	return false
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
