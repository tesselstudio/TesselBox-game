package boss

import (
	"fmt"
	"time"

	"tesselbox/pkg/economy"
	"tesselbox/pkg/items"
)

// BossType represents boss difficulty
type BossType int

const (
	BossNormal BossType = iota
	BossElite
	BossLegendary
)

// BossPhase represents current boss phase
type BossPhase int

const (
	PhaseNormal BossPhase = iota
	PhaseEnraged
	PhaseDesperate
)

// BossInstance represents an active boss fight
type BossInstance struct {
	ID      string `json:"id"`
	BossID  string `json:"boss_id"`
	WorldID string `json:"world_id"`

	// Location
	X float64 `json:"x"`
	Y float64 `json:"y"`

	// State
	Health    float64   `json:"health"`
	MaxHealth float64   `json:"max_health"`
	Phase     BossPhase `json:"phase"`

	// Combat
	Players []string           `json:"players"` // Participating players
	Damage  map[string]float64 `json:"damage"`  // Player -> damage dealt

	// Stats
	SpawnedAt time.Time  `json:"spawned_at"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
	WinnerID  *string    `json:"winner_id,omitempty"` // Player who got last hit

	// Status
	Active   bool `json:"active"`
	Defeated bool `json:"defeated"`
}

// NewBossInstance creates a boss fight
func NewBossInstance(id, bossID, worldID string, bossType BossType, x, y float64) *BossInstance {
	health := 1000.0
	switch bossType {
	case BossElite:
		health = 5000.0
	case BossLegendary:
		health = 20000.0
	}

	return &BossInstance{
		ID:        id,
		BossID:    bossID,
		WorldID:   worldID,
		X:         x,
		Y:         y,
		Health:    health,
		MaxHealth: health,
		Phase:     PhaseNormal,
		Players:   make([]string, 0),
		Damage:    make(map[string]float64),
		SpawnedAt: time.Now(),
		Active:    true,
		Defeated:  false,
	}
}

// Join adds a player to the fight
func (bi *BossInstance) Join(playerID string) {
	for _, id := range bi.Players {
		if id == playerID {
			return // Already in
		}
	}
	bi.Players = append(bi.Players, playerID)
	bi.Damage[playerID] = 0
}

// RecordDamage records damage dealt
func (bi *BossInstance) RecordDamage(playerID string, amount float64) {
	bi.Health -= amount
	if bi.Health < 0 {
		bi.Health = 0
	}

	// Track damage
	if _, exists := bi.Damage[playerID]; !exists {
		bi.Join(playerID)
	}
	bi.Damage[playerID] += amount

	// Update phase
	bi.UpdatePhase()

	// Check defeat
	if bi.Health <= 0 {
		bi.Defeat(playerID)
	}
}

// UpdatePhase updates boss phase based on health
func (bi *BossInstance) UpdatePhase() {
	healthPercent := bi.Health / bi.MaxHealth

	switch {
	case healthPercent <= 0.2:
		bi.Phase = PhaseDesperate
	case healthPercent <= 0.5:
		bi.Phase = PhaseEnraged
	default:
		bi.Phase = PhaseNormal
	}
}

// Defeat defeats the boss
func (bi *BossInstance) Defeat(lastHitPlayerID string) {
	if bi.Defeated {
		return
	}

	bi.Defeated = true
	bi.Active = false
	bi.WinnerID = &lastHitPlayerID

	now := time.Now()
	bi.EndedAt = &now
}

// GetTopDamage returns top damage dealers
func (bi *BossInstance) GetTopDamage(count int) []struct {
	PlayerID string
	Damage   float64
} {
	// Convert to slice
	type damageEntry struct {
		PlayerID string
		Damage   float64
	}

	entries := make([]damageEntry, 0, len(bi.Damage))
	for id, dmg := range bi.Damage {
		entries = append(entries, damageEntry{id, dmg})
	}

	// Sort by damage
	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[i].Damage < entries[j].Damage {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}

	if count > len(entries) {
		count = len(entries)
	}

	result := make([]struct {
		PlayerID string
		Damage   float64
	}, count)

	for i := 0; i < count; i++ {
		result[i] = struct {
			PlayerID string
			Damage   float64
		}{entries[i].PlayerID, entries[i].Damage}
	}

	return result
}

// BossDefinition defines a boss type
type BossDefinition struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        BossType `json:"type"`

	// Stats
	BaseHealth float64 `json:"base_health"`
	Damage     float64 `json:"damage"`

	// Rewards
	Rewards BossRewards `json:"rewards"`

	// Spawn
	SpawnInterval  time.Duration `json:"spawn_interval"`
	SpawnLocations []struct {
		WorldID string  `json:"world_id"`
		X       float64 `json:"x"`
		Y       float64 `json:"y"`
	} `json:"spawn_locations"`
}

// BossRewards represents boss loot
type BossRewards struct {
	Guaranteed []items.Item `json:"guaranteed"` // Everyone gets
	TopDamage  []items.Item `json:"top_damage"` // Top 3 damage
	LastHit    []items.Item `json:"last_hit"`   // Who got last hit
	Money      float64      `json:"money"`
}

// BossEventManager manages boss events
type BossEventManager struct {
	definitions map[string]BossDefinition
	instances   map[string]*BossInstance

	instanceCounter int
	walletMgr       *economy.WalletManager
}

// NewBossEventManager creates new manager
func NewBossEventManager(walletMgr *economy.WalletManager) *BossEventManager {
	return &BossEventManager{
		definitions: make(map[string]BossDefinition),
		instances:   make(map[string]*BossInstance),
		walletMgr:   walletMgr,
	}
}

// RegisterBoss registers a boss
func (bm *BossEventManager) RegisterBoss(def BossDefinition) {
	bm.definitions[def.ID] = def
}

// SpawnBoss spawns a boss
func (bm *BossEventManager) SpawnBoss(bossID, worldID string, x, y float64) (*BossInstance, error) {
	def, exists := bm.definitions[bossID]
	if !exists {
		return nil, fmt.Errorf("boss not found")
	}

	bm.instanceCounter++
	instanceID := fmt.Sprintf("boss_%d_%d", bm.instanceCounter, time.Now().Unix())

	instance := NewBossInstance(instanceID, bossID, worldID, def.Type, x, y)
	instance.MaxHealth = def.BaseHealth
	instance.Health = def.BaseHealth

	bm.instances[instanceID] = instance

	return instance, nil
}

// GetInstance gets an instance
func (bm *BossEventManager) GetInstance(instanceID string) (*BossInstance, bool) {
	instance, exists := bm.instances[instanceID]
	return instance, exists
}

// GetActiveBosses returns active boss count
func (bm *BossEventManager) GetActiveBosses() []*BossInstance {
	result := make([]*BossInstance, 0)
	for _, instance := range bm.instances {
		if instance.Active {
			result = append(result, instance)
		}
	}
	return result
}

// RecordDamage records damage to boss
func (bm *BossEventManager) RecordDamage(instanceID, playerID string, amount float64) error {
	instance, exists := bm.GetInstance(instanceID)
	if !exists {
		return fmt.Errorf("boss not found")
	}

	instance.RecordDamage(playerID, amount)

	// If defeated, distribute rewards
	if instance.Defeated {
		bm.distributeRewards(instance)
	}

	return nil
}

// distributeRewards gives out boss rewards
func (bm *BossEventManager) distributeRewards(instance *BossInstance) {
	def, exists := bm.definitions[instance.BossID]
	if !exists {
		return
	}

	// Give guaranteed rewards to all participants
	for _, playerID := range instance.Players {
		wallet := bm.walletMgr.GetOrCreateWallet(playerID)
		wallet.Add(def.Rewards.Money/float64(len(instance.Players)), economy.TransactionEarn, "BOSS", "Boss reward")
	}

	// Give top damage rewards
	topDamage := instance.GetTopDamage(3)
	for i, entry := range topDamage {
		if i < len(def.Rewards.TopDamage) {
			// Give item reward (in real implementation)
			_ = entry
		}
	}

	// Give last hit reward
	if instance.WinnerID != nil && len(def.Rewards.LastHit) > 0 {
		// Give special reward to last hitter
	}
}

// Cleanup removes defeated bosses
func (bm *BossEventManager) Cleanup() {
	for id, instance := range bm.instances {
		if !instance.Active && instance.EndedAt != nil {
			if time.Since(*instance.EndedAt) > 5*time.Minute {
				delete(bm.instances, id)
			}
		}
	}
}

// GetUpcomingSpawn returns when next boss spawns
func (bm *BossEventManager) GetUpcomingSpawn() (bossID string, timeUntil time.Duration) {
	var soonest time.Time
	var nextBoss string

	for id, def := range bm.definitions {
		// In real implementation, track last spawn times
		nextSpawn := time.Now().Add(def.SpawnInterval)
		if soonest.IsZero() || nextSpawn.Before(soonest) {
			soonest = nextSpawn
			nextBoss = id
		}
	}

	if soonest.IsZero() {
		return "", 0
	}

	return nextBoss, time.Until(soonest)
}
