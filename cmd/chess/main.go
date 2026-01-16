package main

import (
	"fmt"

	"github.com/Waleed-Ahmad-dev/Chess-app/internal/game"
)

func main() {
	gameInstance := game.NewGame()

	// Setup a tricky position:
	// White King on e1. Black Rook on e8 (Check!).
	// White has a Bishop on c1 that CAN interpose on e3.
	// White King CAN move to f1.
	fmt.Println("Loading check test position...")
	// 8  r . . . r . . .
	// 7  . . . . . . . .
	// ...
	// 1  . . B . K . . .
	gameInstance.Board.LoadFEN("4r3/8/8/8/8/8/8/2B1K3 w - - 0 1")
	gameInstance.Board.Draw()

	fmt.Printf("Is White in Check? %v\n", gameInstance.Board.InCheck(game.White))

	fmt.Println("Generating LEGAL Moves for White...")
	moves := gameInstance.GenerateLegalMoves()

	for _, m := range moves {
		fmt.Printf("%s moves %s -> %s\n",
			m.Piece,
			game.IndexToCoord(m.From),
			game.IndexToCoord(m.To))
	}
}
