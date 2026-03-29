package entities

import (
	"log"
	"sync"
	"time"
	
	"tesselbox/pkg/items"
	"tesselbox/pkg/organisms"
)

// MigrationManager handles the gradual migration from old to new systems
type MigrationManager struct {
	bridge        *Bridge
	enabled       bool
	migrationStep int
	mutex         sync.RWMutex
	
	// Migration statistics
	oldSystemUsage map[string]int
	newSystemUsage map[string]int
	migrationStats map[string]interface{}
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager() *MigrationManager {
	return &MigrationManager{
		bridge:         NewBridge(),
		enabled:        true,
		migrationStep:  0,
		oldSystemUsage: make(map[string]int),
		newSystemUsage: make(map[string]int),
		migrationStats: make(map[string]interface{}),
	}
}

// GetBridge returns the bridge instance
func (mm *MigrationManager) GetBridge() *Bridge {
	return mm.bridge
}

// SetEnabled enables or disables migration
func (mm *MigrationManager) SetEnabled(enabled bool) {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()
	mm.enabled = enabled
}

// IsEnabled returns whether migration is enabled
func (mm *MigrationManager) IsEnabled() bool {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()
	return mm.enabled
}

// SetMigrationStep sets the current migration step
func (mm *MigrationManager) SetMigrationStep(step int) {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()
	mm.migrationStep = step
	log.Printf("Migration step set to: %d", step)
}

// GetMigrationStep returns the current migration step
func (mm *MigrationManager) GetMigrationStep() int {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()
	return mm.migrationStep
}

// RecordOldSystemUsage records usage of old system
func (mm *MigrationManager) RecordOldSystemUsage(system string) {
	if !mm.enabled {
		return
	}
	
	mm.mutex.Lock()
	defer mm.mutex.Unlock()
	mm.oldSystemUsage[system]++
}

// RecordNewSystemUsage records usage of new system
func (mm *MigrationManager) RecordNewSystemUsage(system string) {
	if !mm.enabled {
		return
	}
	
	mm.mutex.Lock()
	defer mm.mutex.Unlock()
	mm.newSystemUsage[system]++
}

// GetMigrationProgress returns migration progress statistics
func (mm *MigrationManager) GetMigrationProgress() map[string]interface{} {
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()
	
	totalOldUsage := 0
	for _, count := range mm.oldSystemUsage {
		totalOldUsage += count
	}
	
	totalNewUsage := 0
	for _, count := range mm.newSystemUsage {
		totalNewUsage += count
	}
	
	progress := map[string]interface{}{
		"migration_step":    mm.migrationStep,
		"enabled":           mm.enabled,
		"old_system_usage":  mm.oldSystemUsage,
		"new_system_usage":  mm.newSystemUsage,
		"total_old_usage":   totalOldUsage,
		"total_new_usage":   totalNewUsage,
		"migration_ratio":  float64(totalNewUsage) / float64(totalOldUsage+totalNewUsage),
		"last_updated":      time.Now(),
	}
	
	return progress
}

// ShouldUseNewSystem determines if new system should be used based on migration step
func (mm *MigrationManager) ShouldUseNewSystem(feature string) bool {
	if !mm.enabled {
		return false
	}
	
	mm.mutex.RLock()
	defer mm.mutex.RUnlock()
	
	// Migration phases
	switch mm.migrationStep {
	case 0: // Old system only
		return false
	case 1: // Hybrid - blocks only
		return feature == "blocks"
	case 2: // Hybrid - blocks and items
		return feature == "blocks" || feature == "items"
	case 3: // Hybrid - blocks, items, and organisms
		return feature == "blocks" || feature == "items" || feature == "organisms"
	case 4: // New system primary with fallback
		return true
	case 5: // New system only
		return true
	default:
		return false
	}
}

// Update updates the migration manager
func (mm *MigrationManager) Update(deltaTime float64) {
	if !mm.enabled {
		return
	}
	
	// Update the bridge
	mm.bridge.Update(deltaTime)
	
	// Periodically log migration progress
	if time.Now().Unix()%60 == 0 { // Every minute
		progress := mm.GetMigrationProgress()
		log.Printf("Migration Progress: Step %d, Ratio %.2f%%", 
			progress["migration_step"], 
			progress["migration_ratio"].(float64)*100)
	}
}

// Shutdown shuts down the migration manager
func (mm *MigrationManager) Shutdown() error {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()
	
	// Log final migration statistics
	progress := mm.GetMigrationProgress()
	log.Printf("Final Migration Statistics: %+v", progress)
	
	return mm.bridge.Shutdown()
}

// Global migration manager instance
var globalMigrationManager *MigrationManager
var migrationManagerOnce sync.Once

// GetGlobalMigrationManager returns the global migration manager
func GetGlobalMigrationManager() *MigrationManager {
	migrationManagerOnce.Do(func() {
		globalMigrationManager = NewMigrationManager()
	})
	return globalMigrationManager
}

// ============================================================================
// Migration Helper Functions
// ============================================================================

// MigrateBlockColor migrates block color queries
func MigrateBlockColor(blockType string) (interface{}, error) {
	mm := GetGlobalMigrationManager()
	
	if mm.ShouldUseNewSystem("blocks") {
		mm.RecordNewSystemUsage("block_color")
		return mm.bridge.GetBlockColor(blockType), nil
	}
	
	mm.RecordOldSystemUsage("block_color")
	// Fallback to old system would be here
	return nil, nil
}

// MigrateItemProperties migrates item property queries
func MigrateItemProperties(itemType interface{}) (map[string]interface{}, error) {
	mm := GetGlobalMigrationManager()
	
	if mm.ShouldUseNewSystem("items") {
		mm.RecordNewSystemUsage("item_properties")
		if item, ok := itemType.(items.ItemType); ok {
			return mm.bridge.GetItemProperties(item), nil
		}
	}
	
	mm.RecordOldSystemUsage("item_properties")
	// Fallback to old system would be here
	return nil, nil
}

// MigrateOrganismProperties migrates organism property queries
func MigrateOrganismProperties(orgType interface{}) (map[string]interface{}, error) {
	mm := GetGlobalMigrationManager()
	
	if mm.ShouldUseNewSystem("organisms") {
		mm.RecordNewSystemUsage("organism_properties")
		// Convert to proper type
		if org, ok := orgType.(int); ok {
			return mm.bridge.GetOrganismProperties(organisms.OrganismType(org)), nil
		}
	}
	
	mm.RecordOldSystemUsage("organism_properties")
	// Fallback to old system would be here
	return nil, nil
}

// CreateEntityUsingNewSystem creates entities using the new system
func CreateEntityUsingNewSystem(entityType, entityID string, x, y, z float64, quantity int, playerID string) error {
	mm := GetGlobalMigrationManager()
	
	if !mm.ShouldUseNewSystem(entityType) {
		return nil // Don't create if migration step doesn't allow
	}
	
	switch entityType {
	case "blocks":
		mm.RecordNewSystemUsage("create_block")
		return mm.bridge.CreateBlock(entityID, x, y, z)
	case "items":
		mm.RecordNewSystemUsage("create_item")
		return mm.bridge.CreateItem(entityID, quantity, playerID)
	case "organisms":
		mm.RecordNewSystemUsage("create_organism")
		return mm.bridge.CreateOrganism(entityID, x, y, z)
	default:
		return nil
	}
}

// ============================================================================
// Migration Commands
// ============================================================================

// SetMigrationStepCommand sets the migration step
func SetMigrationStepCommand(step int) {
	mm := GetGlobalMigrationManager()
	mm.SetMigrationStep(step)
	log.Printf("Migration step set to %d", step)
}

// GetMigrationStatusCommand returns migration status
func GetMigrationStatusCommand() map[string]interface{} {
	mm := GetGlobalMigrationManager()
	return mm.GetMigrationProgress()
}

// ToggleMigrationCommand toggles migration on/off
func ToggleMigrationCommand() {
	mm := GetGlobalMigrationManager()
	mm.SetEnabled(!mm.IsEnabled())
	log.Printf("Migration enabled: %v", mm.IsEnabled())
}

// ============================================================================
// Integration Hooks
// ============================================================================

// InitializeMigration initializes the migration system
func InitializeMigration() {
	mm := GetGlobalMigrationManager()
	log.Printf("Migration system initialized - Step: %d, Enabled: %v", 
		mm.GetMigrationStep(), mm.IsEnabled())
}

// UpdateMigration updates the migration system
func UpdateMigration(deltaTime float64) {
	mm := GetGlobalMigrationManager()
	mm.Update(deltaTime)
}

// ShutdownMigration shuts down the migration system
func ShutdownMigration() error {
	mm := GetGlobalMigrationManager()
	return mm.Shutdown()
}
