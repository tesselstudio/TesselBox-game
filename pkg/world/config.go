package world

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// DifficultyLevel represents world difficulty
type DifficultyLevel int

const (
	DifficultyPeaceful DifficultyLevel = iota // No hostile mobs, no hunger
	DifficultyEasy                            // Reduced damage, fewer mobs
	DifficultyNormal                          // Standard gameplay
	DifficultyHard                            // Increased damage, more mobs
)

// String returns the string representation
func (d DifficultyLevel) String() string {
	switch d {
	case DifficultyPeaceful:
		return "peaceful"
	case DifficultyEasy:
		return "easy"
	case DifficultyNormal:
		return "normal"
	case DifficultyHard:
		return "hard"
	}
	return "unknown"
}

// ParseDifficulty parses a difficulty string
func ParseDifficulty(s string) DifficultyLevel {
	switch s {
	case "peaceful":
		return DifficultyPeaceful
	case "easy":
		return DifficultyEasy
	case "normal":
		return DifficultyNormal
	case "hard":
		return DifficultyHard
	}
	return DifficultyNormal
}

// DifficultyModifiers contains multipliers for difficulty settings
type DifficultyModifiers struct {
	MobHealthMultiplier     float64
	MobDamageMultiplier     float64
	MobSpeedMultiplier      float64
	MobSpawnMultiplier      float64
	MobAggroRangeMultiplier float64
	HungerDrainMultiplier   float64
	DamageTakenMultiplier   float64 // For player
}

// GetModifiers returns the modifiers for this difficulty
func (d DifficultyLevel) GetModifiers() DifficultyModifiers {
	switch d {
	case DifficultyPeaceful:
		return DifficultyModifiers{
			MobHealthMultiplier:     0.5,
			MobDamageMultiplier:     0.0, // No mob damage
			MobSpeedMultiplier:      0.8,
			MobSpawnMultiplier:      0.0, // No hostile spawns
			MobAggroRangeMultiplier: 0.0,
			HungerDrainMultiplier:   0.0, // No hunger
			DamageTakenMultiplier:   0.5,
		}
	case DifficultyEasy:
		return DifficultyModifiers{
			MobHealthMultiplier:     0.7,
			MobDamageMultiplier:     0.5,
			MobSpeedMultiplier:      0.9,
			MobSpawnMultiplier:      0.7,
			MobAggroRangeMultiplier: 0.8,
			HungerDrainMultiplier:   0.7,
			DamageTakenMultiplier:   0.7,
		}
	case DifficultyNormal:
		return DifficultyModifiers{
			MobHealthMultiplier:     1.0,
			MobDamageMultiplier:     1.0,
			MobSpeedMultiplier:      1.0,
			MobSpawnMultiplier:      1.0,
			MobAggroRangeMultiplier: 1.0,
			HungerDrainMultiplier:   1.0,
			DamageTakenMultiplier:   1.0,
		}
	case DifficultyHard:
		return DifficultyModifiers{
			MobHealthMultiplier:     1.5,
			MobDamageMultiplier:     1.5,
			MobSpeedMultiplier:      1.2,
			MobSpawnMultiplier:      1.5,
			MobAggroRangeMultiplier: 1.3,
			HungerDrainMultiplier:   1.2,
			DamageTakenMultiplier:   1.3,
		}
	}
	return DifficultyNormal.GetModifiers()
}

// GameRules contains configurable game rules
type GameRules struct {
	// Command settings
	CommandFeedback     bool `json:"command_feedback"`
	LogAdminCommands    bool `json:"log_admin_commands"`
	SendCommandFeedback bool `json:"send_command_feedback"`

	// World behavior
	PvPEnabled          bool `json:"pvp_enabled"`
	FireSpread          bool `json:"fire_spread"`
	MobGriefing         bool `json:"mob_griefing"`   // Can mobs break blocks
	KeepInventory       bool `json:"keep_inventory"` // On death
	DaylightCycle       bool `json:"daylight_cycle"`
	WeatherCycle        bool `json:"weather_cycle"`
	DoMobSpawning       bool `json:"do_mob_spawning"`
	DoEntityDrops       bool `json:"do_entity_drops"`
	DoTileDrops         bool `json:"do_tile_drops"` // Block drops
	NaturalRegeneration bool `json:"natural_regeneration"`

	// Game mode
	AllowCreative bool   `json:"allow_creative"`
	AllowSurvival bool   `json:"allow_survival"`
	ForceGameMode string `json:"force_game_mode,omitempty"` // "", "creative", "survival"

	// Restrictions
	BlockedCommands []string `json:"blocked_commands,omitempty"`
	BlockedPlugins  []string `json:"blocked_plugins,omitempty"`

	// Access
	WhitelistEnabled bool     `json:"whitelist_enabled"`
	Whitelist        []string `json:"whitelist,omitempty"`

	// Spawn
	SpawnRadius     int `json:"spawn_radius"`     // Spawn area radius
	SpawnProtection int `json:"spawn_protection"` // Unbreakable radius

	// Performance
	RandomTickSpeed   int `json:"random_tick_speed"`   // Default 3
	MaxEntityCramming int `json:"max_entity_cramming"` // Default 24
}

// DefaultGameRules returns default game rules
func DefaultGameRules() GameRules {
	return GameRules{
		CommandFeedback:     true,
		LogAdminCommands:    true,
		SendCommandFeedback: true,
		PvPEnabled:          true,
		FireSpread:          true,
		MobGriefing:         true,
		KeepInventory:       false,
		DaylightCycle:       true,
		WeatherCycle:        true,
		DoMobSpawning:       true,
		DoEntityDrops:       true,
		DoTileDrops:         true,
		NaturalRegeneration: true,
		AllowCreative:       true,
		AllowSurvival:       true,
		ForceGameMode:       "",
		BlockedCommands:     []string{},
		BlockedPlugins:      []string{},
		WhitelistEnabled:    false,
		Whitelist:           []string{},
		SpawnRadius:         10,
		SpawnProtection:     5, // 5 chunks radius = 80 blocks
		RandomTickSpeed:     3,
		MaxEntityCramming:   24,
	}
}

// EconomySettings contains per-world economy configuration
type EconomySettings struct {
	Enabled          bool    `json:"enabled"`
	StartingBalance  float64 `json:"starting_balance"`
	MaxTransaction   float64 `json:"max_transaction"`
	MinTransaction   float64 `json:"min_transaction"`
	TradeTaxRate     float64 `json:"trade_tax_rate"` // 0.05 = 5%
	ShopTaxRate      float64 `json:"shop_tax_rate"`
	AuctionTaxRate   float64 `json:"auction_tax_rate"`
	DeathPenaltyRate float64 `json:"death_penalty_rate"` // 0.10 = 10%
	EnableLoans      bool    `json:"enable_loans"`
	MaxLoanAmount    float64 `json:"max_loan_amount"`
	LoanInterestRate float64 `json:"loan_interest_rate"` // Daily rate
}

// DefaultEconomySettings returns default economy settings
func DefaultEconomySettings() EconomySettings {
	return EconomySettings{
		Enabled:          true,
		StartingBalance:  100.0,
		MaxTransaction:   1000000.0,
		MinTransaction:   0.01,
		TradeTaxRate:     0.0,  // No tax on player trades
		ShopTaxRate:      0.05, // 5% shop tax
		AuctionTaxRate:   0.10, // 10% auction fee
		DeathPenaltyRate: 0.10, // Lose 10% on death
		EnableLoans:      true,
		MaxLoanAmount:    10000.0,
		LoanInterestRate: 0.05, // 5% daily
	}
}

// WorldBorder defines the world border settings
type WorldBorder struct {
	Enabled         bool    `json:"enabled"`
	CenterX         float64 `json:"center_x"`
	CenterY         float64 `json:"center_y"`
	Radius          float64 `json:"radius"`
	DamagePerSec    float64 `json:"damage_per_sec"`
	WarningDistance float64 `json:"warning_distance"`
}

// DefaultWorldBorder returns default world border (disabled)
func DefaultWorldBorder() WorldBorder {
	return WorldBorder{
		Enabled:         false,
		CenterX:         0,
		CenterY:         0,
		Radius:          10000,
		DamagePerSec:    1.0,
		WarningDistance: 50.0,
	}
}

// WorldConfig contains all per-world settings
type WorldConfig struct {
	// Identity
	WorldID   string    `json:"world_id"`
	WorldName string    `json:"world_name"`
	OwnerID   string    `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`

	// Type
	WorldType string `json:"world_type"` // "survival", "creative", "adventure", "minigame"

	// Difficulty
	Difficulty DifficultyLevel `json:"difficulty"`

	// Game Rules
	GameRules GameRules `json:"game_rules"`

	// Economy
	Economy EconomySettings `json:"economy"`

	// Border
	WorldBorder WorldBorder `json:"world_border"`

	// Spawn
	SpawnX float64 `json:"spawn_x"`
	SpawnY float64 `json:"spawn_y"`

	// Inventory isolation
	SeparateInventories bool `json:"separate_inventories"`
	ConfigVersion       int  `json:"config_version"`

	// Metadata
	LastPlayed   time.Time `json:"last_played"`
	PlayCount    int       `json:"play_count"`
	TotalPlayers int       `json:"total_players"`

	// Plugin settings
	EnabledPlugins []string                          `json:"enabled_plugins,omitempty"`
	PluginConfigs  map[string]map[string]interface{} `json:"plugin_configs,omitempty"`
}

// NewWorldConfig creates a new world configuration with defaults
func NewWorldConfig(worldID, worldName, ownerID string) *WorldConfig {
	now := time.Now()
	return &WorldConfig{
		WorldID:             worldID,
		WorldName:           worldName,
		OwnerID:             ownerID,
		CreatedAt:           now,
		CreatedBy:           ownerID,
		WorldType:           "survival",
		Difficulty:          DifficultyNormal,
		GameRules:           DefaultGameRules(),
		Economy:             DefaultEconomySettings(),
		WorldBorder:         DefaultWorldBorder(),
		SpawnX:              0,
		SpawnY:              0,
		SeparateInventories: false,
		LastPlayed:          now,
		PlayCount:           0,
		TotalPlayers:        0,
		EnabledPlugins:      []string{},
		PluginConfigs:       make(map[string]map[string]interface{}),
	}
}

// SetDifficulty changes the difficulty
func (wc *WorldConfig) SetDifficulty(d DifficultyLevel) {
	wc.Difficulty = d
}

// GetDifficultyModifiers returns current difficulty modifiers
func (wc *WorldConfig) GetDifficultyModifiers() DifficultyModifiers {
	return wc.Difficulty.GetModifiers()
}

// IsCommandBlocked checks if a command is blocked
func (wc *WorldConfig) IsCommandBlocked(cmd string) bool {
	for _, blocked := range wc.GameRules.BlockedCommands {
		if blocked == cmd {
			return true
		}
	}
	return false
}

// BlockCommand blocks a command
func (wc *WorldConfig) BlockCommand(cmd string) {
	if !wc.IsCommandBlocked(cmd) {
		wc.GameRules.BlockedCommands = append(wc.GameRules.BlockedCommands, cmd)
	}
}

// UnblockCommand unblocks a command
func (wc *WorldConfig) UnblockCommand(cmd string) {
	for i, blocked := range wc.GameRules.BlockedCommands {
		if blocked == cmd {
			wc.GameRules.BlockedCommands = append(
				wc.GameRules.BlockedCommands[:i],
				wc.GameRules.BlockedCommands[i+1:]...,
			)
			return
		}
	}
}

// IsPluginEnabled checks if a plugin is enabled
func (wc *WorldConfig) IsPluginEnabled(pluginID string) bool {
	for _, enabled := range wc.EnabledPlugins {
		if enabled == pluginID {
			return true
		}
	}
	return false
}

// EnablePlugin enables a plugin for this world
func (wc *WorldConfig) EnablePlugin(pluginID string) {
	if !wc.IsPluginEnabled(pluginID) {
		wc.EnabledPlugins = append(wc.EnabledPlugins, pluginID)
	}
}

// DisablePlugin disables a plugin for this world
func (wc *WorldConfig) DisablePlugin(pluginID string) {
	for i, enabled := range wc.EnabledPlugins {
		if enabled == pluginID {
			wc.EnabledPlugins = append(
				wc.EnabledPlugins[:i],
				wc.EnabledPlugins[i+1:]...,
			)
			return
		}
	}
}

// SetPluginConfig sets a plugin configuration value
func (wc *WorldConfig) SetPluginConfig(pluginID, key string, value interface{}) {
	if wc.PluginConfigs[pluginID] == nil {
		wc.PluginConfigs[pluginID] = make(map[string]interface{})
	}
	wc.PluginConfigs[pluginID][key] = value
}

// GetPluginConfig gets a plugin configuration value
func (wc *WorldConfig) GetPluginConfig(pluginID, key string) (interface{}, bool) {
	if config, exists := wc.PluginConfigs[pluginID]; exists {
		val, ok := config[key]
		return val, ok
	}
	return nil, false
}

// IsWhitelisted checks if a player is whitelisted
func (wc *WorldConfig) IsWhitelisted(playerID string) bool {
	if !wc.GameRules.WhitelistEnabled {
		return true // No whitelist = everyone allowed
	}
	for _, id := range wc.GameRules.Whitelist {
		if id == playerID {
			return true
		}
	}
	return false
}

// AddToWhitelist adds a player to whitelist
func (wc *WorldConfig) AddToWhitelist(playerID string) {
	if !wc.IsWhitelisted(playerID) {
		wc.GameRules.Whitelist = append(wc.GameRules.Whitelist, playerID)
	}
}

// RemoveFromWhitelist removes a player from whitelist
func (wc *WorldConfig) RemoveFromWhitelist(playerID string) {
	for i, id := range wc.GameRules.Whitelist {
		if id == playerID {
			wc.GameRules.Whitelist = append(
				wc.GameRules.Whitelist[:i],
				wc.GameRules.Whitelist[i+1:]...,
			)
			return
		}
	}
}

// RecordPlay records a play session
func (wc *WorldConfig) RecordPlay() {
	wc.LastPlayed = time.Now()
	wc.PlayCount++
}

// IsSpawnProtected checks if a location is within spawn protection
func (wc *WorldConfig) IsSpawnProtected(x, y float64) bool {
	if wc.GameRules.SpawnProtection <= 0 {
		return false
	}

	// Calculate distance from spawn in chunks
	chunkX := int(x / GetChunkWidth())
	chunkY := int(y / GetChunkHeight())
	spawnChunkX := int(wc.SpawnX / GetChunkWidth())
	spawnChunkY := int(wc.SpawnY / GetChunkHeight())

	dx := chunkX - spawnChunkX
	dy := chunkY - spawnChunkY
	distance := dx*dx + dy*dy

	return distance <= wc.GameRules.SpawnProtection*wc.GameRules.SpawnProtection
}

// IsWithinWorldBorder checks if a location is within world border
func (wc *WorldConfig) IsWithinWorldBorder(x, y float64) bool {
	if !wc.WorldBorder.Enabled {
		return true
	}

	dx := x - wc.WorldBorder.CenterX
	dy := y - wc.WorldBorder.CenterY
	distance := dx*dx + dy*dy

	return distance <= wc.WorldBorder.Radius*wc.WorldBorder.Radius
}

// GetDistanceToWorldBorder returns distance to border (negative if outside)
func (wc *WorldConfig) GetDistanceToWorldBorder(x, y float64) float64 {
	if !wc.WorldBorder.Enabled {
		return 999999 // Effectively infinite
	}

	dx := x - wc.WorldBorder.CenterX
	dy := y - wc.WorldBorder.CenterY
	distance := sqrt(dx*dx + dy*dy)

	return wc.WorldBorder.Radius - distance
}

func sqrt(x float64) float64 {
	// Simple square root - use math.Sqrt in production
	if x <= 0 {
		return 0
	}
	z := 1.0
	for i := 0; i < 10; i++ {
		z = (z + x/z) / 2
	}
	return z
}

// Save saves the world configuration
func (wc *WorldConfig) Save(path string) error {
	data, err := json.MarshalIndent(wc, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal world config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write world config: %w", err)
	}

	return nil
}

// LoadWorldConfig loads a world configuration
func LoadWorldConfig(path string) (*WorldConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read world config: %w", err)
	}

	var config WorldConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal world config: %w", err)
	}

	return &config, nil
}

// GetConfigPath returns the default config path for a world
func GetConfigPath(worldName string) string {
	return filepath.Join("saves", worldName, "config.json")
}
