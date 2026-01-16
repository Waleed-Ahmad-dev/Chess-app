package main

import (
	"fmt"

	"github.com/Waleed-Ahmad-dev/Chess-app/internal/game"
)

func main() {
	// 1. Create a new board
	var board game.Board

	// 2. Load the starting position
	fmt.Println("Loading starting position...")
	board.LoadFEN(game.StartFEN)

	// 3. Draw the board
	board.Draw()
}
