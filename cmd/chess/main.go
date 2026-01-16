package main

import (
	"fmt"

	"github.com/Waleed-Ahmad-dev/Chess-app/internal/game"
)

func main() {
	var board game.Board
	fmt.Println("Loading starting position...")
	board.LoadFEN(game.StartFEN)
	board.Draw()

	// Test Move Generation
	fmt.Println("Generating Moves for White...")
	moves := board.GeneratePseudoLegalMoves(game.White)

	for _, m := range moves {
		// Use our new Helper function to print readable moves
		fmt.Printf("%s moves %s -> %s\n",
			m.Piece,
			game.IndexToCoord(m.From),
			game.IndexToCoord(m.To))
	}
}
