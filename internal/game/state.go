package game

// MakeMove executes a move on the board and updates game state
func (g *Game) MakeMove(m Move) MoveResult {
	// --- Capture State Snapshot ---
	var lastMove Move
	if len(g.History) > 0 {
		lastMove = g.History[len(g.History)-1]
	}

	snapshot := StateSnapshot{
		Board:           g.Board,
		Castling:        g.Castling,
		EnPassantTarget: g.EnPassantTarget,
		Turn:            g.Turn,
		LastMove:        lastMove,
	}
	g.StateHistory = append(g.StateHistory, snapshot)
	// ----------------------------------------

	// Store original pieces for sound detection
	targetPiece := g.Board[m.To]

	wasCapture := targetPiece.Type != Empty
	wasCastle := m.MoveType == MoveCastling
	wasPromotion := m.Promotion != Empty

	// Execute the move (existing code remains the same)
	capturedPiece := g.Board[m.To]
	movingPiece := g.Board[m.From]
	g.Board[m.From] = Piece{Type: Empty}

	if m.Promotion != Empty {
		movingPiece.Type = m.Promotion
	}
	g.Board[m.To] = movingPiece

	// --- Special Move Logic ---
	if m.MoveType == MoveEnPassant {
		var captureSq int
		if g.Turn == White {
			captureSq = m.To - 8
		} else {
			captureSq = m.To + 8
		}
		g.Board[captureSq] = Piece{Type: Empty}
		wasCapture = true
	}

	if m.MoveType == MoveCastling {
		switch m.To {
		case 62:
			rook := g.Board[63]
			g.Board[63] = Piece{Type: Empty}
			g.Board[61] = rook
		case 58:
			rook := g.Board[56]
			g.Board[56] = Piece{Type: Empty}
			g.Board[59] = rook
		case 6:
			rook := g.Board[7]
			g.Board[7] = Piece{Type: Empty}
			g.Board[5] = rook
		case 2:
			rook := g.Board[0]
			g.Board[0] = Piece{Type: Empty}
			g.Board[3] = rook
		}
	}

	// Update game state
	g.EnPassantTarget = -1
	if movingPiece.Type == Pawn {
		diff := m.To - m.From
		if diff == 16 || diff == -16 {
			g.EnPassantTarget = m.From + (diff / 2)
		}
	}

	// Update Castling Rights
	if movingPiece.Type == King {
		if g.Turn == White {
			g.Castling.WhiteKingSide = false
			g.Castling.WhiteQueenSide = false
		} else {
			g.Castling.BlackKingSide = false
			g.Castling.BlackQueenSide = false
		}
	}

	if movingPiece.Type == Rook {
		if m.From == 63 {
			g.Castling.WhiteKingSide = false
		}
		if m.From == 56 {
			g.Castling.WhiteQueenSide = false
		}
		if m.From == 7 {
			g.Castling.BlackKingSide = false
		}
		if m.From == 0 {
			g.Castling.BlackQueenSide = false
		}
	}

	if capturedPiece.Type == Rook {
		switch m.To {
		case 63:
			g.Castling.WhiteKingSide = false
		case 56:
			g.Castling.WhiteQueenSide = false
		case 7:
			g.Castling.BlackKingSide = false
		case 0:
			g.Castling.BlackQueenSide = false
		}
	}

	// Record history
	g.History = append(g.History, m)

	// Check for check/checkmate after move
	opponentColor := Black
	if g.Turn == Black {
		opponentColor = White
	}

	wasCheck := g.Board.InCheck(opponentColor)
	wasCheckmate := false

	if wasCheck {
		// Generate opponent's moves to check if it's checkmate
		tempGame := &Game{
			Board:           g.Board,
			Turn:            opponentColor,
			Castling:        g.Castling,
			EnPassantTarget: g.EnPassantTarget,
		}
		opponentMoves := tempGame.GenerateLegalMoves()
		wasCheckmate = len(opponentMoves) == 0
	}

	// Create move result
	result := MoveResult{
		Move:         m,
		WasCapture:   wasCapture,
		WasCheck:     wasCheck,
		WasCheckmate: wasCheckmate,
		WasCastle:    wasCastle,
		WasPromotion: wasPromotion,
		WasIllegal:   false,
	}

	g.MoveResults = append(g.MoveResults, result)

	// Switch turn
	if g.Turn == White {
		g.Turn = Black
	} else {
		g.Turn = White
	}

	return result
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

// GetLastMoveResult returns the last move result
func (g *Game) GetLastMoveResult() *MoveResult {
	if len(g.MoveResults) == 0 {
		return nil
	}
	return &g.MoveResults[len(g.MoveResults)-1]
}
