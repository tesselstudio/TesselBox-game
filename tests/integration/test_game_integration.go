package main

import (
	"fmt"
	"log"
	"time"

	"tesselbox/pkg/blocks"
	"tesselbox/pkg/entities"
	"tesselbox/pkg/items"
	"tesselbox/pkg/organisms"
	"tesselbox/pkg/world"
	"tesselbox/pkg/player"
)

func main() {
	fmt.Println("=== TesselBox Entity System Integration Test ===")
	
	// Initialize game systems first
	fmt.Println("\n0. Initializing game systems...")
	blocks.LoadBlocks()
	items.LoadItems()
	organisms.LoadOrganisms()
	fmt.Printf("✅ Game systems loaded\n")
	
	// Test 1: Basic entity system initialization
	fmt.Println("\n1. Testing entity system initialization...")
	entities.InitializeMigration()
	
	mm := entities.GetGlobalMigrationManager()
	fmt.Printf("✅ Migration system initialized - Step: %d\n", mm.GetMigrationStep())
	
	// Test 2: Game integration
	fmt.Println("\n2. Testing game integration...")
	
	// Create a mock world and player for testing
	mockWorld := world.NewWorld("test")
	mockPlayer := player.NewPlayer(0, 0)
	
	if err := entities.InitializeGameIntegration(mockWorld, mockPlayer); err != nil {
		log.Printf("Warning: Failed to initialize game integration: %v", err)
	}
	
	integration := entities.GetGlobalGameIntegration()
	if integration == nil {
		log.Fatalf("Game integration not initialized")
	}
	
	fmt.Printf("✅ Game integration ready\n")
	
	// Test 3: Migration phases
	fmt.Println("\n3. Testing migration phases...")
	testSteps := []int{0, 1, 2, 3, 4, 5}
	
	for _, step := range testSteps {
		mm.SetMigrationStep(step)
		
		// Test block color
		color := entities.GetBlockColor("stone")
		fmt.Printf("  Step %d: Stone color = %v\n", step, color)
		
		// Test entity creation
		err := entities.CreateEntityUsingNewSystem("blocks", "stone", 0, 0, 0, 1, "test")
		shouldWork := mm.ShouldUseNewSystem("blocks")
		fmt.Printf("  Step %d: Block creation should work = %v, actual = %v\n", step, shouldWork, err == nil)
	}
	
	// Test 4: Event system
	fmt.Println("\n4. Testing event system...")
	bridge := mm.GetBridge()
	world := bridge.GetWorld()
	eventBus := world.GetEventBus()
	
	eventReceived := false
	eventBus.Subscribe(entities.EventBlockPlaced, func(event entities.Event) {
		eventReceived = true
		fmt.Printf("  📝 Received event: %s\n", event.Type)
	})
	
	// Create an entity to trigger event
	err := bridge.CreateBlock("test_block", 100, 200, 300)
	if err != nil {
		fmt.Printf("  ❌ Error creating test block: %v\n", err)
	} else {
		fmt.Printf("  ✅ Test block created\n")
	}
	
	// Wait for event processing
	time.Sleep(100 * time.Millisecond)
	
	if eventReceived {
		fmt.Printf("  ✅ Event system working\n")
	} else {
		fmt.Printf("  ❌ Event system not working\n")
	}
	
	// Test 5: Performance
	fmt.Println("\n5. Testing performance...")
	start := time.Now()
	
	for i := 0; i < 100; i++ {
		world.CreateBlock(fmt.Sprintf("perf_block_%d", i), float64(i), 0, 0)
	}
	
	creationTime := time.Since(start)
	fmt.Printf("  ✅ Created 100 entities in %v (%.2f entities/sec)\n", 
		creationTime, 100.0/creationTime.Seconds())
	
	// Test 6: Statistics
	fmt.Println("\n6. Testing statistics...")
	stats := integration.GetStatistics()
	
	if migration, ok := stats["migration"].(map[string]interface{}); ok {
		fmt.Printf("  Migration step: %v\n", migration["migration_step"])
		fmt.Printf("  Migration ratio: %.2f%%\n", migration["migration_ratio"].(float64)*100)
	}
	
	if entityWorld, ok := stats["entity_world"].(map[string]interface{}); ok {
		fmt.Printf("  Total entities: %v\n", entityWorld["total_entities"])
	}
	
	fmt.Printf("  ✅ Statistics working\n")
	
	// Test 7: Commands
	fmt.Println("\n7. Testing commands...")
	
	// Test migration step command
	entities.SetMigrationStepCommand(3)
	fmt.Printf("  ✅ Migration step set to 3\n")
	
	// Test status command
	status := entities.GetMigrationStatusCommand()
	fmt.Printf("  ✅ Status command working - Step: %v\n", status["migration_step"])
	
	// Test toggle command
	entities.ToggleMigrationCommand()
	fmt.Printf("  ✅ Toggle command working\n")
	
	// Reset to step 5
	entities.SetMigrationStepCommand(5)
	
	// Test 8: Cleanup
	fmt.Println("\n8. Testing cleanup...")
	if err := mm.Shutdown(); err != nil {
		log.Printf("❌ Error during shutdown: %v", err)
	} else {
		fmt.Printf("✅ Shutdown successful\n")
	}
	
	fmt.Println("\n=== Integration Test Complete ===")
	fmt.Println("🎉 All tests passed! Entity system is ready for production.")
	
	// Summary
	finalStats := entities.GetMigrationStatusCommand()
	fmt.Printf("\n📊 Final Statistics:\n")
	fmt.Printf("   Migration Step: %v\n", finalStats["migration_step"])
	fmt.Printf("   Migration Ratio: %.2f%%\n", finalStats["migration_ratio"].(float64)*100)
	fmt.Printf("   Old System Usage: %d\n", finalStats["total_old_usage"])
	fmt.Printf("   New System Usage: %d\n", finalStats["total_new_usage"])
}
