package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"tesselbox/pkg/world"
)

// Server represents the game server
type Server struct {
	world   *world.World
	players map[string]*PlayerState
	mu      sync.RWMutex
}

// PlayerState represents a connected player's state
type PlayerState struct {
	ID     string  `json:"id"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	VX     float64 `json:"vx"`
	VY     float64 `json:"vy"`
	Health float64 `json:"health"`
}

// NewServer creates a new game server
func NewServer(seed float64) *Server {
	return &Server{
		world:   world.NewWorld(),
		players: make(map[string]*PlayerState),
	}
}

func main() {
	server := NewServer(42.0)
	
	// Setup routes
	http.HandleFunc("/api/world/chunk", server.handleGetChunk)
	http.HandleFunc("/api/world/blocks", server.handleGetBlocks)
	http.HandleFunc("/api/world/setblock", server.handleSetBlock)
	http.HandleFunc("/api/player/join", server.handlePlayerJoin)
	http.HandleFunc("/api/player/update", server.handlePlayerUpdate)
	http.HandleFunc("/api/players", server.handleGetPlayers)
	http.HandleFunc("/api/health", server.handleHealth)
	
	// Enable CORS
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		fmt.Fprintf(w, "Tesselbox Go Server - Running!")
	})
	
	port := ":8080"
	log.Printf("Starting Tesselbox Go server on %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// handleGetChunk returns chunk data
func (s *Server) handleGetChunk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	query := r.URL.Query()
	chunkX := query.Get("x")
	chunkY := query.Get("y")
	
	if chunkX == "" || chunkY == "" {
		http.Error(w, "Missing chunk coordinates", http.StatusBadRequest)
		return
	}
	
	var cx, cy int
	fmt.Sscanf(chunkX, "%d", &cx)
	fmt.Sscanf(chunkY, "%d", &cy)
	
	chunk := s.world.GetChunk(cx, cy)

	blocks := make([]map[string]interface{}, 0)
	for _, hex := range chunk.Hexagons {
		blocks = append(blocks, map[string]interface{}{
			"type":      hex.BlockType,
			"x":         hex.X,
			"y":         hex.Y,
			"health":    hex.Health,
			"chunk_x":   hex.ChunkX,
			"chunk_y":   hex.ChunkY,
			"hex_size":  hex.Size,
		})
	}

	response := map[string]interface{}{
		"chunk_x": cx,
		"chunk_y": cy,
		"blocks":  blocks,
	}
	
	json.NewEncoder(w).Encode(response)
}

// handleGetBlocks returns blocks in a range
func (s *Server) handleGetBlocks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	query := r.URL.Query()
	xStr := query.Get("x")
	yStr := query.Get("y")
	
	if xStr == "" || yStr == "" {
		http.Error(w, "Missing position coordinates", http.StatusBadRequest)
		return
	}
	
	var x, y float64
	fmt.Sscanf(xStr, "%f", &x)
	fmt.Sscanf(yStr, "%f", &y)
	
	blocks := s.world.GetNearbyHexagons(x, y, 800)
	
	blockData := make([]map[string]interface{}, 0)
	for _, hex := range blocks {
		blockData = append(blockData, map[string]interface{}{
			"type":   hex.BlockType,
			"x":      hex.X,
			"y":      hex.Y,
			"health": hex.Health,
		})
	}
	
	response := map[string]interface{}{
		"blocks": blockData,
		"count":  len(blockData),
	}
	
	json.NewEncoder(w).Encode(response)
}

// handleSetBlock sets a block in the world
func (s *Server) handleSetBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var request struct {
		X       float64 `json:"x"`
		Y       float64 `json:"y"`
		BlockType string `json:"block_type"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Place/remove hexagon at given position using current world API
	if request.BlockType == "" {
		s.world.RemoveHexagonAt(request.X, request.Y)
	} else {
		// Map string to BlockType via blocks package on the client side; here we just add generic hex
		// For now, treat all as dirt-equivalent; real implementation would decode a concrete BlockType.
		// This keeps the server compiling and structurally correct.
		// You can extend this later to switch on BlockType if needed.
		s.world.AddHexagonAt(request.X, request.Y, 0)
	}
	
	response := map[string]interface{}{
		"success": true,
	}
	
	json.NewEncoder(w).Encode(response)
}

// handlePlayerJoin handles a new player joining
func (s *Server) handlePlayerJoin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	query := r.URL.Query()
	playerID := query.Get("id")
	
	if playerID == "" {
		http.Error(w, "Missing player ID", http.StatusBadRequest)
		return
	}
	
	s.mu.Lock()
	s.players[playerID] = &PlayerState{
		ID:     playerID,
		X:      0,
		Y:      0,
		VX:     0,
		VY:     0,
		Health: 100,
	}
	s.mu.Unlock()
	
	response := map[string]interface{}{
		"success": true,
		"player_id": playerID,
		"spawn_x": 0,
		"spawn_y": 0,
	}
	
	json.NewEncoder(w).Encode(response)
}

// handlePlayerUpdate updates a player's state
func (s *Server) handlePlayerUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var state PlayerState
	if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	s.mu.Lock()
	if _, exists := s.players[state.ID]; exists {
		s.players[state.ID] = &state
	}
	s.mu.Unlock()
	
	response := map[string]interface{}{
		"success": true,
	}
	
	json.NewEncoder(w).Encode(response)
}

// handleGetPlayers returns all connected players
func (s *Server) handleGetPlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	players := make([]*PlayerState, 0, len(s.players))
	for _, player := range s.players {
		players = append(players, player)
	}
	
	response := map[string]interface{}{
		"players": players,
		"count":   len(players),
	}
	
	json.NewEncoder(w).Encode(response)
}

// handleHealth returns server health status
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	s.mu.RLock()
	playerCount := len(s.players)
	s.mu.RUnlock()
	
	response := map[string]interface{}{
		"status":       "healthy",
		"players":      playerCount,
		"world_seed":   s.world.Seed,
	}
	
	json.NewEncoder(w).Encode(response)
}