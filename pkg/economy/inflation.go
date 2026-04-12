package economy

import (
	"fmt"
	"math"
	"time"
)

// EconomicHealth represents the current state of the economy
type EconomicHealth int

const (
	EconomyHealthy EconomicHealth = iota    // Balanced
	EconomyInflation                        // Too much money
	EconomyDeflation                        // Too little money
	EconomyCrisis                           // Severe imbalance
)

// MonetaryPolicy contains policy targets
type MonetaryPolicy struct {
	// Money targets
	TargetMoneyPerPlayer float64  // Ideal money per active player (default: 1000)
	MinHealthyRatio        float64 // 0.5 - below this is deflation
	MaxHealthyRatio        float64 // 1.5 - above this is inflation
	
	// Rate caps
	MaxInflationRate       float64 // 0.10 = 10% daily max
	MaxDeflationRate       float64 // 0.05 = 5% daily max
	
	// Adjustment factors
	InflationAdjustment    float64 // How fast to adjust (0.1 = slow, 0.5 = fast)
	DeflationAdjustment    float64
}

// DefaultMonetaryPolicy returns default policy settings
func DefaultMonetaryPolicy() MonetaryPolicy {
	return MonetaryPolicy{
		TargetMoneyPerPlayer: 1000.0,
		MinHealthyRatio:      0.5,
		MaxHealthyRatio:      1.5,
		MaxInflationRate:     0.10,
		MaxDeflationRate:     0.05,
		InflationAdjustment:  0.2,
		DeflationAdjustment:  0.1,
	}
}

// EconomyEngine manages the dynamic economy
type EconomyEngine struct {
	// State
	TotalCurrency      float64        // All money in circulation
	ActivePlayers      int            // Players online in last 24h
	MoneyVelocity      float64        // Average transactions per player per day
	
	// Calculated rates
	InflationRate      float64        // Current inflation (0.0 - 0.10)
	DeflationRate      float64        // Current deflation (0.0 - 0.05)
	PriceMultiplier    float64        // Applied to all prices (0.95 - 1.10)
	
	// Health
	HealthStatus       EconomicHealth
	HealthScore        int            // 0-100
	
	// Policy
	Policy             MonetaryPolicy
	
	// History
	History            []EconomicSnapshot
	maxHistoryDays     int
	
	// Last update
	LastCalculation    time.Time
	
	// Sink/Source modifiers (adjusted based on economy)
	DeathPenaltyMod    float64        // Modifier for death penalty
	ShopTaxMod         float64        // Modifier for shop tax
	MiningRewardMod    float64        // Modifier for mining rewards
	MobBountyMod       float64        // Modifier for mob kills
	
	// Callbacks
	OnInflationAlert   func(rate float64)
	OnDeflationAlert   func(rate float64)
	OnRecessionAlert   func()
}

// EconomicSnapshot represents economy state at a point in time
type EconomicSnapshot struct {
	Timestamp       time.Time      `json:"timestamp"`
	TotalCurrency   float64        `json:"total_currency"`
	ActivePlayers   int            `json:"active_players"`
	InflationRate   float64        `json:"inflation_rate"`
	DeflationRate   float64        `json:"deflation_rate"`
	PriceMultiplier float64        `json:"price_multiplier"`
	HealthStatus    EconomicHealth `json:"health_status"`
	HealthScore     int            `json:"health_score"`
}

// NewEconomyEngine creates a new economy engine
func NewEconomyEngine() *EconomyEngine {
	return &EconomyEngine{
		TotalCurrency:      0,
		ActivePlayers:      1,
		MoneyVelocity:      0,
		InflationRate:      0,
		DeflationRate:      0,
		PriceMultiplier:    1.0,
		HealthStatus:       EconomyHealthy,
		HealthScore:        100,
		Policy:             DefaultMonetaryPolicy(),
		History:            make([]EconomicSnapshot, 0),
		maxHistoryDays:     30,
		LastCalculation:    time.Now(),
		DeathPenaltyMod:    1.0,
		ShopTaxMod:         1.0,
		MiningRewardMod:    1.0,
		MobBountyMod:       1.0,
	}
}

// UpdateState updates the economy state with current data
func (e *EconomyEngine) UpdateState(totalCurrency float64, activePlayers int, velocity float64) {
	e.TotalCurrency = totalCurrency
	e.ActivePlayers = activePlayers
	if e.ActivePlayers < 1 {
		e.ActivePlayers = 1 // Prevent division by zero
	}
	e.MoneyVelocity = velocity
}

// Calculate runs the economic calculation
func (e *EconomyEngine) Calculate() {
	now := time.Now()
	
	// Only recalculate every hour (or force with CalculateNow)
	if now.Sub(e.LastCalculation) < time.Hour {
		return
	}
	
	e.CalculateNow()
}

// CalculateNow forces immediate calculation
func (e *EconomyEngine) CalculateNow() {
	// Calculate ideal money supply
	idealTotal := float64(e.ActivePlayers) * e.Policy.TargetMoneyPerPlayer
	
	// Calculate ratio
	ratio := e.TotalCurrency / idealTotal
	
	// Determine economic state and calculate rates
	if ratio > e.Policy.MaxHealthyRatio {
		// INFLATION - too much money
		e.calculateInflation(ratio, idealTotal)
	} else if ratio < e.Policy.MinHealthyRatio {
		// DEFLATION - too little money
		e.calculateDeflation(ratio, idealTotal)
	} else {
		// HEALTHY - normalize toward 1.0
		e.normalize(ratio)
	}
	
	// Calculate health score
	e.calculateHealthScore(ratio)
	
	// Record snapshot
	e.recordSnapshot()
	
	// Trim history
	e.trimHistory()
	
	e.LastCalculation = time.Now()
}

// calculateInflation handles inflation state
func (e *EconomyEngine) calculateInflation(ratio, idealTotal float64) {
	e.HealthStatus = EconomyInflation
	
	// Calculate excess money
	excess := e.TotalCurrency - (idealTotal * e.Policy.MaxHealthyRatio)
	
	// Calculate inflation rate (capped)
	rawInflation := (excess / idealTotal) * e.Policy.InflationAdjustment
	e.InflationRate = math.Min(rawInflation, e.Policy.MaxInflationRate)
	e.DeflationRate = 0
	
	// Price multiplier (things cost more)
	e.PriceMultiplier = 1.0 + e.InflationRate
	
	// Adjust sinks (increase to remove money)
	e.DeathPenaltyMod = 1.0 + e.InflationRate
	e.ShopTaxMod = 1.0 + (e.InflationRate * 0.5)
	
	// Adjust sources (decrease to reduce money injection)
	e.MiningRewardMod = 1.0 - (e.InflationRate * 0.5)
	e.MobBountyMod = 1.0 - (e.InflationRate * 0.5)
	
	// Alert if severe
	if e.InflationRate > 0.05 && e.OnInflationAlert != nil {
		e.OnInflationAlert(e.InflationRate)
	}
}

// calculateDeflation handles deflation state
func (e *EconomyEngine) calculateDeflation(ratio, idealTotal float64) {
	e.HealthStatus = EconomyDeflation
	
	// Calculate shortage
	shortage := (idealTotal * e.Policy.MinHealthyRatio) - e.TotalCurrency
	
	// Calculate deflation rate (capped)
	rawDeflation := (shortage / idealTotal) * e.Policy.DeflationAdjustment
	e.DeflationRate = math.Min(rawDeflation, e.Policy.MaxDeflationRate)
	e.InflationRate = 0
	
	// Price multiplier (things cost less)
	e.PriceMultiplier = 1.0 - e.DeflationRate
	
	// Adjust sinks (decrease)
	e.DeathPenaltyMod = 1.0 - e.DeflationRate
	e.ShopTaxMod = 1.0 - (e.DeflationRate * 0.5)
	
	// Adjust sources (increase)
	e.MiningRewardMod = 1.0 + (e.DeflationRate * 2.0)
	e.MobBountyMod = 1.0 + (e.DeflationRate * 2.0)
	
	// Alert if severe
	if e.DeflationRate > 0.03 && e.OnDeflationAlert != nil {
		e.OnDeflationAlert(e.DeflationRate)
	}
}

// normalize brings economy toward balance
func (e *EconomyEngine) normalize(ratio float64) {
	e.HealthStatus = EconomyHealthy
	
	// Gradually normalize rates
	e.InflationRate = e.InflationRate * 0.8
	e.DeflationRate = e.DeflationRate * 0.8
	
	// Normalize price multiplier toward 1.0
	if e.PriceMultiplier > 1.0 {
		e.PriceMultiplier = 1.0 + (e.PriceMultiplier-1.0)*0.5
	} else if e.PriceMultiplier < 1.0 {
		e.PriceMultiplier = 1.0 - (1.0-e.PriceMultiplier)*0.5
	}
	
	// Normalize modifiers
	e.DeathPenaltyMod = 1.0 + (e.DeathPenaltyMod-1.0)*0.5
	e.ShopTaxMod = 1.0 + (e.ShopTaxMod-1.0)*0.5
	e.MiningRewardMod = 1.0 + (e.MiningRewardMod-1.0)*0.5
	e.MobBountyMod = 1.0 + (e.MobBountyMod-1.0)*0.5
}

// calculateHealthScore calculates 0-100 health score
func (e *EconomyEngine) calculateHealthScore(ratio float64) {
	// Ideal ratio is 1.0
	distance := math.Abs(ratio - 1.0)
	
	// Score decreases as distance from ideal increases
	score := 100 - int(distance*100)
	
	// Penalty for velocity (too fast = instability)
	if e.MoneyVelocity > 100 {
		score -= int((e.MoneyVelocity - 100) / 10)
	}
	
	// Clamp
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}
	
	e.HealthScore = score
	
	// Check for crisis
	if score < 20 && e.HealthStatus != EconomyHealthy {
		e.HealthStatus = EconomyCrisis
		if e.OnRecessionAlert != nil {
			e.OnRecessionAlert()
		}
	}
}

// recordSnapshot records current state
func (e *EconomyEngine) recordSnapshot() {
	snapshot := EconomicSnapshot{
		Timestamp:       time.Now(),
		TotalCurrency:   e.TotalCurrency,
		ActivePlayers:   e.ActivePlayers,
		InflationRate:   e.InflationRate,
		DeflationRate:   e.DeflationRate,
		PriceMultiplier: e.PriceMultiplier,
		HealthStatus:    e.HealthStatus,
		HealthScore:     e.HealthScore,
	}
	
	e.History = append(e.History, snapshot)
}

// trimHistory removes old snapshots
func (e *EconomyEngine) trimHistory() {
	maxEntries := e.maxHistoryDays * 24 // Hourly snapshots
	if len(e.History) > maxEntries {
		e.History = e.History[len(e.History)-maxEntries:]
	}
}

// GetEffectivePrice applies price multiplier to a base price
func (e *EconomyEngine) GetEffectivePrice(basePrice float64) float64 {
	return basePrice * e.PriceMultiplier
}

// GetEffectiveMiningReward applies reward modifier
func (e *EconomyEngine) GetEffectiveMiningReward(baseReward float64) float64 {
	return baseReward * e.MiningRewardMod
}

// GetEffectiveMobBounty applies bounty modifier
func (e *EconomyEngine) GetEffectiveMobBounty(baseBounty float64) float64 {
	return baseBounty * e.MobBountyMod
}

// GetEffectiveDeathPenalty applies penalty modifier
func (e *EconomyEngine) GetEffectiveDeathPenalty(basePenalty float64) float64 {
	return basePenalty * e.DeathPenaltyMod
}

// GetEffectiveShopTax applies tax modifier
func (e *EconomyEngine) GetEffectiveShopTax(baseTax float64) float64 {
	return baseTax * e.ShopTaxMod
}

// GetHealthColor returns a color code for the health status
func (e *EconomyEngine) GetHealthColor() string {
	switch e.HealthStatus {
	case EconomyHealthy:
		return "#00FF00" // Green
	case EconomyInflation:
		if e.InflationRate > 0.05 {
			return "#FF0000" // Red (severe)
		}
		return "#FFFF00" // Yellow (mild)
	case EconomyDeflation:
		if e.DeflationRate > 0.03 {
			return "#FF0000" // Red (severe)
		}
		return "#FFFF00" // Yellow (mild)
	case EconomyCrisis:
		return "#FF0000" // Red
	}
	return "#FFFFFF"
}

// GetHealthDescription returns a human-readable description
func (e *EconomyEngine) GetHealthDescription() string {
	switch e.HealthStatus {
	case EconomyHealthy:
		return fmt.Sprintf("Healthy (%d/100)", e.HealthScore)
	case EconomyInflation:
		return fmt.Sprintf("Inflation %.1f%% (+%.1f%% prices)", 
			e.InflationRate*100, (e.PriceMultiplier-1.0)*100)
	case EconomyDeflation:
		return fmt.Sprintf("Deflation %.1f%% (-%.1f%% prices)", 
			e.DeflationRate*100, (1.0-e.PriceMultiplier)*100)
	case EconomyCrisis:
		return "ECONOMIC CRISIS - Intervention needed"
	}
	return "Unknown"
}

// GetTrend returns the economic trend based on history
func (e *EconomyEngine) GetTrend() string {
	if len(e.History) < 2 {
		return "stable"
	}
	
	// Compare last 24h to previous 24h
	recent := e.History[len(e.History)-24:]
	if len(recent) < 24 {
		recent = e.History
	}
	
	previous := e.History[len(e.History)-48:len(e.History)-24]
	if len(previous) < 24 {
		return "stable"
	}
	
	recentAvg := 0.0
	for _, s := range recent {
		recentAvg += s.PriceMultiplier
	}
	recentAvg /= float64(len(recent))
	
	previousAvg := 0.0
	for _, s := range previous {
		previousAvg += s.PriceMultiplier
	}
	previousAvg /= float64(len(previous))
	
	diff := recentAvg - previousAvg
	
	if diff > 0.02 {
		return "rising"
	} else if diff < -0.02 {
		return "falling"
	}
	return "stable"
}

// GetHistory returns economic history
func (e *EconomyEngine) GetHistory(days int) []EconomicSnapshot {
	hours := days * 24
	if hours > len(e.History) {
		return e.History
	}
	return e.History[len(e.History)-hours:]
}

// InjectCurrency adds money to economy (for stimulus)
func (e *EconomyEngine) InjectCurrency(amount float64) {
	e.TotalCurrency += amount
	// Force recalculation
	e.CalculateNow()
}

// RemoveCurrency removes money from economy
func (e *EconomyEngine) RemoveCurrency(amount float64) bool {
	if amount > e.TotalCurrency {
		return false
	}
	e.TotalCurrency -= amount
	e.CalculateNow()
	return true
}
