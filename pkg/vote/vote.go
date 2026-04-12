package vote

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"tesselbox/pkg/economy"
)

// VoteSite represents a voting website
type VoteSite int

const (
	SiteA VoteSite = iota
	SiteB
	SiteC
)

// String returns site name
func (v VoteSite) String() string {
	switch v {
	case SiteA:
		return "SiteA"
	case SiteB:
		return "SiteB"
	case SiteC:
		return "SiteC"
	}
	return "Unknown"
}

// PlayerVote represents a player's vote
type PlayerVote struct {
	PlayerID    string    `json:"player_id"`
	Site        VoteSite  `json:"site"`
	VotedAt     time.Time `json:"voted_at"`
	Claimed     bool      `json:"claimed"`
	ClaimedAt   *time.Time `json:"claimed_at,omitempty"`
	IP          string    `json:"ip,omitempty"`
}

// VoteReward represents rewards for voting
type VoteReward struct {
	Money      float64 `json:"money"`
	XP         int     `json:"xp"`
	Items      []string `json:"items,omitempty"`
	Keys       int      `json:"keys"` // Vote crate keys
}

// VoteStreak tracks consecutive votes
type VoteStreak struct {
	PlayerID    string    `json:"player_id"`
	CurrentStreak int     `json:"current_streak"`
	BestStreak    int     `json:"best_streak"`
	LastVoteAt    time.Time `json:"last_vote_at"`
}

// VoteManager manages voting
type VoteManager struct {
	votes       map[string]PlayerVote // key: "playerID_site"
	streaks     map[string]*VoteStreak
	
	// Configuration
	rewards     VoteReward
	streakBonus map[int]float64 // Streak length -> multiplier
	voteCooldown time.Duration
	
	walletMgr   *economy.WalletManager
	
	storagePath string
}

// NewVoteManager creates new manager
func NewVoteManager(walletMgr *economy.WalletManager, storageDir string) *VoteManager {
	return &VoteManager{
		votes:        make(map[string]PlayerVote),
		streaks:      make(map[string]*VoteStreak),
		rewards:      VoteReward{Money: 100, XP: 50, Keys: 1},
		streakBonus:  map[int]float64{2: 1.5, 5: 2.0, 10: 3.0, 30: 5.0},
		voteCooldown: 24 * time.Hour,
		walletMgr:    walletMgr,
		storagePath:  filepath.Join(storageDir, "votes.json"),
	}
}

// CanVote checks if player can vote on a site
func (vm *VoteManager) CanVote(playerID string, site VoteSite) bool {
	key := fmt.Sprintf("%s_%d", playerID, site)
	
	if vote, exists := vm.votes[key]; exists {
		return time.Since(vote.VotedAt) >= vm.voteCooldown
	}
	
	return true
}

// RecordVote records a vote
func (vm *VoteManager) RecordVote(playerID string, site VoteSite, ip string) error {
	if !vm.CanVote(playerID, site) {
		return fmt.Errorf("already voted today")
	}
	
	key := fmt.Sprintf("%s_%d", playerID, site)
	
	vote := PlayerVote{
		PlayerID: playerID,
		Site:     site,
		VotedAt:  time.Now(),
		Claimed:  false,
		IP:       ip,
	}
	
	vm.votes[key] = vote
	
	// Update streak
	vm.updateStreak(playerID)
	
	return nil
}

// updateStreak updates player's vote streak
func (vm *VoteManager) updateStreak(playerID string) {
	streak, exists := vm.streaks[playerID]
	if !exists {
		streak = &VoteStreak{
			PlayerID:      playerID,
			CurrentStreak: 0,
			BestStreak:    0,
		}
		vm.streaks[playerID] = streak
	}
	
	// Check if streak continues (voted within 48 hours)
	if time.Since(streak.LastVoteAt) < 48*time.Hour {
		streak.CurrentStreak++
	} else {
		streak.CurrentStreak = 1 // Reset but count this vote
	}
	
	if streak.CurrentStreak > streak.BestStreak {
		streak.BestStreak = streak.CurrentStreak
	}
	
	streak.LastVoteAt = time.Now()
}

// ClaimRewards claims vote rewards
func (vm *VoteManager) ClaimRewards(playerID string, site VoteSite) (float64, error) {
	key := fmt.Sprintf("%s_%d", playerID, site)
	
	vote, exists := vm.votes[key]
	if !exists {
		return 0, fmt.Errorf("no vote found")
	}
	
	if vote.Claimed {
		return 0, fmt.Errorf("already claimed")
	}
	
	// Calculate reward with streak bonus
	streak := vm.streaks[playerID]
	multiplier := 1.0
	
	for streakLen, bonus := range vm.streakBonus {
		if streak.CurrentStreak >= streakLen && bonus > multiplier {
			multiplier = bonus
		}
	}
	
	reward := vm.rewards.Money * multiplier
	
	// Give reward
	wallet := vm.walletMgr.GetOrCreateWallet(playerID)
	wallet.Add(reward, economy.TransactionEarn, "VOTE", fmt.Sprintf("Vote reward from %s", site.String()))
	
	// Mark claimed
	vote.Claimed = true
	now := time.Now()
	vote.ClaimedAt = &now
	vm.votes[key] = vote
	
	return reward, nil
}

// GetVoteStatus gets player's vote status for all sites
func (vm *VoteManager) GetVoteStatus(playerID string) map[VoteSite]struct {
	CanVote   bool
	Voted     bool
	Claimed   bool
	TimeLeft  time.Duration
} {
	result := make(map[VoteSite]struct {
		CanVote   bool
		Voted     bool
		Claimed   bool
		TimeLeft  time.Duration
	})
	
	for _, site := range []VoteSite{SiteA, SiteB, SiteC} {
		key := fmt.Sprintf("%s_%d", playerID, site)
		
		status := struct {
			CanVote   bool
			Voted     bool
			Claimed   bool
			TimeLeft  time.Duration
		}{}
		
		if vote, exists := vm.votes[key]; exists {
			status.Voted = true
			status.Claimed = vote.Claimed
			
			timeSince := time.Since(vote.VotedAt)
			if timeSince < vm.voteCooldown {
				status.TimeLeft = vm.voteCooldown - timeSince
			} else {
				status.CanVote = true
			}
		} else {
			status.CanVote = true
		}
		
		result[site] = status
	}
	
	return result
}

// GetStreak gets player's vote streak
func (vm *VoteManager) GetStreak(playerID string) *VoteStreak {
	if streak, exists := vm.streaks[playerID]; exists {
		return streak
	}
	return &VoteStreak{PlayerID: playerID}
}

// GetTopVoters returns top voters
func (vm *VoteManager) GetTopVoters(count int) []struct {
	PlayerID     string
	TotalVotes   int
	CurrentStreak int
} {
	// Count votes per player
	voteCounts := make(map[string]int)
	for _, vote := range vm.votes {
		voteCounts[vote.PlayerID]++
	}
	
	// Convert to slice
	type voter struct {
		PlayerID      string
		TotalVotes    int
		CurrentStreak int
	}
	
	voters := make([]voter, 0, len(voteCounts))
	for playerID, total := range voteCounts {
		streak := 0
		if s, exists := vm.streaks[playerID]; exists {
			streak = s.CurrentStreak
		}
		voters = append(voters, voter{playerID, total, streak})
	}
	
	// Sort by total votes
	for i := 0; i < len(voters); i++ {
		for j := i + 1; j < len(voters); j++ {
			if voters[i].TotalVotes < voters[j].TotalVotes {
				voters[i], voters[j] = voters[j], voters[i]
			}
		}
	}
	
	if count > len(voters) {
		count = len(voters)
	}
	
	result := make([]struct {
		PlayerID      string
		TotalVotes    int
		CurrentStreak int
	}, count)
	
	for i := 0; i < count; i++ {
		result[i] = struct {
			PlayerID      string
			TotalVotes    int
			CurrentStreak int
		}{voters[i].PlayerID, voters[i].TotalVotes, voters[i].CurrentStreak}
	}
	
	return result
}

// Save saves vote data
func (vm *VoteManager) Save() error {
	data := struct {
		Votes   map[string]PlayerVote `json:"votes"`
		Streaks map[string]*VoteStreak `json:"streaks"`
	}{
		Votes:   vm.votes,
		Streaks: vm.streaks,
	}
	
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}
	
	if err := os.WriteFile(vm.storagePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	
	return nil
}

// Load loads vote data
func (vm *VoteManager) Load() error {
	data, err := os.ReadFile(vm.storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read: %w", err)
	}
	
	var loaded struct {
		Votes   map[string]PlayerVote `json:"votes"`
		Streaks map[string]*VoteStreak `json:"streaks"`
	}
	
	if err := json.Unmarshal(data, &loaded); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}
	
	vm.votes = loaded.Votes
	if vm.votes == nil {
		vm.votes = make(map[string]PlayerVote)
	}
	
	vm.streaks = loaded.Streaks
	if vm.streaks == nil {
		vm.streaks = make(map[string]*VoteStreak)
	}
	
	return nil
}

// GetTotalVotes returns total vote count
func (vm *VoteManager) GetTotalVotes() int {
	return len(vm.votes)
}

// GetTodayVotes returns votes in last 24 hours
func (vm *VoteManager) GetTodayVotes() int {
	count := 0
	cutoff := time.Now().Add(-24 * time.Hour)
	
	for _, vote := range vm.votes {
		if vote.VotedAt.After(cutoff) {
			count++
		}
	}
	
	return count
}
