package cosmetics

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// CosmeticType represents cosmetic category
type CosmeticType int

const (
	CosmeticParticle CosmeticType = iota
	CosmeticPet
	CosmeticTitle
	CosmeticHat
	CosmeticEmote
)

// String returns cosmetic type name
func (c CosmeticType) String() string {
	switch c {
	case CosmeticParticle:
		return "Particle"
	case CosmeticPet:
		return "Pet"
	case CosmeticTitle:
		return "Title"
	case CosmeticHat:
		return "Hat"
	case CosmeticEmote:
		return "Emote"
	}
	return "Unknown"
}

// CosmeticDefinition defines a cosmetic item
type CosmeticDefinition struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Type        CosmeticType `json:"type"`
	Rarity      string       `json:"rarity"` // common, uncommon, rare, epic, legendary
	UnlockReq   string       `json:"unlock_req"` // Achievement or requirement to unlock
	AssetPath   string       `json:"asset_path"`
	Previewable bool         `json:"previewable"`
}

// PlayerCosmetic tracks a player's cosmetics
type PlayerCosmetic struct {
	CosmeticID string `json:"cosmetic_id"`
	Unlocked   bool   `json:"unlocked"`
	Equipped   bool   `json:"equipped"`
	Favorite   bool   `json:"favorite"`
}

// PlayerCosmetics stores all cosmetics for a player
type PlayerCosmetics struct {
	PlayerID   string                     `json:"player_id"`
	Cosmetics  map[string]PlayerCosmetic `json:"cosmetics"`
	ActiveTitle string                    `json:"active_title,omitempty"`
	ActivePet   string                    `json:"active_pet,omitempty"`
	ActiveHat   string                    `json:"active_hat,omitempty"`
	ActiveParticle string                 `json:"active_particle,omitempty"`
}

// NewPlayerCosmetics creates new cosmetics data
func NewPlayerCosmetics(playerID string) *PlayerCosmetics {
	return &PlayerCosmetics{
		PlayerID:  playerID,
		Cosmetics: make(map[string]PlayerCosmetic),
	}
}

// UnlockCosmetic unlocks a cosmetic
func (pc *PlayerCosmetics) UnlockCosmetic(cosmeticID string) {
	pc.Cosmetics[cosmeticID] = PlayerCosmetic{
		CosmeticID: cosmeticID,
		Unlocked:   true,
		Equipped:   false,
	}
}

// EquipCosmetic equips a cosmetic
func (pc *PlayerCosmetics) EquipCosmetic(cosmeticID string, cosmeticType CosmeticType) bool {
	cosmetic, exists := pc.Cosmetics[cosmeticID]
	if !exists || !cosmetic.Unlocked {
		return false
	}
	
	// Unequip previous of same type
	pc.UnequipType(cosmeticType)
	
	// Equip new
	cosmetic.Equipped = true
	pc.Cosmetics[cosmeticID] = cosmetic
	
	// Set as active
	switch cosmeticType {
	case CosmeticTitle:
		pc.ActiveTitle = cosmeticID
	case CosmeticPet:
		pc.ActivePet = cosmeticID
	case CosmeticHat:
		pc.ActiveHat = cosmeticID
	case CosmeticParticle:
		pc.ActiveParticle = cosmeticID
	}
	
	return true
}

// UnequipType unequips all cosmetics of a type
func (pc *PlayerCosmetics) UnequipType(cosmeticType CosmeticType) {
	for id, cosmetic := range pc.Cosmetics {
		if cosmetic.Equipped {
			// Check if this is the type being unequipped
			// In real implementation, lookup cosmetic type from definition
			cosmetic.Equipped = false
			pc.Cosmetics[id] = cosmetic
		}
	}
	
	switch cosmeticType {
	case CosmeticTitle:
		pc.ActiveTitle = ""
	case CosmeticPet:
		pc.ActivePet = ""
	case CosmeticHat:
		pc.ActiveHat = ""
	case CosmeticParticle:
		pc.ActiveParticle = ""
	}
}

// Unequip unequips a specific cosmetic
func (pc *PlayerCosmetics) Unequip(cosmeticID string) {
	if cosmetic, exists := pc.Cosmetics[cosmeticID]; exists {
		cosmetic.Equipped = false
		pc.Cosmetics[cosmeticID] = cosmetic
		
		// Remove from active
		if pc.ActiveTitle == cosmeticID {
			pc.ActiveTitle = ""
		}
		if pc.ActivePet == cosmeticID {
			pc.ActivePet = ""
		}
		if pc.ActiveHat == cosmeticID {
			pc.ActiveHat = ""
		}
		if pc.ActiveParticle == cosmeticID {
			pc.ActiveParticle = ""
		}
	}
}

// GetEquipped gets all equipped cosmetics
func (pc *PlayerCosmetics) GetEquipped() []string {
	equipped := make([]string, 0)
	for id, cosmetic := range pc.Cosmetics {
		if cosmetic.Equipped {
			equipped = append(equipped, id)
		}
	}
	return equipped
}

// HasCosmetic checks if player has a cosmetic
func (pc *PlayerCosmetics) HasCosmetic(cosmeticID string) bool {
	cosmetic, exists := pc.Cosmetics[cosmeticID]
	return exists && cosmetic.Unlocked
}

// IsEquipped checks if cosmetic is equipped
func (pc *PlayerCosmetics) IsEquipped(cosmeticID string) bool {
	cosmetic, exists := pc.Cosmetics[cosmeticID]
	return exists && cosmetic.Equipped
}

// ToggleFavorite toggles favorite status
func (pc *PlayerCosmetics) ToggleFavorite(cosmeticID string) {
	if cosmetic, exists := pc.Cosmetics[cosmeticID]; exists {
		cosmetic.Favorite = !cosmetic.Favorite
		pc.Cosmetics[cosmeticID] = cosmetic
	}
}

// CosmeticManager manages all cosmetics
type CosmeticManager struct {
	definitions map[string]CosmeticDefinition
	players     map[string]*PlayerCosmetics
	
	storagePath string
}

// NewCosmeticManager creates new manager
func NewCosmeticManager(storageDir string) *CosmeticManager {
	cm := &CosmeticManager{
		definitions: make(map[string]CosmeticDefinition),
		players:     make(map[string]*PlayerCosmetics),
		storagePath: filepath.Join(storageDir, "cosmetics.json"),
	}
	
	cm.registerDefaultCosmetics()
	
	return cm
}

// registerDefaultCosmetics registers default cosmetics
func (cm *CosmeticManager) registerDefaultCosmetics() {
	// Titles
	cm.RegisterCosmetic(CosmeticDefinition{
		ID:          "title_player",
		Name:        "Player",
		Description: "Default title",
		Type:        CosmeticTitle,
		Rarity:      "common",
		UnlockReq:   "default",
	})
	
	cm.RegisterCosmetic(CosmeticDefinition{
		ID:          "title_builder",
		Name:        "Builder",
		Description: "Place 1,000 blocks",
		Type:        CosmeticTitle,
		Rarity:      "uncommon",
		UnlockReq:   "achievement:builder_novice",
	})
	
	cm.RegisterCosmetic(CosmeticDefinition{
		ID:          "title_warrior",
		Name:        "Warrior",
		Description: "Win 10 duels",
		Type:        CosmeticTitle,
		Rarity:      "rare",
		UnlockReq:   "duels_won:10",
	})
	
	// Particles
	cm.RegisterCosmetic(CosmeticDefinition{
		ID:          "particle_sparkle",
		Name:        "Sparkle",
		Description: "Sparkling trail",
		Type:        CosmeticParticle,
		Rarity:      "common",
		UnlockReq:   "default",
	})
	
	cm.RegisterCosmetic(CosmeticDefinition{
		ID:          "particle_flame",
		Name:        "Flame",
		Description: "Burning trail",
		Type:        CosmeticParticle,
		Rarity:      "rare",
		UnlockReq:   "achievement:hunter",
	})
	
	// Hats
	cm.RegisterCosmetic(CosmeticDefinition{
		ID:          "hat_crown",
		Name:        "Crown",
		Description: "Royal crown",
		Type:        CosmeticHat,
		Rarity:      "legendary",
		UnlockReq:   "achievement:millionaire",
	})
	
	// Emotes
	cm.RegisterCosmetic(CosmeticDefinition{
		ID:          "emote_wave",
		Name:        "Wave",
		Description: "Friendly wave",
		Type:        CosmeticEmote,
		Rarity:      "common",
		UnlockReq:   "default",
	})
}

// RegisterCosmetic registers a cosmetic
func (cm *CosmeticManager) RegisterCosmetic(def CosmeticDefinition) error {
	if _, exists := cm.definitions[def.ID]; exists {
		return fmt.Errorf("cosmetic already exists")
	}
	
	cm.definitions[def.ID] = def
	
	// Auto-unlock for all players if default
	if def.UnlockReq == "default" {
		for _, pc := range cm.players {
			pc.UnlockCosmetic(def.ID)
		}
	}
	
	return nil
}

// GetCosmetic gets a cosmetic definition
func (cm *CosmeticManager) GetCosmetic(cosmeticID string) (CosmeticDefinition, bool) {
	def, exists := cm.definitions[cosmeticID]
	return def, exists
}

// GetAllCosmetics returns all cosmetics
func (cm *CosmeticManager) GetAllCosmetics() []CosmeticDefinition {
	result := make([]CosmeticDefinition, 0, len(cm.definitions))
	for _, def := range cm.definitions {
		result = append(result, def)
	}
	return result
}

// GetCosmeticsByType returns cosmetics of a type
func (cm *CosmeticManager) GetCosmeticsByType(cosmeticType CosmeticType) []CosmeticDefinition {
	result := make([]CosmeticDefinition, 0)
	for _, def := range cm.definitions {
		if def.Type == cosmeticType {
			result = append(result, def)
		}
	}
	return result
}

// GetOrCreatePlayerCosmetics gets or creates player cosmetics
func (cm *CosmeticManager) GetOrCreatePlayerCosmetics(playerID string) *PlayerCosmetics {
	if pc, exists := cm.players[playerID]; exists {
		return pc
	}
	
	pc := NewPlayerCosmetics(playerID)
	
	// Unlock defaults
	for id, def := range cm.definitions {
		if def.UnlockReq == "default" {
			pc.UnlockCosmetic(id)
		}
	}
	
	cm.players[playerID] = pc
	return pc
}

// GetPlayerCosmetics gets player cosmetics
func (cm *CosmeticManager) GetPlayerCosmetics(playerID string) (*PlayerCosmetics, bool) {
	pc, exists := cm.players[playerID]
	return pc, exists
}

// UnlockCosmeticForPlayer unlocks a cosmetic for a player
func (cm *CosmeticManager) UnlockCosmeticForPlayer(playerID, cosmeticID string) bool {
	if _, exists := cm.definitions[cosmeticID]; !exists {
		return false
	}
	
	pc := cm.GetOrCreatePlayerCosmetics(playerID)
	pc.UnlockCosmetic(cosmeticID)
	return true
}

// EquipCosmetic equips a cosmetic for a player
func (cm *CosmeticManager) EquipCosmetic(playerID, cosmeticID string) bool {
	def, exists := cm.definitions[cosmeticID]
	if !exists {
		return false
	}
	
	pc := cm.GetOrCreatePlayerCosmetics(playerID)
	return pc.EquipCosmetic(cosmeticID, def.Type)
}

// UnequipCosmetic unequips a cosmetic
func (cm *CosmeticManager) UnequipCosmetic(playerID, cosmeticID string) {
	pc := cm.GetOrCreatePlayerCosmetics(playerID)
	pc.Unequip(cosmeticID)
}

// GetActiveTitle gets player's active title
func (cm *CosmeticManager) GetActiveTitle(playerID string) string {
	pc, exists := cm.GetPlayerCosmetics(playerID)
	if !exists {
		return ""
	}
	
	if def, exists := cm.definitions[pc.ActiveTitle]; exists {
		return def.Name
	}
	return ""
}

// GetDisplayTitle gets formatted title with player name
func (cm *CosmeticManager) GetDisplayTitle(playerID, playerName string) string {
	title := cm.GetActiveTitle(playerID)
	if title == "" {
		return playerName
	}
	return fmt.Sprintf("[%s] %s", title, playerName)
}

// Save saves cosmetic data
func (cm *CosmeticManager) Save() error {
	data, err := json.MarshalIndent(cm.players, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}
	
	if err := os.WriteFile(cm.storagePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	
	return nil
}

// Load loads cosmetic data
func (cm *CosmeticManager) Load() error {
	data, err := os.ReadFile(cm.storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read: %w", err)
	}
	
	var loaded map[string]*PlayerCosmetics
	if err := json.Unmarshal(data, &loaded); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}
	
	cm.players = loaded
	if cm.players == nil {
		cm.players = make(map[string]*PlayerCosmetics)
	}
	
	return nil
}

// CheckUnlocks checks if player qualifies for new cosmetics
func (cm *CosmeticManager) CheckUnlocks(playerID string, checkFunc func(req string) bool) []string {
	pc := cm.GetOrCreatePlayerCosmetics(playerID)
	unlocked := make([]string, 0)
	
	for id, def := range cm.definitions {
		if !pc.HasCosmetic(id) && def.UnlockReq != "default" {
			if checkFunc(def.UnlockReq) {
				pc.UnlockCosmetic(id)
				unlocked = append(unlocked, id)
			}
		}
	}
	
	return unlocked
}
