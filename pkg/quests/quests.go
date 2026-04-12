package quests

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"tesselbox/pkg/items"
)

// QuestType represents quest category
type QuestType int

const (
	QuestKill QuestType = iota
	QuestGather
	QuestDeliver
	QuestExplore
	QuestEscort
)

// QuestStatus represents quest state
type QuestStatus int

const (
	QuestAvailable QuestStatus = iota
	QuestActive
	QuestCompleted
	QuestTurnedIn
)

// QuestObjective represents a quest objective
type QuestObjective struct {
	Type        string `json:"type"`   // "kill", "gather", "deliver", "reach"
	Target      string `json:"target"` // mob type, item type, location, etc
	Amount      int    `json:"amount"`
	Current     int    `json:"current"`
	Completed   bool   `json:"completed"`
	Description string `json:"description"`
}

// QuestReward represents quest rewards
type QuestReward struct {
	Money      float64      `json:"money"`
	XP         int          `json:"xp"`
	Items      []items.Item `json:"items,omitempty"`
	Reputation string       `json:"reputation,omitempty"` // Faction rep gain
}

// QuestDefinition defines a quest
type QuestDefinition struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        QuestType `json:"type"`
	LevelReq    int       `json:"level_req"`
	PrereqQuest string    `json:"prereq_quest,omitempty"`

	// Content
	Objectives []QuestObjective `json:"objectives"`
	Reward     QuestReward      `json:"reward"`

	// Meta
	Repeatable bool          `json:"repeatable"`
	Cooldown   time.Duration `json:"cooldown"`
	TimeLimit  time.Duration `json:"time_limit"`
}

// PlayerQuest tracks a player's quest progress
type PlayerQuest struct {
	QuestID     string           `json:"quest_id"`
	Status      QuestStatus      `json:"status"`
	Progress    []QuestObjective `json:"progress"`
	AcceptedAt  time.Time        `json:"accepted_at"`
	ExpiresAt   *time.Time       `json:"expires_at,omitempty"`
	CompletedAt *time.Time       `json:"completed_at,omitempty"`
	TurnedInAt  *time.Time       `json:"turned_in_at,omitempty"`
}

// QuestManager manages quests
type QuestManager struct {
	definitions  map[string]QuestDefinition
	playerQuests map[string]map[string]PlayerQuest // playerID -> questID -> quest

	storagePath string
}

// NewQuestManager creates new manager
func NewQuestManager(storageDir string) *QuestManager {
	qm := &QuestManager{
		definitions:  make(map[string]QuestDefinition),
		playerQuests: make(map[string]map[string]PlayerQuest),
		storagePath:  filepath.Join(storageDir, "quests.json"),
	}

	qm.registerDefaultQuests()

	return qm
}

// registerDefaultQuests registers built-in quests
func (qm *QuestManager) registerDefaultQuests() {
	// Tutorial quest
	qm.RegisterQuest(QuestDefinition{
		ID:          "tutorial_beginner",
		Name:        "Getting Started",
		Description: "Break 10 blocks and place 10 blocks to learn the basics",
		Type:        QuestGather,
		Objectives: []QuestObjective{
			{Type: "break", Target: "any", Amount: 10, Description: "Break 10 blocks"},
			{Type: "place", Target: "any", Amount: 10, Description: "Place 10 blocks"},
		},
		Reward: QuestReward{Money: 50, XP: 100},
	})

	// Kill quest
	qm.RegisterQuest(QuestDefinition{
		ID:          "hunter_initiation",
		Name:        "Hunter's Initiation",
		Description: "Kill 5 zombies to prove your combat skills",
		Type:        QuestKill,
		LevelReq:    5,
		Objectives: []QuestObjective{
			{Type: "kill", Target: "zombie", Amount: 5, Description: "Kill 5 zombies"},
		},
		Reward: QuestReward{Money: 100, XP: 200, Items: []items.Item{
			{Type: 1, Quantity: 1},
		}},
	})

	// Gather quest
	qm.RegisterQuest(QuestDefinition{
		ID:          "miner_prospector",
		Name:        "Prospector",
		Description: "Gather valuable ores from the mines",
		Type:        QuestGather,
		LevelReq:    10,
		PrereqQuest: "tutorial_beginner",
		Objectives: []QuestObjective{
			{Type: "gather", Target: "coal_ore", Amount: 20, Description: "Mine 20 coal ore"},
			{Type: "gather", Target: "iron_ore", Amount: 10, Description: "Mine 10 iron ore"},
		},
		Reward: QuestReward{Money: 250, XP: 500},
	})

	// Daily repeatable quest
	qm.RegisterQuest(QuestDefinition{
		ID:          "daily_defense",
		Name:        "Daily Defense",
		Description: "Defend the village by killing 20 mobs",
		Type:        QuestKill,
		Objectives: []QuestObjective{
			{Type: "kill", Target: "any", Amount: 20, Description: "Kill 20 mobs"},
		},
		Reward:     QuestReward{Money: 200, XP: 300},
		Repeatable: true,
		Cooldown:   24 * time.Hour,
	})
}

// RegisterQuest registers a quest
func (qm *QuestManager) RegisterQuest(quest QuestDefinition) error {
	if _, exists := qm.definitions[quest.ID]; exists {
		return fmt.Errorf("quest already exists")
	}

	qm.definitions[quest.ID] = quest
	return nil
}

// GetQuest gets a quest definition
func (qm *QuestManager) GetQuest(questID string) (QuestDefinition, bool) {
	quest, exists := qm.definitions[questID]
	return quest, exists
}

// GetAvailableQuests gets quests available to a player
func (qm *QuestManager) GetAvailableQuests(playerID string, playerLevel int) []QuestDefinition {
	available := make([]QuestDefinition, 0)
	playerData := qm.getPlayerQuests(playerID)

	for _, quest := range qm.definitions {
		// Check level requirement
		if playerLevel < quest.LevelReq {
			continue
		}

		// Check if already active or completed (non-repeatable)
		if playerQuest, exists := playerData[quest.ID]; exists {
			if playerQuest.Status == QuestActive {
				continue
			}
			if playerQuest.Status == QuestTurnedIn && !quest.Repeatable {
				continue
			}
			// Check cooldown for repeatable
			if quest.Repeatable && playerQuest.TurnedInAt != nil {
				if time.Since(*playerQuest.TurnedInAt) < quest.Cooldown {
					continue
				}
			}
		}

		// Check prerequisites
		if quest.PrereqQuest != "" {
			if prereq, exists := playerData[quest.PrereqQuest]; !exists || prereq.Status != QuestTurnedIn {
				continue
			}
		}

		available = append(available, quest)
	}

	return available
}

// getPlayerQuests gets or creates player quest data
func (qm *QuestManager) getPlayerQuests(playerID string) map[string]PlayerQuest {
	if _, exists := qm.playerQuests[playerID]; !exists {
		qm.playerQuests[playerID] = make(map[string]PlayerQuest)
	}
	return qm.playerQuests[playerID]
}

// AcceptQuest accepts a quest
func (qm *QuestManager) AcceptQuest(playerID, questID string) (*PlayerQuest, error) {
	quest, exists := qm.GetQuest(questID)
	if !exists {
		return nil, fmt.Errorf("quest not found")
	}

	playerData := qm.getPlayerQuests(playerID)

	// Check if already active
	if existing, exists := playerData[questID]; exists && existing.Status == QuestActive {
		return nil, fmt.Errorf("quest already active")
	}

	// Create quest instance
	playerQuest := PlayerQuest{
		QuestID:    questID,
		Status:     QuestActive,
		Progress:   make([]QuestObjective, len(quest.Objectives)),
		AcceptedAt: time.Now(),
	}

	// Copy objectives
	for i, obj := range quest.Objectives {
		playerQuest.Progress[i] = obj
		playerQuest.Progress[i].Current = 0
	}

	// Set expiration if time limited
	if quest.TimeLimit > 0 {
		expires := playerQuest.AcceptedAt.Add(quest.TimeLimit)
		playerQuest.ExpiresAt = &expires
	}

	playerData[questID] = playerQuest

	return &playerQuest, nil
}

// UpdateProgress updates quest progress
func (qm *QuestManager) UpdateProgress(playerID, objectiveType, target string, amount int) []PlayerQuest {
	updated := make([]PlayerQuest, 0)
	playerData := qm.getPlayerQuests(playerID)

	for questID, playerQuest := range playerData {
		if playerQuest.Status != QuestActive {
			continue
		}

		updatedAny := false
		for i := range playerQuest.Progress {
			obj := &playerQuest.Progress[i]
			if obj.Type == objectiveType && (obj.Target == target || obj.Target == "any") {
				if !obj.Completed {
					obj.Current += amount
					if obj.Current >= obj.Amount {
						obj.Current = obj.Amount
						obj.Completed = true
					}
					updatedAny = true
				}
			}
		}

		if updatedAny {
			// Check if all objectives complete
			allComplete := true
			for _, obj := range playerQuest.Progress {
				if !obj.Completed {
					allComplete = false
					break
				}
			}

			if allComplete {
				playerQuest.Status = QuestCompleted
				now := time.Now()
				playerQuest.CompletedAt = &now
			}

			playerData[questID] = playerQuest
			updated = append(updated, playerQuest)
		}
	}

	return updated
}

// CompleteQuest completes a quest and gives rewards
func (qm *QuestManager) CompleteQuest(playerID, questID string) (*QuestReward, error) {
	playerData := qm.getPlayerQuests(playerID)

	playerQuest, exists := playerData[questID]
	if !exists {
		return nil, fmt.Errorf("quest not found")
	}

	if playerQuest.Status != QuestCompleted {
		return nil, fmt.Errorf("quest not completed")
	}

	quest, exists := qm.GetQuest(questID)
	if !exists {
		return nil, fmt.Errorf("quest definition not found")
	}

	// Mark as turned in
	playerQuest.Status = QuestTurnedIn
	now := time.Now()
	playerQuest.TurnedInAt = &now
	playerData[questID] = playerQuest

	return &quest.Reward, nil
}

// GetActiveQuests gets player's active quests
func (qm *QuestManager) GetActiveQuests(playerID string) []PlayerQuest {
	playerData := qm.getPlayerQuests(playerID)

	active := make([]PlayerQuest, 0)
	for _, quest := range playerData {
		if quest.Status == QuestActive {
			active = append(active, quest)
		}
	}

	return active
}

// AbandonQuest abandons a quest
func (qm *QuestManager) AbandonQuest(playerID, questID string) error {
	playerData := qm.getPlayerQuests(playerID)

	playerQuest, exists := playerData[questID]
	if !exists {
		return fmt.Errorf("quest not found")
	}

	if playerQuest.Status != QuestActive {
		return fmt.Errorf("quest not active")
	}

	delete(playerData, questID)
	return nil
}

// CheckExpired checks and removes expired quests
func (qm *QuestManager) CheckExpired() {
	for playerID, playerData := range qm.playerQuests {
		for questID, quest := range playerData {
			if quest.Status == QuestActive && quest.ExpiresAt != nil {
				if time.Now().After(*quest.ExpiresAt) {
					delete(playerData, questID)
				}
			}
		}
		qm.playerQuests[playerID] = playerData
	}
}

// Save saves quest data
func (qm *QuestManager) Save() error {
	data, err := json.MarshalIndent(qm.playerQuests, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	if err := os.WriteFile(qm.storagePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	return nil
}

// Load loads quest data
func (qm *QuestManager) Load() error {
	data, err := os.ReadFile(qm.storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read: %w", err)
	}

	var loaded map[string]map[string]PlayerQuest
	if err := json.Unmarshal(data, &loaded); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	qm.playerQuests = loaded
	if qm.playerQuests == nil {
		qm.playerQuests = make(map[string]map[string]PlayerQuest)
	}

	return nil
}
