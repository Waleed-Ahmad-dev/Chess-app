package game

// MakeMove executes a move on the board and updates game state
func (g *Game) MakeMove(m Move) {
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
		case 6: // White Short: h1->f1
			rook := g.Board[7]
			g.Board[7] = Piece{Type: Empty}
			g.Board[5] = rook
		case 2: // White Long: a1->d1
			rook := g.Board[0]
			g.Board[0] = Piece{Type: Empty}
			g.Board[3] = rook
		case 62: // Black Short: h8->f8
			rook := g.Board[63]
			g.Board[63] = Piece{Type: Empty}
			g.Board[61] = rook
		case 58: // Black Long: a8->d8
			rook := g.Board[56]
			g.Board[56] = Piece{Type: Empty}
			g.Board[59] = rook
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
		if m.From == 7 {
			g.Castling.WhiteKingSide = false
		} // h1
		if m.From == 0 {
			g.Castling.WhiteQueenSide = false
		} // a1
		if m.From == 63 {
			g.Castling.BlackKingSide = false
		} // h8
		if m.From == 56 {
			g.Castling.BlackQueenSide = false
		} // a8
	}

	// If a rook is captured, remove castling rights for that side
	if capturedPiece.Type == Rook {
		switch m.To {
		case 7: // h1
			g.Castling.WhiteKingSide = false
		case 0: // a1
			g.Castling.WhiteQueenSide = false
		case 63: // h8
			g.Castling.BlackKingSide = false
		case 56: // a8
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

// UndoMove reverts the last move
// Note: This is now significantly harder because we need to restore
// CastlingRights and EnPassantTarget.
// For now, we will just revert the physical move as in Phase 4,
// but note that "Takeback" logic needs a full snapshot history in the future.
func (g *Game) UndoMove() {
	// Implementation intentionally left simple as per Phase 4 requirements.
	// A robust undo requires storing previous Castling/EP states in the History struct.
	if len(g.History) == 0 {
		return
	}
	lastMove := g.History[len(g.History)-1]
	g.History = g.History[:len(g.History)-1]

	movedPiece := g.Board[lastMove.To]
	if lastMove.Promotion != Empty {
		movedPiece.Type = Pawn
	}
	g.Board[lastMove.From] = movedPiece
	g.Board[lastMove.To] = Piece{Type: Empty}

	// TODO: Restore captured piece and castling rights
	// This requires storing more information in the Move struct

	if g.Turn == White {
		g.Turn = Black
	} else {
		g.Turn = White
	}
}
