package game

import (
	"fmt"
	"strings"
)

// ParseMove takes a UCI string (e.g., "e2e4" or "a7a8q") and a list of legal moves.
// It returns the matching Move struct if found, and a boolean indicating success.
func ParseMove(input string, legalMoves []Move) (Move, error) {
	cleanInput := strings.TrimSpace(input)

	// Basic validation: UCI moves are 4 or 5 chars (e.g., "e2e4", "e7e8q")
	if len(cleanInput) < 4 || len(cleanInput) > 5 {
		return Move{}, fmt.Errorf("invalid format. Use UCI format like 'e2e4' or 'a7a8q'")
	}

	fromStr := cleanInput[:2]
	toStr := cleanInput[2:4]

	fromIdx := CoordToIndex(fromStr)
	toIdx := CoordToIndex(toStr)

	if fromIdx == -1 || toIdx == -1 {
		return Move{}, fmt.Errorf("invalid coordinates")
	}

	// Filter legal moves to find candidates that match From and To
	var candidates []Move
	for _, m := range legalMoves {
		if m.From == fromIdx && m.To == toIdx {
			candidates = append(candidates, m)
		}
	}

	if len(candidates) == 0 {
		return Move{}, fmt.Errorf("illegal move")
	}

	// If only one candidate, it's a standard move (or unique move)
	if len(candidates) == 1 {
		return candidates[0], nil
	}

	// If multiple candidates, it must be a promotion (Q, R, B, N available)
	// We need the 5th character to distinguish.
	if len(cleanInput) != 5 {
		return Move{}, fmt.Errorf("promotion detected. Please specify piece (q, r, b, n). Example: %sq", cleanInput)
	}

	promotionChar := rune(cleanInput[4])
	var targetType PieceType

	switch strings.ToLower(string(promotionChar)) {
	case "q":
		targetType = Queen
	case "r":
		targetType = Rook
	case "b":
		targetType = Bishop
	case "n":
		targetType = Knight
	default:
		return Move{}, fmt.Errorf("invalid promotion piece '%c'. Use q, r, b, or n", promotionChar)
	}

	// Find the specific promotion move
	for _, m := range candidates {
		if m.Promotion == targetType {
			return m, nil
		}
	}

	return Move{}, fmt.Errorf("could not match promotion move")
}
