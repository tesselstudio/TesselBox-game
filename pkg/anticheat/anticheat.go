package anticheat

import (
	"fmt"
	"math"
	"time"
)

// ViolationType represents the type of cheat violation
type ViolationType int

const (
	ViolationSpeed ViolationType = iota
	ViolationFly
	ViolationReach
	ViolationFastBreak
	ViolationFastPlace
	ViolationSpam
	ViolationCommandSpam
	ViolationAimbot
)

// String returns human-readable violation name
func (v ViolationType) String() string {
	switch v {
	case ViolationSpeed:
		return "Speed Hack"
	case ViolationFly:
		return "Fly Hack"
	case ViolationReach:
		return "Reach Hack"
	case ViolationFastBreak:
		return "Fast Break"
	case ViolationFastPlace:
		return "Fast Place"
	case ViolationSpam:
		return "Chat Spam"
	case ViolationCommandSpam:
		return "Command Spam"
	case ViolationAimbot:
		return "Aimbot"
	}
	return "Unknown"
}

// ViolationLevel represents severity level
type ViolationLevel int

const (
	LevelNotice ViolationLevel = iota    // Log only
	LevelWarning                         // Alert moderators
	LevelKick                            // Kick player
	LevelTempBan                         // Temporary ban
	LevelBan                             // Permanent ban
)

// Violation represents a single violation
type Violation struct {
	ID          string         `json:"id"`
	PlayerID    string         `json:"player_id"`
	Type        ViolationType  `json:"type"`
	Level       ViolationLevel `json:"level"`
	Description string         `json:"description"`
	Evidence    string         `json:"evidence"` // Data supporting the violation
	Timestamp   time.Time      `json:"timestamp"`
	Location    Position       `json:"location,omitempty"`
}

// Position represents a 3D position
type Position struct {
	X, Y  float64
	World string
}

// PlayerACData stores anti-cheat data for a player
type PlayerACData struct {
	PlayerID         string
	LastUpdate       time.Time
	
	// Position history for speed check
	PositionHistory  []TimedPosition
	VelocityHistory  []VelocitySample
	
	// Block interaction tracking
	BlocksPlaced     int
	BlocksBroken     int
	LastBlockPlace   time.Time
	LastBlockBreak   time.Time
	
	// Chat/Command tracking
	MessagesSent     int
	CommandsSent     int
	LastMessage      time.Time
	LastCommand      time.Time
	MessageHistory   []string // Recent messages for spam detection
	
	// Combat tracking
	AttacksMade      int
	LastAttack       time.Time
	AttackTargets    []string
	
	// Trust score (0-100, higher = more trusted)
	TrustScore       float64
	ViolationCount   int
	LastViolation    time.Time
}

// TimedPosition represents a position at a specific time
type TimedPosition struct {
	X, Y   float64
	Time   time.Time
	OnGround bool
}

// VelocitySample represents velocity at a specific time
type VelocitySample struct {
	VX, VY float64
	Time   time.Time
}

// NewPlayerACData creates new anti-cheat data for a player
func NewPlayerACData(playerID string) *PlayerACData {
	now := time.Now()
	return &PlayerACData{
		PlayerID:        playerID,
		LastUpdate:      now,
		PositionHistory: make([]TimedPosition, 0),
		VelocityHistory: make([]VelocitySample, 0),
		TrustScore:      100.0,
		LastMessage:     now,
		LastCommand:     now,
		LastBlockPlace:  now,
		LastBlockBreak:  now,
		LastAttack:      now,
	}
}

// UpdatePosition updates player position
func (pac *PlayerACData) UpdatePosition(x, y float64, onGround bool) {
	now := time.Now()
	
	// Calculate velocity
	var vx, vy float64
	if len(pac.PositionHistory) > 0 {
		last := pac.PositionHistory[len(pac.PositionHistory)-1]
		dt := now.Sub(last.Time).Seconds()
		if dt > 0 {
			vx = (x - last.X) / dt
			vy = (y - last.Y) / dt
		}
	}
	
	// Add to history
	pac.PositionHistory = append(pac.PositionHistory, TimedPosition{
		X:        x,
		Y:        y,
		Time:     now,
		OnGround: onGround,
	})
	
	// Add velocity to history
	pac.VelocityHistory = append(pac.VelocityHistory, VelocitySample{
		VX:   vx,
		VY:   vy,
		Time: now,
	})
	
	// Trim histories
	if len(pac.PositionHistory) > 100 {
		pac.PositionHistory = pac.PositionHistory[len(pac.PositionHistory)-100:]
	}
	if len(pac.VelocityHistory) > 100 {
		pac.VelocityHistory = pac.VelocityHistory[len(pac.VelocityHistory)-100:]
	}
	
	pac.LastUpdate = now
}

// RecordBlockPlace records a block placement
func (pac *PlayerACData) RecordBlockPlace() {
	pac.BlocksPlaced++
	pac.LastBlockPlace = time.Now()
}

// RecordBlockBreak records a block break
func (pac *PlayerACData) RecordBlockBreak() {
	pac.BlocksBroken++
	pac.LastBlockBreak = time.Now()
}

// RecordMessage records a chat message
func (pac *PlayerACData) RecordMessage(msg string) {
	pac.MessagesSent++
	pac.LastMessage = time.Now()
	pac.MessageHistory = append(pac.MessageHistory, msg)
	
	// Keep last 10 messages
	if len(pac.MessageHistory) > 10 {
		pac.MessageHistory = pac.MessageHistory[len(pac.MessageHistory)-10:]
	}
}

// RecordCommand records a command
func (pac *PlayerACData) RecordCommand() {
	pac.CommandsSent++
	pac.LastCommand = time.Now()
}

// RecordAttack records an attack
func (pac *PlayerACData) RecordAttack(targetID string) {
	pac.AttacksMade++
	pac.LastAttack = time.Now()
	pac.AttackTargets = append(pac.AttackTargets, targetID)
	
	// Keep last 20 targets
	if len(pac.AttackTargets) > 20 {
		pac.AttackTargets = pac.AttackTargets[len(pac.AttackTargets)-20:]
	}
}

// AddViolation reduces trust score
func (pac *PlayerACData) AddViolation() {
	pac.ViolationCount++
	pac.LastViolation = time.Now()
	
	// Reduce trust score
	pac.TrustScore = pac.TrustScore * 0.9 // Lose 10% per violation
	if pac.TrustScore < 0 {
		pac.TrustScore = 0
	}
}

// RecoverTrust slowly recovers trust over time
func (pac *PlayerACData) RecoverTrust() {
	// Recover 1 point per hour without violations
	if time.Since(pac.LastViolation) > time.Hour {
		pac.TrustScore = math.Min(pac.TrustScore+1, 100)
	}
}

// GetBlockPlaceRate returns blocks placed per second
func (pac *PlayerACData) GetBlockPlaceRate(window time.Duration) float64 {
	if time.Since(pac.LastBlockPlace) > window {
		return 0
	}
	return float64(pac.BlocksPlaced) / window.Seconds()
}

// GetBlockBreakRate returns blocks broken per second
func (pac *PlayerACData) GetBlockBreakRate(window time.Duration) float64 {
	if time.Since(pac.LastBlockBreak) > window {
		return 0
	}
	return float64(pac.BlocksBroken) / window.Seconds()
}

// GetMessageRate returns messages per second
func (pac *PlayerACData) GetMessageRate(window time.Duration) float64 {
	if time.Since(pac.LastMessage) > window {
		return 0
	}
	return float64(pac.MessagesSent) / window.Seconds()
}

// GetCommandRate returns commands per second
func (pac *PlayerACData) GetCommandRate(window time.Duration) float64 {
	if time.Since(pac.LastCommand) > window {
		return 0
	}
	return float64(pac.CommandsSent) / window.Seconds()
}

// GetMaxHorizontalSpeed returns max horizontal speed in the last window
func (pac *PlayerACData) GetMaxHorizontalSpeed(window time.Duration) float64 {
	maxSpeed := 0.0
	cutoff := time.Now().Add(-window)
	
	for _, vel := range pac.VelocityHistory {
		if vel.Time.After(cutoff) {
			speed := math.Sqrt(vel.VX*vel.VX + vel.VY*vel.VY)
			if speed > maxSpeed {
				maxSpeed = speed
			}
		}
	}
	
	return maxSpeed
}

// ACRules contains anti-cheat detection thresholds
type ACRules struct {
	// Movement
	MaxSpeedGround    float64 // Max horizontal speed on ground
	MaxSpeedAir       float64 // Max horizontal speed in air
	MaxSpeedSprint    float64 // Max sprint speed
	MaxVerticalSpeed  float64 // Max vertical speed (jump/fall)
	MaxJumpHeight     float64 // Max jump height
	
	// Block interaction
	MaxBlocksPerSec   int     // Max blocks placed/broken per second
	MaxReachDistance  float64 // Max interaction distance
	
	// Chat/Commands
	MaxMessagesPerSec int     // Max chat messages per second
	MaxCommandsPerSec int     // Max commands per second
	MaxMessageLength  int     // Max message length
	MaxCapsPercent    float64 // Max percentage of caps in message
	
	// Combat
	MaxAttacksPerSec  int     // Max attacks per second
	MaxReachAttack    float64 // Max attack reach
	MaxCPS            int     // Max clicks per second
}

// DefaultACRules returns default anti-cheat rules
func DefaultACRules() ACRules {
	return ACRules{
		MaxSpeedGround:    50.0,  // blocks per second
		MaxSpeedAir:       40.0,
		MaxSpeedSprint:    66.0,
		MaxVerticalSpeed:  78.4,  // Terminal velocity
		MaxJumpHeight:     2.5,
		MaxBlocksPerSec:   20,
		MaxReachDistance:  6.0,
		MaxMessagesPerSec: 5,
		MaxCommandsPerSec: 10,
		MaxMessageLength:  200,
		MaxCapsPercent:    0.7,   // 70%
		MaxAttacksPerSec:  15,
		MaxReachAttack:    6.0,
		MaxCPS:            20,
	}
}

// AntiCheat is the main anti-cheat manager
type AntiCheat struct {
	players     map[string]*PlayerACData
	rules       ACRules
	violations  []Violation
	
	// Callbacks
	OnViolation func(v *Violation)
	OnKick      func(playerID string, reason string)
	OnBan       func(playerID string, duration time.Duration, reason string)
}

// NewAntiCheat creates a new anti-cheat system
func NewAntiCheat() *AntiCheat {
	return &AntiCheat{
		players:    make(map[string]*PlayerACData),
		rules:      DefaultACRules(),
		violations: make([]Violation, 0),
	}
}

// SetRules sets anti-cheat rules
func (ac *AntiCheat) SetRules(rules ACRules) {
	ac.rules = rules
}

// GetPlayerData gets or creates player AC data
func (ac *AntiCheat) GetPlayerData(playerID string) *PlayerACData {
	if data, exists := ac.players[playerID]; exists {
		return data
	}
	
	data := NewPlayerACData(playerID)
	ac.players[playerID] = data
	return data
}

// UpdatePlayer updates player position data
func (ac *AntiCheat) UpdatePlayer(playerID string, x, y float64, onGround bool) {
	data := ac.GetPlayerData(playerID)
	data.UpdatePosition(x, y, onGround)
	data.RecoverTrust()
	
	// Check for speed violations
	ac.checkSpeed(playerID, data)
	
	// Check for fly violations
	ac.checkFly(playerID, data)
}

// checkSpeed checks for speed hacks
func (ac *AntiCheat) checkSpeed(playerID string, data *PlayerACData) {
	maxSpeed := data.GetMaxHorizontalSpeed(1 * time.Second)
	
	var limit float64
	// Determine appropriate limit based on recent movement
	recent := data.PositionHistory
	if len(recent) > 0 && recent[len(recent)-1].OnGround {
		limit = ac.rules.MaxSpeedGround
	} else {
		limit = ac.rules.MaxSpeedAir
	}
	
	if maxSpeed > limit*1.5 { // 50% tolerance for lag
		ac.recordViolation(playerID, ViolationSpeed, LevelWarning, 
			fmt.Sprintf("Speed: %.2f blocks/sec (limit: %.2f)", maxSpeed, limit),
			fmt.Sprintf("max_speed=%.2f, limit=%.2f", maxSpeed, limit))
	}
}

// checkFly checks for fly hacks
func (ac *AntiCheat) checkFly(playerID string, data *PlayerACData) {
	if len(data.PositionHistory) < 2 {
		return
	}
	
	// Check for sustained vertical movement without being on ground
	// This is simplified - real fly detection needs more sophisticated logic
	recent := data.PositionHistory[len(data.PositionHistory)-10:]
	if len(recent) < 5 {
		return
	}
	
	gainedHeight := 0.0
	lastOnGround := false
	
	for _, pos := range recent {
		if pos.OnGround {
			lastOnGround = true
			gainedHeight = 0
			continue
		}
		
		if lastOnGround {
			// Check if gained height without jumping (simplified)
			if gainedHeight > ac.rules.MaxJumpHeight && !pos.OnGround {
				ac.recordViolation(playerID, ViolationFly, LevelKick,
					fmt.Sprintf("Fly: Gained %.2f blocks without ground contact", gainedHeight),
					fmt.Sprintf("height_gained=%.2f", gainedHeight))
				return
			}
		}
	}
}

// CheckBlockPlace checks block placement for fast place
func (ac *AntiCheat) CheckBlockPlace(playerID string, x, y float64) {
	data := ac.GetPlayerData(playerID)
	data.RecordBlockPlace()
	
	rate := data.GetBlockPlaceRate(1 * time.Second)
	if rate > float64(ac.rules.MaxBlocksPerSec) {
		ac.recordViolation(playerID, ViolationFastPlace, LevelWarning,
			fmt.Sprintf("Fast place: %.2f blocks/sec (limit: %d)", rate, ac.rules.MaxBlocksPerSec),
			fmt.Sprintf("rate=%.2f, limit=%d", rate, ac.rules.MaxBlocksPerSec))
	}
}

// CheckBlockBreak checks block breaking for fast break
func (ac *AntiCheat) CheckBlockBreak(playerID string, x, y float64) {
	data := ac.GetPlayerData(playerID)
	data.RecordBlockBreak()
	
	rate := data.GetBlockBreakRate(1 * time.Second)
	if rate > float64(ac.rules.MaxBlocksPerSec) {
		ac.recordViolation(playerID, ViolationFastBreak, LevelWarning,
			fmt.Sprintf("Fast break: %.2f blocks/sec (limit: %d)", rate, ac.rules.MaxBlocksPerSec),
			fmt.Sprintf("rate=%.2f, limit=%d", rate, ac.rules.MaxBlocksPerSec))
	}
}

// CheckMessage checks chat message for spam
func (ac *AntiCheat) CheckMessage(playerID string, message string) {
	data := ac.GetPlayerData(playerID)
	data.RecordMessage(message)
	
	// Check rate
	rate := data.GetMessageRate(5 * time.Second)
	if rate > float64(ac.rules.MaxMessagesPerSec) {
		ac.recordViolation(playerID, ViolationSpam, LevelWarning,
			fmt.Sprintf("Chat spam: %.2f msg/sec (limit: %d)", rate, ac.rules.MaxMessagesPerSec),
			fmt.Sprintf("rate=%.2f", rate))
		return
	}
	
	// Check length
	if len(message) > ac.rules.MaxMessageLength {
		ac.recordViolation(playerID, ViolationSpam, LevelNotice,
			fmt.Sprintf("Oversized message: %d chars (limit: %d)", len(message), ac.rules.MaxMessageLength),
			fmt.Sprintf("length=%d", len(message)))
	}
	
	// Check caps
	capsCount := 0
	for _, c := range message {
		if c >= 'A' && c <= 'Z' {
			capsCount++
		}
	}
	if len(message) > 5 {
		capsPercent := float64(capsCount) / float64(len(message))
		if capsPercent > ac.rules.MaxCapsPercent {
			ac.recordViolation(playerID, ViolationSpam, LevelNotice,
				fmt.Sprintf("CAPS spam: %.0f%% caps", capsPercent*100),
				fmt.Sprintf("caps_percent=%.2f", capsPercent))
		}
	}
}

// CheckCommand checks command for spam
func (ac *AntiCheat) CheckCommand(playerID string, command string) {
	data := ac.GetPlayerData(playerID)
	data.RecordCommand()
	
	rate := data.GetCommandRate(5 * time.Second)
	if rate > float64(ac.rules.MaxCommandsPerSec) {
		ac.recordViolation(playerID, ViolationCommandSpam, LevelWarning,
			fmt.Sprintf("Command spam: %.2f cmd/sec (limit: %d)", rate, ac.rules.MaxCommandsPerSec),
			fmt.Sprintf("rate=%.2f", rate))
	}
}

// CheckAttack checks combat for hacks
func (ac *AntiCheat) CheckAttack(playerID string, targetID string, distance float64) {
	data := ac.GetPlayerData(playerID)
	data.RecordAttack(targetID)
	
	// Check reach
	if distance > ac.rules.MaxReachAttack {
		ac.recordViolation(playerID, ViolationReach, LevelWarning,
			fmt.Sprintf("Reach: %.2f blocks (limit: %.2f)", distance, ac.rules.MaxReachAttack),
			fmt.Sprintf("distance=%.2f", distance))
	}
	
	// Check rate
	if data.AttacksMade > ac.rules.MaxAttacksPerSec {
		ac.recordViolation(playerID, ViolationAimbot, LevelWarning,
			fmt.Sprintf("Attack rate: %d attacks (limit: %d/sec)", data.AttacksMade, ac.rules.MaxAttacksPerSec),
			fmt.Sprintf("count=%d", data.AttacksMade))
	}
}

// recordViolation records a violation and determines response
func (ac *AntiCheat) recordViolation(playerID string, vType ViolationType, baseLevel ViolationLevel, description, evidence string) {
	data := ac.GetPlayerData(playerID)
	
	// Adjust level based on trust score
	level := ac.calculateViolationLevel(data.TrustScore, baseLevel)
	
	violation := Violation{
		ID:          fmt.Sprintf("v_%d", time.Now().UnixNano()),
		PlayerID:    playerID,
		Type:        vType,
		Level:       level,
		Description: description,
		Evidence:    evidence,
		Timestamp:   time.Now(),
	}
	
	ac.violations = append(ac.violations, violation)
	data.AddViolation()
	
	// Trigger callbacks
	if ac.OnViolation != nil {
		ac.OnViolation(&violation)
	}
	
	// Execute punishment
	ac.executePunishment(playerID, level, vType.String())
}

// calculateViolationLevel determines final level based on trust
func (ac *AntiCheat) calculateViolationLevel(trust float64, base ViolationLevel) ViolationLevel {
	if trust > 90 && base == LevelWarning {
		return LevelNotice
	}
	if trust < 50 && base < LevelKick {
		return LevelKick
	}
	if trust < 20 && base < LevelTempBan {
		return LevelTempBan
	}
	if trust < 5 {
		return LevelBan
	}
	return base
}

// executePunishment executes the appropriate punishment
func (ac *AntiCheat) executePunishment(playerID string, level ViolationLevel, reason string) {
	switch level {
	case LevelKick:
		if ac.OnKick != nil {
			ac.OnKick(playerID, reason)
		}
	case LevelTempBan:
		if ac.OnBan != nil {
			ac.OnBan(playerID, 1*time.Hour, reason)
		}
	case LevelBan:
		if ac.OnBan != nil {
			ac.OnBan(playerID, 0, reason) // 0 = permanent
		}
	}
}

// GetViolations returns all violations for a player
func (ac *AntiCheat) GetViolations(playerID string) []Violation {
	result := make([]Violation, 0)
	for _, v := range ac.violations {
		if v.PlayerID == playerID {
			result = append(result, v)
		}
	}
	return result
}

// GetRecentViolations returns recent violations
func (ac *AntiCheat) GetRecentViolations(count int) []Violation {
	if count > len(ac.violations) {
		count = len(ac.violations)
	}
	start := len(ac.violations) - count
	if start < 0 {
		start = 0
	}
	return ac.violations[start:]
}

// ClearPlayerData removes player data (on logout)
func (ac *AntiCheat) ClearPlayerData(playerID string) {
	delete(ac.players, playerID)
}

// GetPlayerTrustScore returns a player's trust score
func (ac *AntiCheat) GetPlayerTrustScore(playerID string) float64 {
	if data, exists := ac.players[playerID]; exists {
		return data.TrustScore
	}
	return 100 // Default trust for new players
}
