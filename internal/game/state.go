package game

// MakeMove executes a move on the board and updates game state
func (g *Game) MakeMove(m Move) {
	// 1. Identify the moving piece
	movingPiece := g.Board[m.From]

	// 2. Clear the old square
	g.Board[m.From] = Piece{Type: Empty}

	// 3. Update the new square
	// Handle Promotion: if the move specifies a promotion, change piece type
	if m.Promotion != Empty {
		movingPiece.Type = m.Promotion
	}
	g.Board[m.To] = movingPiece

	// 4. Record the move in history
	g.History = append(g.History, m)

	// 5. Switch Turn
	if g.Turn == White {
		g.Turn = Black
	} else {
		g.Turn = White
	}
}

// UndoMove reverts the last move (basic implementation)
func (g *Game) UndoMove() {
	if len(g.History) == 0 {
		return
	}

	// 1. Get last move
	lastMove := g.History[len(g.History)-1]
	g.History = g.History[:len(g.History)-1]

	// 2. Move piece back
	// Note: This is simplified. To do this properly we need to capture
	// what piece was taken. For now, we just reverse the movement.
	movedPiece := g.Board[lastMove.To]

	// Revert promotion if necessary
	if lastMove.Promotion != Empty {
		movedPiece.Type = Pawn // It was a pawn before promotion
	}

	g.Board[lastMove.From] = movedPiece
	g.Board[lastMove.To] = Piece{Type: Empty} // Warning: This erases captured pieces!
	// (We will fix captured piece restoration in the next phase)

	// 3. Switch Turn back
	if g.Turn == White {
		g.Turn = Black
	} else {
		g.Turn = White
	}
}
