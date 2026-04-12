package stats

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// AchievementRarity represents how rare an achievement is
type AchievementRarity int

const (
	RarityCommon AchievementRarity = iota
	RarityUncommon
	RarityRare
	RarityEpic
	RarityLegendary
)

// String returns rarity name
func (r AchievementRarity) String() string {
	switch r {
	case RarityCommon:
		return "Common"
	case RarityUncommon:
		return "Uncommon"
	case RarityRare:
		return "Rare"
	case RarityEpic:
		return "Epic"
	case RarityLegendary:
		return "Legendary"
	}
	return "Unknown"
}

// AchievementReward represents rewards for completing an achievement
type AchievementReward struct {
	Money    float64 `json:"money,omitempty"`
	XP       int     `json:"xp,omitempty"`
	Title    string  `json:"title,omitempty"`
	ItemType string  `json:"item_type,omitempty"`
}

// AchievementDefinition defines an achievement
type AchievementDefinition struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Category    string              `json:"category"`
	Rarity      AchievementRarity   `json:"rarity"`
	Secret      bool                `json:"secret"` // Hidden until earned
	
	// Requirements
	TargetValue int64               `json:"target_value"` // Value needed to complete
	
	// Rewards
	Reward      AchievementReward   `json:"reward"`
	
	// Prerequisites
	Requires    []string            `json:"requires,omitempty"` // Other achievement IDs needed first
}

// PlayerAchievement represents a player's progress on an achievement
type PlayerAchievement struct {
	AchievementID string    `json:"achievement_id"`
	CurrentValue  int64     `json:"current_value"`
	Completed     bool      `json:"completed"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
}

// AchievementManager manages achievements
type AchievementManager struct {
	definitions map[string]AchievementDefinition
	playerData  map[string]map[string]PlayerAchievement // playerID -> achievementID -> data
	
	storagePath string
}

// NewAchievementManager creates a new achievement manager
func NewAchievementManager(storageDir string) *AchievementManager {
	am := &AchievementManager{
		definitions: make(map[string]AchievementDefinition),
		playerData:  make(map[string]map[string]PlayerAchievement),
		storagePath: filepath.Join(storageDir, "achievements.json"),
	}
	
	// Register default achievements
	am.registerDefaultAchievements()
	
	return am
}

// registerDefaultAchievements registers built-in achievements
func (am *AchievementManager) registerDefaultAchievements() {
	// Mining achievements
	am.RegisterAchievement(AchievementDefinition{
		ID:          "miner_novice",
		Name:        "Novice Miner",
		Description: "Mine 100 blocks",
		Category:    "mining",
		Rarity:      RarityCommon,
		TargetValue: 100,
		Reward:      AchievementReward{Money: 50, XP: 100},
	})
	
	am.RegisterAchievement(AchievementDefinition{
		ID:          "miner_apprentice",
		Name:        "Apprentice Miner",
		Description: "Mine 1,000 blocks",
		Category:    "mining",
		Rarity:      RarityUncommon,
		TargetValue: 1000,
		Requires:    []string{"miner_novice"},
		Reward:      AchievementReward{Money: 250, XP: 500},
	})
	
	am.RegisterAchievement(AchievementDefinition{
		ID:          "miner_master",
		Name:        "Master Miner",
		Description: "Mine 10,000 blocks",
		Category:    "mining",
		Rarity:      RarityRare,
		TargetValue: 10000,
		Requires:    []string{"miner_apprentice"},
		Reward:      AchievementReward{Money: 1000, XP: 2000, Title: "The Miner"},
	})
	
	// Building achievements
	am.RegisterAchievement(AchievementDefinition{
		ID:          "builder_novice",
		Name:        "Novice Builder",
		Description: "Place 100 blocks",
		Category:    "building",
		Rarity:      RarityCommon,
		TargetValue: 100,
		Reward:      AchievementReward{Money: 50, XP: 100},
	})
	
	am.RegisterAchievement(AchievementDefinition{
		ID:          "builder_architect",
		Name:        "Architect",
		Description: "Place 10,000 blocks",
		Category:    "building",
		Rarity:      RarityRare,
		TargetValue: 10000,
		Requires:    []string{"builder_novice"},
		Reward:      AchievementReward{Money: 1000, XP: 2000, Title: "The Architect"},
	})
	
	// Combat achievements
	am.RegisterAchievement(AchievementDefinition{
		ID:          "zombie_slayer",
		Name:        "Zombie Slayer",
		Description: "Kill 100 zombies",
		Category:    "combat",
		Rarity:      RarityCommon,
		TargetValue: 100,
		Reward:      AchievementReward{Money: 100, XP: 200},
	})
	
	am.RegisterAchievement(AchievementDefinition{
		ID:          "hunter",
		Name:        "Hunter",
		Description: "Kill 1,000 mobs",
		Category:    "combat",
		Rarity:      RarityUncommon,
		TargetValue: 1000,
		Requires:    []string{"zombie_slayer"},
		Reward:      AchievementReward{Money: 500, XP: 1000},
	})
	
	// Economy achievements
	am.RegisterAchievement(AchievementDefinition{
		ID:          "first_steps",
		Name:        "First Steps",
		Description: "Earn your first 100 coins",
		Category:    "economy",
		Rarity:      RarityCommon,
		TargetValue: 100,
		Reward:      AchievementReward{Money: 25, XP: 50},
	})
	
	am.RegisterAchievement(AchievementDefinition{
		ID:          "wealthy",
		Name:        "Wealthy",
		Description: "Have 10,000 coins at once",
		Category:    "economy",
		Rarity:      RarityUncommon,
		TargetValue: 10000,
		Reward:      AchievementReward{Money: 500, XP: 1000},
	})
	
	am.RegisterAchievement(AchievementDefinition{
		ID:          "millionaire",
		Name:        "Millionaire",
		Description: "Have 1,000,000 coins at once",
		Category:    "economy",
		Rarity:      RarityLegendary,
		Secret:      true,
		TargetValue: 1000000,
		Requires:    []string{"wealthy"},
		Reward:      AchievementReward{Money: 10000, XP: 5000, Title: "Millionaire"},
	})
	
	// Exploration achievements
	am.RegisterAchievement(AchievementDefinition{
		ID:          "traveler",
		Name:        "Traveler",
		Description: "Walk 10,000 blocks",
		Category:    "exploration",
		Rarity:      RarityCommon,
		TargetValue: 10000,
		Reward:      AchievementReward{Money: 100, XP: 200},
	})
	
	am.RegisterAchievement(AchievementDefinition{
		ID:          "explorer",
		Name:        "Explorer",
		Description: "Walk 100,000 blocks",
		Category:    "exploration",
		Rarity:      RarityRare,
		TargetValue: 100000,
		Requires:    []string{"traveler"},
		Reward:      AchievementReward{Money: 1000, XP: 2000, Title: "The Explorer"},
	})
	
	// Social achievements
	am.RegisterAchievement(AchievementDefinition{
		ID:          "friendly",
		Name:        "Friendly",
		Description: "Make 10 friends",
		Category:    "social",
		Rarity:      RarityCommon,
		TargetValue: 10,
		Reward:      AchievementReward{Money: 100, XP: 200},
	})
	
	am.RegisterAchievement(AchievementDefinition{
		ID:          "socialite",
		Name:        "Socialite",
		Description: "Make 50 friends",
		Category:    "social",
		Rarity:      RarityRare,
		TargetValue: 50,
		Requires:    []string{"friendly"},
		Reward:      AchievementReward{Money: 500, XP: 1000, Title: "Socialite"},
	})
	
	// Playtime achievements
	am.RegisterAchievement(AchievementDefinition{
		ID:          "dedicated",
		Name:        "Dedicated",
		Description: "Play for 1 hour",
		Category:    "playtime",
		Rarity:      RarityCommon,
		TargetValue: 60, // minutes
		Reward:      AchievementReward{Money: 50, XP: 100},
	})
	
	am.RegisterAchievement(AchievementDefinition{
		ID:          "committed",
		Name:        "Committed",
		Description: "Play for 24 hours",
		Category:    "playtime",
		Rarity:      RarityUncommon,
		TargetValue: 1440, // minutes
		Requires:    []string{"dedicated"},
		Reward:      AchievementReward{Money: 500, XP: 1000},
	})
	
	am.RegisterAchievement(AchievementDefinition{
		ID:          "veteran",
		Name:        "Veteran",
		Description: "Play for 100 hours",
		Category:    "playtime",
		Rarity:      RarityEpic,
		TargetValue: 6000, // minutes
		Requires:    []string{"committed"},
		Reward:      AchievementReward{Money: 2000, XP: 3000, Title: "Veteran"},
	})
}

// RegisterAchievement registers a new achievement
func (am *AchievementManager) RegisterAchievement(def AchievementDefinition) error {
	if _, exists := am.definitions[def.ID]; exists {
		return fmt.Errorf("achievement with ID '%s' already exists", def.ID)
	}
	
	am.definitions[def.ID] = def
	return nil
}

// GetAchievement gets an achievement definition
func (am *AchievementManager) GetAchievement(achievementID string) (AchievementDefinition, bool) {
	def, exists := am.definitions[achievementID]
	return def, exists
}

// GetAllAchievements returns all achievement definitions
func (am *AchievementManager) GetAllAchievements() []AchievementDefinition {
	result := make([]AchievementDefinition, 0, len(am.definitions))
	for _, def := range am.definitions {
		result = append(result, def)
	}
	return result
}

// GetAchievementsByCategory returns achievements in a category
func (am *AchievementManager) GetAchievementsByCategory(category string) []AchievementDefinition {
	result := make([]AchievementDefinition, 0)
	for _, def := range am.definitions {
		if def.Category == category {
			result = append(result, def)
		}
	}
	return result
}

// GetPlayerAchievement gets a player's achievement data
func (am *AchievementManager) GetPlayerAchievement(playerID, achievementID string) (PlayerAchievement, bool) {
	if playerAchievements, exists := am.playerData[playerID]; exists {
		if data, exists := playerAchievements[achievementID]; exists {
			return data, true
		}
	}
	return PlayerAchievement{}, false
}

// UpdateProgress updates progress on an achievement
func (am *AchievementManager) UpdateProgress(playerID, achievementID string, value int64) (completed bool, reward *AchievementReward, err error) {
	def, exists := am.GetAchievement(achievementID)
	if !exists {
		return false, nil, fmt.Errorf("achievement not found")
	}
	
	// Initialize player data if needed
	if _, exists := am.playerData[playerID]; !exists {
		am.playerData[playerID] = make(map[string]PlayerAchievement)
	}
	
	// Get current data
	data, exists := am.playerData[playerID][achievementID]
	if !exists {
		data = PlayerAchievement{
			AchievementID: achievementID,
			CurrentValue:  0,
			Completed:     false,
		}
	}
	
	// Check if already completed
	if data.Completed {
		return false, nil, nil
	}
	
	// Check prerequisites
	for _, prereqID := range def.Requires {
		prereq, exists := am.GetPlayerAchievement(playerID, prereqID)
		if !exists || !prereq.Completed {
			return false, nil, fmt.Errorf("prerequisites not met")
		}
	}
	
	// Update value
	if value > data.CurrentValue {
		data.CurrentValue = value
	}
	
	// Check completion
	if data.CurrentValue >= def.TargetValue {
		now := time.Now()
		data.Completed = true
		data.CompletedAt = &now
		
		am.playerData[playerID][achievementID] = data
		return true, &def.Reward, nil
	}
	
	am.playerData[playerID][achievementID] = data
	return false, nil, nil
}

// IncrementProgress increments progress by 1
func (am *AchievementManager) IncrementProgress(playerID, achievementID string) (completed bool, reward *AchievementReward, err error) {
	data, exists := am.GetPlayerAchievement(playerID, achievementID)
	if !exists {
		return am.UpdateProgress(playerID, achievementID, 1)
	}
	return am.UpdateProgress(playerID, achievementID, data.CurrentValue+1)
}

// GetPlayerAchievements returns all achievements for a player
func (am *AchievementManager) GetPlayerAchievements(playerID string) []PlayerAchievement {
	if playerAchievements, exists := am.playerData[playerID]; exists {
		result := make([]PlayerAchievement, 0, len(playerAchievements))
		for _, data := range playerAchievements {
			result = append(result, data)
		}
		return result
	}
	return []PlayerAchievement{}
}

// GetCompletedAchievements returns completed achievements for a player
func (am *AchievementManager) GetCompletedAchievements(playerID string) []PlayerAchievement {
	all := am.GetPlayerAchievements(playerID)
	completed := make([]PlayerAchievement, 0)
	
	for _, data := range all {
		if data.Completed {
			completed = append(completed, data)
		}
	}
	
	return completed
}

// GetCompletedCount returns number of completed achievements
func (am *AchievementManager) GetCompletedCount(playerID string) int {
	return len(am.GetCompletedAchievements(playerID))
}

// GetCompletionPercentage returns % of achievements completed
func (am *AchievementManager) GetCompletionPercentage(playerID string) float64 {
	total := len(am.definitions)
	if total == 0 {
		return 0
	}
	
	completed := am.GetCompletedCount(playerID)
	return float64(completed) / float64(total) * 100
}

// GetVisibleAchievements returns achievements visible to player
func (am *AchievementManager) GetVisibleAchievements(playerID string) []struct {
	Definition AchievementDefinition
	Progress   PlayerAchievement
} {
	result := make([]struct {
		Definition AchievementDefinition
		Progress   PlayerAchievement
	}, 0)
	
	for _, def := range am.definitions {
		// Show if not secret or already completed
		progress, _ := am.GetPlayerAchievement(playerID, def.ID)
		
		if !def.Secret || progress.Completed {
			result = append(result, struct {
				Definition AchievementDefinition
				Progress   PlayerAchievement
			}{def, progress})
		}
	}
	
	return result
}

// GetTitles returns earned titles for a player
func (am *AchievementManager) GetTitles(playerID string) []string {
	titles := make([]string, 0)
	
	completed := am.GetCompletedAchievements(playerID)
	for _, data := range completed {
		if def, exists := am.GetAchievement(data.AchievementID); exists {
			if def.Reward.Title != "" {
				titles = append(titles, def.Reward.Title)
			}
		}
	}
	
	return titles
}

// Save saves achievement data
func (am *AchievementManager) Save() error {
	data, err := json.MarshalIndent(am.playerData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}
	
	if err := os.WriteFile(am.storagePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	
	return nil
}

// Load loads achievement data
func (am *AchievementManager) Load() error {
	data, err := os.ReadFile(am.storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read: %w", err)
	}
	
	var loaded map[string]map[string]PlayerAchievement
	if err := json.Unmarshal(data, &loaded); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}
	
	am.playerData = loaded
	if am.playerData == nil {
		am.playerData = make(map[string]map[string]PlayerAchievement)
	}
	
	return nil
}

// GetLeaderboard returns top players by completed achievements
func (am *AchievementManager) GetLeaderboard(count int) []struct {
	PlayerID string
	Count    int
	Percent  float64
} {
	// Count achievements per player
	playerCounts := make(map[string]int)
	for playerID := range am.playerData {
		playerCounts[playerID] = am.GetCompletedCount(playerID)
	}
	
	// Convert to slice
	result := make([]struct {
		PlayerID string
		Count    int
		Percent  float64
	}, 0, len(playerCounts))
	
	totalAchievements := len(am.definitions)
	
	for id, count := range playerCounts {
		percent := 0.0
		if totalAchievements > 0 {
			percent = float64(count) / float64(totalAchievements) * 100
		}
		
		result = append(result, struct {
			PlayerID string
			Count    int
			Percent  float64
		}{id, count, percent})
	}
	
	// Sort by count
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Count < result[j].Count {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	
	if count > len(result) {
		count = len(result)
	}
	
	return result[:count]
}
