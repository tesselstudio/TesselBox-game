package pvp

import (
	"fmt"
	"time"

	"tesselbox/pkg/economy"
)

// ArenaType represents arena format
type ArenaType int

const (
	ArenaFreeForAll ArenaType = iota // Everyone vs everyone
	ArenaTeam                        // Team vs team
	Arena1v1                         // 1v1 matches
	ArenaTournament                  // Bracket tournament
)

// ArenaStatus represents arena state
type ArenaStatus int

const (
	ArenaWaiting ArenaStatus = iota
	ArenaInProgress
	ArenaCooldown
)

// ArenaPlayer represents a player in arena
type ArenaPlayer struct {
	PlayerID   string    `json:"player_id"`
	PlayerName string    `json:"player_name"`
	Kills      int       `json:"kills"`
	Deaths     int       `json:"deaths"`
	Score      int       `json:"score"`
	Team       int       `json:"team,omitempty"`
	JoinedAt   time.Time `json:"joined_at"`
	IsAlive    bool      `json:"is_alive"`
}

// Arena represents a PvP arena
type Arena struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	WorldID string    `json:"world_id"`
	Type    ArenaType `json:"type"`

	// Boundaries
	MinX float64 `json:"min_x"`
	MinY float64 `json:"min_y"`
	MaxX float64 `json:"max_x"`
	MaxY float64 `json:"max_y"`

	// Players
	Players    []ArenaPlayer `json:"players"`
	MaxPlayers int           `json:"max_players"`
	MinPlayers int           `json:"min_players"`

	// Status
	Status ArenaStatus `json:"status"`

	// Game settings
	TimeLimit      time.Duration `json:"time_limit"`
	ScoreLimit     int           `json:"score_limit"`
	LivesPerPlayer int           `json:"lives_per_player"`

	// Rewards
	EntryFee  float64 `json:"entry_fee"`
	PrizePool float64 `json:"prize_pool"`

	// Stats
	MatchesPlayed int        `json:"matches_played"`
	CreatedAt     time.Time  `json:"created_at"`
	LastMatch     *time.Time `json:"last_match,omitempty"`

	// Cooldown
	CooldownEnd *time.Time `json:"cooldown_end,omitempty"`
}

// NewArena creates a new arena
func NewArena(id, name, worldID string, arenaType ArenaType) *Arena {
	return &Arena{
		ID:             id,
		Name:           name,
		WorldID:        worldID,
		Type:           arenaType,
		Players:        make([]ArenaPlayer, 0),
		MaxPlayers:     16,
		MinPlayers:     2,
		Status:         ArenaWaiting,
		TimeLimit:      10 * time.Minute,
		ScoreLimit:     20,
		LivesPerPlayer: 3,
		EntryFee:       0,
		PrizePool:      0,
		CreatedAt:      time.Now(),
	}
}

// IsInBounds checks if position is in arena
func (a *Arena) IsInBounds(x, y float64) bool {
	return x >= a.MinX && x <= a.MaxX && y >= a.MinY && y <= a.MaxY
}

// Join adds a player to arena
func (a *Arena) Join(playerID, playerName string, team int) error {
	if a.Status != ArenaWaiting {
		return fmt.Errorf("arena not accepting players")
	}

	if len(a.Players) >= a.MaxPlayers {
		return fmt.Errorf("arena is full")
	}

	// Check if already in
	for _, p := range a.Players {
		if p.PlayerID == playerID {
			return fmt.Errorf("already in arena")
		}
	}

	player := ArenaPlayer{
		PlayerID:   playerID,
		PlayerName: playerName,
		Team:       team,
		JoinedAt:   time.Now(),
		IsAlive:    true,
	}

	a.Players = append(a.Players, player)

	return nil
}

// Leave removes a player
func (a *Arena) Leave(playerID string) error {
	for i, p := range a.Players {
		if p.PlayerID == playerID {
			a.Players = append(a.Players[:i], a.Players[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("player not in arena")
}

// IsInArena checks if player is in arena
func (a *Arena) IsInArena(playerID string) bool {
	for _, p := range a.Players {
		if p.PlayerID == playerID {
			return true
		}
	}
	return false
}

// GetPlayer gets arena player data
func (a *Arena) GetPlayer(playerID string) *ArenaPlayer {
	for i := range a.Players {
		if a.Players[i].PlayerID == playerID {
			return &a.Players[i]
		}
	}
	return nil
}

// Start starts arena match
func (a *Arena) Start() error {
	if a.Status != ArenaWaiting {
		return fmt.Errorf("arena not ready")
	}

	if len(a.Players) < a.MinPlayers {
		return fmt.Errorf("not enough players")
	}

	now := time.Now()
	a.Status = ArenaInProgress
	a.LastMatch = &now
	a.MatchesPlayed++

	return nil
}

// End ends the match
func (a *Arena) End() {
	if a.Status != ArenaInProgress {
		return
	}

	a.Status = ArenaCooldown
	cooldownEnd := time.Now().Add(30 * time.Second)
	a.CooldownEnd = &cooldownEnd

	// Clear players after cooldown
	a.Players = make([]ArenaPlayer, 0)
}

// RecordKill records a kill
func (a *Arena) RecordKill(killerID, victimID string) {
	if a.Status != ArenaInProgress {
		return
	}

	killer := a.GetPlayer(killerID)
	victim := a.GetPlayer(victimID)

	if killer != nil {
		killer.Kills++
		killer.Score += 10
	}

	if victim != nil {
		victim.Deaths++
		victim.IsAlive = false

		// Check if out of lives
		if victim.Deaths >= a.LivesPerPlayer {
			// Player eliminated
		}
	}

	// Check win conditions
	a.checkWinConditions()
}

// checkWinConditions checks if match should end
func (a *Arena) checkWinConditions() {
	// Check time limit
	if a.LastMatch != nil && time.Since(*a.LastMatch) > a.TimeLimit {
		a.End()
		return
	}

	// Check score limit
	for _, p := range a.Players {
		if p.Score >= a.ScoreLimit {
			a.End()
			return
		}
	}

	// Check if only one player/team alive
	aliveCount := 0
	for _, p := range a.Players {
		if p.IsAlive {
			aliveCount++
		}
	}

	if aliveCount == 1 {
		a.End()
		return
	}

	// For team mode, check if one team remains
	if a.Type == ArenaTeam {
		teamAlive := make(map[int]bool)
		for _, p := range a.Players {
			if p.IsAlive {
				teamAlive[p.Team] = true
			}
		}
		if len(teamAlive) == 1 {
			a.End()
		}
	}
}

// GetWinner gets the winner(s)
func (a *Arena) GetWinner() *ArenaPlayer {
	if a.Status != ArenaCooldown {
		return nil
	}

	// Find highest score
	var winner *ArenaPlayer
	for i := range a.Players {
		if winner == nil || a.Players[i].Score > winner.Score {
			winner = &a.Players[i]
		}
	}

	return winner
}

// ArenaManager manages arenas
type ArenaManager struct {
	arenas  map[string]*Arena
	byWorld map[string][]string

	arenaCounter int
	walletMgr    *economy.WalletManager
}

// NewArenaManager creates new manager
func NewArenaManager(walletMgr *economy.WalletManager) *ArenaManager {
	return &ArenaManager{
		arenas:       make(map[string]*Arena),
		byWorld:      make(map[string][]string),
		arenaCounter: 0,
		walletMgr:    walletMgr,
	}
}

// CreateArena creates a new arena
func (am *ArenaManager) CreateArena(name, worldID string, arenaType ArenaType, minX, minY, maxX, maxY float64) *Arena {
	am.arenaCounter++
	arenaID := fmt.Sprintf("arena_%d_%d", am.arenaCounter, time.Now().Unix())

	arena := NewArena(arenaID, name, worldID, arenaType)
	arena.MinX = minX
	arena.MinY = minY
	arena.MaxX = maxX
	arena.MaxY = maxY

	am.arenas[arenaID] = arena
	am.byWorld[worldID] = append(am.byWorld[worldID], arenaID)

	return arena
}

// GetArena gets an arena
func (am *ArenaManager) GetArena(arenaID string) (*Arena, bool) {
	arena, exists := am.arenas[arenaID]
	return arena, exists
}

// GetArenasByWorld gets arenas in a world
func (am *ArenaManager) GetArenasByWorld(worldID string) []*Arena {
	arenaIDs := am.byWorld[worldID]
	arenas := make([]*Arena, 0, len(arenaIDs))

	for _, id := range arenaIDs {
		if arena, exists := am.arenas[id]; exists {
			arenas = append(arenas, arena)
		}
	}

	return arenas
}

// GetAvailableArenas gets arenas accepting players
func (am *ArenaManager) GetAvailableArenas() []*Arena {
	result := make([]*Arena, 0)
	for _, arena := range am.arenas {
		if arena.Status == ArenaWaiting {
			result = append(result, arena)
		}
	}
	return result
}

// JoinArena adds player to arena with entry fee
func (am *ArenaManager) JoinArena(arenaID, playerID, playerName string, team int) error {
	arena, exists := am.GetArena(arenaID)
	if !exists {
		return fmt.Errorf("arena not found")
	}

	// Take entry fee
	if arena.EntryFee > 0 {
		wallet := am.walletMgr.GetWallet(playerID)
		if wallet == nil || !wallet.CanAfford(arena.EntryFee) {
			return fmt.Errorf("cannot afford entry fee")
		}

		_, success := wallet.Remove(arena.EntryFee, economy.TransactionSpend, "ARENA", "Arena entry fee")
		if !success {
			return fmt.Errorf("payment failed")
		}

		arena.PrizePool += arena.EntryFee
	}

	return arena.Join(playerID, playerName, team)
}

// StartArena starts a match
func (am *ArenaManager) StartArena(arenaID string) error {
	arena, exists := am.GetArena(arenaID)
	if !exists {
		return fmt.Errorf("arena not found")
	}

	return arena.Start()
}

// EndArena ends a match and pays winner
func (am *ArenaManager) EndArena(arenaID string) error {
	arena, exists := am.GetArena(arenaID)
	if !exists {
		return fmt.Errorf("arena not found")
	}

	arena.End()

	// Pay winner
	winner := arena.GetWinner()
	if winner != nil && arena.PrizePool > 0 {
		winnerWallet := am.walletMgr.GetOrCreateWallet(winner.PlayerID)
		winnerWallet.Add(arena.PrizePool, economy.TransactionEarn, "ARENA", "Arena winnings")
		arena.PrizePool = 0
	}

	return nil
}

// RecordKill records arena kill
func (am *ArenaManager) RecordKill(arenaID, killerID, victimID string) error {
	arena, exists := am.GetArena(arenaID)
	if !exists {
		return fmt.Errorf("arena not found")
	}

	arena.RecordKill(killerID, victimID)

	// Auto-end if conditions met
	if arena.Status == ArenaCooldown {
		am.EndArena(arenaID)
	}

	return nil
}

// Update updates all arenas
func (am *ArenaManager) Update() {
	for _, arena := range am.arenas {
		// Check time limit
		if arena.Status == ArenaInProgress && arena.LastMatch != nil {
			if time.Since(*arena.LastMatch) > arena.TimeLimit {
				am.EndArena(arena.ID)
			}
		}

		// Reset cooldown
		if arena.Status == ArenaCooldown && arena.CooldownEnd != nil {
			if time.Now().After(*arena.CooldownEnd) {
				arena.Status = ArenaWaiting
				arena.CooldownEnd = nil
				arena.Players = make([]ArenaPlayer, 0)
			}
		}
	}
}

// GetActiveMatches returns active arena count
func (am *ArenaManager) GetActiveMatches() int {
	count := 0
	for _, arena := range am.arenas {
		if arena.Status == ArenaInProgress {
			count++
		}
	}
	return count
}
