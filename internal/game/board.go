package game

import "fmt"

const StartFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// Draw prints the board to the console
func (b *Board) Draw() {
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
}
