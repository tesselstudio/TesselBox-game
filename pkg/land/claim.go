package land

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ChunkCoord represents a chunk coordinate
type ChunkCoord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// String returns string representation
func (c ChunkCoord) String() string {
	return fmt.Sprintf("%d,%d", c.X, c.Y)
}

// MemberPermission represents permission levels for claim members
type MemberPermission int

const (
	PermNone   MemberPermission = iota
	PermVisit                   // Can enter, look around
	PermBuild                   // Can place/break blocks
	PermManage                  // Can trust/untrust others
	PermOwner                   // Full control
)

// String returns human-readable permission name
func (p MemberPermission) String() string {
	switch p {
	case PermNone:
		return "none"
	case PermVisit:
		return "visit"
	case PermBuild:
		return "build"
	case PermManage:
		return "manage"
	case PermOwner:
		return "owner"
	}
	return "unknown"
}

// ClaimFlags contains togglable claim settings
type ClaimFlags struct {
	PvPAllowed   bool   `json:"pvp_allowed"`
	MobSpawning  bool   `json:"mob_spawning"`
	FireSpread   bool   `json:"fire_spread"`
	Explosions   bool   `json:"explosions"`
	ItemPickup   bool   `json:"item_pickup"`
	PublicAccess bool   `json:"public_access"`
	EntryMessage string `json:"entry_message,omitempty"`
	ExitMessage  string `json:"exit_message,omitempty"`
}

// DefaultClaimFlags returns default flags
func DefaultClaimFlags() ClaimFlags {
	return ClaimFlags{
		PvPAllowed:   false,
		MobSpawning:  true,
		FireSpread:   false,
		Explosions:   false,
		ItemPickup:   true,
		PublicAccess: false,
		EntryMessage: "",
		ExitMessage:  "",
	}
}

// LandClaim represents a claimed land area
type LandClaim struct {
	ID      string `json:"id"`
	OwnerID string `json:"owner_id"`
	WorldID string `json:"world_id"`

	// Area
	MainChunk ChunkCoord   `json:"main_chunk"`
	SubClaims []SubClaim   `json:"sub_claims,omitempty"`
	Adjacent  []ChunkCoord `json:"adjacent,omitempty"` // Expanded claims

	// Members
	Members map[string]MemberPermission `json:"members,omitempty"`
	Banned  []string                    `json:"banned,omitempty"`

	// Settings
	Flags ClaimFlags `json:"flags"`

	// Meta
	CreatedAt  time.Time  `json:"created_at"`
	LastActive time.Time  `json:"last_active"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"` // For rentals

	// Cost
	PurchasePrice float64 `json:"purchase_price"`
	UpkeepCost    float64 `json:"upkeep_cost"` // Monthly
}

// SubClaim represents a sub-division within a claim
type SubClaim struct {
	Name       string
	MinX, MinY float64
	MaxX, MaxY float64
	Members    map[string]MemberPermission
	Flags      ClaimFlags
}

// NewLandClaim creates a new land claim
func NewLandClaim(id, ownerID, worldID string, chunk ChunkCoord, price float64) *LandClaim {
	now := time.Now()
	return &LandClaim{
		ID:            id,
		OwnerID:       ownerID,
		WorldID:       worldID,
		MainChunk:     chunk,
		SubClaims:     make([]SubClaim, 0),
		Adjacent:      make([]ChunkCoord, 0),
		Members:       make(map[string]MemberPermission),
		Banned:        make([]string, 0),
		Flags:         DefaultClaimFlags(),
		CreatedAt:     now,
		LastActive:    now,
		PurchasePrice: price,
		UpkeepCost:    price * 0.1, // 10% monthly upkeep
	}
}

// RecordActivity updates last active time
func (lc *LandClaim) RecordActivity() {
	lc.LastActive = time.Now()
}

// IsOwner checks if player is the owner
func (lc *LandClaim) IsOwner(playerID string) bool {
	return lc.OwnerID == playerID
}

// GetPlayerPermission gets permission level for a player
func (lc *LandClaim) GetPlayerPermission(playerID string) MemberPermission {
	// Owner has full access
	if lc.IsOwner(playerID) {
		return PermOwner
	}

	// Check members
	if perm, exists := lc.Members[playerID]; exists {
		return perm
	}

	// Public access
	if lc.Flags.PublicAccess {
		return PermVisit
	}

	return PermNone
}

// CanBuild checks if player can build
func (lc *LandClaim) CanBuild(playerID string) bool {
	perm := lc.GetPlayerPermission(playerID)
	return perm >= PermBuild
}

// CanInteract checks if player can interact (doors, chests, etc)
func (lc *LandClaim) CanInteract(playerID string) bool {
	perm := lc.GetPlayerPermission(playerID)
	return perm >= PermVisit
}

// CanManage checks if player can manage the claim
func (lc *LandClaim) CanManage(playerID string) bool {
	perm := lc.GetPlayerPermission(playerID)
	return perm >= PermManage
}

// IsBanned checks if player is banned
func (lc *LandClaim) IsBanned(playerID string) bool {
	for _, id := range lc.Banned {
		if id == playerID {
			return true
		}
	}
	return false
}

// Trust adds a member with permission
func (lc *LandClaim) Trust(playerID string, perm MemberPermission) {
	lc.Members[playerID] = perm
	lc.RecordActivity()
}

// Untrust removes a member
func (lc *LandClaim) Untrust(playerID string) {
	delete(lc.Members, playerID)
	lc.RecordActivity()
}

// Ban bans a player
func (lc *LandClaim) Ban(playerID string) {
	// Remove from members if present
	delete(lc.Members, playerID)

	// Add to banned
	if !lc.IsBanned(playerID) {
		lc.Banned = append(lc.Banned, playerID)
	}

	lc.RecordActivity()
}

// Unban unbans a player
func (lc *LandClaim) Unban(playerID string) {
	for i, id := range lc.Banned {
		if id == playerID {
			lc.Banned = append(lc.Banned[:i], lc.Banned[i+1:]...)
			break
		}
	}
	lc.RecordActivity()
}

// Expand adds an adjacent chunk
func (lc *LandClaim) Expand(chunk ChunkCoord) bool {
	// Check if already claimed
	if lc.MainChunk == chunk {
		return false
	}
	for _, adj := range lc.Adjacent {
		if adj == chunk {
			return false
		}
	}

	lc.Adjacent = append(lc.Adjacent, chunk)
	lc.RecordActivity()
	return true
}

// Contract removes an adjacent chunk
func (lc *LandClaim) Contract(chunk ChunkCoord) bool {
	for i, adj := range lc.Adjacent {
		if adj == chunk {
			lc.Adjacent = append(lc.Adjacent[:i], lc.Adjacent[i+1:]...)
			lc.RecordActivity()
			return true
		}
	}
	return false
}

// IsExpired checks if claim has expired
func (lc *LandClaim) IsExpired() bool {
	if lc.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*lc.ExpiresAt)
}

// IsAbandoned checks if claim has been inactive too long
func (lc *LandClaim) IsAbandoned(duration time.Duration) bool {
	return time.Since(lc.LastActive) > duration
}

// ChunkCount returns total number of chunks
func (lc *LandClaim) ChunkCount() int {
	return 1 + len(lc.Adjacent)
}

// Contains checks if a location is within the claim
func (lc *LandClaim) Contains(x, y float64, chunkWidth, chunkHeight float64) bool {
	// Convert to chunk coordinates
	chunkX := int(x / chunkWidth)
	chunkY := int(y / chunkHeight)
	coord := ChunkCoord{X: chunkX, Y: chunkY}

	// Check main chunk
	if lc.MainChunk == coord {
		return true
	}

	// Check adjacent chunks
	for _, adj := range lc.Adjacent {
		if adj == coord {
			return true
		}
	}

	return false
}

// LandManager manages all land claims
type LandManager struct {
	claims  map[string]*LandClaim // By claim ID
	byChunk map[string]string     // Chunk key -> claim ID
	byOwner map[string][]string   // Owner ID -> claim IDs

	chunkWidth  float64
	chunkHeight float64
	claimCost   float64 // Base cost per chunk

	storagePath string
}

// NewLandManager creates a new land manager
func NewLandManager(storageDir string) *LandManager {
	return &LandManager{
		claims:      make(map[string]*LandClaim),
		byChunk:     make(map[string]string),
		byOwner:     make(map[string][]string),
		chunkWidth:  256.0, // Default chunk size
		chunkHeight: 256.0,
		claimCost:   100.0,
		storagePath: filepath.Join(storageDir, "land_claims.json"),
	}
}

// SetChunkSize sets the chunk dimensions
func (lm *LandManager) SetChunkSize(width, height float64) {
	lm.chunkWidth = width
	lm.chunkHeight = height
}

// chunkKey generates a unique key for a chunk coordinate
func (lm *LandManager) chunkKey(c ChunkCoord) string {
	return fmt.Sprintf("%d,%d", c.X, c.Y)
}

// Claim attempts to claim land for a player
func (lm *LandManager) Claim(claimID, ownerID, worldID string, chunk ChunkCoord, price float64) (*LandClaim, error) {
	// Check if chunk is already claimed
	key := lm.chunkKey(chunk)
	if _, exists := lm.byChunk[key]; exists {
		return nil, fmt.Errorf("chunk %v is already claimed", chunk)
	}

	// Create claim
	claim := NewLandClaim(claimID, ownerID, worldID, chunk, price)

	// Register
	lm.claims[claimID] = claim
	lm.byChunk[key] = claimID
	lm.byOwner[ownerID] = append(lm.byOwner[ownerID], claimID)

	return claim, nil
}

// GetClaim gets a claim by ID
func (lm *LandManager) GetClaim(claimID string) (*LandClaim, bool) {
	claim, exists := lm.claims[claimID]
	return claim, exists
}

// GetClaimAt gets claim at a specific location
func (lm *LandManager) GetClaimAt(x, y float64) (*LandClaim, bool) {
	chunkX := int(x / lm.chunkWidth)
	chunkY := int(y / lm.chunkHeight)
	key := lm.chunkKey(ChunkCoord{X: chunkX, Y: chunkY})

	claimID, exists := lm.byChunk[key]
	if !exists {
		return nil, false
	}

	return lm.GetClaim(claimID)
}

// GetClaimsByOwner gets all claims for a player
func (lm *LandManager) GetClaimsByOwner(ownerID string) []*LandClaim {
	claimIDs := lm.byOwner[ownerID]
	claims := make([]*LandClaim, 0, len(claimIDs))

	for _, id := range claimIDs {
		if claim, exists := lm.claims[id]; exists {
			claims = append(claims, claim)
		}
	}

	return claims
}

// CanBuildAt checks if a player can build at a location
func (lm *LandManager) CanBuildAt(x, y float64, playerID string) bool {
	claim, exists := lm.GetClaimAt(x, y)
	if !exists {
		return true // Unclaimed land - anyone can build
	}

	return claim.CanBuild(playerID)
}

// CanInteractAt checks if a player can interact at a location
func (lm *LandManager) CanInteractAt(x, y float64, playerID string) bool {
	claim, exists := lm.GetClaimAt(x, y)
	if !exists {
		return true // Unclaimed land
	}

	if claim.IsBanned(playerID) {
		return false
	}

	return claim.CanInteract(playerID)
}

// Unclaim removes a claim
func (lm *LandManager) Unclaim(claimID string) error {
	claim, exists := lm.claims[claimID]
	if !exists {
		return fmt.Errorf("claim not found")
	}

	// Remove from chunk map
	delete(lm.byChunk, lm.chunkKey(claim.MainChunk))
	for _, adj := range claim.Adjacent {
		delete(lm.byChunk, lm.chunkKey(adj))
	}

	// Remove from owner map
	ownerClaims := lm.byOwner[claim.OwnerID]
	for i, id := range ownerClaims {
		if id == claimID {
			lm.byOwner[claim.OwnerID] = append(ownerClaims[:i], ownerClaims[i+1:]...)
			break
		}
	}

	// Remove claim
	delete(lm.claims, claimID)

	return nil
}

// Transfer transfers claim ownership
func (lm *LandManager) Transfer(claimID, newOwnerID string) error {
	claim, exists := lm.claims[claimID]
	if !exists {
		return fmt.Errorf("claim not found")
	}

	oldOwnerID := claim.OwnerID

	// Remove from old owner
	oldClaims := lm.byOwner[oldOwnerID]
	for i, id := range oldClaims {
		if id == claimID {
			lm.byOwner[oldOwnerID] = append(oldClaims[:i], oldClaims[i+1:]...)
			break
		}
	}

	// Update claim
	claim.OwnerID = newOwnerID

	// Add to new owner
	lm.byOwner[newOwnerID] = append(lm.byOwner[newOwnerID], claimID)

	return nil
}

// ExpandClaim expands a claim to include another chunk
func (lm *LandManager) ExpandClaim(claimID string, chunk ChunkCoord) error {
	claim, exists := lm.claims[claimID]
	if !exists {
		return fmt.Errorf("claim not found")
	}

	// Check if chunk is free
	key := lm.chunkKey(chunk)
	if _, occupied := lm.byChunk[key]; occupied {
		return fmt.Errorf("chunk %v is already claimed", chunk)
	}

	// Add to claim
	if !claim.Expand(chunk) {
		return fmt.Errorf("chunk already in claim")
	}

	// Register chunk
	lm.byChunk[key] = claimID

	return nil
}

// GetClaimCost returns the cost to claim a chunk
func (lm *LandManager) GetClaimCost() float64 {
	return lm.claimCost
}

// SetClaimCost sets the base claim cost
func (lm *LandManager) SetClaimCost(cost float64) {
	lm.claimCost = cost
}

// GetTotalClaimedChunks returns total number of claimed chunks
func (lm *LandManager) GetTotalClaimedChunks() int {
	total := 0
	for _, claim := range lm.claims {
		total += claim.ChunkCount()
	}
	return total
}

// GetOwnerChunkCount returns how many chunks a player owns
func (lm *LandManager) GetOwnerChunkCount(ownerID string) int {
	claims := lm.GetClaimsByOwner(ownerID)
	count := 0
	for _, claim := range claims {
		count += claim.ChunkCount()
	}
	return count
}

// GetAbandonedClaims returns claims abandoned for longer than duration
func (lm *LandManager) GetAbandonedClaims(duration time.Duration) []*LandClaim {
	abandoned := make([]*LandClaim, 0)

	for _, claim := range lm.claims {
		if claim.IsAbandoned(duration) {
			abandoned = append(abandoned, claim)
		}
	}

	return abandoned
}

// Save saves all claims to disk
func (lm *LandManager) Save() error {
	data := struct {
		Claims []*LandClaim `json:"claims"`
		Count  int          `json:"count"`
	}{
		Claims: make([]*LandClaim, 0, len(lm.claims)),
		Count:  len(lm.claims),
	}

	for _, claim := range lm.claims {
		data.Claims = append(data.Claims, claim)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal claims: %w", err)
	}

	if err := os.WriteFile(lm.storagePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write claims: %w", err)
	}

	return nil
}

// Load loads claims from disk
func (lm *LandManager) Load() error {
	data, err := os.ReadFile(lm.storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No existing data
		}
		return fmt.Errorf("failed to read claims: %w", err)
	}

	var loaded struct {
		Claims []*LandClaim `json:"claims"`
	}

	if err := json.Unmarshal(data, &loaded); err != nil {
		return fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	// Clear current state
	lm.claims = make(map[string]*LandClaim)
	lm.byChunk = make(map[string]string)
	lm.byOwner = make(map[string][]string)

	// Load claims
	for _, claim := range loaded.Claims {
		// Fix nil maps
		if claim.Members == nil {
			claim.Members = make(map[string]MemberPermission)
		}

		lm.claims[claim.ID] = claim

		// Register chunks
		lm.byChunk[lm.chunkKey(claim.MainChunk)] = claim.ID
		for _, adj := range claim.Adjacent {
			lm.byChunk[lm.chunkKey(adj)] = claim.ID
		}

		// Register owner
		lm.byOwner[claim.OwnerID] = append(lm.byOwner[claim.OwnerID], claim.ID)
	}

	return nil
}
