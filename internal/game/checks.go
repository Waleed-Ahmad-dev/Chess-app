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

	// Should not happen in a valid game, but safeguard
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
	rank := sq / 8
	file := sq % 8

	// 1. Check for Pawn attacks
	pawnDir := -1
	if attackerColor == White {
		pawnDir = 1 // White pawns attack upward (from lower ranks)
	}

	checkRank := rank - pawnDir
	if checkRank >= 0 && checkRank < 8 {
		for _, fileOffset := range []int{-1, 1} {
			checkFile := file + fileOffset
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
	knightMoves := [][2]int{
		{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2},
		{1, -2}, {1, 2}, {2, -1}, {2, 1},
	}
	for _, move := range knightMoves {
		checkRank := rank + move[0]
		checkFile := file + move[1]
		if checkRank >= 0 && checkRank < 8 && checkFile >= 0 && checkFile < 8 {
			idx := checkRank*8 + checkFile
			target := b[idx]
			if target.Type == Knight && target.Color == attackerColor {
				return true
			}
		}
	}

	// 3. Check for Sliding attacks (Rook/Queen) - Orthogonal
	directions := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for _, dir := range directions {
		if b.checkSlidingAttack(sq, dir, attackerColor, true) {
			return true
		}
	}

	// 4. Check for Sliding attacks (Bishop/Queen) - Diagonal
	directions = [][2]int{{-1, -1}, {-1, 1}, {1, -1}, {1, 1}}
	for _, dir := range directions {
		if b.checkSlidingAttack(sq, dir, attackerColor, false) {
			return true
		}
	}

	// 5. Check for King attacks (1 square in all directions)
	kingMoves := [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}
	for _, move := range kingMoves {
		checkRank := rank + move[0]
		checkFile := file + move[1]
		if checkRank >= 0 && checkRank < 8 && checkFile >= 0 && checkFile < 8 {
			idx := checkRank*8 + checkFile
			target := b[idx]
			if target.Type == King && target.Color == attackerColor {
				return true
			}
		}
	}

	return false
}

// checkSlidingAttack checks for sliding piece attacks in a given direction
// isOrthogonal: true for rook-like moves, false for bishop-like moves
func (b *Board) checkSlidingAttack(startSq int, direction [2]int, attackerColor Color, isOrthogonal bool) bool {
	rank := startSq / 8
	file := startSq % 8

	dRank := direction[0]
	dFile := direction[1]

	for i := 1; i < 8; i++ {
		checkRank := rank + (dRank * i)
		checkFile := file + (dFile * i)

		// Out of bounds
		if checkRank < 0 || checkRank > 7 || checkFile < 0 || checkFile > 7 {
			break
		}

		idx := checkRank*8 + checkFile
		target := b[idx]

		// Empty square, continue sliding
		if target.Type == Empty {
			continue
		}

		// Found a piece
		if target.Color == attackerColor {
			// Check if it's the right type of attacker
			if isOrthogonal {
				// Orthogonal: Rook or Queen
				if target.Type == Rook || target.Type == Queen {
					return true
				}
			} else {
				// Diagonal: Bishop or Queen
				if target.Type == Bishop || target.Type == Queen {
					return true
				}
			}
		}

		// Blocked by any piece (same color or opponent)
		break
	}

	return false
}
