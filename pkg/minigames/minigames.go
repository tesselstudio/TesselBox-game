package minigames

import (
	"fmt"
	"time"
)

// MinigameType represents the type of minigame
type MinigameType int

const (
	MinigameSpleef MinigameType = iota
	MinigameParkour
	MinigameCTF
	MinigameTNTRun
	MinigameMobArena
	MinigameDuel
)

// String returns minigame name
func (m MinigameType) String() string {
	switch m {
	case MinigameSpleef:
		return "Spleef"
	case MinigameParkour:
		return "Parkour"
	case MinigameCTF:
		return "Capture The Flag"
	case MinigameTNTRun:
		return "TNT Run"
	case MinigameMobArena:
		return "Mob Arena"
	case MinigameDuel:
		return "Duel"
	}
	return "Unknown"
}

// GameStatus represents the status of a game
type GameStatus int

const (
	GameWaiting GameStatus = iota
	GameStarting
	GameActive
	GameEnded
)

// PlayerScore represents a player's score in a minigame
type PlayerScore struct {
	PlayerID   string    `json:"player_id"`
	PlayerName string    `json:"player_name"`
	Score      int       `json:"score"`
	Kills      int       `json:"kills"`
	Deaths     int       `json:"deaths"`
	Place      int       `json:"place"` // 1st, 2nd, 3rd, etc
}

// Minigame represents a minigame instance
type Minigame struct {
	ID          string       `json:"id"`
	Type        MinigameType `json:"type"`
	Name        string       `json:"name"`
	WorldID     string       `json:"world_id"`
	
	// Players
	Players     []string     `json:"players"` // Player IDs
	Scores      []PlayerScore `json:"scores"`
	
	// State
	Status      GameStatus   `json:"status"`
	StartTime   time.Time    `json:"start_time"`
	EndTime     *time.Time   `json:"end_time,omitempty"`
	
	// Settings
	MaxPlayers  int          `json:"max_players"`
	MinPlayers  int          `json:"min_players"`
	Duration    time.Duration `json:"duration"`
	
	// Meta
	CreatedBy   string       `json:"created_by"`
	CreatedAt   time.Time    `json:"created_at"`
	IsPublic    bool         `json:"is_public"`
}

// NewMinigame creates a new minigame
func NewMinigame(id string, gameType MinigameType, name, worldID, createdBy string) *Minigame {
	return &Minigame{
		ID:         id,
		Type:       gameType,
		Name:       name,
		WorldID:    worldID,
		Players:    make([]string, 0),
		Scores:     make([]PlayerScore, 0),
		Status:     GameWaiting,
		MaxPlayers: 8,
		MinPlayers: 2,
		Duration:   10 * time.Minute,
		CreatedBy:  createdBy,
		CreatedAt:  time.Now(),
		IsPublic:   true,
	}
}

// AddPlayer adds a player to the game
func (m *Minigame) AddPlayer(playerID, playerName string) error {
	if m.Status != GameWaiting && m.Status != GameStarting {
		return fmt.Errorf("game has already started")
	}
	
	// Check if full
	if len(m.Players) >= m.MaxPlayers {
		return fmt.Errorf("game is full")
	}
	
	// Check if already in
	for _, id := range m.Players {
		if id == playerID {
			return fmt.Errorf("already in game")
		}
	}
	
	m.Players = append(m.Players, playerID)
	
	// Add score entry
	m.Scores = append(m.Scores, PlayerScore{
		PlayerID:   playerID,
		PlayerName: playerName,
		Score:      0,
		Kills:      0,
		Deaths:     0,
	})
	
	return nil
}

// RemovePlayer removes a player
func (m *Minigame) RemovePlayer(playerID string) error {
	for i, id := range m.Players {
		if id == playerID {
			m.Players = append(m.Players[:i], m.Players[i+1:]...)
			
			// Remove from scores
			for j, score := range m.Scores {
				if score.PlayerID == playerID {
					m.Scores = append(m.Scores[:j], m.Scores[j+1:]...)
					break
				}
			}
			
			return nil
		}
	}
	return fmt.Errorf("player not in game")
}

// IsPlayerIn checks if player is in game
func (m *Minigame) IsPlayerIn(playerID string) bool {
	for _, id := range m.Players {
		if id == playerID {
			return true
		}
	}
	return false
}

// Start starts the game
func (m *Minigame) Start() error {
	if m.Status != GameWaiting {
		return fmt.Errorf("game cannot be started")
	}
	
	if len(m.Players) < m.MinPlayers {
		return fmt.Errorf("not enough players (need %d, have %d)", m.MinPlayers, len(m.Players))
	}
	
	m.Status = GameStarting
	
	// Could add countdown here
	
	m.Status = GameActive
	m.StartTime = time.Now()
	
	return nil
}

// End ends the game
func (m *Minigame) End() {
	if m.Status == GameEnded {
		return
	}
	
	now := time.Now()
	m.Status = GameEnded
	m.EndTime = &now
	
	// Sort scores and assign places
	m.sortScores()
}

// sortScores sorts players by score and assigns places
func (m *Minigame) sortScores() {
	// Bubble sort by score
	for i := 0; i < len(m.Scores); i++ {
		for j := i + 1; j < len(m.Scores); j++ {
			if m.Scores[i].Score < m.Scores[j].Score {
				m.Scores[i], m.Scores[j] = m.Scores[j], m.Scores[i]
			}
		}
	}
	
	// Assign places
	for i := range m.Scores {
		m.Scores[i].Place = i + 1
	}
}

// UpdateScore updates a player's score
func (m *Minigame) UpdateScore(playerID string, scoreDelta, killsDelta, deathsDelta int) {
	for i := range m.Scores {
		if m.Scores[i].PlayerID == playerID {
			m.Scores[i].Score += scoreDelta
			m.Scores[i].Kills += killsDelta
			m.Scores[i].Deaths += deathsDelta
			return
		}
	}
}

// GetWinner returns the winner (1st place)
func (m *Minigame) GetWinner() *PlayerScore {
	if len(m.Scores) == 0 {
		return nil
	}
	return &m.Scores[0]
}

// GetPlayerScore gets a player's score
func (m *Minigame) GetPlayerScore(playerID string) *PlayerScore {
	for i := range m.Scores {
		if m.Scores[i].PlayerID == playerID {
			return &m.Scores[i]
		}
	}
	return nil
}

// TimeRemaining returns time until game ends
func (m *Minigame) TimeRemaining() time.Duration {
	if m.Status != GameActive {
		return 0
	}
	
	elapsed := time.Since(m.StartTime)
	remaining := m.Duration - elapsed
	
	if remaining < 0 {
		return 0
	}
	return remaining
}

// MinigameManager manages minigames
type MinigameManager struct {
	games       map[string]*Minigame
	byType      map[MinigameType][]string
	byPlayer    map[string]string // PlayerID -> GameID
	
	gameCounter int
}

// NewMinigameManager creates a new minigame manager
func NewMinigameManager() *MinigameManager {
	return &MinigameManager{
		games:       make(map[string]*Minigame),
		byType:      make(map[MinigameType][]string),
		byPlayer:    make(map[string]string),
		gameCounter: 0,
	}
}

// CreateGame creates a new minigame
func (mm *MinigameManager) CreateGame(gameType MinigameType, name, worldID, createdBy string) *Minigame {
	mm.gameCounter++
	gameID := fmt.Sprintf("game_%d_%d", mm.gameCounter, time.Now().Unix())
	
	game := NewMinigame(gameID, gameType, name, worldID, createdBy)
	
	mm.games[gameID] = game
	mm.byType[gameType] = append(mm.byType[gameType], gameID)
	
	return game
}

// GetGame gets a game by ID
func (mm *MinigameManager) GetGame(gameID string) (*Minigame, bool) {
	game, exists := mm.games[gameID]
	return game, exists
}

// GetGamesByType gets games of a specific type
func (mm *MinigameManager) GetGamesByType(gameType MinigameType) []*Minigame {
	gameIDs := mm.byType[gameType]
	games := make([]*Minigame, 0, len(gameIDs))
	
	for _, id := range gameIDs {
		if game, exists := mm.games[id]; exists {
			games = append(games, game)
		}
	}
	
	return games
}

// GetWaitingGames gets games waiting for players
func (mm *MinigameManager) GetWaitingGames() []*Minigame {
	result := make([]*Minigame, 0)
	for _, game := range mm.games {
		if game.Status == GameWaiting {
			result = append(result, game)
		}
	}
	return result
}

// GetPlayerGame gets the game a player is in
func (mm *MinigameManager) GetPlayerGame(playerID string) (*Minigame, bool) {
	gameID, exists := mm.byPlayer[playerID]
	if !exists {
		return nil, false
	}
	return mm.GetGame(gameID)
}

// IsInGame checks if player is in a game
func (mm *MinigameManager) IsInGame(playerID string) bool {
	_, exists := mm.byPlayer[playerID]
	return exists
}

// JoinGame adds a player to a game
func (mm *MinigameManager) JoinGame(gameID, playerID, playerName string) error {
	game, exists := mm.GetGame(gameID)
	if !exists {
		return fmt.Errorf("game not found")
	}
	
	if mm.IsInGame(playerID) {
		return fmt.Errorf("already in a game")
	}
	
	if err := game.AddPlayer(playerID, playerName); err != nil {
		return err
	}
	
	mm.byPlayer[playerID] = gameID
	
	return nil
}

// LeaveGame removes a player from their game
func (mm *MinigameManager) LeaveGame(playerID string) error {
	game, exists := mm.GetPlayerGame(playerID)
	if !exists {
		return fmt.Errorf("not in a game")
	}
	
	if err := game.RemovePlayer(playerID); err != nil {
		return err
	}
	
	delete(mm.byPlayer, playerID)
	
	// Clean up empty games
	if len(game.Players) == 0 {
		mm.DeleteGame(game.ID)
	}
	
	return nil
}

// DeleteGame deletes a game
func (mm *MinigameManager) DeleteGame(gameID string) {
	game, exists := mm.GetGame(gameID)
	if !exists {
		return
	}
	
	// Remove players from byPlayer
	for _, playerID := range game.Players {
		delete(mm.byPlayer, playerID)
	}
	
	// Remove from byType
	for i, id := range mm.byType[game.Type] {
		if id == gameID {
			mm.byType[game.Type] = append(mm.byType[game.Type][:i], mm.byType[game.Type][i+1:]...)
			break
		}
	}
	
	// Remove game
	delete(mm.games, gameID)
}

// Update updates all active games
func (mm *MinigameManager) Update() {
	now := time.Now()
	
	for _, game := range mm.games {
		if game.Status == GameActive {
			// Check if game should end
			if now.After(game.StartTime.Add(game.Duration)) {
				game.End()
			}
		}
	}
}

// GetActiveGames returns all active games
func (mm *MinigameManager) GetActiveGames() []*Minigame {
	result := make([]*Minigame, 0)
	for _, game := range mm.games {
		if game.Status == GameActive {
			result = append(result, game)
		}
	}
	return result
}

// GetStats returns minigame statistics
func (mm *MinigameManager) GetStats() map[string]int {
	stats := map[string]int{
		"total":     0,
		"waiting":   0,
		"starting":  0,
		"active":    0,
		"ended":     0,
		"players":   0,
	}
	
	for _, game := range mm.games {
		stats["total"]++
		stats["players"] += len(game.Players)
		
		switch game.Status {
		case GameWaiting:
			stats["waiting"]++
		case GameStarting:
			stats["starting"]++
		case GameActive:
			stats["active"]++
		case GameEnded:
			stats["ended"]++
		}
	}
	
	return stats
}
