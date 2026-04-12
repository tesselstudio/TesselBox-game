package pvp

import (
	"fmt"
	"time"

	"tesselbox/pkg/economy"
)

// DuelStatus represents duel state
type DuelStatus int

const (
	DuelInvited DuelStatus = iota
	DuelAccepted
	DuelInProgress
	DuelEnded
	DuelCancelled
)

// DuelType represents duel format
type DuelType int

const (
	DuelClassic DuelType = iota // Standard fight
	DuelBow                     // Bow only
	DuelSword                   // Sword only
	DuelFist                    // No weapons
	DuelPot                     // Potions allowed
)

// Duel represents a 1v1 duel
type Duel struct {
	ID           string `json:"id"`
	ChallengerID string `json:"challenger_id"`
	TargetID     string `json:"target_id"`

	// Settings
	Type    DuelType `json:"type"`
	Wager   float64  `json:"wager"`
	WorldID string   `json:"world_id"`

	// Status
	Status   DuelStatus `json:"status"`
	WinnerID *string    `json:"winner_id,omitempty"`

	// Scores
	ChallengerKills int `json:"challenger_kills"`
	TargetKills     int `json:"target_kills"`
	RoundsToWin     int `json:"rounds_to_win"`

	// Location
	ArenaX float64 `json:"arena_x"`
	ArenaY float64 `json:"arena_y"`

	// Timing
	CreatedAt time.Time  `json:"created_at"`
	StartedAt *time.Time `json:"started_at,omitempty"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`

	// Callbacks (not serialized)
	OnStart func()
	OnEnd   func(winnerID string)
}

// NewDuel creates a new duel
func NewDuel(id, challengerID, targetID, worldID string, duelType DuelType, wager float64) *Duel {
	return &Duel{
		ID:              id,
		ChallengerID:    challengerID,
		TargetID:        targetID,
		Type:            duelType,
		Wager:           wager,
		WorldID:         worldID,
		Status:          DuelInvited,
		ChallengerKills: 0,
		TargetKills:     0,
		RoundsToWin:     3, // Best of 5
		CreatedAt:       time.Now(),
	}
}

// Accept accepts the duel
func (d *Duel) Accept() error {
	if d.Status != DuelInvited {
		return fmt.Errorf("duel not in invited state")
	}

	d.Status = DuelAccepted
	return nil
}

// Start starts the duel
func (d *Duel) Start() error {
	if d.Status != DuelAccepted {
		return fmt.Errorf("duel not accepted")
	}

	now := time.Now()
	d.Status = DuelInProgress
	d.StartedAt = &now

	if d.OnStart != nil {
		d.OnStart()
	}

	return nil
}

// RecordKill records a kill in the duel
func (d *Duel) RecordKill(killerID string) {
	if d.Status != DuelInProgress {
		return
	}

	if killerID == d.ChallengerID {
		d.ChallengerKills++
	} else if killerID == d.TargetID {
		d.TargetKills++
	}

	// Check for winner
	if d.ChallengerKills >= d.RoundsToWin {
		d.End(&d.ChallengerID)
	} else if d.TargetKills >= d.RoundsToWin {
		d.End(&d.TargetID)
	}
}

// End ends the duel
func (d *Duel) End(winnerID *string) {
	if d.Status == DuelEnded {
		return
	}

	now := time.Now()
	d.Status = DuelEnded
	d.WinnerID = winnerID
	d.EndedAt = &now

	if d.OnEnd != nil && winnerID != nil {
		d.OnEnd(*winnerID)
	}
}

// Cancel cancels the duel
func (d *Duel) Cancel() {
	if d.Status == DuelEnded || d.Status == DuelInProgress {
		return
	}

	d.Status = DuelCancelled
}

// Forfeit forfeits the duel
func (d *Duel) Forfeit(playerID string) {
	if d.Status != DuelInProgress {
		return
	}

	var winnerID string
	if playerID == d.ChallengerID {
		winnerID = d.TargetID
	} else {
		winnerID = d.ChallengerID
	}

	d.End(&winnerID)
}

// GetOpponent gets the opponent of a player
func (d *Duel) GetOpponent(playerID string) string {
	if playerID == d.ChallengerID {
		return d.TargetID
	}
	return d.ChallengerID
}

// IsParticipant checks if player is in this duel
func (d *Duel) IsParticipant(playerID string) bool {
	return playerID == d.ChallengerID || playerID == d.TargetID
}

// GetScore returns score for a player
func (d *Duel) GetScore(playerID string) int {
	if playerID == d.ChallengerID {
		return d.ChallengerKills
	}
	return d.TargetKills
}

// DuelManager manages duels
type DuelManager struct {
	duels    map[string]*Duel
	byPlayer map[string]string // PlayerID -> DuelID

	duelCounter int
	walletMgr   *economy.WalletManager
}

// NewDuelManager creates a new duel manager
func NewDuelManager(walletMgr *economy.WalletManager) *DuelManager {
	return &DuelManager{
		duels:       make(map[string]*Duel),
		byPlayer:    make(map[string]string),
		duelCounter: 0,
		walletMgr:   walletMgr,
	}
}

// Challenge creates a new duel challenge
func (dm *DuelManager) Challenge(challengerID, targetID, worldID string, duelType DuelType, wager float64) (*Duel, error) {
	// Check if already dueling
	if dm.IsDueling(challengerID) {
		return nil, fmt.Errorf("already in a duel")
	}
	if dm.IsDueling(targetID) {
		return nil, fmt.Errorf("target is already in a duel")
	}

	// Check wager
	if wager > 0 {
		challengerWallet := dm.walletMgr.GetWallet(challengerID)
		if challengerWallet == nil || !challengerWallet.CanAfford(wager) {
			return nil, fmt.Errorf("cannot afford wager")
		}
	}

	dm.duelCounter++
	duelID := fmt.Sprintf("duel_%d_%d", dm.duelCounter, time.Now().Unix())

	duel := NewDuel(duelID, challengerID, targetID, worldID, duelType, wager)

	dm.duels[duelID] = duel
	dm.byPlayer[challengerID] = duelID
	dm.byPlayer[targetID] = duelID

	return duel, nil
}

// GetDuel gets a duel
func (dm *DuelManager) GetDuel(duelID string) (*Duel, bool) {
	duel, exists := dm.duels[duelID]
	return duel, exists
}

// GetPlayerDuel gets a player's current duel
func (dm *DuelManager) GetPlayerDuel(playerID string) (*Duel, bool) {
	duelID, exists := dm.byPlayer[playerID]
	if !exists {
		return nil, false
	}
	return dm.GetDuel(duelID)
}

// IsDueling checks if player is in a duel
func (dm *DuelManager) IsDueling(playerID string) bool {
	duel, exists := dm.GetPlayerDuel(playerID)
	if !exists {
		return false
	}
	return duel.Status == DuelInvited || duel.Status == DuelAccepted || duel.Status == DuelInProgress
}

// AcceptDuel accepts a duel
func (dm *DuelManager) AcceptDuel(duelID, playerID string) error {
	duel, exists := dm.GetDuel(duelID)
	if !exists {
		return fmt.Errorf("duel not found")
	}

	if duel.TargetID != playerID {
		return fmt.Errorf("not the challenged player")
	}

	// Check target can afford wager
	if duel.Wager > 0 {
		targetWallet := dm.walletMgr.GetWallet(playerID)
		if targetWallet == nil || !targetWallet.CanAfford(duel.Wager) {
			return fmt.Errorf("cannot afford wager")
		}
	}

	return duel.Accept()
}

// StartDuel starts a duel
func (dm *DuelManager) StartDuel(duelID string) error {
	duel, exists := dm.GetDuel(duelID)
	if !exists {
		return fmt.Errorf("duel not found")
	}

	// Take wagers
	if duel.Wager > 0 {
		challengerWallet := dm.walletMgr.GetWallet(duel.ChallengerID)
		targetWallet := dm.walletMgr.GetWallet(duel.TargetID)

		_, success1 := challengerWallet.Remove(duel.Wager, economy.TransactionSpend, "DUEL", "Duel wager")
		_, success2 := targetWallet.Remove(duel.Wager, economy.TransactionSpend, "DUEL", "Duel wager")

		if !success1 || !success2 {
			// Refund if failed (in real implementation)
			return fmt.Errorf("wager collection failed")
		}
	}

	return duel.Start()
}

// RecordDuelKill records a kill
func (dm *DuelManager) RecordDuelKill(duelID, killerID string) error {
	duel, exists := dm.GetDuel(duelID)
	if !exists {
		return fmt.Errorf("duel not found")
	}

	duel.RecordKill(killerID)

	// Check if duel ended
	if duel.Status == DuelEnded && duel.WinnerID != nil {
		dm.payWinner(duel)
	}

	return nil
}

// payWinner pays the winner
func (dm *DuelManager) payWinner(duel *Duel) {
	if duel.WinnerID == nil || duel.Wager == 0 {
		return
	}

	winnerWallet := dm.walletMgr.GetOrCreateWallet(*duel.WinnerID)
	winnerWallet.Add(duel.Wager*2, economy.TransactionEarn, "DUEL", "Duel winnings")
}

// ForfeitDuel forfeits a duel
func (dm *DuelManager) ForfeitDuel(duelID, playerID string) error {
	duel, exists := dm.GetDuel(duelID)
	if !exists {
		return fmt.Errorf("duel not found")
	}

	if !duel.IsParticipant(playerID) {
		return fmt.Errorf("not in this duel")
	}

	duel.Forfeit(playerID)

	// Pay winner
	if duel.WinnerID != nil {
		dm.payWinner(duel)
	}

	// Cleanup
	dm.cleanupDuel(duelID)

	return nil
}

// CancelDuel cancels a duel
func (dm *DuelManager) CancelDuel(duelID, playerID string) error {
	duel, exists := dm.GetDuel(duelID)
	if !exists {
		return fmt.Errorf("duel not found")
	}

	if duel.ChallengerID != playerID {
		return fmt.Errorf("only challenger can cancel")
	}

	if duel.Status == DuelInProgress {
		return fmt.Errorf("cannot cancel in-progress duel")
	}

	duel.Cancel()
	dm.cleanupDuel(duelID)

	return nil
}

// cleanupDuel removes a duel
func (dm *DuelManager) cleanupDuel(duelID string) {
	duel, exists := dm.duels[duelID]
	if !exists {
		return
	}

	delete(dm.byPlayer, duel.ChallengerID)
	delete(dm.byPlayer, duel.TargetID)
	delete(dm.duels, duelID)
}

// GetActiveDuels returns active duel count
func (dm *DuelManager) GetActiveDuels() []*Duel {
	result := make([]*Duel, 0)
	for _, duel := range dm.duels {
		if duel.Status == DuelInProgress {
			result = append(result, duel)
		}
	}
	return result
}
