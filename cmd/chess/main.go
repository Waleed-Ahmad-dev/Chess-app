package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/Waleed-Ahmad-dev/Chess-app/internal/game"
	"github.com/Waleed-Ahmad-dev/Chess-app/internal/server"
	"github.com/Waleed-Ahmad-dev/Chess-app/internal/sound"
)

func main() {
	// Check if user wants web mode
	if len(os.Args) > 1 && os.Args[1] == "web" {
		fmt.Println("Starting Chess Web Server...")
		fmt.Println("Open your browser and go to: http://localhost:8080")
		server.StartServer()
		return
	}

	// Check for sound disable flag
	soundEnabled := true
	for _, arg := range os.Args[1:] {
		if arg == "--no-sound" || arg == "-n" {
			soundEnabled = false
			break
		}
	}

	if soundEnabled {
		sound.EnableSound()
	} else {
		sound.DisableSound()
		fmt.Println("Sound disabled. Use --no-sound or -n to disable sound.")
	}

	// Original CLI mode
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
				if soundEnabled {
					sound.PlaySound(sound.SoundCheckmate)
				}
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
			if soundEnabled {
				sound.PlaySound(sound.SoundIllegal)
			}
			reader.ReadString('\n') // Wait for user to read error
			continue
		}

		// 6. Execute Move
		result := gameInstance.MakeMove(move)

		// 7. Play appropriate sound
		if soundEnabled {
			playMoveSound(result)
		}
	}
}

// Helper to clear the terminal screen
func clearScreen() {
	// Standard ANSI escape code to clear screen
	fmt.Print("\033[H\033[2J")

	// Fallback for Windows
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

// playMoveSound plays the appropriate sound for a move
func playMoveSound(result game.MoveResult) {
	if result.WasCheckmate {
		sound.PlaySound(sound.SoundCheckmate)
		return
	}

	if result.WasCheck {
		sound.PlaySound(sound.SoundCheck)
		return
	}

	if result.WasCastle {
		sound.PlaySound(sound.SoundCastle)
		return
	}

	if result.WasPromotion {
		sound.PlaySound(sound.SoundPromote)
		return
	}

	if result.WasCapture {
		sound.PlaySound(sound.SoundCapture)
		return
	}

	// Default move sound
	sound.PlaySound(sound.SoundMove)
}
