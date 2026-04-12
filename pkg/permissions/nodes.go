package permissions

// PermissionNode represents a single permission node
type PermissionNode string

// Core build permissions
const (
	PermBuildPlace   PermissionNode = "build.place"
	PermBuildBreak   PermissionNode = "build.break"
	PermBuildInteract PermissionNode = "build.interact"
	PermBuildAdmin   PermissionNode = "build.admin"
)

// Game mode permissions
const (
	PermGameModeCreative  PermissionNode = "gamemode.creative"
	PermGameModeSurvival  PermissionNode = "gamemode.survival"
	PermGameModeSpectator PermissionNode = "gamemode.spectator"
	PermGameModeSwitch    PermissionNode = "gamemode.switch"
)

// Command permissions
const (
	PermCmdGive      PermissionNode = "commands.give"
	PermCmdTP        PermissionNode = "commands.tp"
	PermCmdTPA       PermissionNode = "commands.tpa"
	PermCmdPlugin    PermissionNode = "commands.plugin"
	PermCmdKick      PermissionNode = "commands.kick"
	PermCmdBan       PermissionNode = "commands.ban"
	PermCmdTempBan   PermissionNode = "commands.tempban"
	PermCmdMute      PermissionNode = "commands.mute"
	PermCmdWarn      PermissionNode = "commands.warn"
	PermCmdHistory   PermissionNode = "commands.history"
	PermCmdMoney     PermissionNode = "commands.money"
	PermCmdPay       PermissionNode = "commands.pay"
	PermCmdBalance   PermissionNode = "commands.balance"
	PermCmdShop      PermissionNode = "commands.shop"
	PermCmdAuction   PermissionNode = "commands.auction"
	PermCmdTrade     PermissionNode = "commands.trade"
	PermCmdClaim     PermissionNode = "commands.claim"
	PermCmdUnclaim   PermissionNode = "commands.unclaim"
	PermCmdTrust     PermissionNode = "commands.trust"
	PermCmdHome      PermissionNode = "commands.home"
	PermCmdSetHome   PermissionNode = "commands.sethome"
	PermCmdWarp      PermissionNode = "commands.warp"
	PermCmdMail      PermissionNode = "commands.mail"
	PermCmdParty     PermissionNode = "commands.party"
	PermCmdGuild     PermissionNode = "commands.guild"
	PermCmdQuest     PermissionNode = "commands.quest"
	PermCmdDuel      PermissionNode = "commands.duel"
	PermCmdBounty    PermissionNode = "commands.bounty"
)

// Admin permissions
const (
	PermAdminWorldSettings PermissionNode = "admin.world.settings"
	PermAdminWorldDelete   PermissionNode = "admin.world.delete"
	PermAdminWorldBackup   PermissionNode = "admin.world.backup"
	PermAdminPluginEnable  PermissionNode = "admin.plugins.enable"
	PermAdminPluginInstall PermissionNode = "admin.plugins.install"
	PermAdminPluginConfig  PermissionNode = "admin.plugins.config"
	PermAdminPlayerPromote PermissionNode = "admin.player.promote"
	PermAdminPlayerDemote  PermissionNode = "admin.player.demote"
	PermAdminPlayerReset   PermissionNode = "admin.player.reset"
	PermAdminEcoReset      PermissionNode = "admin.economy.reset"
	PermAdminEcoMint       PermissionNode = "admin.economy.mint"
	PermAdminEcoAdjust     PermissionNode = "admin.economy.adjust"
	PermAdminLandOverride  PermissionNode = "admin.land.override"
	PermAdminCheatBypass   PermissionNode = "admin.cheat.bypass"
	PermAdminVanish        PermissionNode = "admin.vanish"
	PermAdminBypassMute    PermissionNode = "admin.bypass.mute"
	PermAdminBypassKick    PermissionNode = "admin.bypass.kick"
	PermAdminBypassBan     PermissionNode = "admin.bypass.ban"
)

// Economy permissions
const (
	PermEcoReceive     PermissionNode = "economy.receive"
	PermEcoSpend       PermissionNode = "economy.spend"
	PermEcoTrade       PermissionNode = "economy.trade"
	PermEcoShopCreate  PermissionNode = "economy.shop.create"
	PermEcoShopDelete  PermissionNode = "economy.shop.delete"
	PermEcoShopAdmin   PermissionNode = "economy.shop.admin"
	PermEcoAuctionCreate PermissionNode = "economy.auction.create"
	PermEcoAuctionBid    PermissionNode = "economy.auction.bid"
	PermEcoBankWithdraw  PermissionNode = "economy.bank.withdraw"
	PermEcoBankDeposit   PermissionNode = "economy.bank.deposit"
	PermEcoBankLoan      PermissionNode = "economy.bank.loan"
	PermEcoJobJoin       PermissionNode = "economy.job.join"
	PermEcoJobAdmin      PermissionNode = "economy.job.admin"
)

// World access permissions
const (
	PermWorldAccess PermissionNode = "world.access"
	PermWorldBuild  PermissionNode = "world.build"
	PermWorldPvP    PermissionNode = "world.pvp"
	PermWorldEco    PermissionNode = "world.economy"
	PermWorldTrade  PermissionNode = "world.trade"
)

// Land permissions
const (
	PermLandClaim      PermissionNode = "land.claim"
	PermLandUnclaim    PermissionNode = "land.unclaim"
	PermLandExpand     PermissionNode = "land.expand"
	PermLandSubdivide  PermissionNode = "land.subdivide"
	PermLandTrustAdd   PermissionNode = "land.trust.add"
	PermLandTrustRemove PermissionNode = "land.trust.remove"
	PermLandFlagSet    PermissionNode = "land.flag.set"
)

// Social permissions
const (
	PermChatGlobal    PermissionNode = "social.chat.global"
	PermChatWhisper   PermissionNode = "social.chat.whisper"
	PermChatParty     PermissionNode = "social.chat.party"
	PermChatGuild     PermissionNode = "social.chat.guild"
	PermFriendAdd     PermissionNode = "social.friend.add"
	PermFriendRemove  PermissionNode = "social.friend.remove"
	PermFriendView    PermissionNode = "social.friend.view"
	PermPartyCreate   PermissionNode = "social.party.create"
	PermPartyJoin     PermissionNode = "social.party.join"
	PermPartyLead     PermissionNode = "social.party.lead"
	PermGuildCreate   PermissionNode = "social.guild.create"
	PermGuildJoin     PermissionNode = "social.guild.join"
	PermGuildManage   PermissionNode = "social.guild.manage"
	PermMailSend      PermissionNode = "social.mail.send"
	PermMailReceive   PermissionNode = "social.mail.receive"
)

// Minigame permissions
const (
	PermMinigameJoin   PermissionNode = "minigame.join"
	PermMinigameCreate PermissionNode = "minigame.create"
	PermMinigameAdmin  PermissionNode = "minigame.admin"
	PermMinigameSpleef PermissionNode = "minigame.spleef"
	PermMinigameParkour PermissionNode = "minigame.parkour"
	PermMinigameCTF    PermissionNode = "minigame.ctf"
	PermMinigameTNTRun PermissionNode = "minigame.tntrun"
	PermMinigameMobArena PermissionNode = "minigame.mobarena"
)

// Cosmetic permissions
const (
	PermCosmeticParticles PermissionNode = "cosmetic.particles"
	PermCosmeticPets      PermissionNode = "cosmetic.pets"
	PermCosmeticTitles    PermissionNode = "cosmetic.titles"
	PermCosmeticHats      PermissionNode = "cosmetic.hats"
	PermCosmeticEmotes    PermissionNode = "cosmetic.emotes"
	PermCosmeticTrails    PermissionNode = "cosmetic.trails"
)

// AllNodes returns all defined permission nodes
func AllNodes() []PermissionNode {
	return []PermissionNode{
		// Build
		PermBuildPlace, PermBuildBreak, PermBuildInteract, PermBuildAdmin,
		// Game mode
		PermGameModeCreative, PermGameModeSurvival, PermGameModeSpectator, PermGameModeSwitch,
		// Commands
		PermCmdGive, PermCmdTP, PermCmdTPA, PermCmdPlugin,
		PermCmdKick, PermCmdBan, PermCmdTempBan, PermCmdMute, PermCmdWarn, PermCmdHistory,
		PermCmdMoney, PermCmdPay, PermCmdBalance, PermCmdShop, PermCmdAuction, PermCmdTrade,
		PermCmdClaim, PermCmdUnclaim, PermCmdTrust,
		PermCmdHome, PermCmdSetHome, PermCmdWarp, PermCmdMail,
		PermCmdParty, PermCmdGuild, PermCmdQuest, PermCmdDuel, PermCmdBounty,
		// Admin
		PermAdminWorldSettings, PermAdminWorldDelete, PermAdminWorldBackup,
		PermAdminPluginEnable, PermAdminPluginInstall, PermAdminPluginConfig,
		PermAdminPlayerPromote, PermAdminPlayerDemote, PermAdminPlayerReset,
		PermAdminEcoReset, PermAdminEcoMint, PermAdminEcoAdjust,
		PermAdminLandOverride, PermAdminCheatBypass, PermAdminVanish,
		PermAdminBypassMute, PermAdminBypassKick, PermAdminBypassBan,
		// Economy
		PermEcoReceive, PermEcoSpend, PermEcoTrade,
		PermEcoShopCreate, PermEcoShopDelete, PermEcoShopAdmin,
		PermEcoAuctionCreate, PermEcoAuctionBid,
		PermEcoBankWithdraw, PermEcoBankDeposit, PermEcoBankLoan,
		PermEcoJobJoin, PermEcoJobAdmin,
		// World
		PermWorldAccess, PermWorldBuild, PermWorldPvP, PermWorldEco, PermWorldTrade,
		// Land
		PermLandClaim, PermLandUnclaim, PermLandExpand, PermLandSubdivide,
		PermLandTrustAdd, PermLandTrustRemove, PermLandFlagSet,
		// Social
		PermChatGlobal, PermChatWhisper, PermChatParty, PermChatGuild,
		PermFriendAdd, PermFriendRemove, PermFriendView,
		PermPartyCreate, PermPartyJoin, PermPartyLead,
		PermGuildCreate, PermGuildJoin, PermGuildManage,
		PermMailSend, PermMailReceive,
		// Minigames
		PermMinigameJoin, PermMinigameCreate, PermMinigameAdmin,
		PermMinigameSpleef, PermMinigameParkour, PermMinigameCTF,
		PermMinigameTNTRun, PermMinigameMobArena,
		// Cosmetics
		PermCosmeticParticles, PermCosmeticPets, PermCosmeticTitles,
		PermCosmeticHats, PermCosmeticEmotes, PermCosmeticTrails,
	}
}

// IsValid checks if a permission node is valid
func IsValid(node PermissionNode) bool {
	for _, n := range AllNodes() {
		if n == node {
			return true
		}
	}
	return false
}
