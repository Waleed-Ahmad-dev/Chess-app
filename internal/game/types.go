package game

// import "fmt"

// Color represents the side (White or Black)
type Color int

const (
	White Color = iota
	Black
)

// String returns the string representation of the color
func (c Color) String() string {
	if c == White {
		return "White"
	}
	return "Black"
}

// PieceType represents the rank of the piece
type PieceType int

const (
	Empty PieceType = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

func (p PieceType) String() string {
	switch p {
	case Pawn:
		return "Pawn"
	case Knight:
		return "Knight"
	case Bishop:
		return "Bishop"
	case Rook:
		return "Rook"
	case Queen:
		return "Queen"
	case King:
		return "King"
	default:
		return "Empty"
	}
}

// Piece represents a piece on the board
type Piece struct {
	Type  PieceType
	Color Color
}

// Board represents the 8x8 chess board
// We use a 1D array of size 64.
// Index 0 = a1, 7 = h1, ..., 63 = h8
type Board [64]Piece

// Move represents a single move
type Move struct {
	From      int       // Square index (0-63)
	To        int       // Square index (0-63)
	Piece     PieceType // The piece being moved
	Promotion PieceType // If pawn promotion, what type? (Empty otherwise)
}

// String provides a simple debug print for the piece
func (p Piece) String() string {
	if p.Type == Empty {
		return "."
	}
	// Simple ASCII representation
	// We will make a better renderer later
	switch p.Type {
	case Pawn:
		if p.Color == White {
			return "P"
		} else {
			return "p"
		}
	case Knight:
		if p.Color == White {
			return "N"
		} else {
			return "n"
		}
	case Bishop:
		if p.Color == White {
			return "B"
		} else {
			return "b"
		}
	case Rook:
		if p.Color == White {
			return "R"
		} else {
			return "r"
		}
	case Queen:
		if p.Color == White {
			return "Q"
		} else {
			return "q"
		}
	case King:
		if p.Color == White {
			return "K"
		} else {
			return "k"
		}
	}
	return "?"
}
