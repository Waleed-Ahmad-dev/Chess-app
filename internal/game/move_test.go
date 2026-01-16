package game

import "testing"

func TestSimplePawnMove(t *testing.T) {
	// 1. Setup Game
	g := NewGame()

	// 2. Define a move: White Pawn e2 -> e4
	from := CoordToIndex("e2")
	to := CoordToIndex("e4")
	move := Move{From: from, To: to, Piece: Pawn}

	// 3. Execute Move
	g.MakeMove(move)

	// 4. Assertions
	// Check e2 is empty
	if g.Board[from].Type != Empty {
		t.Errorf("Expected e2 to be empty, got %v", g.Board[from])
	}

	// Check e4 has a White Pawn
	pieceAtDest := g.Board[to]
	if pieceAtDest.Type != Pawn || pieceAtDest.Color != White {
		t.Errorf("Expected White Pawn at e4, got %v", pieceAtDest)
	}

	// Check Turn Switched
	if g.Turn != Black {
		t.Errorf("Expected Turn to be Black, got %v", g.Turn)
	}
}
