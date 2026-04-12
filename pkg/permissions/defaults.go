package permissions

// CreateDefaultRoles creates the built-in role definitions
func CreateDefaultRoles() map[string]*Role {
	roles := make(map[string]*Role)
	
	// Owner - Full access, highest weight
	owner := NewRoleBuilder("owner", "Owner").
		WithDescription("Server owner with full access to all commands and settings").
		WithWeight(100).
		WithPrefix("[OWNER]").
		Wildcard("*"). // All permissions
		Build()
	roles["owner"] = owner
	
	// Co-Owner - Almost full access, can manage most things except server ownership
	coOwner := NewRoleBuilder("co-owner", "Co-Owner").
		WithDescription("Co-owner with extensive administrative access").
		WithWeight(95).
		WithPrefix("[CO-OWNER]").
		Inherit("owner").
		Deny(PermAdminPlayerReset). // Can't reset owner
		Build()
	roles["co-owner"] = coOwner
	
	// Admin - Can manage players, worlds, and moderate
	admin := NewRoleBuilder("admin", "Admin").
		WithDescription("Administrator with player management and world control").
		WithWeight(90).
		WithPrefix("[ADMIN]").
		GrantWildcard("admin.*").
		GrantWildcard("commands.*").
		GrantWildcard("social.*").
		GrantWildcard("cosmetic.*").
		GrantWildcard("minigame.*").
		GrantWildcard("land.*").
		GrantWildcard("economy.*").
		GrantWildcard("build.*").
		GrantWildcard("gamemode.*").
		GrantWildcard("world.*").
		Grant(PermAdminVanish).
		Grant(PermAdminCheatBypass).
		Build()
	roles["admin"] = admin
	
	// Moderator - Can moderate chat and kick/mute players
	moderator := NewRoleBuilder("moderator", "Moderator").
		WithDescription("Moderator with chat control and player moderation powers").
		WithWeight(70).
		WithPrefix("[MOD]").
		Inherit("player").
		Grant(PermCmdKick).
		Grant(PermCmdMute).
		Grant(PermCmdWarn).
		Grant(PermCmdHistory).
		Grant(PermAdminBypassMute).
		Grant(PermChatParty).
		Grant(PermChatGuild).
		Grant(PermAdminVanish).
		Build()
	roles["moderator"] = moderator
	
	// Helper - Junior moderator with limited powers
	helper := NewRoleBuilder("helper", "Helper").
		WithDescription("Helper with basic moderation and player assistance").
		WithWeight(60).
		WithPrefix("[HELPER]").
		Inherit("player").
		Grant(PermCmdWarn).
		Grant(PermCmdHistory).
		Grant(PermChatParty).
		Build()
	roles["helper"] = helper
	
	// Player - Standard player with full gameplay access
	player := NewRoleBuilder("player", "Player").
		WithDescription("Standard player with full gameplay access").
		WithWeight(50).
		Grant(PermBuildPlace).
		Grant(PermBuildBreak).
		Grant(PermBuildInteract).
		Grant(PermGameModeSurvival).
		Grant(PermGameModeSwitch).
		Grant(PermWorldAccess).
		Grant(PermWorldBuild).
		Grant(PermWorldPvP).
		Grant(PermWorldEco).
		Grant(PermWorldTrade).
		Grant(PermLandClaim).
		Grant(PermLandUnclaim).
		Grant(PermLandExpand).
		Grant(PermLandTrustAdd).
		Grant(PermLandTrustRemove).
		Grant(PermCmdTPA).
		Grant(PermCmdPay).
		Grant(PermCmdBalance).
		Grant(PermCmdShop).
		Grant(PermCmdTrade).
		Grant(PermCmdClaim).
		Grant(PermCmdUnclaim).
		Grant(PermCmdTrust).
		Grant(PermCmdHome).
		Grant(PermCmdSetHome).
		Grant(PermCmdWarp).
		Grant(PermCmdMail).
		Grant(PermCmdParty).
		Grant(PermCmdGuild).
		Grant(PermCmdQuest).
		Grant(PermCmdDuel).
		Grant(PermCmdBounty).
		Grant(PermEcoReceive).
		Grant(PermEcoSpend).
		Grant(PermEcoTrade).
		Grant(PermEcoShopCreate).
		Grant(PermEcoShopDelete).
		Grant(PermEcoAuctionBid).
		Grant(PermEcoBankWithdraw).
		Grant(PermEcoBankDeposit).
		Grant(PermEcoJobJoin).
		Grant(PermChatGlobal).
		Grant(PermChatWhisper).
		Grant(PermFriendAdd).
		Grant(PermFriendRemove).
		Grant(PermFriendView).
		Grant(PermPartyCreate).
		Grant(PermPartyJoin).
		Grant(PermGuildCreate).
		Grant(PermGuildJoin).
		Grant(PermMailSend).
		Grant(PermMailReceive).
		Grant(PermMinigameJoin).
		Grant(PermMinigameSpleef).
		Grant(PermMinigameParkour).
		Grant(PermMinigameCTF).
		Grant(PermMinigameTNTRun).
		Grant(PermMinigameMobArena).
		Grant(PermCosmeticParticles).
		Grant(PermCosmeticPets).
		Grant(PermCosmeticTitles).
		Grant(PermCosmeticHats).
		Grant(PermCosmeticEmotes).
		Grant(PermCosmeticTrails).
		Build()
	roles["player"] = player
	
	// Shopkeeper - Player focused on trading and economy
	shopkeeper := NewRoleBuilder("shopkeeper", "Shopkeeper").
		WithDescription("Player with enhanced economy and shop management abilities").
		WithWeight(55).
		WithPrefix("[SHOP]").
		Inherit("player").
		Grant(PermEcoShopAdmin).
		Grant(PermEcoAuctionCreate).
		Grant(PermEcoBankLoan).
		Grant(PermEcoJobAdmin).
		Grant(PermCmdAuction).
		Build()
	roles["shopkeeper"] = shopkeeper
	
	// Builder - Player with enhanced building permissions
	builder := NewRoleBuilder("builder", "Builder").
		WithDescription("Player with enhanced building and creative abilities").
		WithWeight(55).
		WithPrefix("[BUILDER]").
		Inherit("player").
		Grant(PermGameModeCreative).
		Grant(PermBuildAdmin).
		Grant(PermLandSubdivide).
		Grant(PermLandFlagSet).
		Build()
	roles["builder"] = builder
	
	// Visitor - Very limited, view-only access
	visitor := NewRoleBuilder("visitor", "Visitor").
		WithDescription("Visitor with limited access, cannot build or interact with most systems").
		WithWeight(10).
		WithPrefix("[VISITOR]").
		Grant(PermWorldAccess).
		Grant(PermChatGlobal).
		Grant(PermChatWhisper).
		Grant(PermFriendView).
		Grant(PermMinigameJoin).
		Grant(PermMinigameSpleef).
		Grant(PermMinigameParkour).
		Build()
	roles["visitor"] = visitor
	
	return roles
}

// GetDefaultRole returns the default role for new players
func GetDefaultRole() string {
	return "player"
}

// GetVisitorRole returns the role for visitors
func GetVisitorRole() string {
	return "visitor"
}

// IsValidRole checks if a role ID is valid
func IsValidRole(roleID string, roles map[string]*Role) bool {
	_, exists := roles[roleID]
	return exists
}

// GetRoleWeight returns the weight of a role (for comparison)
func GetRoleWeight(roleID string, roles map[string]*Role) int {
	if role, exists := roles[roleID]; exists {
		return role.Weight
	}
	return 0
}

// RoleHierarchy returns roles in order of weight (highest first)
func RoleHierarchy(roles map[string]*Role) []*Role {
	result := make([]*Role, 0, len(roles))
	for _, role := range roles {
		result = append(result, role)
	}
	
	// Simple bubble sort by weight (descending)
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Weight < result[j].Weight {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	
	return result
}
