package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/Waleed-Ahmad-dev/Chess-app/internal/game"
)

var globalGame *game.Game

// Structures to match server.go response format
type GameStateResponse struct {
	Board          []string `json:"board"`
	Turn           string   `json:"turn"`
	InCheck        bool     `json:"inCheck"`
	LegalMoves     []string `json:"legalMoves"`
	GameOver       bool     `json:"gameOver"`
	Winner         string   `json:"winner,omitempty"`
	IsStalemate    bool     `json:"isStalemate"`
	LastMove       string   `json:"lastMove,omitempty"`
	CapturedPieces []string `json:"capturedPieces"`
}

type MoveResponse struct {
	Success bool              `json:"success"`
	Error   string            `json:"error,omitempty"`
	State   GameStateResponse `json:"state"`
}

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("goNewGame", js.FuncOf(goNewGame))
	js.Global().Set("goMakeMove", js.FuncOf(goMakeMove))
	js.Global().Set("goGetState", js.FuncOf(goGetState))

	fmt.Println("Wasm Go Chess Initialized")
	<-c
}

func goNewGame(this js.Value, args []js.Value) interface{} {
	globalGame = game.NewGame()

	res := map[string]interface{}{
		"sessionId": "local_wasm_session",
		"state":     getGameState(globalGame),
	}

	encoded, _ := json.Marshal(res)
	return string(encoded)
}

func goMakeMove(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return errorResponse("Missing move argument")
	}
	moveStr := args[0].String()

	if globalGame == nil {
		return errorResponse("Game not initialized")
	}

	legalMoves := globalGame.GenerateLegalMoves()
	move, err := game.ParseMove(moveStr, legalMoves)
	if err != nil {
		return errorResponse(err.Error())
	}

	globalGame.MakeMove(move)

	response := MoveResponse{
		Success: true,
		State:   getGameState(globalGame),
	}
	encoded, _ := json.Marshal(response)
	return string(encoded)
}

func goGetState(this js.Value, args []js.Value) interface{} {
	if globalGame == nil {
		return errorResponse("Game not initialized")
	}
	encoded, _ := json.Marshal(getGameState(globalGame))
	return string(encoded)
}

func errorResponse(msg string) string {
	resp := MoveResponse{
		Success: false,
		Error:   msg,
	}
	if globalGame != nil {
		resp.State = getGameState(globalGame)
	}
	encoded, _ := json.Marshal(resp)
	return string(encoded)
}

// Logic copied/adapted from server.go to ensure consistent JSON structure
func getGameState(g *game.Game) GameStateResponse {
	board := make([]string, 64)
	for i := 0; i < 64; i++ {
		board[i] = g.Board[i].String()
	}

	legalMoves := g.GenerateLegalMoves()
	legalMovesStr := make([]string, len(legalMoves))
	for i, m := range legalMoves {
		legalMovesStr[i] = fmt.Sprintf("%s%s", game.IndexToCoord(m.From), game.IndexToCoord(m.To))
	}

	gameOver := len(legalMoves) == 0
	winner := ""
	isStalemate := false

	if gameOver {
		if g.Board.InCheck(g.Turn) {
			if g.Turn == game.White {
				winner = "Black"
			} else {
				winner = "White"
			}
		} else {
			isStalemate = true
		}
	}

	lastMove := ""
	if len(g.History) > 0 {
		m := g.History[len(g.History)-1]
		lastMove = fmt.Sprintf("%s%s", game.IndexToCoord(m.From), game.IndexToCoord(m.To))
	}

	capturedPieces := getCapturedPieces(g)

	return GameStateResponse{
		Board:          board,
		Turn:           g.Turn.String(),
		InCheck:        g.Board.InCheck(g.Turn),
		LegalMoves:     legalMovesStr,
		GameOver:       gameOver,
		Winner:         winner,
		IsStalemate:    isStalemate,
		LastMove:       lastMove,
		CapturedPieces: capturedPieces,
	}
}

func getCapturedPieces(g *game.Game) []string {
	captured := []string{}

	// Count pieces on board
	whitePieces := map[game.PieceType]int{
		game.Pawn: 8, game.Knight: 2, game.Bishop: 2,
		game.Rook: 2, game.Queen: 1, game.King: 1,
	}
	blackPieces := map[game.PieceType]int{
		game.Pawn: 8, game.Knight: 2, game.Bishop: 2,
		game.Rook: 2, game.Queen: 1, game.King: 1,
	}

	for i := 0; i < 64; i++ {
		piece := g.Board[i]
		if piece.Type != game.Empty {
			if piece.Color == game.White {
				whitePieces[piece.Type]--
			} else {
				blackPieces[piece.Type]--
			}
		}
	}

	// Add captured pieces
	for pType, count := range whitePieces {
		for i := 0; i < count; i++ {
			captured = append(captured, "w"+string(game.Piece{Type: pType, Color: game.White}.String()))
		}
	}
	for pType, count := range blackPieces {
		for i := 0; i < count; i++ {
			captured = append(captured, "b"+string(game.Piece{Type: pType, Color: game.Black}.String()))
		}
	}

	return captured
}
