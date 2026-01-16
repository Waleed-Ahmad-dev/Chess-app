package game

import "fmt"

// IndexToCoord converts 0 -> "a1", 63 -> "h8"
func IndexToCoord(index int) string {
	rank := index / 8
	file := index % 8
	// 'a' is ASCII 97, '1' is ASCII 49
	return fmt.Sprintf("%c%d", 'a'+file, rank+1)
}

// CoordToIndex converts "a1" -> 0
func CoordToIndex(coord string) int {
	if len(coord) != 2 {
		return -1
	}
	file := int(coord[0] - 'a')
	rank := int(coord[1] - '1')
	return rank*8 + file
}

// IsOnBoard checks if a square is valid (0-63)
func IsOnBoard(sq int) bool {
	return sq >= 0 && sq < 64
}
