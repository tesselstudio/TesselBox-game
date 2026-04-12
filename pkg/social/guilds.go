package social

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// GuildRank represents a guild member rank
type GuildRank int

const (
	GuildRecruit GuildRank = iota
	RankMember
	GuildVeteran
	GuildOfficer
	GuildCoLeader
	GuildLeader
)

// String returns rank name
func (g GuildRank) String() string {
	switch g {
	case GuildRecruit:
		return "Recruit"
	case RankMember:
		return "Member"
	case GuildVeteran:
		return "Veteran"
	case GuildOfficer:
		return "Officer"
	case GuildCoLeader:
		return "Co-Leader"
	case GuildLeader:
		return "Leader"
	}
	return "Unknown"
}

// GuildMemberData represents a guild member
type GuildMemberData struct {
	PlayerID     string    `json:"player_id"`
	PlayerName   string    `json:"player_name"`
	Rank         GuildRank `json:"rank"`
	JoinedAt     time.Time `json:"joined_at"`
	LastActive   time.Time `json:"last_active"`
	Contribution float64   `json:"contribution"` // Money contributed
}

// CanInvite checks if member can invite others
func (gm *GuildMemberData) CanInvite() bool {
	return gm.Rank >= RankMember
}

// CanKick checks if member can kick others
func (gm *GuildMemberData) CanKick() bool {
	return gm.Rank >= GuildOfficer
}

// CanPromote checks if member can promote others
func (gm *GuildMemberData) CanPromote() bool {
	return gm.Rank >= GuildCoLeader
}

// CanManage checks if member can manage guild settings
func (gm *GuildMemberData) CanManage() bool {
	return gm.Rank >= GuildOfficer
}

// CanDisband checks if member can disband guild
func (gm *GuildMemberData) CanDisband() bool {
	return gm.Rank >= GuildLeader
}

// GuildRelation represents relation to another guild
type GuildRelation struct {
	GuildID   string    `json:"guild_id"`
	GuildName string    `json:"guild_name"`
	Relation  string    `json:"relation"` // "ally", "enemy", "neutral"
	Since     time.Time `json:"since"`
}

// Guild represents a player guild
type Guild struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Tag         string `json:"tag"` // [TAG] prefix
	Description string `json:"description"`
	WorldID     string `json:"world_id"`

	// Leadership
	LeaderID  string   `json:"leader_id"`
	CoLeaders []string `json:"co_leaders,omitempty"`
	Officers  []string `json:"officers,omitempty"`

	// Members
	Members    []GuildMemberData `json:"members"`
	MaxMembers int               `json:"max_members"`

	// Bank
	BankBalance float64 `json:"bank_balance"`
	TaxRate     float64 `json:"tax_rate"` // % of member earnings

	// Stats
	Level       int `json:"level"`
	Experience  int `json:"experience"`
	TotalKills  int `json:"total_kills"`
	TotalDeaths int `json:"total_deaths"`
	WarsWon     int `json:"wars_won"`
	WarsLost    int `json:"wars_lost"`

	// Relations
	Allies  []GuildRelation `json:"allies,omitempty"`
	Enemies []GuildRelation `json:"enemies,omitempty"`

	// Settings
	RecruitmentOpen bool    `json:"recruitment_open"`
	MinLevelToJoin  int     `json:"min_level_to_join"`
	EntryFee        float64 `json:"entry_fee"`

	// Territory
	LandClaims []string `json:"land_claims,omitempty"` // Claim IDs

	// Meta
	CreatedAt  time.Time `json:"created_at"`
	LastActive time.Time `json:"last_active"`
}

// NewGuild creates a new guild
func NewGuild(id, name, tag, leaderID, leaderName, worldID string) *Guild {
	now := time.Now()
	return &Guild{
		ID:              id,
		Name:            name,
		Tag:             tag,
		LeaderID:        leaderID,
		Members:         make([]GuildMemberData, 0),
		MaxMembers:      50,
		BankBalance:     0,
		TaxRate:         0.05, // 5% default
		Level:           1,
		Experience:      0,
		CoLeaders:       make([]string, 0),
		Officers:        make([]string, 0),
		Allies:          make([]GuildRelation, 0),
		Enemies:         make([]GuildRelation, 0),
		LandClaims:      make([]string, 0),
		RecruitmentOpen: true,
		MinLevelToJoin:  0,
		EntryFee:        0,
		CreatedAt:       now,
		LastActive:      now,
	}
}

// AddMember adds a member to the guild
func (g *Guild) AddMember(playerID, playerName string, rank GuildRank) error {
	// Check if full
	if len(g.Members) >= g.MaxMembers {
		return fmt.Errorf("guild is full")
	}

	// Check if already member
	for _, m := range g.Members {
		if m.PlayerID == playerID {
			return fmt.Errorf("already a member")
		}
	}

	member := GuildMemberData{
		PlayerID:   playerID,
		PlayerName: playerName,
		Rank:       rank,
		JoinedAt:   time.Now(),
		LastActive: time.Now(),
	}

	g.Members = append(g.Members, member)
	g.LastActive = time.Now()

	return nil
}

// RemoveMember removes a member
func (g *Guild) RemoveMember(playerID string) error {
	for i, m := range g.Members {
		if m.PlayerID == playerID {
			g.Members = append(g.Members[:i], g.Members[i+1:]...)
			g.LastActive = time.Now()

			// Remove from leadership lists
			g.removeFromLeaders(playerID)

			return nil
		}
	}
	return fmt.Errorf("member not found")
}

// removeFromLeaders removes player from leadership lists
func (g *Guild) removeFromLeaders(playerID string) {
	// Remove from co-leaders
	for i, id := range g.CoLeaders {
		if id == playerID {
			g.CoLeaders = append(g.CoLeaders[:i], g.CoLeaders[i+1:]...)
			break
		}
	}

	// Remove from officers
	for i, id := range g.Officers {
		if id == playerID {
			g.Officers = append(g.Officers[:i], g.Officers[i+1:]...)
			break
		}
	}
}

// GetMember gets a guild member
func (g *Guild) GetMember(playerID string) *GuildMemberData {
	for i := range g.Members {
		if g.Members[i].PlayerID == playerID {
			return &g.Members[i]
		}
	}
	return nil
}

// IsMember checks if player is in guild
func (g *Guild) IsMember(playerID string) bool {
	return g.GetMember(playerID) != nil
}

// IsLeader checks if player is leader
func (g *Guild) IsLeader(playerID string) bool {
	return g.LeaderID == playerID
}

// SetRank sets a member's rank
func (g *Guild) SetRank(playerID string, rank GuildRank) error {
	member := g.GetMember(playerID)
	if member == nil {
		return fmt.Errorf("member not found")
	}

	oldRank := member.Rank
	member.Rank = rank
	member.LastActive = time.Now()

	// Update leadership lists
	if rank == GuildCoLeader {
		g.addCoLeader(playerID)
	} else if oldRank == GuildCoLeader {
		g.removeFromLeaders(playerID)
	}

	if rank == GuildOfficer {
		g.addOfficer(playerID)
	} else if oldRank == GuildOfficer {
		g.removeFromLeaders(playerID)
	}

	g.LastActive = time.Now()
	return nil
}

// addCoLeader adds a co-leader
func (g *Guild) addCoLeader(playerID string) {
	for _, id := range g.CoLeaders {
		if id == playerID {
			return // Already exists
		}
	}
	g.CoLeaders = append(g.CoLeaders, playerID)
}

// addOfficer adds an officer
func (g *Guild) addOfficer(playerID string) {
	for _, id := range g.Officers {
		if id == playerID {
			return // Already exists
		}
	}
	g.Officers = append(g.Officers, playerID)
}

// TransferLeadership transfers leadership to another member
func (g *Guild) TransferLeadership(newLeaderID string) error {
	if !g.IsMember(newLeaderID) {
		return fmt.Errorf("new leader must be a member")
	}

	// Demote old leader to member (or remove from leadership)
	oldLeader := g.GetMember(g.LeaderID)
	if oldLeader != nil {
		oldLeader.Rank = GuildVeteran
	}

	// Promote new leader
	newLeader := g.GetMember(newLeaderID)
	newLeader.Rank = GuildLeader
	g.LeaderID = newLeaderID

	g.LastActive = time.Now()
	return nil
}

// AddToBank adds money to guild bank
func (g *Guild) AddToBank(amount float64) {
	g.BankBalance += amount
	g.LastActive = time.Now()
}

// WithdrawFromBank withdraws money from guild bank
func (g *Guild) WithdrawFromBank(amount float64) bool {
	if g.BankBalance < amount {
		return false
	}
	g.BankBalance -= amount
	g.LastActive = time.Now()
	return true
}

// AddAlly adds an ally guild
func (g *Guild) AddAlly(guildID, guildName string) {
	// Check if already ally
	for _, a := range g.Allies {
		if a.GuildID == guildID {
			return
		}
	}

	g.Allies = append(g.Allies, GuildRelation{
		GuildID:   guildID,
		GuildName: guildName,
		Relation:  "ally",
		Since:     time.Now(),
	})
	g.LastActive = time.Now()
}

// AddEnemy adds an enemy guild
func (g *Guild) AddEnemy(guildID, guildName string) {
	// Check if already enemy
	for _, e := range g.Enemies {
		if e.GuildID == guildID {
			return
		}
	}

	g.Enemies = append(g.Enemies, GuildRelation{
		GuildID:   guildID,
		GuildName: guildName,
		Relation:  "enemy",
		Since:     time.Now(),
	})
	g.LastActive = time.Now()
}

// RemoveRelation removes an ally or enemy
func (g *Guild) RemoveRelation(guildID string) {
	// Remove from allies
	for i, a := range g.Allies {
		if a.GuildID == guildID {
			g.Allies = append(g.Allies[:i], g.Allies[i+1:]...)
			return
		}
	}

	// Remove from enemies
	for i, e := range g.Enemies {
		if e.GuildID == guildID {
			g.Enemies = append(g.Enemies[:i], g.Enemies[i+1:]...)
			return
		}
	}
}

// RecordWarResult records war outcome
func (g *Guild) RecordWarResult(won bool) {
	if won {
		g.WarsWon++
	} else {
		g.WarsLost++
	}
	g.LastActive = time.Now()
}

// GetOnlineCount returns number of online members
func (g *Guild) GetOnlineCount() int {
	count := 0
	for _, m := range g.Members {
		// In real implementation, check online status
		_ = m
		count++
	}
	return count
}

// GuildManager manages all guilds
type GuildManager struct {
	guilds   map[string]*Guild
	byName   map[string]string // Name -> ID
	byPlayer map[string]string // PlayerID -> GuildID

	guildCounter int

	storagePath string
}

// NewGuildManager creates a new guild manager
func NewGuildManager(storageDir string) *GuildManager {
	return &GuildManager{
		guilds:       make(map[string]*Guild),
		byName:       make(map[string]string),
		byPlayer:     make(map[string]string),
		guildCounter: 0,
		storagePath:  filepath.Join(storageDir, "guilds.json"),
	}
}

// CreateGuild creates a new guild
func (gm *GuildManager) CreateGuild(name, tag, leaderID, leaderName, worldID string) (*Guild, error) {
	// Check if name taken
	if _, exists := gm.byName[name]; exists {
		return nil, fmt.Errorf("guild name already taken")
	}

	// Check if leader already in guild
	if _, exists := gm.byPlayer[leaderID]; exists {
		return nil, fmt.Errorf("already in a guild")
	}

	gm.guildCounter++
	guildID := fmt.Sprintf("guild_%d_%d", gm.guildCounter, time.Now().Unix())

	guild := NewGuild(guildID, name, tag, leaderID, leaderName, worldID)
	guild.AddMember(leaderID, leaderName, GuildLeader)

	gm.guilds[guildID] = guild
	gm.byName[name] = guildID
	gm.byPlayer[leaderID] = guildID

	return guild, nil
}

// GetGuild gets a guild by ID
func (gm *GuildManager) GetGuild(guildID string) (*Guild, bool) {
	guild, exists := gm.guilds[guildID]
	return guild, exists
}

// GetGuildByName gets a guild by name
func (gm *GuildManager) GetGuildByName(name string) (*Guild, bool) {
	guildID, exists := gm.byName[name]
	if !exists {
		return nil, false
	}
	return gm.GetGuild(guildID)
}

// GetPlayerGuild gets the guild a player is in
func (gm *GuildManager) GetPlayerGuild(playerID string) (*Guild, bool) {
	guildID, exists := gm.byPlayer[playerID]
	if !exists {
		return nil, false
	}
	return gm.GetGuild(guildID)
}

// IsInGuild checks if player is in a guild
func (gm *GuildManager) IsInGuild(playerID string) bool {
	_, exists := gm.byPlayer[playerID]
	return exists
}

// JoinGuild adds a player to a guild
func (gm *GuildManager) JoinGuild(guildID, playerID, playerName string) error {
	guild, exists := gm.GetGuild(guildID)
	if !exists {
		return fmt.Errorf("guild not found")
	}

	if gm.IsInGuild(playerID) {
		return fmt.Errorf("already in a guild")
	}

	if err := guild.AddMember(playerID, playerName, GuildRecruit); err != nil {
		return err
	}

	gm.byPlayer[playerID] = guildID
	return nil
}

// LeaveGuild removes a player from their guild
func (gm *GuildManager) LeaveGuild(playerID string) error {
	guild, exists := gm.GetPlayerGuild(playerID)
	if !exists {
		return fmt.Errorf("not in a guild")
	}

	// Check if leader
	if guild.IsLeader(playerID) {
		return fmt.Errorf("leader must transfer leadership before leaving")
	}

	if err := guild.RemoveMember(playerID); err != nil {
		return err
	}

	delete(gm.byPlayer, playerID)
	return nil
}

// KickFromGuild kicks a player from a guild
func (gm *GuildManager) KickFromGuild(guildID, kickerID, targetID string) error {
	guild, exists := gm.GetGuild(guildID)
	if !exists {
		return fmt.Errorf("guild not found")
	}

	kicker := guild.GetMember(kickerID)
	if kicker == nil || !kicker.CanKick() {
		return fmt.Errorf("permission denied")
	}

	if guild.IsLeader(targetID) {
		return fmt.Errorf("cannot kick the leader")
	}

	if err := guild.RemoveMember(targetID); err != nil {
		return err
	}

	delete(gm.byPlayer, targetID)
	return nil
}

// DisbandGuild disbands a guild
func (gm *GuildManager) DisbandGuild(guildID, disbanderID string) error {
	guild, exists := gm.GetGuild(guildID)
	if !exists {
		return fmt.Errorf("guild not found")
	}

	if !guild.IsLeader(disbanderID) {
		return fmt.Errorf("only leader can disband")
	}

	// Remove all members from byPlayer
	for _, m := range guild.Members {
		delete(gm.byPlayer, m.PlayerID)
	}

	// Remove guild
	delete(gm.guilds, guildID)
	delete(gm.byName, guild.Name)

	return nil
}

// GetAllGuilds returns all guilds
func (gm *GuildManager) GetAllGuilds() []*Guild {
	guilds := make([]*Guild, 0, len(gm.guilds))
	for _, guild := range gm.guilds {
		guilds = append(guilds, guild)
	}
	return guilds
}

// GetRecruitingGuilds returns guilds open for recruitment
func (gm *GuildManager) GetRecruitingGuilds() []*Guild {
	result := make([]*Guild, 0)
	for _, guild := range gm.guilds {
		if guild.RecruitmentOpen && len(guild.Members) < guild.MaxMembers {
			result = append(result, guild)
		}
	}
	return result
}

// GetTopGuilds returns top guilds by level/experience
func (gm *GuildManager) GetTopGuilds(count int) []*Guild {
	allGuilds := gm.GetAllGuilds()

	// Sort by level then experience
	for i := 0; i < len(allGuilds); i++ {
		for j := i + 1; j < len(allGuilds); j++ {
			if allGuilds[i].Level < allGuilds[j].Level ||
				(allGuilds[i].Level == allGuilds[j].Level && allGuilds[i].Experience < allGuilds[j].Experience) {
				allGuilds[i], allGuilds[j] = allGuilds[j], allGuilds[i]
			}
		}
	}

	if count > len(allGuilds) {
		count = len(allGuilds)
	}

	return allGuilds[:count]
}

// Save saves guilds to disk
func (gm *GuildManager) Save() error {
	data, err := json.MarshalIndent(gm.guilds, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal guilds: %w", err)
	}

	if err := os.WriteFile(gm.storagePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write guilds: %w", err)
	}

	return nil
}

// Load loads guilds from disk
func (gm *GuildManager) Load() error {
	data, err := os.ReadFile(gm.storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read guilds: %w", err)
	}

	if err := json.Unmarshal(data, &gm.guilds); err != nil {
		return fmt.Errorf("failed to unmarshal guilds: %w", err)
	}

	// Rebuild indexes
	gm.byName = make(map[string]string)
	gm.byPlayer = make(map[string]string)

	for _, guild := range gm.guilds {
		gm.byName[guild.Name] = guild.ID
		for _, m := range guild.Members {
			gm.byPlayer[m.PlayerID] = guild.ID
		}
	}

	return nil
}
