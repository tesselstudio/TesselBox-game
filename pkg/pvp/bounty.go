package pvp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// BountyStatus represents the status of a bounty
type BountyStatus int

const (
	BountyActive BountyStatus = iota
	BountyClaimed
	BountyExpired
	BountyCancelled
)

// BountyClaim represents a claim on a bounty
type BountyClaim struct {
	HunterID   string    `json:"hunter_id"`
	ClaimedAt  time.Time `json:"claimed_at"`
	Proof      []string  `json:"proof,omitempty"` // Screenshot paths, etc
	Verified   bool      `json:"verified"`
	Paid       bool      `json:"paid"`
}

// Bounty represents a player bounty
type Bounty struct {
	ID          string         `json:"id"`
	TargetID    string         `json:"target_id"`
	TargetName  string         `json:"target_name"`
	IssuerID    string         `json:"issuer_id"`
	IssuerName  string         `json:"issuer_name"`
	Amount      float64        `json:"amount"`
	Reason      string         `json:"reason"`
	
	// Status
	Status      BountyStatus   `json:"status"`
	
	// Claims
	Claims      []BountyClaim  `json:"claims,omitempty"`
	WinningClaim *BountyClaim  `json:"winning_claim,omitempty"`
	
	// Settings
	Anonymous   bool           `json:"anonymous"`
	ExpiresAt   time.Time      `json:"expires_at"`
	
	// Meta
	CreatedAt   time.Time      `json:"created_at"`
	Views       int            `json:"views"`
}

// NewBounty creates a new bounty
func NewBounty(id, targetID, targetName, issuerID, issuerName string, amount float64, reason string, anonymous bool, duration time.Duration) *Bounty {
	now := time.Now()
	return &Bounty{
		ID:         id,
		TargetID:   targetID,
		TargetName: targetName,
		IssuerID:   issuerID,
		IssuerName: issuerName,
		Amount:     amount,
		Reason:     reason,
		Status:     BountyActive,
		Claims:     make([]BountyClaim, 0),
		Anonymous:  anonymous,
		ExpiresAt:  now.Add(duration),
		CreatedAt:  now,
		Views:      0,
	}
}

// IsExpired checks if bounty has expired
func (b *Bounty) IsExpired() bool {
	return time.Now().After(b.ExpiresAt)
}

// RecordView increments view count
func (b *Bounty) RecordView() {
	b.Views++
}

// SubmitClaim submits a claim on this bounty
func (b *Bounty) SubmitClaim(hunterID string, proof []string) error {
	if b.Status != BountyActive {
		return fmt.Errorf("bounty is not active")
	}
	
	// Check if already claimed by this hunter
	for _, claim := range b.Claims {
		if claim.HunterID == hunterID {
			return fmt.Errorf("already submitted a claim")
		}
	}
	
	claim := BountyClaim{
		HunterID:  hunterID,
		ClaimedAt: time.Now(),
		Proof:     proof,
		Verified:  false,
		Paid:      false,
	}
	
	b.Claims = append(b.Claims, claim)
	return nil
}

// ApproveClaim approves a claim and pays the hunter
func (b *Bounty) ApproveClaim(hunterID string) error {
	if b.Status != BountyActive {
		return fmt.Errorf("bounty is not active")
	}
	
	// Find the claim
	for i := range b.Claims {
		if b.Claims[i].HunterID == hunterID {
			b.Claims[i].Verified = true
			b.Claims[i].Paid = true
			b.WinningClaim = &b.Claims[i]
			b.Status = BountyClaimed
			return nil
		}
	}
	
	return fmt.Errorf("claim not found")
}

// GetIssuersDisplayName returns the name to display for issuer
func (b *Bounty) GetIssuersDisplayName() string {
	if b.Anonymous {
		return "Anonymous"
	}
	return b.IssuerName
}

// BountyBoard manages all bounties
type BountyBoard struct {
	bounties      map[string]*Bounty
	byTarget      map[string][]string // TargetID -> Bounty IDs
	byIssuer      map[string][]string // IssuerID -> Bounty IDs
	
	bountyCounter int
	
	// Cooldown tracking
	killCooldowns map[string]time.Time // TargetID -> Last kill time
	cooldownDuration time.Duration
	
	storagePath   string
}

// NewBountyBoard creates a new bounty board
func NewBountyBoard(storageDir string) *BountyBoard {
	return &BountyBoard{
		bounties:         make(map[string]*Bounty),
		byTarget:         make(map[string][]string),
		byIssuer:         make(map[string][]string),
		bountyCounter:    0,
		killCooldowns:    make(map[string]time.Time),
		cooldownDuration: 1 * time.Hour, // 1 hour cooldown
		storagePath:      filepath.Join(storageDir, "bounties.json"),
	}
}

// CreateBounty creates a new bounty
func (bb *BountyBoard) CreateBounty(targetID, targetName, issuerID, issuerName string, amount float64, reason string, anonymous bool, duration time.Duration) (*Bounty, error) {
	// Check if already has active bounty from this issuer on this target
	issuerBounties := bb.GetBountiesByIssuer(issuerID)
	for _, b := range issuerBounties {
		if b.TargetID == targetID && b.Status == BountyActive {
			return nil, fmt.Errorf("already have an active bounty on this target")
		}
	}
	
	bb.bountyCounter++
	bountyID := fmt.Sprintf("bounty_%d_%d", bb.bountyCounter, time.Now().Unix())
	
	if duration == 0 {
		duration = 7 * 24 * time.Hour // Default 7 days
	}
	
	bounty := NewBounty(bountyID, targetID, targetName, issuerID, issuerName, amount, reason, anonymous, duration)
	
	bb.bounties[bountyID] = bounty
	bb.byTarget[targetID] = append(bb.byTarget[targetID], bountyID)
	bb.byIssuer[issuerID] = append(bb.byIssuer[issuerID], bountyID)
	
	return bounty, nil
}

// GetBounty gets a bounty by ID
func (bb *BountyBoard) GetBounty(bountyID string) (*Bounty, bool) {
	bounty, exists := bb.bounties[bountyID]
	return bounty, exists
}

// GetBountiesByTarget gets bounties on a target
func (bb *BountyBoard) GetBountiesByTarget(targetID string) []*Bounty {
	bountyIDs := bb.byTarget[targetID]
	bounties := make([]*Bounty, 0, len(bountyIDs))
	
	for _, id := range bountyIDs {
		if bounty, exists := bb.bounties[id]; exists {
			bounties = append(bounties, bounty)
		}
	}
	
	return bounties
}

// GetActiveBountiesOnTarget gets active bounties on a target
func (bb *BountyBoard) GetActiveBountiesOnTarget(targetID string) []*Bounty {
	all := bb.GetBountiesByTarget(targetID)
	active := make([]*Bounty, 0)
	
	for _, b := range all {
		if b.Status == BountyActive && !b.IsExpired() {
			active = append(active, b)
		}
	}
	
	return active
}

// GetBountiesByIssuer gets bounties issued by a player
func (bb *BountyBoard) GetBountiesByIssuer(issuerID string) []*Bounty {
	bountyIDs := bb.byIssuer[issuerID]
	bounties := make([]*Bounty, 0, len(bountyIDs))
	
	for _, id := range bountyIDs {
		if bounty, exists := bb.bounties[id]; exists {
			bounties = append(bounties, bounty)
		}
	}
	
	return bounties
}

// GetActiveBounties gets all active bounties
func (bb *BountyBoard) GetActiveBounties() []*Bounty {
	result := make([]*Bounty, 0)
	for _, bounty := range bb.bounties {
		if bounty.Status == BountyActive && !bounty.IsExpired() {
			result = append(result, bounty)
		}
	}
	return result
}

// GetTopBounties returns bounties sorted by amount
func (bb *BountyBoard) GetTopBounties(count int) []*Bounty {
	bounties := bb.GetActiveBounties()
	
	// Sort by amount (bubble sort)
	for i := 0; i < len(bounties); i++ {
		for j := i + 1; j < len(bounties); j++ {
			if bounties[i].Amount < bounties[j].Amount {
				bounties[i], bounties[j] = bounties[j], bounties[i]
			}
		}
	}
	
	if count > len(bounties) {
		count = len(bounties)
	}
	
	return bounties[:count]
}

// ClaimBounty processes a bounty claim
func (bb *BountyBoard) ClaimBounty(bountyID, hunterID string, proof []string) error {
	bounty, exists := bb.GetBounty(bountyID)
	if !exists {
		return fmt.Errorf("bounty not found")
	}
	
	if bounty.Status != BountyActive {
		return fmt.Errorf("bounty is not active")
	}
	
	// Check cooldown
	if lastKill, exists := bb.killCooldowns[bounty.TargetID]; exists {
		if time.Since(lastKill) < bb.cooldownDuration {
			return fmt.Errorf("target was recently killed, please wait")
		}
	}
	
	// Check if hunter is target (can't claim own bounty)
	if hunterID == bounty.TargetID {
		return fmt.Errorf("cannot claim bounty on yourself")
	}
	
	// Check if hunter is issuer (can't claim own bounty)
	if hunterID == bounty.IssuerID {
		return fmt.Errorf("cannot claim your own bounty")
	}
	
	return bounty.SubmitClaim(hunterID, proof)
}

// ApproveBountyClaim approves a claim and pays the hunter
func (bb *BountyBoard) ApproveBountyClaim(bountyID, hunterID string) error {
	bounty, exists := bb.GetBounty(bountyID)
	if !exists {
		return fmt.Errorf("bounty not found")
	}
	
	if err := bounty.ApproveClaim(hunterID); err != nil {
		return err
	}
	
	// Record cooldown
	bb.killCooldowns[bounty.TargetID] = time.Now()
	
	return nil
}

// CancelBounty cancels a bounty (issuer only)
func (bb *BountyBoard) CancelBounty(bountyID, cancellerID string) error {
	bounty, exists := bb.GetBounty(bountyID)
	if !exists {
		return fmt.Errorf("bounty not found")
	}
	
	if bounty.IssuerID != cancellerID {
		return fmt.Errorf("only issuer can cancel")
	}
	
	if bounty.Status != BountyActive {
		return fmt.Errorf("bounty is not active")
	}
	
	bounty.Status = BountyCancelled
	return nil
}

// GetTotalBountyOnTarget returns total active bounty amount on a target
func (bb *BountyBoard) GetTotalBountyOnTarget(targetID string) float64 {
	bounties := bb.GetActiveBountiesOnTarget(targetID)
	total := 0.0
	for _, b := range bounties {
		total += b.Amount
	}
	return total
}

// GetHunterRankings returns hunters by successful claims
func (bb *BountyBoard) GetHunterRankings(count int) []struct {
	HunterID     string
	Claims       int
	TotalEarned  float64
} {
	// Count claims by hunter
	hunterStats := make(map[string]struct {
		Claims      int
		TotalEarned float64
	})
	
	for _, bounty := range bb.bounties {
		if bounty.Status == BountyClaimed && bounty.WinningClaim != nil {
			stats := hunterStats[bounty.WinningClaim.HunterID]
			stats.Claims++
			stats.TotalEarned += bounty.Amount
			hunterStats[bounty.WinningClaim.HunterID] = stats
		}
	}
	
	// Convert to slice
	result := make([]struct {
		HunterID    string
		Claims      int
		TotalEarned float64
	}, 0, len(hunterStats))
	
	for id, stats := range hunterStats {
		result = append(result, struct {
			HunterID    string
			Claims      int
			TotalEarned float64
		}{id, stats.Claims, stats.TotalEarned})
	}
	
	// Sort by claims (bubble sort)
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Claims < result[j].Claims {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	
	if count > len(result) {
		count = len(result)
	}
	
	return result[:count]
}

// Update processes all bounties (cleanup expired)
func (bb *BountyBoard) Update() {
	for _, bounty := range bb.bounties {
		if bounty.Status == BountyActive && bounty.IsExpired() {
			bounty.Status = BountyExpired
		}
	}
}

// Save saves bounties to disk
func (bb *BountyBoard) Save() error {
	data, err := json.MarshalIndent(bb.bounties, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}
	
	if err := os.WriteFile(bb.storagePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	
	return nil
}

// Load loads bounties from disk
func (bb *BountyBoard) Load() error {
	data, err := os.ReadFile(bb.storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read: %w", err)
	}
	
	var loaded map[string]*Bounty
	if err := json.Unmarshal(data, &loaded); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}
	
	bb.bounties = loaded
	if bb.bounties == nil {
		bb.bounties = make(map[string]*Bounty)
	}
	
	// Rebuild indexes
	bb.byTarget = make(map[string][]string)
	bb.byIssuer = make(map[string][]string)
	
	for _, bounty := range bb.bounties {
		bb.byTarget[bounty.TargetID] = append(bb.byTarget[bounty.TargetID], bounty.ID)
		bb.byIssuer[bounty.IssuerID] = append(bb.byIssuer[bounty.IssuerID], bounty.ID)
	}
	
	return nil
}

// GetStats returns bounty statistics
func (bb *BountyBoard) GetStats() (active, claimed, expired, cancelled int, totalValue float64) {
	for _, bounty := range bb.bounties {
		switch bounty.Status {
		case BountyActive:
			active++
			totalValue += bounty.Amount
		case BountyClaimed:
			claimed++
		case BountyExpired:
			expired++
		case BountyCancelled:
			cancelled++
		}
	}
	return
}
