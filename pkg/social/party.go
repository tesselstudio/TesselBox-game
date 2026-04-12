package social

import (
	"fmt"
	"time"
)

// PartyLootMode represents how party loot is distributed
type PartyLootMode int

const (
	LootFreeForAll PartyLootMode = iota // Everyone can pick up
	LootRoundRobin                      // Rotate between members
	LootLeaderOnly                      // Only leader can pick up
	LootNeedBeforeGreed                 // Roll for items
)

// String returns loot mode name
func (p PartyLootMode) String() string {
	switch p {
	case LootFreeForAll:
		return "freeforall"
	case LootRoundRobin:
		return "roundrobin"
	case LootLeaderOnly:
		return "leaderonly"
	case LootNeedBeforeGreed:
		return "needbeforegreed"
	}
	return "unknown"
}

// PartyMember represents a party member
type PartyMember struct {
	PlayerID   string    `json:"player_id"`
	PlayerName string    `json:"player_name"`
	Role       string    `json:"role"` // "leader", "member"
	JoinedAt   time.Time `json:"joined_at"`
	LastActive time.Time `json:"last_active"`
	Online     bool      `json:"online"`
}

// IsLeader returns true if member is party leader
func (pm *PartyMember) IsLeader() bool {
	return pm.Role == "leader"
}

// Party represents a player party
type Party struct {
	ID           string        `json:"id"`
	LeaderID     string        `json:"leader_id"`
	Members      []PartyMember `json:"members"`
	MaxSize      int           `json:"max_size"`
	
	// Settings
	LootMode     PartyLootMode `json:"loot_mode"`
	SharedXP     bool          `json:"shared_xp"`
	FriendlyFire bool          `json:"friendly_fire"` // Can hurt party members?
	
	// State
	WorldID      string        `json:"world_id"`
	CreatedAt    time.Time     `json:"created_at"`
	LastActive   time.Time     `json:"last_active"`
	Disbanded    bool          `json:"disbanded"`
	DisbandedAt  *time.Time    `json:"disbanded_at,omitempty"`
	
	// Loot tracking
	LootIndex    int           `json:"loot_index"` // For round-robin
}

// NewParty creates a new party
func NewParty(id, leaderID, leaderName, worldID string) *Party {
	now := time.Now()
	return &Party{
		ID:           id,
		LeaderID:     leaderID,
		Members:      make([]PartyMember, 0),
		MaxSize:      8,
		LootMode:     LootFreeForAll,
		SharedXP:     true,
		FriendlyFire: false,
		WorldID:      worldID,
		CreatedAt:    now,
		LastActive:   now,
		Disbanded:    false,
		LootIndex:    0,
	}
}

// AddMember adds a member to the party
func (p *Party) AddMember(playerID, playerName string) error {
	if p.Disbanded {
		return fmt.Errorf("party has been disbanded")
	}
	
	if len(p.Members) >= p.MaxSize {
		return fmt.Errorf("party is full")
	}
	
	// Check if already in party
	for _, m := range p.Members {
		if m.PlayerID == playerID {
			return fmt.Errorf("already in party")
		}
	}
	
	member := PartyMember{
		PlayerID:   playerID,
		PlayerName: playerName,
		Role:       "member",
		JoinedAt:   time.Now(),
		LastActive: time.Now(),
		Online:     true,
	}
	
	p.Members = append(p.Members, member)
	p.LastActive = time.Now()
	
	return nil
}

// RemoveMember removes a member from the party
func (p *Party) RemoveMember(playerID string) error {
	if p.Disbanded {
		return fmt.Errorf("party has been disbanded")
	}
	
	for i, m := range p.Members {
		if m.PlayerID == playerID {
			p.Members = append(p.Members[:i], p.Members[i+1:]...)
			
			// If leader left, promote next member
			if m.IsLeader() && len(p.Members) > 0 {
				p.Members[0].Role = "leader"
				p.LeaderID = p.Members[0].PlayerID
			}
			
			// If empty, disband
			if len(p.Members) == 0 {
				p.Disband()
			}
			
			p.LastActive = time.Now()
			return nil
		}
	}
	
	return fmt.Errorf("member not found")
}

// PromoteMember promotes a member to leader
func (p *Party) PromoteMember(playerID string) error {
	if p.Disbanded {
		return fmt.Errorf("party has been disbanded")
	}
	
	// Find the member
	found := false
	for i := range p.Members {
		if p.Members[i].PlayerID == playerID {
			p.Members[i].Role = "leader"
			found = true
		} else {
			p.Members[i].Role = "member"
		}
	}
	
	if !found {
		return fmt.Errorf("member not found")
	}
	
	p.LeaderID = playerID
	p.LastActive = time.Now()
	
	return nil
}

// IsMember checks if player is in party
func (p *Party) IsMember(playerID string) bool {
	for _, m := range p.Members {
		if m.PlayerID == playerID {
			return true
		}
	}
	return false
}

// IsLeader checks if player is party leader
func (p *Party) IsLeader(playerID string) bool {
	return p.LeaderID == playerID
}

// GetMember gets a party member
func (p *Party) GetMember(playerID string) *PartyMember {
	for i := range p.Members {
		if p.Members[i].PlayerID == playerID {
			return &p.Members[i]
		}
	}
	return nil
}

// UpdateMemberActivity updates member activity
func (p *Party) UpdateMemberActivity(playerID string) {
	for i := range p.Members {
		if p.Members[i].PlayerID == playerID {
			p.Members[i].LastActive = time.Now()
			p.Members[i].Online = true
			p.LastActive = time.Now()
			return
		}
	}
}

// SetMemberOffline marks a member as offline
func (p *Party) SetMemberOffline(playerID string) {
	for i := range p.Members {
		if p.Members[i].PlayerID == playerID {
			p.Members[i].Online = false
			return
		}
	}
}

// Disband disbands the party
func (p *Party) Disband() {
	if p.Disbanded {
		return
	}
	
	now := time.Now()
	p.Disbanded = true
	p.DisbandedAt = &now
	p.Members = make([]PartyMember, 0)
}

// GetOnlineCount returns number of online members
func (p *Party) GetOnlineCount() int {
	count := 0
	for _, m := range p.Members {
		if m.Online {
			count++
		}
	}
	return count
}

// GetMemberIDs returns all member IDs
func (p *Party) GetMemberIDs() []string {
	ids := make([]string, 0, len(p.Members))
	for _, m := range p.Members {
		ids = append(ids, m.PlayerID)
	}
	return ids
}

// GetOnlineMemberIDs returns online member IDs
func (p *Party) GetOnlineMemberIDs() []string {
	ids := make([]string, 0)
	for _, m := range p.Members {
		if m.Online {
			ids = append(ids, m.PlayerID)
		}
	}
	return ids
}

// NextLootRecipient gets the next recipient for round-robin loot
func (p *Party) NextLootRecipient() string {
	if len(p.Members) == 0 {
		return ""
	}
	
	// Find next online member
	for i := 0; i < len(p.Members); i++ {
		idx := (p.LootIndex + i) % len(p.Members)
		if p.Members[idx].Online {
			p.LootIndex = (idx + 1) % len(p.Members)
			return p.Members[idx].PlayerID
		}
	}
	
	return ""
}

// PartyManager manages parties
type PartyManager struct {
	parties    map[string]*Party
	byPlayer   map[string]string // PlayerID -> PartyID
	
	partyCounter int
}

// NewPartyManager creates a new party manager
func NewPartyManager() *PartyManager {
	return &PartyManager{
		parties:      make(map[string]*Party),
		byPlayer:     make(map[string]string),
		partyCounter: 0,
	}
}

// CreateParty creates a new party
func (pm *PartyManager) CreateParty(leaderID, leaderName, worldID string) *Party {
	pm.partyCounter++
	partyID := fmt.Sprintf("party_%d_%d", pm.partyCounter, time.Now().Unix())
	
	party := NewParty(partyID, leaderID, leaderName, worldID)
	party.AddMember(leaderID, leaderName) // Leader is first member
	party.GetMember(leaderID).Role = "leader"
	
	pm.parties[partyID] = party
	pm.byPlayer[leaderID] = partyID
	
	return party
}

// GetParty gets a party by ID
func (pm *PartyManager) GetParty(partyID string) (*Party, bool) {
	party, exists := pm.parties[partyID]
	return party, exists
}

// GetPlayerParty gets the party a player is in
func (pm *PartyManager) GetPlayerParty(playerID string) (*Party, bool) {
	partyID, exists := pm.byPlayer[playerID]
	if !exists {
		return nil, false
	}
	return pm.GetParty(partyID)
}

// IsInParty checks if player is in a party
func (pm *PartyManager) IsInParty(playerID string) bool {
	_, exists := pm.byPlayer[playerID]
	return exists
}

// InviteToParty invites a player to a party
func (pm *PartyManager) InviteToParty(partyID, inviterID, inviteeID string) error {
	party, exists := pm.GetParty(partyID)
	if !exists {
		return fmt.Errorf("party not found")
	}
	
	if party.Disbanded {
		return fmt.Errorf("party has been disbanded")
	}
	
	if !party.IsMember(inviterID) {
		return fmt.Errorf("not in party")
	}
	
	if pm.IsInParty(inviteeID) {
		return fmt.Errorf("player is already in a party")
	}
	
	// Note: Actual invitation system would be implemented separately
	// This just validates the party can accept the invitee
	
	return nil
}

// JoinParty adds a player to a party
func (pm *PartyManager) JoinParty(partyID, playerID, playerName string) error {
	party, exists := pm.GetParty(partyID)
	if !exists {
		return fmt.Errorf("party not found")
	}
	
	if party.Disbanded {
		return fmt.Errorf("party has been disbanded")
	}
	
	if pm.IsInParty(playerID) {
		return fmt.Errorf("already in a party")
	}
	
	if err := party.AddMember(playerID, playerName); err != nil {
		return err
	}
	
	pm.byPlayer[playerID] = partyID
	
	return nil
}

// LeaveParty removes a player from their party
func (pm *PartyManager) LeaveParty(playerID string) error {
	party, exists := pm.GetPlayerParty(playerID)
	if !exists {
		return fmt.Errorf("not in a party")
	}
	
	if err := party.RemoveMember(playerID); err != nil {
		return err
	}
	
	delete(pm.byPlayer, playerID)
	
	// Clean up disbanded parties
	if party.Disbanded {
		delete(pm.parties, party.ID)
	}
	
	return nil
}

// KickFromParty kicks a player from a party
func (pm *PartyManager) KickFromParty(partyID, kickerID, targetID string) error {
	party, exists := pm.GetParty(partyID)
	if !exists {
		return fmt.Errorf("party not found")
	}
	
	if !party.IsLeader(kickerID) {
		return fmt.Errorf("only leader can kick")
	}
	
	if !party.IsMember(targetID) {
		return fmt.Errorf("player is not in party")
	}
	
	if err := party.RemoveMember(targetID); err != nil {
		return err
	}
	
	delete(pm.byPlayer, targetID)
	
	return nil
}

// PromoteInParty promotes a member to leader
func (pm *PartyManager) PromoteInParty(partyID, promoterID, targetID string) error {
	party, exists := pm.GetParty(partyID)
	if !exists {
		return fmt.Errorf("party not found")
	}
	
	if !party.IsLeader(promoterID) {
		return fmt.Errorf("only leader can promote")
	}
	
	if !party.IsMember(targetID) {
		return fmt.Errorf("player is not in party")
	}
	
	return party.PromoteMember(targetID)
}

// DisbandParty disbands a party
func (pm *PartyManager) DisbandParty(partyID, disbanderID string) error {
	party, exists := pm.GetParty(partyID)
	if !exists {
		return fmt.Errorf("party not found")
	}
	
	if !party.IsLeader(disbanderID) {
		return fmt.Errorf("only leader can disband")
	}
	
	// Remove all members
	for _, m := range party.Members {
		delete(pm.byPlayer, m.PlayerID)
	}
	
	party.Disband()
	delete(pm.parties, partyID)
	
	return nil
}

// GetActiveParties returns all active (non-disbanded) parties
func (pm *PartyManager) GetActiveParties() []*Party {
	result := make([]*Party, 0)
	for _, party := range pm.parties {
		if !party.Disbanded {
			result = append(result, party)
		}
	}
	return result
}

// CleanupInactive removes parties that have been inactive too long
func (pm *PartyManager) CleanupInactive(maxInactive time.Duration) {
	now := time.Now()
	
	for id, party := range pm.parties {
		if party.Disbanded {
			continue
		}
		
		// Check if all members offline for too long
		allOffline := true
		for _, m := range party.Members {
			if m.Online || now.Sub(m.LastActive) < maxInactive {
				allOffline = false
				break
			}
		}
		
		if allOffline {
			party.Disband()
			for _, m := range party.Members {
				delete(pm.byPlayer, m.PlayerID)
			}
			delete(pm.parties, id)
		}
	}
}
