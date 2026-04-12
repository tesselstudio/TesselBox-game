package social

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"
)

// KarmaLevel represents karma tier
type KarmaLevel int

const (
	KarmaEvil KarmaLevel = iota
	KarmaBad
	KarmaNeutral
	KarmaGood
	KarmaHero
)

// String returns karma level name
func (k KarmaLevel) String() string {
	switch k {
	case KarmaEvil:
		return "Evil"
	case KarmaBad:
		return "Bad"
	case KarmaNeutral:
		return "Neutral"
	case KarmaGood:
		return "Good"
	case KarmaHero:
		return "Hero"
	}
	return "Unknown"
}

// ReputationAction represents actions that affect karma
type ReputationAction int

const (
	ActionKillMob ReputationAction = iota
	ActionKillPlayer
	ActionHelpPlayer
	ActionTrade
	ActionGive
	ActionSteal
	ActionGrief
	ActionHelpNewbie
	ActionVote
	ActionReport
)

// KarmaChange represents karma values for actions
var KarmaValues = map[ReputationAction]float64{
	ActionKillMob:     1,
	ActionKillPlayer:  -10,
	ActionHelpPlayer:  5,
	ActionTrade:       2,
	ActionGive:        3,
	ActionSteal:       -20,
	ActionGrief:       -50,
	ActionHelpNewbie:  10,
	ActionVote:        2,
	ActionReport:      1,
}

// PlayerReputation tracks a player's reputation
type PlayerReputation struct {
	PlayerID       string                 `json:"player_id"`
	Karma          float64                `json:"karma"`
	KarmaLevel     KarmaLevel             `json:"karma_level"`
	Reputation     float64                `json:"reputation"` // 0-100
	
	// Stats
	GoodDeeds      int                    `json:"good_deeds"`
	BadDeeds       int                    `json:"bad_deeds"`
	PlayersHelped  int                    `json:"players_helped"`
	PlayersKilled  int                    `json:"players_killed"`
	ItemsGiven     int64                  `json:"items_given"`
	MoneyGiven     float64                `json:"money_given"`
	
	// History
	RecentActions  []ReputationActionEntry `json:"recent_actions,omitempty"`
	
	// Ratings from other players
	PlayerRatings  map[string]int         `json:"player_ratings,omitempty"` // PlayerID -> rating (1-5)
	
	// Timestamps
	LastUpdated    time.Time              `json:"last_updated"`
}

// ReputationActionEntry tracks a single action
type ReputationActionEntry struct {
	Action    ReputationAction `json:"action"`
	TargetID  string           `json:"target_id,omitempty"`
	Value     float64          `json:"value"`
	Timestamp time.Time        `json:"timestamp"`
}

// NewPlayerReputation creates new reputation data
func NewPlayerReputation(playerID string) *PlayerReputation {
	return &PlayerReputation{
		PlayerID:      playerID,
		Karma:         0,
		KarmaLevel:    KarmaNeutral,
		Reputation:    50, // Start neutral
		GoodDeeds:     0,
		BadDeeds:      0,
		RecentActions: make([]ReputationActionEntry, 0),
		PlayerRatings: make(map[string]int),
		LastUpdated:   time.Now(),
	}
}

// RecordAction records a reputation action
func (pr *PlayerReputation) RecordAction(action ReputationAction, targetID string) {
	value := KarmaValues[action]
	pr.Karma += value
	
	// Track good/bad deeds
	if value > 0 {
		pr.GoodDeeds++
	} else if value < 0 {
		pr.BadDeeds++
	}
	
	// Track specific stats
	switch action {
	case ActionKillPlayer:
		pr.PlayersKilled++
	case ActionHelpPlayer, ActionHelpNewbie:
		pr.PlayersHelped++
	case ActionGive:
		pr.ItemsGiven++
	}
	
	// Add to history
	entry := ReputationActionEntry{
		Action:    action,
		TargetID:  targetID,
		Value:     value,
		Timestamp: time.Now(),
	}
	pr.RecentActions = append(pr.RecentActions, entry)
	
	// Keep last 50 actions
	if len(pr.RecentActions) > 50 {
		pr.RecentActions = pr.RecentActions[len(pr.RecentActions)-50:]
	}
	
	pr.updateKarmaLevel()
	pr.LastUpdated = time.Now()
}

// updateKarmaLevel updates karma tier
func (pr *PlayerReputation) updateKarmaLevel() {
	switch {
	case pr.Karma <= -100:
		pr.KarmaLevel = KarmaEvil
	case pr.Karma <= -25:
		pr.KarmaLevel = KarmaBad
	case pr.Karma >= 100:
		pr.KarmaLevel = KarmaHero
	case pr.Karma >= 25:
		pr.KarmaLevel = KarmaGood
	default:
		pr.KarmaLevel = KarmaNeutral
	}
}

// RatePlayer rates another player
func (pr *PlayerReputation) RatePlayer(targetID string, rating int) {
	if rating < 1 {
		rating = 1
	}
	if rating > 5 {
		rating = 5
	}
	
	pr.PlayerRatings[targetID] = rating
}

// GetAverageRating returns average rating from other players
func (pr *PlayerReputation) GetAverageRating() float64 {
	if len(pr.PlayerRatings) == 0 {
		return 0
	}
	
	total := 0
	for _, rating := range pr.PlayerRatings {
		total += rating
	}
	
	return float64(total) / float64(len(pr.PlayerRatings))
}

// CalculateReputation calculates overall reputation score (0-100)
func (pr *PlayerReputation) CalculateReputation() {
	// Base from karma (-100 to 100 karma -> 0 to 100 score)
	karmaScore := (pr.Karma + 100) / 2
	if karmaScore < 0 {
		karmaScore = 0
	}
	if karmaScore > 100 {
		karmaScore = 100
	}
	
	// Rating factor (0-5 stars -> 0-25 points)
	ratingScore := pr.GetAverageRating() * 5
	
	// Combined score
	pr.Reputation = (karmaScore*0.7) + (ratingScore*0.3)
	
	if pr.Reputation < 0 {
		pr.Reputation = 0
	}
	if pr.Reputation > 100 {
		pr.Reputation = 100
	}
}

// GetTitle returns reputation title
func (pr *PlayerReputation) GetTitle() string {
	switch pr.KarmaLevel {
	case KarmaEvil:
		return "Villain"
	case KarmaBad:
		return "Outlaw"
	case KarmaNeutral:
		if pr.Reputation > 60 {
			return "Citizen"
		}
		return "Neutral"
	case KarmaGood:
		return "Helper"
	case KarmaHero:
		return "Hero"
	}
	return "Unknown"
}

// ReputationManager manages all player reputations
type ReputationManager struct {
	reputations map[string]*PlayerReputation
	
	storagePath string
}

// NewReputationManager creates a new manager
func NewReputationManager(storageDir string) *ReputationManager {
	return &ReputationManager{
		reputations: make(map[string]*PlayerReputation),
		storagePath: filepath.Join(storageDir, "reputation.json"),
	}
}

// GetOrCreateReputation gets or creates reputation
func (rm *ReputationManager) GetOrCreateReputation(playerID string) *PlayerReputation {
	if rep, exists := rm.reputations[playerID]; exists {
		return rep
	}
	
	rep := NewPlayerReputation(playerID)
	rm.reputations[playerID] = rep
	return rep
}

// GetReputation gets reputation
func (rm *ReputationManager) GetReputation(playerID string) (*PlayerReputation, bool) {
	rep, exists := rm.reputations[playerID]
	return rep, exists
}

// RecordAction records an action for a player
func (rm *ReputationManager) RecordAction(playerID string, action ReputationAction, targetID string) {
	rep := rm.GetOrCreateReputation(playerID)
	rep.RecordAction(action, targetID)
}

// GetLeaderboard returns top players by karma
func (rm *ReputationManager) GetLeaderboard(count int) []*PlayerReputation {
	all := make([]*PlayerReputation, 0, len(rm.reputations))
	for _, rep := range rm.reputations {
		all = append(all, rep)
	}
	
	// Sort by karma
	for i := 0; i < len(all); i++ {
		for j := i + 1; j < len(all); j++ {
			if all[i].Karma < all[j].Karma {
				all[i], all[j] = all[j], all[i]
			}
		}
	}
	
	if count > len(all) {
		count = len(all)
	}
	
	return all[:count]
}

// GetVillains returns most evil players
func (rm *ReputationManager) GetVillains(count int) []*PlayerReputation {
	all := make([]*PlayerReputation, 0, len(rm.reputations))
	for _, rep := range rm.reputations {
		if rep.Karma < 0 {
			all = append(all, rep)
		}
	}
	
	// Sort by karma (lowest first)
	for i := 0; i < len(all); i++ {
		for j := i + 1; j < len(all); j++ {
			if all[i].Karma > all[j].Karma {
				all[i], all[j] = all[j], all[i]
			}
		}
	}
	
	if count > len(all) {
		count = len(all)
	}
	
	return all[:count]
}

// GetReputationBonus returns bonus multiplier based on reputation
func (rm *ReputationManager) GetReputationBonus(playerID string) float64 {
	rep, exists := rm.GetReputation(playerID)
	if !exists {
		return 1.0
	}
	
	// Good karma = bonus, bad karma = penalty
	return 1.0 + (rep.Karma / 1000.0)
}

// CanUseFeature checks if player can use a feature based on karma
func (rm *ReputationManager) CanUseFeature(playerID string, minKarma float64) bool {
	rep, exists := rm.GetReputation(playerID)
	if !exists {
		return minKarma <= 0
	}
	
	return rep.Karma >= minKarma
}

// Save saves reputation data
func (rm *ReputationManager) Save() error {
	data, err := json.MarshalIndent(rm.reputations, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}
	
	if err := os.WriteFile(rm.storagePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	
	return nil
}

// Load loads reputation data
func (rm *ReputationManager) Load() error {
	data, err := os.ReadFile(rm.storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read: %w", err)
	}
	
	var loaded map[string]*PlayerReputation
	if err := json.Unmarshal(data, &loaded); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}
	
	rm.reputations = loaded
	if rm.reputations == nil {
		rm.reputations = make(map[string]*PlayerReputation)
	}
	
	return nil
}

// DecayKarma decays karma over time (slowly returns toward 0)
func (rm *ReputationManager) DecayKarma() {
	for _, rep := range rm.reputations {
		if rep.Karma > 0 {
			rep.Karma = math.Max(0, rep.Karma-0.1)
		} else if rep.Karma < 0 {
			rep.Karma = math.Min(0, rep.Karma+0.1)
		}
		rep.updateKarmaLevel()
	}
}

// GetGlobalStats returns global reputation stats
func (rm *ReputationManager) GetGlobalStats() (avgKarma, avgReputation float64, heroes, villains int) {
	if len(rm.reputations) == 0 {
		return 0, 0, 0, 0
	}
	
	totalKarma := 0.0
	totalRep := 0.0
	
	for _, rep := range rm.reputations {
		totalKarma += rep.Karma
		totalRep += rep.Reputation
		
		if rep.KarmaLevel == KarmaHero {
			heroes++
		} else if rep.KarmaLevel == KarmaEvil {
			villains++
		}
	}
	
	return totalKarma / float64(len(rm.reputations)),
		totalRep / float64(len(rm.reputations)),
		heroes, villains
}
