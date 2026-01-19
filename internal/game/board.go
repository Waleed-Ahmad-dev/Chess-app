package game

import "fmt"

const StartFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// Draw prints the board to the console based on the perspective
func (b *Board) Draw(perspective Color) {
	if perspective == White {
		fmt.Println("\n   A B C D E F G H")
		fmt.Println("  -----------------")

		// Loop from Rank 8 (top) down to Rank 1 (bottom)
		for rank := 7; rank >= 0; rank-- {
			fmt.Printf("%d |", rank+1)
			for file := 0; file < 8; file++ {
				index := rank*8 + file
				piece := b[index]
				fmt.Printf(" %s", piece.String())
			}
			fmt.Printf(" | %d\n", rank+1)
		}
		fmt.Println("  -----------------")
		fmt.Println("   A B C D E F G H")
	} else {
		// Black Perspective: Rank 1 at top, File H on left
		fmt.Println("\n   H G F E D C B A")
		fmt.Println("  -----------------")

		// Loop from Rank 1 (top visual) up to Rank 8
		for rank := 0; rank < 8; rank++ {
			fmt.Printf("%d |", rank+1)
			// Loop from File H (7) down to File A (0)
			for file := 7; file >= 0; file-- {
				index := rank*8 + file
				piece := b[index]
				fmt.Printf(" %s", piece.String())
			}
			fmt.Printf(" | %d\n", rank+1)
		}
		fmt.Println("  -----------------")
		fmt.Println("   H G F E D C B A")
	}
}
