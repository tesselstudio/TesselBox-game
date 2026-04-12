package social

import (
	"fmt"
	"time"

	"tesselbox/pkg/economy"
)

// GuildWarStatus represents war state
type GuildWarStatus int

const (
	WarDeclared GuildWarStatus = iota
	WarActive
	WarEnded
	WarCancelled
)

// GuildWar represents a war between two guilds
type GuildWar struct {
	ID         string `json:"id"`
	AttackerID string `json:"attacker_id"`
	DefenderID string `json:"defender_id"`

	// Status
	Status   GuildWarStatus `json:"status"`
	WinnerID *string        `json:"winner_id,omitempty"`

	// Terms
	Wager     float64       `json:"wager"`
	MaxKills  int           `json:"max_kills"`
	TimeLimit time.Duration `json:"time_limit"`

	// Scores
	AttackerKills int `json:"attacker_kills"`
	DefenderKills int `json:"defender_kills"`

	// Active players
	AttackerPlayers []string `json:"attacker_players,omitempty"`
	DefenderPlayers []string `json:"defender_players,omitempty"`

	// Timing
	DeclaredAt time.Time  `json:"declared_at"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	EndedAt    *time.Time `json:"ended_at,omitempty"`
}

// NewGuildWar creates a new war
func NewGuildWar(id, attackerID, defenderID string, wager float64, maxKills int, timeLimit time.Duration) *GuildWar {
	return &GuildWar{
		ID:              id,
		AttackerID:      attackerID,
		DefenderID:      defenderID,
		Status:          WarDeclared,
		Wager:           wager,
		MaxKills:        maxKills,
		TimeLimit:       timeLimit,
		AttackerKills:   0,
		DefenderKills:   0,
		AttackerPlayers: make([]string, 0),
		DefenderPlayers: make([]string, 0),
		DeclaredAt:      time.Now(),
	}
}

// Accept accepts the war declaration
func (gw *GuildWar) Accept() {
	if gw.Status == WarDeclared {
		gw.Status = WarActive
		now := time.Now()
		gw.StartedAt = &now
	}
}

// Decline declines the war
func (gw *GuildWar) Decline() {
	if gw.Status == WarDeclared {
		gw.Status = WarCancelled
	}
}

// RecordKill records a kill in the war
func (gw *GuildWar) RecordKill(killerGuildID string) {
	if gw.Status != WarActive {
		return
	}

	if killerGuildID == gw.AttackerID {
		gw.AttackerKills++
	} else if killerGuildID == gw.DefenderID {
		gw.DefenderKills++
	}

	// Check win condition
	if gw.MaxKills > 0 {
		if gw.AttackerKills >= gw.MaxKills {
			gw.End(&gw.AttackerID)
		} else if gw.DefenderKills >= gw.MaxKills {
			gw.End(&gw.DefenderID)
		}
	}
}

// End ends the war
func (gw *GuildWar) End(winnerID *string) {
	if gw.Status != WarActive {
		return
	}

	now := time.Now()
	gw.Status = WarEnded
	gw.WinnerID = winnerID
	gw.EndedAt = &now
}

// IsExpired checks if war has expired
func (gw *GuildWar) IsExpired() bool {
	if gw.StartedAt == nil {
		return false
	}
	return time.Since(*gw.StartedAt) > gw.TimeLimit
}

// GetScore gets score for a guild
func (gw *GuildWar) GetScore(guildID string) int {
	if guildID == gw.AttackerID {
		return gw.AttackerKills
	}
	return gw.DefenderKills
}

// GuildWarManager manages guild wars
type GuildWarManager struct {
	wars    map[string]*GuildWar
	byGuild map[string][]string // GuildID -> War IDs

	warCounter int
	guildMgr   *GuildManager
	walletMgr  *economy.WalletManager
}

// NewGuildWarManager creates a new manager
func NewGuildWarManager(guildMgr *GuildManager, walletMgr *economy.WalletManager) *GuildWarManager {
	return &GuildWarManager{
		wars:       make(map[string]*GuildWar),
		byGuild:    make(map[string][]string),
		warCounter: 0,
		guildMgr:   guildMgr,
		walletMgr:  walletMgr,
	}
}

// DeclareWar declares war on another guild
func (gwm *GuildWarManager) DeclareWar(attackerID, defenderID string, wager float64, maxKills int) (*GuildWar, error) {
	// Get guilds
	attacker, exists := gwm.guildMgr.GetGuild(attackerID)
	if !exists {
		return nil, fmt.Errorf("attacking guild not found")
	}

	_, exists = gwm.guildMgr.GetGuild(defenderID)
	if !exists {
		return nil, fmt.Errorf("defending guild not found")
	}

	// Check if already at war
	if gwm.AreAtWar(attackerID, defenderID) {
		return nil, fmt.Errorf("already at war")
	}

	// Check if allies
	for _, ally := range attacker.Allies {
		if ally.GuildID == defenderID {
			return nil, fmt.Errorf("cannot declare war on ally")
		}
	}

	// Check wager
	if wager > 0 {
		// In real implementation, check guild bank
		_ = attacker.BankBalance
	}

	gwm.warCounter++
	warID := fmt.Sprintf("war_%d_%d", gwm.warCounter, time.Now().Unix())

	war := NewGuildWar(warID, attackerID, defenderID, wager, maxKills, 24*time.Hour)

	gwm.wars[warID] = war
	gwm.byGuild[attackerID] = append(gwm.byGuild[attackerID], warID)
	gwm.byGuild[defenderID] = append(gwm.byGuild[defenderID], warID)

	return war, nil
}

// AcceptWar accepts a war declaration
func (gwm *GuildWarManager) AcceptWar(warID, guildID string) error {
	war, exists := gwm.GetWar(warID)
	if !exists {
		return fmt.Errorf("war not found")
	}

	if war.DefenderID != guildID {
		return fmt.Errorf("only defender can accept")
	}

	if war.Status != WarDeclared {
		return fmt.Errorf("war not in declared state")
	}

	war.Accept()
	return nil
}

// DeclineWar declines a war declaration
func (gwm *GuildWarManager) DeclineWar(warID, guildID string) error {
	war, exists := gwm.GetWar(warID)
	if !exists {
		return fmt.Errorf("war not found")
	}

	if war.DefenderID != guildID {
		return fmt.Errorf("only defender can decline")
	}

	war.Decline()
	return nil
}

// GetWar gets a war
func (gwm *GuildWarManager) GetWar(warID string) (*GuildWar, bool) {
	war, exists := gwm.wars[warID]
	return war, exists
}

// GetGuildWars gets all wars for a guild
func (gwm *GuildWarManager) GetGuildWars(guildID string) []*GuildWar {
	warIDs := gwm.byGuild[guildID]
	wars := make([]*GuildWar, 0, len(warIDs))

	for _, id := range warIDs {
		if war, exists := gwm.wars[id]; exists {
			wars = append(wars, war)
		}
	}

	return wars
}

// GetActiveWars gets active wars for a guild
func (gwm *GuildWarManager) GetActiveWars(guildID string) []*GuildWar {
	all := gwm.GetGuildWars(guildID)
	active := make([]*GuildWar, 0)

	for _, war := range all {
		if war.Status == WarActive {
			active = append(active, war)
		}
	}

	return active
}

// AreAtWar checks if two guilds are at war
func (gwm *GuildWarManager) AreAtWar(guildA, guildB string) bool {
	wars := gwm.GetGuildWars(guildA)

	for _, war := range wars {
		if (war.AttackerID == guildB || war.DefenderID == guildB) && war.Status == WarActive {
			return true
		}
	}

	return false
}

// RecordKill records a kill in a war
func (gwm *GuildWarManager) RecordKill(killerGuildID, victimGuildID string) {
	// Find active war between these guilds
	wars := gwm.GetActiveWars(killerGuildID)

	for _, war := range wars {
		if war.AttackerID == victimGuildID || war.DefenderID == victimGuildID {
			war.RecordKill(killerGuildID)

			// Update guild stats
			if killer, exists := gwm.guildMgr.GetGuild(killerGuildID); exists {
				killer.TotalKills++
			}
			if victim, exists := gwm.guildMgr.GetGuild(victimGuildID); exists {
				victim.TotalDeaths++
			}

			break
		}
	}
}

// EndWar ends a war
func (gwm *GuildWarManager) EndWar(warID string, winnerID *string) error {
	war, exists := gwm.GetWar(warID)
	if !exists {
		return fmt.Errorf("war not found")
	}

	war.End(winnerID)

	// Pay wager
	if winnerID != nil && war.Wager > 0 {
		if winner, exists := gwm.guildMgr.GetGuild(*winnerID); exists {
			winner.AddToBank(war.Wager * 2)

			winner.RecordWarResult(true)
			if loser, exists := gwm.guildMgr.GetGuild(war.AttackerID); exists && loser.ID != *winnerID {
				loser.RecordWarResult(false)
			}
			if loser, exists := gwm.guildMgr.GetGuild(war.DefenderID); exists && loser.ID != *winnerID {
				loser.RecordWarResult(false)
			}
		}
	}

	return nil
}

// Update updates all wars
func (gwm *GuildWarManager) Update() {
	for _, war := range gwm.wars {
		if war.Status == WarActive && war.IsExpired() {
			// Time expired - defender wins by default
			gwm.EndWar(war.ID, &war.DefenderID)
		}
	}
}

// GetAllActiveWars returns all active wars
func (gwm *GuildWarManager) GetAllActiveWars() []*GuildWar {
	result := make([]*GuildWar, 0)
	for _, war := range gwm.wars {
		if war.Status == WarActive {
			result = append(result, war)
		}
	}
	return result
}

// GetWarStats returns war statistics
func (gwm *GuildWarManager) GetWarStats() (declared, active, ended int) {
	for _, war := range gwm.wars {
		switch war.Status {
		case WarDeclared:
			declared++
		case WarActive:
			active++
		case WarEnded:
			ended++
		}
	}
	return
}
