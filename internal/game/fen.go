package game

import (
	"strings"
	"unicode"
)

// LoadFEN parses a FEN string and updates the Game state
func (g *Game) LoadFEN(fen string) {
	parts := strings.Split(fen, " ")

	// 1. Piece Placement
	piecePlacement := parts[0]
	g.Board = Board{} // Clear board

	rank := 7
	file := 0

	for _, char := range piecePlacement {
		if char == '/' {
			rank--
			file = 0
			continue
		}

		if unicode.IsDigit(char) {
			emptyCount := int(char - '0')
			file += emptyCount
		} else {
			index := rank*8 + file
			g.Board[index] = charToPiece(char)
			file++
		}
	}

	// 2. Turn
	if len(parts) > 1 {
		if parts[1] == "w" {
			g.Turn = White
		} else {
			g.Turn = Black
		}
	}

	// 3. Castling Rights
	// Default to false, then enable based on string
	g.Castling = CastlingRights{}
	if len(parts) > 2 {
		c := parts[2]
		if c != "-" {
			for _, char := range c {
				switch char {
				case 'K':
					g.Castling.WhiteKingSide = true
				case 'Q':
					g.Castling.WhiteQueenSide = true
				case 'k':
					g.Castling.BlackKingSide = true
				case 'q':
					g.Castling.BlackQueenSide = true
				}
			}
		}
	}

	// 4. En Passant
	g.EnPassantTarget = -1
	if len(parts) > 3 {
		ep := parts[3]
		if ep != "-" {
			g.EnPassantTarget = CoordToIndex(ep)
		}
	}
}

func charToPiece(char rune) Piece {
	switch char {
	case 'P':
		return Piece{Type: Pawn, Color: White}
	case 'N':
		return Piece{Type: Knight, Color: White}
	case 'B':
		return Piece{Type: Bishop, Color: White}
	case 'R':
		return Piece{Type: Rook, Color: White}
	case 'Q':
		return Piece{Type: Queen, Color: White}
	case 'K':
		return Piece{Type: King, Color: White}
	case 'p':
		return Piece{Type: Pawn, Color: Black}
	case 'n':
		return Piece{Type: Knight, Color: Black}
	case 'b':
		return Piece{Type: Bishop, Color: Black}
	case 'r':
		return Piece{Type: Rook, Color: Black}
	case 'q':
		return Piece{Type: Queen, Color: Black}
	case 'k':
		return Piece{Type: King, Color: Black}
	}
	return Piece{Type: Empty}
}
