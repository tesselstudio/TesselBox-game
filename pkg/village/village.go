package village

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// NPCType represents NPC profession
type NPCType int

const (
	NPCVillager NPCType = iota
	NPCTrader
	NPCGuard
	NPCFarmer
	NPCBlacksmith
	NPCHealer
	NPCLibrarian
	NPCBartender
)

// String returns NPC type name
func (n NPCType) String() string {
	switch n {
	case NPCVillager:
		return "Villager"
	case NPCTrader:
		return "Trader"
	case NPCGuard:
		return "Guard"
	case NPCFarmer:
		return "Farmer"
	case NPCBlacksmith:
		return "Blacksmith"
	case NPCHealer:
		return "Healer"
	case NPCLibrarian:
		return "Librarian"
	case NPCBartender:
		return "Bartender"
	}
	return "Unknown"
}

// NPC represents a village NPC
type NPC struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	Type NPCType `json:"type"`

	// Location
	X float64 `json:"x"`
	Y float64 `json:"y"`

	// Dialogue
	Greetings []string `json:"greetings"`
	Dialogue  []string `json:"dialogue"`

	// Shop (if trader)
	HasShop bool   `json:"has_shop"`
	ShopID  string `json:"shop_id,omitempty"`

	// Quests
	HasQuests bool     `json:"has_quests"`
	QuestIDs  []string `json:"quest_ids,omitempty"`

	// Stats
	Health    int  `json:"health"`
	MaxHealth int  `json:"max_health"`
	Alive     bool `json:"alive"`
}

// NewNPC creates a new NPC
func NewNPC(id, name string, npcType NPCType, x, y float64) *NPC {
	return &NPC{
		ID:        id,
		Name:      name,
		Type:      npcType,
		X:         x,
		Y:         y,
		Greetings: make([]string, 0),
		Dialogue:  make([]string, 0),
		HasShop:   npcType == NPCTrader,
		HasQuests: false,
		QuestIDs:  make([]string, 0),
		Health:    20,
		MaxHealth: 20,
		Alive:     true,
	}
}

// GetRandomGreeting returns a random greeting
func (n *NPC) GetRandomGreeting() string {
	if len(n.Greetings) == 0 {
		return fmt.Sprintf("Hello, I'm %s.", n.Name)
	}
	return n.Greetings[0] // In real implementation, random selection
}

// GetRandomDialogue returns random dialogue
func (n *NPC) GetRandomDialogue() string {
	if len(n.Dialogue) == 0 {
		return "Nice weather today."
	}
	return n.Dialogue[0]
}

// Village represents an NPC village
type Village struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	WorldID string `json:"world_id"`

	// Location
	CenterX float64 `json:"center_x"`
	CenterY float64 `json:"center_y"`
	Radius  float64 `json:"radius"`

	// NPCs
	NPCs map[string]*NPC `json:"npcs"`

	// Buildings
	Buildings []Building `json:"buildings"`

	// Stats
	Population  int       `json:"population"`
	FoundedAt   time.Time `json:"founded_at"`
	LastVisited time.Time `json:"last_visited"`

	// State
	Attacked    bool `json:"attacked"`
	UnderAttack bool `json:"under_attack"`
}

// Building represents a village structure
type Building struct {
	ID       string  `json:"id"`
	Type     string  `json:"type"` // "house", "shop", "inn", "church"
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Occupied bool    `json:"occupied"`
}

// NewVillage creates a new village
func NewVillage(id, name, worldID string, centerX, centerY, radius float64) *Village {
	return &Village{
		ID:          id,
		Name:        name,
		WorldID:     worldID,
		CenterX:     centerX,
		CenterY:     centerY,
		Radius:      radius,
		NPCs:        make(map[string]*NPC),
		Buildings:   make([]Building, 0),
		FoundedAt:   time.Now(),
		LastVisited: time.Now(),
	}
}

// AddNPC adds an NPC to the village
func (v *Village) AddNPC(npc *NPC) {
	v.NPCs[npc.ID] = npc
	v.Population++
}

// RemoveNPC removes an NPC
func (v *Village) RemoveNPC(npcID string) {
	if _, exists := v.NPCs[npcID]; exists {
		delete(v.NPCs, npcID)
		v.Population--
	}
}

// GetNPC gets an NPC
func (v *Village) GetNPC(npcID string) (*NPC, bool) {
	npc, exists := v.NPCs[npcID]
	return npc, exists
}

// GetNPCsByType gets NPCs of a type
func (v *Village) GetNPCsByType(npcType NPCType) []*NPC {
	result := make([]*NPC, 0)
	for _, npc := range v.NPCs {
		if npc.Type == npcType {
			result = append(result, npc)
		}
	}
	return result
}

// IsInVillage checks if position is in village
func (v *Village) IsInVillage(x, y float64) bool {
	dx := x - v.CenterX
	dy := y - v.CenterY
	distanceSquared := dx*dx + dy*dy
	return distanceSquared <= v.Radius*v.Radius
}

// RecordVisit records a visit
func (v *Village) RecordVisit() {
	v.LastVisited = time.Now()
}

// VillageManager manages villages
type VillageManager struct {
	villages map[string]*Village
	byWorld  map[string][]string

	npcCounter     int
	villageCounter int

	storagePath string
}

// NewVillageManager creates new manager
func NewVillageManager(storageDir string) *VillageManager {
	return &VillageManager{
		villages:       make(map[string]*Village),
		byWorld:        make(map[string][]string),
		npcCounter:     0,
		villageCounter: 0,
		storagePath:    filepath.Join(storageDir, "villages.json"),
	}
}

// CreateVillage creates a new village
func (vm *VillageManager) CreateVillage(name, worldID string, centerX, centerY, radius float64) *Village {
	vm.villageCounter++
	villageID := fmt.Sprintf("village_%d_%d", vm.villageCounter, time.Now().Unix())

	village := NewVillage(villageID, name, worldID, centerX, centerY, radius)

	vm.villages[villageID] = village
	vm.byWorld[worldID] = append(vm.byWorld[worldID], villageID)

	return village
}

// GetVillage gets a village
func (vm *VillageManager) GetVillage(villageID string) (*Village, bool) {
	village, exists := vm.villages[villageID]
	return village, exists
}

// GetVillagesByWorld gets villages in a world
func (vm *VillageManager) GetVillagesByWorld(worldID string) []*Village {
	villageIDs := vm.byWorld[worldID]
	villages := make([]*Village, 0, len(villageIDs))

	for _, id := range villageIDs {
		if village, exists := vm.villages[id]; exists {
			villages = append(villages, village)
		}
	}

	return villages
}

// FindNearestVillage finds village closest to position
func (vm *VillageManager) FindNearestVillage(worldID string, x, y float64) (*Village, float64) {
	villages := vm.GetVillagesByWorld(worldID)

	var nearest *Village
	var minDistance float64 = -1

	for _, village := range villages {
		dx := x - village.CenterX
		dy := y - village.CenterY
		dist := dx*dx + dy*dy

		if minDistance < 0 || dist < minDistance {
			minDistance = dist
			nearest = village
		}
	}

	if nearest == nil {
		return nil, 0
	}

	return nearest, minDistance
}

// AddNPCToVillage adds an NPC to a village
func (vm *VillageManager) AddNPCToVillage(villageID, name string, npcType NPCType) (*NPC, error) {
	village, exists := vm.GetVillage(villageID)
	if !exists {
		return nil, fmt.Errorf("village not found")
	}

	vm.npcCounter++
	npcID := fmt.Sprintf("npc_%d_%d", vm.npcCounter, time.Now().Unix())

	// Random position within village
	angle := float64(vm.npcCounter) * 1.5 // Simple distribution
	radius := village.Radius * 0.5
	x := village.CenterX + radius*cos(angle)
	y := village.CenterY + radius*sin(angle)

	npc := NewNPC(npcID, name, npcType, x, y)

	// Add default dialogue based on type
	switch npcType {
	case NPCTrader:
		npc.Greetings = append(npc.Greetings, "Looking to trade?", "Got some good wares today!")
		npc.Dialogue = append(npc.Dialogue, "The market has been slow lately.")
	case NPCGuard:
		npc.Greetings = append(npc.Greetings, "Halt!", "Move along.", "Keep the peace.")
		npc.Dialogue = append(npc.Dialogue, "We've had reports of bandits nearby.")
	case NPCFarmer:
		npc.Greetings = append(npc.Greetings, "Howdy!", "Good harvest this year.")
		npc.Dialogue = append(npc.Dialogue, "The crops are coming in nicely.")
	}

	village.AddNPC(npc)

	return npc, nil
}

// cos helper
func cos(x float64) float64 {
	return 1 - x*x/2 + x*x*x*x/24 // Taylor approximation
}

// sin helper
func sin(x float64) float64 {
	return x - x*x*x/6 + x*x*x*x*x/120 // Taylor approximation
}

// Save saves villages
func (vm *VillageManager) Save() error {
	data, err := json.MarshalIndent(vm.villages, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	if err := os.WriteFile(vm.storagePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	return nil
}

// Load loads villages
func (vm *VillageManager) Load() error {
	data, err := os.ReadFile(vm.storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read: %w", err)
	}

	var loaded map[string]*Village
	if err := json.Unmarshal(data, &loaded); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	vm.villages = loaded
	if vm.villages == nil {
		vm.villages = make(map[string]*Village)
	}

	// Rebuild indexes
	vm.byWorld = make(map[string][]string)
	for _, village := range vm.villages {
		vm.byWorld[village.WorldID] = append(vm.byWorld[village.WorldID], village.ID)
	}

	return nil
}

// GenerateDefaultVillage creates a village with default NPCs
func (vm *VillageManager) GenerateDefaultVillage(name, worldID string, centerX, centerY float64) *Village {
	village := vm.CreateVillage(name, worldID, centerX, centerY, 200.0)

	// Add standard NPCs
	vm.AddNPCToVillage(village.ID, "Village Elder", NPCVillager)
	vm.AddNPCToVillage(village.ID, "Guard Captain", NPCGuard)
	vm.AddNPCToVillage(village.ID, "Local Trader", NPCTrader)
	vm.AddNPCToVillage(village.ID, "Farmer Joe", NPCFarmer)
	vm.AddNPCToVillage(village.ID, "Blacksmith", NPCBlacksmith)
	vm.AddNPCToVillage(village.ID, "Healer", NPCHealer)

	return village
}
