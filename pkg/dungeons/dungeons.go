package dungeons

import (
	"fmt"
	"time"

	"tesselbox/pkg/items"
)

// DungeonTier represents difficulty level
type DungeonTier int

const (
	TierEasy DungeonTier = iota
	TierNormal
	TierHard
	TierExpert
	TierMaster
)

// String returns tier name
func (d DungeonTier) String() string {
	switch d {
	case TierEasy:
		return "Easy"
	case TierNormal:
		return "Normal"
	case TierHard:
		return "Hard"
	case TierExpert:
		return "Expert"
	case TierMaster:
		return "Master"
	}
	return "Unknown"
}

// DungeonRoom represents a room in dungeon
type DungeonRoom struct {
	ID       string `json:"id"`
	Type     string `json:"type"` // "combat", "puzzle", "treasure", "boss"
	Cleared  bool   `json:"cleared"`
	MobCount int    `json:"mob_count"`
}

// DungeonInstance represents an active dungeon run
type DungeonInstance struct {
	ID        string `json:"id"`
	DungeonID string `json:"dungeon_id"`
	WorldID   string `json:"world_id"`

	// Players
	Players  []string `json:"players"`
	LeaderID string   `json:"leader_id"`

	// State
	Status      DungeonStatus `json:"status"`
	CurrentRoom int           `json:"current_room"`
	Rooms       []DungeonRoom `json:"rooms"`

	// Progress
	MobsKilled    int `json:"mobs_killed"`
	ChestsOpened  int `json:"chests_opened"`
	PuzzlesSolved int `json:"puzzles_solved"`

	// Timing
	StartedAt time.Time     `json:"started_at"`
	EndedAt   *time.Time    `json:"ended_at,omitempty"`
	TimeLimit time.Duration `json:"time_limit"`

	// Rewards
	BonusRewards bool `json:"bonus_rewards"` // No deaths, all rooms cleared
}

// DungeonStatus represents instance state
type DungeonStatus int

const (
	DungeonActive DungeonStatus = iota
	DungeonCompleted
	DungeonFailed
	DungeonAbandoned
)

// NewDungeonInstance creates a new dungeon run
func NewDungeonInstance(id, dungeonID, worldID, leaderID string, tier DungeonTier, playerCount int) *DungeonInstance {
	roomCount := 3 + int(tier)*2 // 3-11 rooms based on tier

	rooms := make([]DungeonRoom, roomCount)
	for i := range rooms {
		roomType := "combat"
		if i == roomCount-1 {
			roomType = "boss"
		} else if i%3 == 2 {
			roomType = "treasure"
		} else if i%2 == 1 {
			roomType = "puzzle"
		}

		rooms[i] = DungeonRoom{
			ID:       fmt.Sprintf("room_%d", i),
			Type:     roomType,
			Cleared:  false,
			MobCount: 2 + int(tier) + (i / 2),
		}
	}

	timeLimit := 30 * time.Minute
	switch tier {
	case TierHard:
		timeLimit = 45 * time.Minute
	case TierExpert:
		timeLimit = 60 * time.Minute
	case TierMaster:
		timeLimit = 90 * time.Minute
	}

	return &DungeonInstance{
		ID:          id,
		DungeonID:   dungeonID,
		WorldID:     worldID,
		Players:     []string{leaderID},
		LeaderID:    leaderID,
		Status:      DungeonActive,
		CurrentRoom: 0,
		Rooms:       rooms,
		StartedAt:   time.Now(),
		TimeLimit:   timeLimit,
	}
}

// Join adds a player to dungeon
func (di *DungeonInstance) Join(playerID string) error {
	if di.Status != DungeonActive {
		return fmt.Errorf("dungeon not active")
	}

	for _, id := range di.Players {
		if id == playerID {
			return fmt.Errorf("already in dungeon")
		}
	}

	di.Players = append(di.Players, playerID)
	return nil
}

// Leave removes a player
func (di *DungeonInstance) Leave(playerID string) {
	for i, id := range di.Players {
		if id == playerID {
			di.Players = append(di.Players[:i], di.Players[i+1:]...)
			break
		}
	}

	// If empty, abandon
	if len(di.Players) == 0 {
		di.Status = DungeonAbandoned
	}
}

// AdvanceRoom moves to next room
func (di *DungeonInstance) AdvanceRoom() bool {
	if di.CurrentRoom < len(di.Rooms)-1 {
		di.CurrentRoom++
		return true
	}
	return false
}

// ClearCurrentRoom marks current room cleared
func (di *DungeonInstance) ClearCurrentRoom() {
	if di.CurrentRoom < len(di.Rooms) {
		di.Rooms[di.CurrentRoom].Cleared = true
	}
}

// Complete completes the dungeon
func (di *DungeonInstance) Complete() {
	di.Status = DungeonCompleted
	now := time.Now()
	di.EndedAt = &now

	// Check bonus conditions
	di.BonusRewards = true
	for _, room := range di.Rooms {
		if !room.Cleared {
			di.BonusRewards = false
			break
		}
	}
}

// Fail fails the dungeon
func (di *DungeonInstance) Fail() {
	di.Status = DungeonFailed
	now := time.Now()
	di.EndedAt = &now
}

// IsExpired checks if time limit exceeded
func (di *DungeonInstance) IsExpired() bool {
	return time.Since(di.StartedAt) > di.TimeLimit
}

// GetProgress returns completion percentage
func (di *DungeonInstance) GetProgress() float64 {
	cleared := 0
	for _, room := range di.Rooms {
		if room.Cleared {
			cleared++
		}
	}
	return float64(cleared) / float64(len(di.Rooms)) * 100
}

// DungeonDefinition defines a dungeon type
type DungeonDefinition struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MinLevel    int    `json:"min_level"`
	MaxPlayers  int    `json:"max_players"`

	// Rewards by tier
	Rewards map[DungeonTier]DungeonReward `json:"rewards"`
}

// DungeonReward represents dungeon completion rewards
type DungeonReward struct {
	Money      float64      `json:"money"`
	XP         int          `json:"xp"`
	Items      []items.Item `json:"items,omitempty"`
	BonusItems []items.Item `json:"bonus_items,omitempty"`
}

// DungeonManager manages dungeons
type DungeonManager struct {
	definitions map[string]DungeonDefinition
	instances   map[string]*DungeonInstance

	instanceCounter int
}

// NewDungeonManager creates new manager
func NewDungeonManager() *DungeonManager {
	return &DungeonManager{
		definitions: make(map[string]DungeonDefinition),
		instances:   make(map[string]*DungeonInstance),
	}
}

// RegisterDungeon registers a dungeon
func (dm *DungeonManager) RegisterDungeon(def DungeonDefinition) {
	dm.definitions[def.ID] = def
}

// GetDungeon gets definition
func (dm *DungeonManager) GetDungeon(dungeonID string) (DungeonDefinition, bool) {
	def, exists := dm.definitions[dungeonID]
	return def, exists
}

// StartDungeon starts a dungeon run
func (dm *DungeonManager) StartDungeon(dungeonID, worldID, leaderID string, tier DungeonTier) (*DungeonInstance, error) {
	def, exists := dm.GetDungeon(dungeonID)
	if !exists {
		return nil, fmt.Errorf("dungeon not found")
	}

	dm.instanceCounter++
	instanceID := fmt.Sprintf("dungeon_%d_%d", dm.instanceCounter, time.Now().Unix())

	instance := NewDungeonInstance(instanceID, dungeonID, worldID, leaderID, tier, 1)
	instance.TimeLimit = time.Duration(def.MaxPlayers) * 10 * time.Minute // Dynamic time limit

	dm.instances[instanceID] = instance

	return instance, nil
}

// GetInstance gets an instance
func (dm *DungeonManager) GetInstance(instanceID string) (*DungeonInstance, bool) {
	instance, exists := dm.instances[instanceID]
	return instance, exists
}

// GetPlayerInstance gets instance player is in
func (dm *DungeonManager) GetPlayerInstance(playerID string) (*DungeonInstance, bool) {
	for _, instance := range dm.instances {
		for _, id := range instance.Players {
			if id == playerID {
				return instance, true
			}
		}
	}
	return nil, false
}

// IsInDungeon checks if player is in dungeon
func (dm *DungeonManager) IsInDungeon(playerID string) bool {
	_, exists := dm.GetPlayerInstance(playerID)
	return exists
}

// CompleteDungeon completes a dungeon
func (dm *DungeonManager) CompleteDungeon(instanceID string) (*DungeonReward, error) {
	instance, exists := dm.GetInstance(instanceID)
	if !exists {
		return nil, fmt.Errorf("instance not found")
	}

	if instance.Status != DungeonActive {
		return nil, fmt.Errorf("dungeon not active")
	}

	instance.Complete()

	// Get rewards
	def, exists := dm.GetDungeon(instance.DungeonID)
	if !exists {
		return nil, fmt.Errorf("dungeon definition not found")
	}

	reward := def.Rewards[TierNormal] // Default tier

	if instance.BonusRewards {
		reward.Items = append(reward.Items, reward.BonusItems...)
	}

	return &reward, nil
}

// Update processes all instances
func (dm *DungeonManager) Update() {
	for _, instance := range dm.instances {
		if instance.Status == DungeonActive {
			if instance.IsExpired() {
				instance.Fail()
			}
		}
	}
}

// GetActiveInstances returns active count
func (dm *DungeonManager) GetActiveInstances() []*DungeonInstance {
	result := make([]*DungeonInstance, 0)
	for _, instance := range dm.instances {
		if instance.Status == DungeonActive {
			result = append(result, instance)
		}
	}
	return result
}
