package game

import (
	"fmt"
	"math/rand"
)

// AI Config
const (
	MaxDepth = 3 // Search depth (3 is fast and decent for casual play)
	Infinity = 1000000
)

// Piece Values (Centipawns)
var pieceValues = map[PieceType]int{
	Pawn:   100,
	Knight: 320,
	Bishop: 330,
	Rook:   500,
	Queen:  900,
	King:   20000,
}

// Simple Position Bonuses (Simplified PST)
// Bonus for controlling center and advancing pieces
var centerBonus = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 5, 5, 5, 5, 5, 5, 0,
	0, 5, 10, 15, 15, 10, 5, 0,
	0, 5, 15, 25, 25, 15, 5, 0,
	0, 5, 15, 25, 25, 15, 5, 0,
	0, 5, 10, 15, 15, 10, 5, 0,
	0, 5, 5, 5, 5, 5, 5, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

// GetBestMove returns the best move for the current turn using Minimax
func (g *Game) GetBestMove(depth int) (Move, error) {
	legalMoves := g.GenerateLegalMoves()
	if len(legalMoves) == 0 {
		return Move{}, fmt.Errorf("no legal moves")
	}

	// Randomize simple moves to avoid identical games
	rand.Shuffle(len(legalMoves), func(i, j int) {
		legalMoves[i], legalMoves[j] = legalMoves[j], legalMoves[i]
	})

	bestMove := legalMoves[0]
	bestScore := -Infinity

	// Maximizing for current player (always treat current turn as maximizing in this logic)
	// We handle color perspective inside Evaluate

	alpha := -Infinity
	beta := Infinity

	for _, move := range legalMoves {
		// Clone game state effectively by using MakeMove/Undo logic manually or copy
		// Since MakeMove alters state, we need a clone or careful Undo.
		// For simplicity/safety, we will use the Clone method.
		tempGame := g.Clone()
		tempGame.MakeMove(move)

		// Call Minimax for opponent
		score := -tempGame.minimax(depth-1, alpha, beta, false)

		if score > bestScore {
			bestScore = score
			bestMove = move
		}
		if score > alpha {
			alpha = score
		}
		if alpha >= beta {
			break
		}
	}

	return bestMove, nil
}

func (g *Game) minimax(depth int, alpha, beta int, isMaximizing bool) int {
	if depth == 0 {
		return g.Evaluate()
	}

	legalMoves := g.GenerateLegalMoves()
	if len(legalMoves) == 0 {
		if g.Board.InCheck(g.Turn) {
			return -Infinity + (MaxDepth - depth) // Prefer faster checkmates
		}
		return 0 // Stalemate
	}

	if isMaximizing {
		maxEval := -Infinity
		for _, move := range legalMoves {
			tempGame := g.Clone()
			tempGame.MakeMove(move)
			eval := tempGame.minimax(depth-1, alpha, beta, false)
			if eval > maxEval {
				maxEval = eval
			}
			if eval > alpha {
				alpha = eval
			}
			if beta <= alpha {
				break
			}
		}
		return maxEval
	} else {
		minEval := Infinity
		for _, move := range legalMoves {
			tempGame := g.Clone()
			tempGame.MakeMove(move)
			// Negamax flip: score returns from opponent perspective, so we don't negate here
			// Standard minimax:
			eval := tempGame.minimax(depth-1, alpha, beta, true)
			if eval < minEval {
				minEval = eval
			}
			if eval < beta {
				beta = eval
			}
			if beta <= alpha {
				break
			}
		}
		return minEval
	}
}

// Evaluate calculates the board score relative to the player whose turn it is
func (g *Game) Evaluate() int {
	whiteScore := 0
	blackScore := 0

	for i, piece := range g.Board {
		if piece.Type == Empty {
			continue
		}

		val := pieceValues[piece.Type]
		// Add positional bonus
		if piece.Type != King {
			val += centerBonus[i]
		}

		if piece.Color == White {
			whiteScore += val
		} else {
			blackScore += val
		}
	}

	score := whiteScore - blackScore
	if g.Turn == Black {
		return -score // Return score relative to current player
	}
	return score
}

// Clone creates a deep copy of the game for the AI calculation
func (g *Game) Clone() *Game {
	newG := &Game{
		Board:           g.Board, // Array copy is by value in Go
		Turn:            g.Turn,
		EnPassantTarget: g.EnPassantTarget,
		Castling:        g.Castling,
	}
	// Copy slices
	newG.History = make([]Move, len(g.History))
	copy(newG.History, g.History)

	newG.StateHistory = make([]StateSnapshot, len(g.StateHistory))
	copy(newG.StateHistory, g.StateHistory)

	return newG
}
