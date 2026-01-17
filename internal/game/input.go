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

	// --- 1. Standard Move Matching ---
	var candidates []Move
	for _, m := range legalMoves {
		if m.From == fromIdx && m.To == toIdx {
			candidates = append(candidates, m)
		}
	}

	// --- 2. Special Castling Handling (King -> Rook) ---
	// If no standard move found, check if user tried to castle by clicking King then Rook
	if len(candidates) == 0 {
		candidates = checkAlternativeCastling(fromIdx, toIdx, legalMoves)
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

// checkAlternativeCastling allows castling by clicking King -> Rook (e.g. e1h1 -> e1g1)
func checkAlternativeCastling(from, to int, legalMoves []Move) []Move {
	// Map of King->Rook moves to their actual King->Dest moves
	// Key: KingFrom * 100 + RookTo, Value: Actual KingTo
	castlingMap := map[int]int{
		6063: 62, // White Short (e1->h1 implies e1->g1)
		6056: 58, // White Long  (e1->a1 implies e1->c1)
		407:  6,  // Black Short (e8->h8 implies e8->g8)
		400:  2,  // Black Long  (e8->a8 implies e8->c8)
	}

	key := from*100 + to
	targetSq, exists := castlingMap[key]

	if !exists {
		return nil
	}

	// Check if the actual castling move exists in legal moves
	var result []Move
	for _, m := range legalMoves {
		if m.From == from && m.To == targetSq && m.MoveType == MoveCastling {
			result = append(result, m)
		}
	}
	return result
}
