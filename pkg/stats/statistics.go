package stats

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// PlayerStats tracks comprehensive player statistics
type PlayerStats struct {
	PlayerID    string `json:"player_id"`
	LastUpdated time.Time `json:"last_updated"`
	
	// Play time
	TotalPlayTime    time.Duration `json:"total_play_time"`
	SessionCount     int           `json:"session_count"`
	LongestSession   time.Duration `json:"longest_session"`
	
	// Building
	BlocksPlaced     int64         `json:"blocks_placed"`
	BlocksBroken     int64         `json:"blocks_broken"`
	BlocksByType     map[string]int64 `json:"blocks_by_type,omitempty"`
	
	// Mining
	OresMined        int64         `json:"ores_mined"`
	StoneMined       int64         `json:"stone_mined"`
	DirtMined        int64         `json:"dirt_mined"`
	
	// Combat
	MobsKilled       int64         `json:"mobs_killed"`
	PlayersKilled    int64         `json:"players_killed"`
	Deaths           int64         `json:"deaths"`
	DamageDealt      float64       `json:"damage_dealt"`
	DamageTaken      float64       `json:"damage_taken"`
	
	// Economy
	MoneyEarned      float64       `json:"money_earned"`
	MoneySpent       float64       `json:"money_spent"`
	ItemsTraded      int64         `json:"items_traded"`
	ShopsCreated     int           `json:"shops_created"`
	ShopSales        float64       `json:"shop_sales"`
	
	// Exploration
	ChunksExplored   int           `json:"chunks_explored"`
	WorldsVisited    int           `json:"worlds_visited"`
	DistanceWalked   float64       `json:"distance_walked"`
	DistanceSprinted float64       `json:"distance_sprinted"`
	DistanceFlown    float64       `json:"distance_flown"`
	
	// Social
	FriendsMade      int           `json:"friends_made"`
	MessagesSent     int64         `json:"messages_sent"`
	PartiesJoined    int           `json:"parties_joined"`
	GuildsJoined     int           `json:"guilds_joined"`
	
	// Land
	ChunksClaimed    int           `json:"chunks_claimed"`
	ClaimsExpanded    int           `json:"claims_expanded"`
	
	// PvP
	DuelsWon         int           `json:"duels_won"`
	DuelsLost        int           `json:"duels_lost"`
	BountiesPlaced   int           `json:"bounties_placed"`
	BountiesClaimed  int           `json:"bounties_claimed"`
	
	// Minigames
	MinigamesPlayed  int           `json:"minigames_played"`
	MinigamesWon     int           `json:"minigames_won"`
	
	// Jobs
	JobLevelsGained  int           `json:"job_levels_gained"`
	JobXPCollected   int64         `json:"job_xp_collected"`
	
	// Inventory
	ItemsCrafted     int64         `json:"items_crafted"`
	ItemsPickedUp    int64         `json:"items_picked_up"`
	ItemsDropped     int64         `json:"items_dropped"`
	
	// Commands
	CommandsUsed     int64         `json:"commands_used"`
	
	// Randomland
	TimesToRandomland int          `json:"times_to_randomland"`
	RandomlandKills   int          `json:"randomland_kills"`
}

// NewPlayerStats creates new player stats
func NewPlayerStats(playerID string) *PlayerStats {
	return &PlayerStats{
		PlayerID:         playerID,
		LastUpdated:      time.Now(),
		BlocksByType:     make(map[string]int64),
	}
}

// RecordBlockPlaced records a block placement
func (ps *PlayerStats) RecordBlockPlaced(blockType string) {
	ps.BlocksPlaced++
	ps.BlocksByType[blockType]++
	ps.LastUpdated = time.Now()
}

// RecordBlockBroken records a block break
func (ps *PlayerStats) RecordBlockBroken(blockType string) {
	ps.BlocksBroken++
	ps.BlocksByType[blockType]++
	ps.LastUpdated = time.Now()
}

// RecordMobKill records a mob kill
func (ps *PlayerStats) RecordMobKill() {
	ps.MobsKilled++
	ps.LastUpdated = time.Now()
}

// RecordPlayerKill records a player kill
func (ps *PlayerStats) RecordPlayerKill() {
	ps.PlayersKilled++
	ps.LastUpdated = time.Now()
}

// RecordDeath records a death
func (ps *PlayerStats) RecordDeath() {
	ps.Deaths++
	ps.LastUpdated = time.Now()
}

// RecordDamageDealt records damage dealt
func (ps *PlayerStats) RecordDamageDealt(amount float64) {
	ps.DamageDealt += amount
	ps.LastUpdated = time.Now()
}

// RecordDamageTaken records damage taken
func (ps *PlayerStats) RecordDamageTaken(amount float64) {
	ps.DamageTaken += amount
	ps.LastUpdated = time.Now()
}

// RecordMoneyEarned records money earned
func (ps *PlayerStats) RecordMoneyEarned(amount float64) {
	ps.MoneyEarned += amount
	ps.LastUpdated = time.Now()
}

// RecordMoneySpent records money spent
func (ps *PlayerStats) RecordMoneySpent(amount float64) {
	ps.MoneySpent += amount
	ps.LastUpdated = time.Now()
}

// RecordDistance records distance traveled
func (ps *PlayerStats) RecordDistance(distance float64, mode string) {
	switch mode {
	case "walk":
		ps.DistanceWalked += distance
	case "sprint":
		ps.DistanceSprinted += distance
	case "fly":
		ps.DistanceFlown += distance
	}
	ps.LastUpdated = time.Now()
}

// RecordMessageSent records a chat message
func (ps *PlayerStats) RecordMessageSent() {
	ps.MessagesSent++
	ps.LastUpdated = time.Now()
}

// RecordCommandUsed records a command
func (ps *PlayerStats) RecordCommandUsed() {
	ps.CommandsUsed++
	ps.LastUpdated = time.Now()
}

// RecordSessionStart starts a new session
func (ps *PlayerStats) RecordSessionStart() {
	ps.SessionCount++
	ps.LastUpdated = time.Now()
}

// RecordSessionEnd ends a session
func (ps *PlayerStats) RecordSessionEnd(duration time.Duration) {
	ps.TotalPlayTime += duration
	if duration > ps.LongestSession {
		ps.LongestSession = duration
	}
	ps.LastUpdated = time.Now()
}

// GetKDRatio returns kill/death ratio
func (ps *PlayerStats) GetKDRatio() float64 {
	if ps.Deaths == 0 {
		if ps.PlayersKilled == 0 {
			return 0
		}
		return float64(ps.PlayersKilled)
	}
	return float64(ps.PlayersKilled) / float64(ps.Deaths)
}

// GetWinRate returns win rate for minigames/duels
func (ps *PlayerStats) GetWinRate() float64 {
	total := ps.DuelsWon + ps.DuelsLost + ps.MinigamesWon + (ps.MinigamesPlayed - ps.MinigamesWon)
	if total == 0 {
		return 0
	}
	wins := ps.DuelsWon + ps.MinigamesWon
	return float64(wins) / float64(total) * 100
}

// GetNetWorth returns net worth (earned - spent)
func (ps *PlayerStats) GetNetWorth() float64 {
	return ps.MoneyEarned - ps.MoneySpent
}

// StatsManager manages player statistics
type StatsManager struct {
	stats       map[string]*PlayerStats
	
	storagePath string
}

// NewStatsManager creates a new stats manager
func NewStatsManager(storageDir string) *StatsManager {
	return &StatsManager{
		stats:       make(map[string]*PlayerStats),
		storagePath: filepath.Join(storageDir, "statistics.json"),
	}
}

// GetOrCreateStats gets or creates player stats
func (sm *StatsManager) GetOrCreateStats(playerID string) *PlayerStats {
	if stats, exists := sm.stats[playerID]; exists {
		return stats
	}
	
	stats := NewPlayerStats(playerID)
	sm.stats[playerID] = stats
	return stats
}

// GetStats gets player stats
func (sm *StatsManager) GetStats(playerID string) (*PlayerStats, bool) {
	stats, exists := sm.stats[playerID]
	return stats, exists
}

// HasStats checks if player has stats
func (sm *StatsManager) HasStats(playerID string) bool {
	_, exists := sm.stats[playerID]
	return exists
}

// GetTopByPlayTime returns top players by play time
func (sm *StatsManager) GetTopByPlayTime(count int) []*PlayerStats {
	allStats := make([]*PlayerStats, 0, len(sm.stats))
	for _, stats := range sm.stats {
		allStats = append(allStats, stats)
	}
	
	// Sort by play time
	for i := 0; i < len(allStats); i++ {
		for j := i + 1; j < len(allStats); j++ {
			if allStats[i].TotalPlayTime < allStats[j].TotalPlayTime {
				allStats[i], allStats[j] = allStats[j], allStats[i]
			}
		}
	}
	
	if count > len(allStats) {
		count = len(allStats)
	}
	
	return allStats[:count]
}

// GetTopByWealth returns top players by money earned
func (sm *StatsManager) GetTopByWealth(count int) []*PlayerStats {
	allStats := make([]*PlayerStats, 0, len(sm.stats))
	for _, stats := range sm.stats {
		allStats = append(allStats, stats)
	}
	
	// Sort by money earned
	for i := 0; i < len(allStats); i++ {
		for j := i + 1; j < len(allStats); j++ {
			if allStats[i].MoneyEarned < allStats[j].MoneyEarned {
				allStats[i], allStats[j] = allStats[j], allStats[i]
			}
		}
	}
	
	if count > len(allStats) {
		count = len(allStats)
	}
	
	return allStats[:count]
}

// GetTopByKills returns top players by PvP kills
func (sm *StatsManager) GetTopByKills(count int) []*PlayerStats {
	allStats := make([]*PlayerStats, 0, len(sm.stats))
	for _, stats := range sm.stats {
		allStats = append(allStats, stats)
	}
	
	// Sort by kills
	for i := 0; i < len(allStats); i++ {
		for j := i + 1; j < len(allStats); j++ {
			if allStats[i].PlayersKilled < allStats[j].PlayersKilled {
				allStats[i], allStats[j] = allStats[j], allStats[i]
			}
		}
	}
	
	if count > len(allStats) {
		count = len(allStats)
	}
	
	return allStats[:count]
}

// GetTopByBlocksPlaced returns top builders
func (sm *StatsManager) GetTopByBlocksPlaced(count int) []*PlayerStats {
	allStats := make([]*PlayerStats, 0, len(sm.stats))
	for _, stats := range sm.stats {
		allStats = append(allStats, stats)
	}
	
	// Sort by blocks placed
	for i := 0; i < len(allStats); i++ {
		for j := i + 1; j < len(allStats); j++ {
			if allStats[i].BlocksPlaced < allStats[j].BlocksPlaced {
				allStats[i], allStats[j] = allStats[j], allStats[i]
			}
		}
	}
	
	if count > len(allStats) {
		count = len(allStats)
	}
	
	return allStats[:count]
}

// GetTopByDistance returns top explorers
func (sm *StatsManager) GetTopByDistance(count int) []*PlayerStats {
	allStats := make([]*PlayerStats, 0, len(sm.stats))
	for _, stats := range sm.stats {
		allStats = append(allStats, stats)
	}
	
	// Sort by total distance
	for i := 0; i < len(allStats); i++ {
		for j := i + 1; j < len(allStats); j++ {
			distI := allStats[i].DistanceWalked + allStats[i].DistanceSprinted + allStats[i].DistanceFlown
			distJ := allStats[j].DistanceWalked + allStats[j].DistanceSprinted + allStats[j].DistanceFlown
			if distI < distJ {
				allStats[i], allStats[j] = allStats[j], allStats[i]
			}
		}
	}
	
	if count > len(allStats) {
		count = len(allStats)
	}
	
	return allStats[:count]
}

// GetGlobalStats returns global statistics
func (sm *StatsManager) GetGlobalStats() struct {
	TotalPlayers     int
	TotalPlayTime    time.Duration
	TotalBlocksPlaced int64
	TotalBlocksBroken int64
	TotalMobsKilled  int64
	TotalMoneyEarned float64
	TotalDistance    float64
} {
	var result struct {
		TotalPlayers     int
		TotalPlayTime    time.Duration
		TotalBlocksPlaced int64
		TotalBlocksBroken int64
		TotalMobsKilled  int64
		TotalMoneyEarned float64
		TotalDistance    float64
	}
	
	for _, stats := range sm.stats {
		result.TotalPlayers++
		result.TotalPlayTime += stats.TotalPlayTime
		result.TotalBlocksPlaced += stats.BlocksPlaced
		result.TotalBlocksBroken += stats.BlocksBroken
		result.TotalMobsKilled += stats.MobsKilled
		result.TotalMoneyEarned += stats.MoneyEarned
		result.TotalDistance += stats.DistanceWalked + stats.DistanceSprinted + stats.DistanceFlown
	}
	
	return result
}

// Save saves all statistics
func (sm *StatsManager) Save() error {
	data, err := json.MarshalIndent(sm.stats, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}
	
	if err := os.WriteFile(sm.storagePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	
	return nil
}

// Load loads all statistics
func (sm *StatsManager) Load() error {
	data, err := os.ReadFile(sm.storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read: %w", err)
	}
	
	var loaded map[string]*PlayerStats
	if err := json.Unmarshal(data, &loaded); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}
	
	sm.stats = loaded
	if sm.stats == nil {
		sm.stats = make(map[string]*PlayerStats)
	}
	
	// Ensure all stats have initialized maps
	for _, stats := range sm.stats {
		if stats.BlocksByType == nil {
			stats.BlocksByType = make(map[string]int64)
		}
	}
	
	return nil
}

// ExportPlayerStats exports a player's stats as JSON
func (sm *StatsManager) ExportPlayerStats(playerID string) ([]byte, error) {
	stats, exists := sm.GetStats(playerID)
	if !exists {
		return nil, fmt.Errorf("player stats not found")
	}
	
	return json.MarshalIndent(stats, "", "  ")
}
