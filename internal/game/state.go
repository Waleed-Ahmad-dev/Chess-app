package game

// MakeMove executes a move on the board and updates game state
func (g *Game) MakeMove(m Move) {
	// --- Phase 2: Capture State Snapshot ---
	// We capture the state *before* the move is executed.
	var lastMove Move
	if len(g.History) > 0 {
		lastMove = g.History[len(g.History)-1]
	}

	snapshot := StateSnapshot{
		Board:           g.Board,    // Arrays are copied by value
		Castling:        g.Castling, // Structs are copied by value
		EnPassantTarget: g.EnPassantTarget,
		Turn:            g.Turn,
		LastMove:        lastMove,
	}
	g.StateHistory = append(g.StateHistory, snapshot)
	// ----------------------------------------

	// Store the captured piece for rook capture check
	capturedPiece := g.Board[m.To]

	// 1. Identify the moving piece
	movingPiece := g.Board[m.From]

	// 2. Clear the old square
	g.Board[m.From] = Piece{Type: Empty}

	// 3. Update the new square
	if m.Promotion != Empty {
		movingPiece.Type = m.Promotion
	}
	g.Board[m.To] = movingPiece

	// --- Special Move Logic ---

	// Handle En Passant Capture (remove the victim pawn)
	if m.MoveType == MoveEnPassant {
		var captureSq int
		if g.Turn == White {
			captureSq = m.To - 8
		} else {
			captureSq = m.To + 8
		}
		g.Board[captureSq] = Piece{Type: Empty}
	}

	// Handle Castling (move the rook)
	if m.MoveType == MoveCastling {
		switch m.To {
		case 62: // White Short: h1->f1
			rook := g.Board[63]
			g.Board[63] = Piece{Type: Empty}
			g.Board[61] = rook
		case 58: // White Long: a1->d1
			rook := g.Board[56]
			g.Board[56] = Piece{Type: Empty}
			g.Board[59] = rook
		case 6: // Black Short: h8->f8
			rook := g.Board[7]
			g.Board[7] = Piece{Type: Empty}
			g.Board[5] = rook
		case 2: // Black Long: a8->d8
			rook := g.Board[0]
			g.Board[0] = Piece{Type: Empty}
			g.Board[3] = rook
		}
	}

	// --- Update Game State (Rights) ---

	// Update En Passant Target
	g.EnPassantTarget = -1 // Default to none
	if movingPiece.Type == Pawn {
		diff := m.To - m.From
		// If moved 2 squares (diff 16 or -16)
		if diff == 16 || diff == -16 {
			// Target is the square skipped
			g.EnPassantTarget = m.From + (diff / 2)
		}
	}

	// Update Castling Rights
	// If King moves, lose both rights
	if movingPiece.Type == King {
		if g.Turn == White {
			g.Castling.WhiteKingSide = false
			g.Castling.WhiteQueenSide = false
		} else {
			g.Castling.BlackKingSide = false
			g.Castling.BlackQueenSide = false
		}
	}

	// If Rook moves, lose specific right
	// We check specific starting squares for rooks
	if movingPiece.Type == Rook {
		if m.From == 63 {
			g.Castling.WhiteKingSide = false
		} // h1
		if m.From == 56 {
			g.Castling.WhiteQueenSide = false
		} // a1
		if m.From == 7 {
			g.Castling.BlackKingSide = false
		} // h8
		if m.From == 0 {
			g.Castling.BlackQueenSide = false
		} // a8
	}

	// If a rook is captured, remove castling rights for that side
	if capturedPiece.Type == Rook {
		switch m.To {
		case 63: // h1
			g.Castling.WhiteKingSide = false
		case 56: // a1
			g.Castling.WhiteQueenSide = false
		case 7: // h8
			g.Castling.BlackKingSide = false
		case 0: // a8
			g.Castling.BlackQueenSide = false
		}
	}

	// 4. Record the move in history
	g.History = append(g.History, m)

	// 5. Switch Turn
	if g.Turn == White {
		g.Turn = Black
	} else {
		g.Turn = White
	}
}

// UndoMove reverts the last move using the StateStack
func (g *Game) UndoMove() {
	// 1. Check if StateHistory is empty
	if len(g.StateHistory) == 0 {
		return
	}

	// 2. Pop the last snapshot
	lastIndex := len(g.StateHistory) - 1
	snapshot := g.StateHistory[lastIndex]

	// 3. Overwrite game state with values from the snapshot
	g.Board = snapshot.Board
	g.Castling = snapshot.Castling
	g.EnPassantTarget = snapshot.EnPassantTarget
	g.Turn = snapshot.Turn

	// 4. Remove the snapshot from the list
	g.StateHistory = g.StateHistory[:lastIndex]

	// 5. Sync the Move History
	// Since the snapshot was taken *before* the move was made,
	// the move currently at the end of g.History is the one we just undid.
	// We must remove it to keep the history in sync with the board.
	if len(g.History) > 0 {
		g.History = g.History[:len(g.History)-1]
	}
}
