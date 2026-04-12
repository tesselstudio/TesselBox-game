package social

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// FriendStatus represents the status of a friendship
type FriendStatus int

const (
	FriendPending FriendStatus = iota
	FriendAccepted
	FriendBlocked
)

// String returns status name
func (f FriendStatus) String() string {
	switch f {
	case FriendPending:
		return "pending"
	case FriendAccepted:
		return "friends"
	case FriendBlocked:
		return "blocked"
	}
	return "unknown"
}

// Friendship represents a friendship between two players
type Friendship struct {
	ID             string       `json:"id"`
	PlayerA        string       `json:"player_a"`      // Always sorted alphabetically
	PlayerB        string       `json:"player_b"`
	Status         FriendStatus `json:"status"`
	InitiatedBy    string       `json:"initiated_by"`  // Who sent the request
	
	// Stats
	FriendsSince   time.Time    `json:"friends_since"`
	LastPlayedTogether time.Time `json:"last_played_together"`
	MinutesTogether int         `json:"minutes_together"`
	
	// Metadata
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

// NewFriendship creates a new friendship
func NewFriendship(playerA, playerB, initiatedBy string) *Friendship {
	now := time.Now()
	
	// Ensure alphabetical ordering for consistency
	id := generateFriendshipID(playerA, playerB)
	
	return &Friendship{
		ID:          id,
		PlayerA:     playerA,
		PlayerB:     playerB,
		Status:      FriendPending,
		InitiatedBy: initiatedBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// generateFriendshipID generates a consistent ID for a friendship
func generateFriendshipID(playerA, playerB string) string {
	// Sort alphabetically
	if playerA > playerB {
		playerA, playerB = playerB, playerA
	}
	return fmt.Sprintf("friend_%s_%s", playerA, playerB)
}

// Accept accepts a pending friend request
func (f *Friendship) Accept() {
	if f.Status == FriendPending {
		f.Status = FriendAccepted
		f.FriendsSince = time.Now()
		f.UpdatedAt = time.Now()
	}
}

// Block blocks the friendship
func (f *Friendship) Block() {
	f.Status = FriendBlocked
	f.UpdatedAt = time.Now()
}

// RecordPlayTime records time spent playing together
func (f *Friendship) RecordPlayTime(minutes int) {
	f.MinutesTogether += minutes
	f.LastPlayedTogether = time.Now()
	f.UpdatedAt = time.Now()
}

// IsFriend returns true if players are friends (accepted)
func (f *Friendship) IsFriend() bool {
	return f.Status == FriendAccepted
}

// IsPending returns true if request is pending
func (f *Friendship) IsPending() bool {
	return f.Status == FriendPending
}

// IsBlocked returns true if friendship is blocked
func (f *Friendship) IsBlocked() bool {
	return f.Status == FriendBlocked
}

// FriendManager manages friendships
type FriendManager struct {
	friendships map[string]*Friendship  // By ID
	byPlayer    map[string][]string      // PlayerID -> friendship IDs
	
	storagePath string
}

// NewFriendManager creates a new friend manager
func NewFriendManager(storageDir string) *FriendManager {
	return &FriendManager{
		friendships: make(map[string]*Friendship),
		byPlayer:    make(map[string][]string),
		storagePath: filepath.Join(storageDir, "friendships.json"),
	}
}

// SendRequest sends a friend request
func (fm *FriendManager) SendRequest(fromID, toID string) (*Friendship, error) {
	if fromID == toID {
		return nil, fmt.Errorf("cannot friend yourself")
	}
	
	// Check if friendship already exists
	id := generateFriendshipID(fromID, toID)
	if _, exists := fm.friendships[id]; exists {
		return nil, fmt.Errorf("friendship already exists")
	}
	
	friendship := NewFriendship(fromID, toID, fromID)
	
	fm.friendships[id] = friendship
	fm.byPlayer[fromID] = append(fm.byPlayer[fromID], id)
	fm.byPlayer[toID] = append(fm.byPlayer[toID], id)
	
	return friendship, nil
}

// AcceptRequest accepts a friend request
func (fm *FriendManager) AcceptRequest(friendshipID, accepterID string) error {
	friendship, exists := fm.friendships[friendshipID]
	if !exists {
		return fmt.Errorf("friendship not found")
	}
	
	// Check if accepter is the target (not the initiator)
	if friendship.InitiatedBy == accepterID {
		return fmt.Errorf("cannot accept your own request")
	}
	
	if !friendship.IsPending() {
		return fmt.Errorf("request is not pending")
	}
	
	friendship.Accept()
	return nil
}

// DeclineRequest declines a friend request
func (fm *FriendManager) DeclineRequest(friendshipID, declinerID string) error {
	friendship, exists := fm.friendships[friendshipID]
	if !exists {
		return fmt.Errorf("friendship not found")
	}
	
	if !friendship.IsPending() {
		return fmt.Errorf("request is not pending")
	}
	
	// Remove the friendship
	fm.removeFriendship(friendshipID)
	return nil
}

// removeFriendship removes a friendship
func (fm *FriendManager) removeFriendship(id string) {
	friendship, exists := fm.friendships[id]
	if !exists {
		return
	}
	
	// Remove from byPlayer lists
	for _, playerID := range []string{friendship.PlayerA, friendship.PlayerB} {
		ids := fm.byPlayer[playerID]
		for i, fid := range ids {
			if fid == id {
				fm.byPlayer[playerID] = append(ids[:i], ids[i+1:]...)
				break
			}
		}
	}
	
	delete(fm.friendships, id)
}

// RemoveFriend removes a friendship
func (fm *FriendManager) RemoveFriend(playerID, friendID string) error {
	id := generateFriendshipID(playerID, friendID)
	friendship, exists := fm.friendships[id]
	if !exists {
		return fmt.Errorf("friendship not found")
	}
	
	// Verify player is part of this friendship
	if friendship.PlayerA != playerID && friendship.PlayerB != playerID {
		return fmt.Errorf("not part of this friendship")
	}
	
	fm.removeFriendship(id)
	return nil
}

// Block blocks a player
func (fm *FriendManager) Block(blockerID, blockedID string) error {
	id := generateFriendshipID(blockerID, blockedID)
	
	friendship, exists := fm.friendships[id]
	if !exists {
		// Create new blocked friendship
		friendship = NewFriendship(blockerID, blockedID, blockerID)
		friendship.Block()
		fm.friendships[id] = friendship
		fm.byPlayer[blockerID] = append(fm.byPlayer[blockerID], id)
		fm.byPlayer[blockedID] = append(fm.byPlayer[blockedID], id)
		return nil
	}
	
	friendship.Block()
	return nil
}

// Unblock unblocks a player
func (fm *FriendManager) Unblock(blockerID, blockedID string) error {
	id := generateFriendshipID(blockerID, blockedID)
	friendship, exists := fm.friendships[id]
	if !exists {
		return fmt.Errorf("friendship not found")
	}
	
	if !friendship.IsBlocked() {
		return fmt.Errorf("player is not blocked")
	}
	
	fm.removeFriendship(id)
	return nil
}

// GetFriendship gets a friendship between two players
func (fm *FriendManager) GetFriendship(playerA, playerB string) (*Friendship, bool) {
	id := generateFriendshipID(playerA, playerB)
	friendship, exists := fm.friendships[id]
	return friendship, exists
}

// GetFriends returns all friends for a player
func (fm *FriendManager) GetFriends(playerID string) []Friendship {
	friendshipIDs := fm.byPlayer[playerID]
	friends := make([]Friendship, 0)
	
	for _, id := range friendshipIDs {
		if friendship, exists := fm.friendships[id]; exists {
			if friendship.IsFriend() {
				friends = append(friends, *friendship)
			}
		}
	}
	
	return friends
}

// GetFriendIDs returns just the friend IDs
func (fm *FriendManager) GetFriendIDs(playerID string) []string {
	friends := fm.GetFriends(playerID)
	ids := make([]string, 0, len(friends))
	
	for _, f := range friends {
		if f.PlayerA == playerID {
			ids = append(ids, f.PlayerB)
		} else {
			ids = append(ids, f.PlayerA)
		}
	}
	
	return ids
}

// GetPendingRequests returns pending friend requests for a player
func (fm *FriendManager) GetPendingRequests(playerID string) []Friendship {
	friendshipIDs := fm.byPlayer[playerID]
	pending := make([]Friendship, 0)
	
	for _, id := range friendshipIDs {
		if friendship, exists := fm.friendships[id]; exists {
			if friendship.IsPending() && friendship.InitiatedBy != playerID {
				// This player is the target of the request
				pending = append(pending, *friendship)
			}
		}
	}
	
	return pending
}

// GetSentRequests returns friend requests sent by a player
func (fm *FriendManager) GetSentRequests(playerID string) []Friendship {
	friendshipIDs := fm.byPlayer[playerID]
	sent := make([]Friendship, 0)
	
	for _, id := range friendshipIDs {
		if friendship, exists := fm.friendships[id]; exists {
			if friendship.IsPending() && friendship.InitiatedBy == playerID {
				sent = append(sent, *friendship)
			}
		}
	}
	
	return sent
}

// GetBlocked returns blocked players
func (fm *FriendManager) GetBlocked(playerID string) []Friendship {
	friendshipIDs := fm.byPlayer[playerID]
	blocked := make([]Friendship, 0)
	
	for _, id := range friendshipIDs {
		if friendship, exists := fm.friendships[id]; exists {
			if friendship.IsBlocked() {
				blocked = append(blocked, *friendship)
			}
		}
	}
	
	return blocked
}

// IsFriend checks if two players are friends
func (fm *FriendManager) IsFriend(playerA, playerB string) bool {
	friendship, exists := fm.GetFriendship(playerA, playerB)
	if !exists {
		return false
	}
	return friendship.IsFriend()
}

// IsBlocked checks if playerA has blocked playerB
func (fm *FriendManager) IsBlocked(blockerID, blockedID string) bool {
	friendship, exists := fm.GetFriendship(blockerID, blockedID)
	if !exists {
		return false
	}
	return friendship.IsBlocked()
}

// CanInteract checks if two players can interact (not blocked)
func (fm *FriendManager) CanInteract(playerA, playerB string) bool {
	// Check if either has blocked the other
	if fm.IsBlocked(playerA, playerB) || fm.IsBlocked(playerB, playerA) {
		return false
	}
	return true
}

// RecordPlayTime records time two friends spent together
func (fm *FriendManager) RecordPlayTime(playerA, playerB string, minutes int) {
	friendship, exists := fm.GetFriendship(playerA, playerB)
	if !exists {
		return
	}
	
	if friendship.IsFriend() {
		friendship.RecordPlayTime(minutes)
	}
}

// GetFriendCount returns number of friends
func (fm *FriendManager) GetFriendCount(playerID string) int {
	return len(fm.GetFriends(playerID))
}

// Save saves friendships to disk
func (fm *FriendManager) Save() error {
	data, err := json.MarshalIndent(fm.friendships, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal friendships: %w", err)
	}
	
	if err := os.WriteFile(fm.storagePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write friendships: %w", err)
	}
	
	return nil
}

// Load loads friendships from disk
func (fm *FriendManager) Load() error {
	data, err := os.ReadFile(fm.storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read friendships: %w", err)
	}
	
	var loaded map[string]Friendship
	if err := json.Unmarshal(data, &loaded); err != nil {
		return fmt.Errorf("failed to unmarshal friendships: %w", err)
	}
	
	// Clear and reload
	fm.friendships = make(map[string]*Friendship)
	fm.byPlayer = make(map[string][]string)
	
	for id, friendship := range loaded {
		f := friendship // Copy
		fm.friendships[id] = &f
		fm.byPlayer[friendship.PlayerA] = append(fm.byPlayer[friendship.PlayerA], id)
		fm.byPlayer[friendship.PlayerB] = append(fm.byPlayer[friendship.PlayerB], id)
	}
	
	return nil
}

// GetTopFriends returns friends sorted by time spent together
func (fm *FriendManager) GetTopFriends(playerID string, count int) []Friendship {
	friends := fm.GetFriends(playerID)
	
	// Sort by minutes together (bubble sort)
	for i := 0; i < len(friends); i++ {
		for j := i + 1; j < len(friends); j++ {
			if friends[i].MinutesTogether < friends[j].MinutesTogether {
				friends[i], friends[j] = friends[j], friends[i]
			}
		}
	}
	
	if count > len(friends) {
		count = len(friends)
	}
	
	return friends[:count]
}
