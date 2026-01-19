package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/Waleed-Ahmad-dev/Chess-app/internal/game"
	"github.com/Waleed-Ahmad-dev/Chess-app/internal/sound"
	"github.com/gorilla/websocket"
)

//go:embed assets
var assetsFS embed.FS

// Room represents a multiplayer lobby
type Room struct {
	ID      string
	Game    *game.Game
	Clients map[*websocket.Conn]game.Color // Map connection to player color
	Mutex   sync.RWMutex
	LastAct time.Time
}

var (
	rooms        = make(map[string]*Room)
	mu           sync.RWMutex // Global mutex for the rooms map
	soundManager = sound.NewWebSoundManager()

	// WebSocket Upgrader
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for simplicity
		},
	}
)

// API Structures

type CreateRoomResponse struct {
	RoomID string            `json:"roomId"`
	State  GameStateResponse `json:"state"`
}

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
	SoundType      string   `json:"soundType,omitempty"`
}

type InitMessage struct {
	Type   string            `json:"type"` // "init"
	Color  string            `json:"color"`
	RoomID string            `json:"roomId"`
	State  GameStateResponse `json:"state"`
}

type MoveRequest struct {
	RoomID string `json:"roomId"`
	Move   string `json:"move"`
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
	// Seed random for room codes
	rand.Seed(time.Now().UnixNano())

	// Start cleanup goroutine
	go cleanupRooms()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", handleConnections) // WebSocket Endpoint

	// New Room API
	http.HandleFunc("/api/create-room", handleCreateRoom)
	http.HandleFunc("/api/game-state", handleGameState)
	http.HandleFunc("/api/make-move", handleMakeMove)
	http.HandleFunc("/api/undo-move", handleUndoMove)
	http.HandleFunc("/api/sounds", handleGetSounds)

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// handleConnections upgrades HTTP to WebSocket and manages Room Joining
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// 1. Get Room ID
	roomID := r.URL.Query().Get("room")
	if roomID == "" {
		http.Error(w, "Room ID required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	room, exists := rooms[roomID]
	mu.Unlock()

	if !exists {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	// 2. Upgrade the connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket Upgrade Error: %v", err)
		return
	}
	defer ws.Close()

	// 3. Register client and Assign Color
	room.Mutex.Lock()
	clientCount := len(room.Clients)
	var assignedColor game.Color

	if clientCount == 0 {
		assignedColor = game.White
	} else if clientCount == 1 {
		assignedColor = game.Black
	} else {
		// Room full (Spectator mode could be added here, but for now we reject or just don't assign a playing color)
		// For this implementation, we allow connection but maybe frontend handles it as spectator
		// defaulting to White for view purposes, or a specific Spectator type
		assignedColor = game.White // Default fallback
	}

	room.Clients[ws] = assignedColor
	room.Mutex.Unlock()

	log.Printf("Client connected to Room: %s as %s", roomID, assignedColor)

	// 4. Send Initial Handshake (Assigned Color + Current State)
	initState := getGameState(room.Game, "")
	initMsg := InitMessage{
		Type:   "init",
		Color:  assignedColor.String(),
		RoomID: roomID,
		State:  initState,
	}
	ws.WriteJSON(initMsg)

	// 5. Keep connection alive / Listen for close
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Client disconnected from Room %s: %v", roomID, err)

			// Remove client from room
			room.Mutex.Lock()
			delete(room.Clients, ws)
			room.Mutex.Unlock()
			break
		}
	}
}

// broadcastState sends the current game state to all clients in the specific room
func broadcastState(room *Room, soundType string) {
	room.Mutex.RLock()
	defer room.Mutex.RUnlock()

	state := getGameState(room.Game, soundType)

	for client := range room.Clients {
		err := client.WriteJSON(state)
		if err != nil {
			log.Printf("WebSocket Write Error: %v", err)
			client.Close()
			// We can't safely delete from map while iterating with RLock,
			// relying on the ReadMessage loop to handle cleanup
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

	// Serve specific MP3 files if requested
	if r.URL.Path == "/move.mp3" || r.URL.Path == "/capture.mp3" ||
		r.URL.Path == "/castle.mp3" || r.URL.Path == "/check.mp3" ||
		r.URL.Path == "/checkmate.mp3" || r.URL.Path == "/illegal.mp3" ||
		r.URL.Path == "/promote.mp3" {

		file, err := assets.Open(r.URL.Path[1:])
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

// handleCreateRoom generates a 4-letter code and creates a new room
func handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	roomID := generateRoomCode()
	g := game.NewGame()

	newRoom := &Room{
		ID:      roomID,
		Game:    g,
		Clients: make(map[*websocket.Conn]game.Color),
		LastAct: time.Now(),
	}

	mu.Lock()
	rooms[roomID] = newRoom
	mu.Unlock()

	response := CreateRoomResponse{
		RoomID: roomID,
		State:  getGameState(g, ""),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleGameState(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("roomId")
	if roomID == "" {
		http.Error(w, "Room ID required", http.StatusBadRequest)
		return
	}

	mu.RLock()
	room, exists := rooms[roomID]
	mu.RUnlock()

	if !exists {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	soundType := r.URL.Query().Get("soundType")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getGameState(room.Game, soundType))
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
	room, exists := rooms[req.RoomID]
	mu.RUnlock()

	if !exists {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MoveResponse{
			Success: false,
			Error:   "Room not found",
		})
		return
	}

	// Lock the room logic
	room.Mutex.Lock()
	g := room.Game
	legalMoves := g.GenerateLegalMoves()
	move, err := game.ParseMove(req.Move, legalMoves)

	if err != nil {
		room.Mutex.Unlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MoveResponse{
			Success: false,
			Error:   err.Error(),
			State:   getGameState(g, "illegal"),
		})
		return
	}

	result := g.MakeMove(move)
	room.LastAct = time.Now()
	room.Mutex.Unlock()

	// Determine sound
	soundType := "move"
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
	}

	// Broadcast to clients in this room
	go broadcastState(room, soundType)

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
		RoomID string `json:"roomId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	mu.RLock()
	room, exists := rooms[req.RoomID]
	mu.RUnlock()

	if !exists {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MoveResponse{
			Success: false,
			Error:   "Room not found",
		})
		return
	}

	room.Mutex.Lock()
	room.Game.UndoMove()
	room.LastAct = time.Now()
	room.Mutex.Unlock()

	go broadcastState(room, "move")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MoveResponse{
		Success: true,
		State:   getGameState(room.Game, ""),
	})
}

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

// Helpers

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
	whitePieces := map[game.PieceType]int{game.Pawn: 8, game.Knight: 2, game.Bishop: 2, game.Rook: 2, game.Queen: 1, game.King: 1}
	blackPieces := map[game.PieceType]int{game.Pawn: 8, game.Knight: 2, game.Bishop: 2, game.Rook: 2, game.Queen: 1, game.King: 1}

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

// generateRoomCode creates a random 4-letter code
func generateRoomCode() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 4)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func cleanupRooms() {
	for {
		time.Sleep(1 * time.Hour)
		mu.Lock()
		now := time.Now()
		for id, room := range rooms {
			if now.Sub(room.LastAct) > 24*time.Hour {
				delete(rooms, id)
			}
		}
		mu.Unlock()
	}
}
