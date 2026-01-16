package game

import (
	"strings"
	"unicode"
)

// LoadFEN parses a FEN string and updates the board
func (b *Board) LoadFEN(fen string) {
	// 1. Split the FEN into parts (Board, Turn, Castling, EnPassant, etc.)
	// We only care about the first part (piece placement) for now.
	parts := strings.Split(fen, " ")
	piecePlacement := parts[0]

	// 2. Clear the board first
	*b = Board{}

	// 3. Loop through the FEN string
	// FEN starts at Rank 8 (index 56) and goes down to Rank 1 (index 0)
	rank := 7 // Rank 8 is index 7 (0-indexed)
	file := 0 // File A is index 0

	for _, char := range piecePlacement {
		if char == '/' {
			rank--
			file = 0
			continue
		}

		if unicode.IsDigit(char) {
			// If it's a number (e.g., '8'), skip that many empty squares
			emptyCount := int(char - '0')
			file += emptyCount
		} else {
			// It's a piece character
			index := rank*8 + file
			b[index] = charToPiece(char)
			file++
		}
	}
}

// Helper to convert FEN character to Piece
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
