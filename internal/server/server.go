package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"sync"
	"time"

	"io"

	"github.com/Waleed-Ahmad-dev/Chess-app/internal/game"
	"github.com/Waleed-Ahmad-dev/Chess-app/internal/sound"
	"github.com/gorilla/websocket"
)

//go:embed assets
var assetsFS embed.FS

var (
	sessions     = make(map[string]*game.Game)
	sessionTimes = make(map[string]time.Time)
	mu           sync.RWMutex
	soundManager = sound.NewWebSoundManager()

	// WebSocket Upgrader
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for simplicity
		},
	}

	// Clients map: Maps a WebSocket connection to a Session ID
	clients = make(map[*websocket.Conn]string)
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
	SoundType      string   `json:"soundType,omitempty"` // Added for sound feedback
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

// SoundData provides sound URLs to frontend
type SoundData struct {
	Move      string `json:"move"`
	Capture   string `json:"capture"`
	Castle    string `json:"castle"`
	Check     string `json:"check"`
	Checkmate string `json:"checkmate"`
	Illegal   string `json:"illegal"`
	Promote   string `json:"promote"`
}

func StartServer() {
	// Start cleanup goroutine
	go cleanupSessions()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", handleConnections) // WebSocket Endpoint
	http.HandleFunc("/api/new-game", handleNewGame)
	http.HandleFunc("/api/game-state", handleGameState)
	http.HandleFunc("/api/make-move", handleMakeMove)
	http.HandleFunc("/api/undo-move", handleUndoMove)
	http.HandleFunc("/api/sounds", handleGetSounds)

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// handleConnections upgrades HTTP to WebSocket
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// 1. Upgrade the connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket Upgrade Error: %v", err)
		return
	}
	defer ws.Close()

	// 2. Get Session ID from query params (e.g., /ws?sessionId=xyz)
	sessionID := r.URL.Query().Get("sessionId")
	if sessionID == "" {
		log.Println("WebSocket connection rejected: No Session ID")
		return
	}

	// 3. Register client
	mu.Lock()
	clients[ws] = sessionID
	mu.Unlock()

	log.Printf("Client connected to session: %s", sessionID)

	// 4. Keep connection alive / Listen for close
	for {
		// We don't expect messages from client via WS for now (using HTTP for moves),
		// but we need to read to handle Close messages / Pings.
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Client disconnected: %v", err)
			mu.Lock()
			delete(clients, ws)
			mu.Unlock()
			break
		}
	}
}

// broadcastState sends the current game state to all clients in the given session
func broadcastState(sessionID string, g *game.Game, soundType string) {
	mu.Lock()
	defer mu.Unlock()

	state := getGameState(g, soundType)

	for client, clientSessionID := range clients {
		if clientSessionID == sessionID {
			err := client.WriteJSON(state)
			if err != nil {
				log.Printf("WebSocket Write Error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	assets, err := fs.Sub(assetsFS, "assets")
	if err != nil {
		log.Printf("Failed to sub-fs assets: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if requesting an MP3 file
	if r.URL.Path == "/move.mp3" || r.URL.Path == "/capture.mp3" ||
		r.URL.Path == "/castle.mp3" || r.URL.Path == "/check.mp3" ||
		r.URL.Path == "/checkmate.mp3" || r.URL.Path == "/illegal.mp3" ||
		r.URL.Path == "/promote.mp3" {
		// Serve the specific sound file
		file, err := assets.Open(r.URL.Path[1:]) // Remove leading slash
		if err != nil {
			http.Error(w, "Sound not found", http.StatusNotFound)
			return
		}
		defer file.Close()

		w.Header().Set("Content-Type", "audio/mpeg")
		http.ServeContent(w, r, r.URL.Path, time.Now(), file.(io.ReadSeeker))
		return
	}

	http.FileServer(http.FS(assets)).ServeHTTP(w, r)
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
		"state":     getGameState(g, ""),
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

	soundType := r.URL.Query().Get("soundType")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getGameState(g, soundType))
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
			State:   getGameState(g, "illegal"),
		})
		return
	}

	result := g.MakeMove(move)

	// Determine sound type based on move result
	soundType := ""
	if result.WasCheckmate {
		soundType = "checkmate"
	} else if result.WasCheck {
		soundType = "check"
	} else if result.WasCastle {
		soundType = "castle"
	} else if result.WasPromotion {
		soundType = "promote"
	} else if result.WasCapture {
		soundType = "capture"
	} else {
		soundType = "move"
	}

	mu.Lock()
	sessionTimes[req.SessionID] = time.Now()
	mu.Unlock()

	// --- Broadcast Update to all clients in this session ---
	go broadcastState(req.SessionID, g, soundType)
	// -----------------------------------------------------

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MoveResponse{
		Success: true,
		State:   getGameState(g, soundType),
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

	// --- Broadcast Update to all clients in this session ---
	go broadcastState(req.SessionID, g, "move")
	// -----------------------------------------------------

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MoveResponse{
		Success: true,
		State:   getGameState(g, ""),
	})
}

// New endpoint to get sound URLs
func handleGetSounds(w http.ResponseWriter, r *http.Request) {
	sounds := SoundData{
		Move:      "/move.mp3",
		Capture:   "/capture.mp3",
		Castle:    "/castle.mp3",
		Check:     "/check.mp3",
		Checkmate: "/checkmate.mp3",
		Illegal:   "/illegal.mp3",
		Promote:   "/promote.mp3",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sounds)
}

func getGameState(g *game.Game, soundType string) GameStateResponse {
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
		SoundType:      soundType,
	}
}

func getCapturedPieces(g *game.Game) []string {
	captured := []string{}

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
