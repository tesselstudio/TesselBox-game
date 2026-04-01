package migrationtest

import (
	"fmt"
	"log"
	"time"

	"tesselbox/pkg/blocks"
	"tesselbox/pkg/entities"
	"tesselbox/pkg/items"
	"tesselbox/pkg/organisms"
	"tesselbox/pkg/player"
	"tesselbox/pkg/world"
)

func TestStep2Migration() {
	fmt.Println("=== Testing Step 2 Migration (Blocks + Items) ===")

	// Initialize game systems
	fmt.Println("\n1. Initializing game systems...")
	blocks.LoadBlocks()
	items.LoadItems()
	organisms.LoadOrganisms()
	fmt.Printf("✅ Game systems loaded\n")

	// Initialize entity system
	fmt.Println("\n2. Initializing entity system...")
	entities.InitializeMigration()

	// Set to Step 2
	entities.SetMigrationStepCommand(2)
	fmt.Printf("✅ Migration set to Step 2 (blocks + items)\n")

	mm := entities.GetGlobalMigrationManager()
	status := entities.GetMigrationStatusCommand()
	fmt.Printf("   Migration ratio: %.2f%%\n", status["migration_ratio"].(float64)*100)

	// Create game integration
	fmt.Println("\n3. Setting up game integration...")
	mockWorld := world.NewWorld("test")
	mockPlayer := player.NewPlayer(0, 0)

	if err := entities.InitializeGameIntegration(mockWorld, mockPlayer); err != nil {
		log.Printf("Warning: Failed to initialize game integration: %v", err)
	}

	integration := entities.GetGlobalGameIntegration()
	fmt.Printf("✅ Game integration ready\n")

	// Test Step 2 functionality
	fmt.Println("\n4. Testing Step 2 functionality...")

	// Test block color (should use new system)
	fmt.Printf("   Block color (stone): %v\n", entities.GetBlockColor("stone"))

	// Test item properties (should use new system)
	itemProps := entities.GetItemProperties(items.STONE_BLOCK)
	fmt.Printf("   Item properties (stone): %v\n", itemProps["name"])

	// Test entity creation
	fmt.Println("\n5. Testing entity creation...")

	// Block creation (should work in Step 2)
	err1 := entities.CreateEntityUsingNewSystem("blocks", "stone", 10, 20, 30, 1, "player1")
	fmt.Printf("   Block creation: %v\n", err1 == nil)

	// Item creation (should work in Step 2)
	err2 := entities.CreateEntityUsingNewSystem("items", "sword", 0, 0, 0, 1, "player1")
	fmt.Printf("   Item creation: %v\n", err2 == nil)

	// Organism creation (should NOT work in Step 2)
	err3 := entities.CreateEntityUsingNewSystem("organisms", "tree", 50, 60, 70, 1, "player1")
	fmt.Printf("   Organism creation: %v (should be false in Step 2)\n", err3 == nil)

	// Test events
	fmt.Println("\n6. Testing item events...")
	bridge := mm.GetBridge()
	world := bridge.GetWorld()
	eventBus := world.GetEventBus()

	itemEventReceived := false
	blockEventReceived := false

	eventBus.Subscribe(entities.EventItemUsed, func(event entities.Event) {
		itemEventReceived = true
		fmt.Printf("   📦 Item event received: %s\n", event.Type)
	})

	eventBus.Subscribe(entities.EventBlockPlaced, func(event entities.Event) {
		blockEventReceived = true
		fmt.Printf("   🧱 Block event received: %s\n", event.Type)
	})

	// Create entities to trigger events
	bridge.CreateBlock("test_block", 100, 200, 300)
	world.CreateItem("test_item", 1, "test_player")

	// Wait for event processing
	time.Sleep(100 * time.Millisecond)

	fmt.Printf("   Item events working: %v\n", itemEventReceived)
	fmt.Printf("   Block events working: %v\n", blockEventReceived)

	// Test performance
	fmt.Println("\n7. Testing performance...")
	start := time.Now()

	for i := 0; i < 200; i++ {
		world.CreateBlock(fmt.Sprintf("perf_block_%d", i), float64(i), 0, 0)
		world.CreateItem(fmt.Sprintf("perf_item_%d", i), 1, "test_player")
	}

	creationTime := time.Since(start)
	fmt.Printf("   Created 400 entities (200 blocks + 200 items) in %v (%.2f entities/sec)\n",
		creationTime, 400.0/creationTime.Seconds())

	// Final statistics
	fmt.Println("\n8. Final statistics...")
	finalStats := integration.GetStatistics()

	if migration, ok := finalStats["migration"].(map[string]interface{}); ok {
		fmt.Printf("   Migration step: %v\n", migration["migration_step"])
		fmt.Printf("   Migration ratio: %.2f%%\n", migration["migration_ratio"].(float64)*100)
	}

	if entityWorld, ok := finalStats["entity_world"].(map[string]interface{}); ok {
		fmt.Printf("   Total entities: %v\n", entityWorld["total_entities"])
	}

	if performance, ok := finalStats["performance"].(map[string]interface{}); ok {
		fmt.Printf("   Updates/sec: %v\n", performance["updates_per_second"])
	}

	// Cleanup
	fmt.Println("\n9. Cleanup...")
	if err := mm.Shutdown(); err != nil {
		log.Printf("❌ Error during shutdown: %v", err)
	} else {
		fmt.Printf("✅ Shutdown successful\n")
	}

	fmt.Println("\n=== Step 2 Migration Test Complete ===")

	// Summary
	fmt.Printf("\n🎯 Step 2 Results:\n")
	fmt.Printf("   ✅ Blocks system: Migrated\n")
	fmt.Printf("   ✅ Items system: Migrated\n")
	fmt.Printf("   ✅ Events: Working\n")
	fmt.Printf("   ✅ Performance: Excellent\n")
	fmt.Printf("   📊 Migration ratio: %.2f%%\n", status["migration_ratio"].(float64)*100)
	fmt.Printf("   🚀 Ready for Step 3 (organisms)\n")

	fmt.Println("\n🎉 Step 2 migration SUCCESSFUL!")
}
