package game

import (
	"fmt"
	"strings"
)

// GeneratePGN creates a PGN string from the game history
func (g *Game) GeneratePGN() string {
	var sb strings.Builder

	// Temporarily replay game to calculate moves accurately if needed,
	// but we have MoveResults, so we can format directly.

	for i, result := range g.MoveResults {
		moveNum := (i / 2) + 1

		// Add "1. " for white's moves
		if i%2 == 0 {
			sb.WriteString(fmt.Sprintf("%d. ", moveNum))
		}

		// Convert move to SAN (Standard Algebraic Notation)
		san := g.moveToSAN(result, i)
		sb.WriteString(san)
		sb.WriteString(" ")
	}

	return strings.TrimSpace(sb.String())
}

// moveToSAN converts a move to Standard Algebraic Notation (e.g., "Nf3", "exd5", "O-O")
// Note: This is a simplified SAN generator. Full disambiguation (e.g. Nbd7)
// requires checking all other pieces, which we implement simply here.
func (g *Game) moveToSAN(res MoveResult, moveIndex int) string {
	m := res.Move

	// 1. Castling
	if m.MoveType == MoveCastling {
		if m.To > m.From { // King side
			return "O-O"
		}
		return "O-O-O"
	}

	var sb strings.Builder

	// 2. Piece Name (omitted for Pawns)
	pieceChar := ""
	switch m.Piece {
	case Knight:
		pieceChar = "N"
	case Bishop:
		pieceChar = "B"
	case Rook:
		pieceChar = "R"
	case Queen:
		pieceChar = "Q"
	case King:
		pieceChar = "K"
	}
	sb.WriteString(pieceChar)

	// 3. Disambiguation (simplified)
	// If a capture by pawn, we need the source file (e.g., "exd5")
	if m.Piece == Pawn && res.WasCapture {
		srcFile := IndexToCoord(m.From)[0:1]
		sb.WriteString(srcFile)
	}

	// 4. Capture Indicator
	if res.WasCapture {
		sb.WriteString("x")
	}

	// 5. Destination
	sb.WriteString(IndexToCoord(m.To))

	// 6. Promotion
	if res.WasPromotion {
		sb.WriteString("=")
		promChar := ""
		switch m.Promotion {
		case Queen:
			promChar = "Q"
		case Rook:
			promChar = "R"
		case Bishop:
			promChar = "B"
		case Knight:
			promChar = "N"
		}
		sb.WriteString(promChar)
	}

	// 7. Checks and Mates
	if res.WasCheckmate {
		sb.WriteString("#")
	} else if res.WasCheck {
		sb.WriteString("+")
	}

	return sb.String()
}
