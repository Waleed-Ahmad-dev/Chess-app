package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Waleed-Ahmad-dev/Chess-app/internal/game"
)

var (
	sessions     = make(map[string]*game.Game)
	sessionTimes = make(map[string]time.Time)
	mu           sync.RWMutex
)

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

type MoveRequest struct {
	SessionID string `json:"sessionId"`
	Move      string `json:"move"`
}

type MoveResponse struct {
	Success bool              `json:"success"`
	Error   string            `json:"error,omitempty"`
	State   GameStateResponse `json:"state"`
}

func StartServer() {
	// Start cleanup goroutine
	go cleanupSessions()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/api/new-game", handleNewGame)
	http.HandleFunc("/api/game-state", handleGameState)
	http.HandleFunc("/api/make-move", handleMakeMove)
	http.HandleFunc("/api/undo-move", handleUndoMove)

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlContent))
}

func handleNewGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionID := generateSessionID()
	g := game.NewGame()

	mu.Lock()
	sessions[sessionID] = g
	sessionTimes[sessionID] = time.Now()
	mu.Unlock()

	response := map[string]interface{}{
		"sessionId": sessionID,
		"state":     getGameState(g),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleGameState(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("sessionId")
	if sessionID == "" {
		http.Error(w, "Session ID required", http.StatusBadRequest)
		return
	}

	mu.RLock()
	g, exists := sessions[sessionID]
	mu.RUnlock()

	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getGameState(g))
}

func handleMakeMove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req MoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	mu.RLock()
	g, exists := sessions[req.SessionID]
	mu.RUnlock()

	if !exists {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MoveResponse{
			Success: false,
			Error:   "Session not found",
		})
		return
	}

	legalMoves := g.GenerateLegalMoves()
	move, err := game.ParseMove(req.Move, legalMoves)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MoveResponse{
			Success: false,
			Error:   err.Error(),
			State:   getGameState(g),
		})
		return
	}

	g.MakeMove(move)

	mu.Lock()
	sessionTimes[req.SessionID] = time.Now()
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MoveResponse{
		Success: true,
		State:   getGameState(g),
	})
}

func handleUndoMove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SessionID string `json:"sessionId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	mu.RLock()
	g, exists := sessions[req.SessionID]
	mu.RUnlock()

	if !exists {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MoveResponse{
			Success: false,
			Error:   "Session not found",
		})
		return
	}

	g.UndoMove()

	mu.Lock()
	sessionTimes[req.SessionID] = time.Now()
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MoveResponse{
		Success: true,
		State:   getGameState(g),
	})
}

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

func generateSessionID() string {
	return fmt.Sprintf("session_%d_%d", len(sessions)+1, time.Now().UnixNano())
}

func cleanupSessions() {
	for {
		time.Sleep(1 * time.Hour)
		mu.Lock()
		now := time.Now()
		for id, lastActive := range sessionTimes {
			if now.Sub(lastActive) > 24*time.Hour {
				delete(sessions, id)
				delete(sessionTimes, id)
			}
		}
		mu.Unlock()
	}
}
