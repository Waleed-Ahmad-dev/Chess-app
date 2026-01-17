package game

// Color represents the side (White or Black)
type Color int

const (
	White Color = iota
	Black
)

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
type Board [64]Piece

// MoveType helps identify special moves
type MoveType int

const (
	MoveNormal MoveType = iota
	MoveCastling
	MoveEnPassant
)

// Move represents a single move
type Move struct {
	From      int
	To        int
	Piece     PieceType
	Promotion PieceType // If pawn promotion, what type?
	MoveType  MoveType  // Normal, Castling, or EnPassant
}

func (p Piece) String() string {
	if p.Type == Empty {
		return "."
	}
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

// CastlingRights tracks permissions
type CastlingRights struct {
	WhiteKingSide  bool
	WhiteQueenSide bool
	BlackKingSide  bool
	BlackQueenSide bool
}

// StateSnapshot captures the game state before a move is made
type StateSnapshot struct {
	Board           Board
	Castling        CastlingRights
	EnPassantTarget int
	Turn            Color
	LastMove        Move
}

type Game struct {
	Board           Board
	Turn            Color
	History         []Move
	StateHistory    []StateSnapshot // Stack for Undo functionality
	Castling        CastlingRights
	EnPassantTarget int
}

// NewGame returns a game with the starting position
func NewGame() *Game {
	g := &Game{
		Turn:            White,
		History:         make([]Move, 0),
		StateHistory:    make([]StateSnapshot, 0),
		EnPassantTarget: -1,
		Castling:        CastlingRights{true, true, true, true},
	}
	g.LoadFEN(StartFEN)
	return g
}
