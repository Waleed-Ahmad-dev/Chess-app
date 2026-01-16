package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/Waleed-Ahmad-dev/Chess-app/internal/game"
)

func main() {
	gameInstance := game.NewGame()
	reader := bufio.NewReader(os.Stdin)

	// --- The Game Loop ---
	for {
		clearScreen()
		gameInstance.Board.Draw()

		// 1. Generate Legal Moves for the current turn
		legalMoves := gameInstance.GenerateLegalMoves()

		// 2. Check Game Over Conditions
		if len(legalMoves) == 0 {
			if gameInstance.Board.InCheck(gameInstance.Turn) {
				fmt.Printf("\nCheckmate! %s wins!\n", getOpponentColor(gameInstance.Turn))
			} else {
				fmt.Println("\nStalemate! The game is a draw.")
			}
			break
		}

		// 3. Print Status
		if gameInstance.Board.InCheck(gameInstance.Turn) {
			fmt.Println("\n⚠️  CHECK! ⚠️")
		}
		fmt.Printf("\n%s's turn to move.\n", gameInstance.Turn)
		fmt.Print("Enter move (e.g., 'e2e4'): ")

		// 4. Read Input
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" || input == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		// 5. Parse & Validate Move
		move, err := game.ParseMove(input, legalMoves)
		if err != nil {
			fmt.Printf("Error: %v\nPress Enter to try again...", err)
			reader.ReadString('\n') // Wait for user to read error
			continue
		}

		// 6. Execute Move
		gameInstance.MakeMove(move)
	}
}

// Helper to clear the terminal screen
func clearScreen() {
	// Standard ANSI escape code to clear screen (Works on Arch Linux/Mac/Most Unix)
	fmt.Print("\033[H\033[2J")

	// Fallback for Windows if needed (though ANSI often works there too now)
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func getOpponentColor(c game.Color) string {
	if c == game.White {
		return "Black"
	}
	return "White"
}
