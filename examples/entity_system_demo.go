package main

import (
	"fmt"
	"log"
	"time"

	"tesselbox/pkg/entities"
)

func main() {
	fmt.Println("=== TesselBox Entity System Demo ===")

	// Create and initialize the game world
	world := entities.NewGameWorld()
	err := world.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize world: %v", err)
	}
	defer world.Shutdown()

	fmt.Println("✅ Game world initialized")

	// Demonstrate entity creation
	fmt.Println("\n--- Entity Creation ---")

	// Create blocks
	stoneBlock, err := world.CreateBlock("stone", 10.0, 20.0, 30.0)
	if err != nil {
		log.Printf("Failed to create stone block: %v", err)
	} else {
		fmt.Printf("✅ Created stone block: %s\n", stoneBlock.ID)
	}

	dirtBlock, err := world.CreateBlock("dirt", 15.0, 25.0, 35.0)
	if err != nil {
		log.Printf("Failed to create dirt block: %v", err)
	} else {
		fmt.Printf("✅ Created dirt block: %s\n", dirtBlock.ID)
	}

	// Create items
	pickaxe, err := world.CreateItem("wooden_pickaxe", 1, "player1")
	if err != nil {
		log.Printf("Failed to create pickaxe: %v", err)
	} else {
		fmt.Printf("✅ Created wooden pickaxe: %s\n", pickaxe.ID)

		// Show tool component
		if toolComp, has := pickaxe.GetComponent("tool"); has {
			tool := toolComp.(*entities.ToolComponent)
			fmt.Printf("   Tool Type: %s, Power: %.1f, Durability: %d\n",
				tool.ToolType, tool.Power, tool.Durability)
		}
	}

	sword, err := world.CreateItem("iron_sword", 1, "player1")
	if err != nil {
		log.Printf("Failed to create sword: %v", err)
	} else {
		fmt.Printf("✅ Created iron sword: %s\n", sword.ID)

		// Show combat component
		if combatComp, has := sword.GetComponent("combat"); has {
			combat := combatComp.(*entities.CombatComponent)
			fmt.Printf("   Weapon Type: %s, Damage: %.1f, Range: %.1f\n",
				combat.WeaponType, combat.Damage, combat.Range)
		}
	}

	// Create organisms
	tree, err := world.CreateOrganism("tree", 50.0, 60.0, 70.0)
	if err != nil {
		log.Printf("Failed to create tree: %v", err)
	} else {
		fmt.Printf("✅ Created tree: %s\n", tree.ID)

		// Show behavior component
		if behaviorComp, has := tree.GetComponent("behavior"); has {
			behavior := behaviorComp.(*entities.BehaviorComponent)
			fmt.Printf("   AI Type: %s, Passive: %t, Hostile: %t\n",
				behavior.AIType, behavior.Passive, behavior.Hostile)
		}
	}

	// Demonstrate system processing
	fmt.Println("\n--- System Processing ---")

	// Run a few update cycles
	for i := 0; i < 5; i++ {
		fmt.Printf("Update cycle %d\n", i+1)
		world.Update(0.016) // 60 FPS
		time.Sleep(100 * time.Millisecond)
	}

	// Demonstrate entity queries
	fmt.Println("\n--- Entity Queries ---")

	// Query by type
	blocks := world.GetEntitiesByType("block")
	fmt.Printf("Found %d block entities\n", len(blocks))

	items := world.GetEntitiesByType("item")
	fmt.Printf("Found %d item entities\n", len(items))

	organisms := world.GetEntitiesByType("organism")
	fmt.Printf("Found %d organism entities\n", len(organisms))

	// Query by component
	renderEntities := world.GetEntitiesByComponent("render")
	fmt.Printf("Found %d entities with render component\n", len(renderEntities))

	toolEntities := world.GetEntitiesByComponent("tool")
	fmt.Printf("Found %d entities with tool component\n", len(toolEntities))

	combatEntities := world.GetEntitiesByComponent("combat")
	fmt.Printf("Found %d entities with combat component\n", len(combatEntities))

	// Query by tag
	collectibleEntities := world.GetEntitiesByTag("collectible")
	fmt.Printf("Found %d collectible entities\n", len(collectibleEntities))

	weaponEntities := world.GetEntitiesByTag("weapon")
	fmt.Printf("Found %d weapon entities\n", len(weaponEntities))

	// Demonstrate event system
	fmt.Println("\n--- Event System ---")

	eventBus := world.GetEventBus()

	// Subscribe to events
	eventBus.Subscribe(entities.EventEntityAdded, func(event entities.Event) {
		if entityEvent, ok := event.Data.(entities.EntityEvent); ok {
			fmt.Printf("📝 Entity Added: %s (%s)\n", entityEvent.Entity.ID, entityEvent.Entity.Type)
		}
	})

	eventBus.Subscribe(entities.EventBlockPlaced, func(event entities.Event) {
		if blockEvent, ok := event.Data.(entities.BlockEvent); ok {
			fmt.Printf("🧱 Block Placed: %s at (%.1f, %.1f, %.1f)\n",
				blockEvent.BlockType, blockEvent.Position.X, blockEvent.Position.Y, blockEvent.Position.Z)
		}
	})

	eventBus.Subscribe(entities.EventItemCrafted, func(event entities.Event) {
		if itemEvent, ok := event.Data.(entities.ItemEvent); ok {
			fmt.Printf("🔨 Item Crafted: %s x%d by %s\n",
				itemEvent.ItemType, itemEvent.Quantity, itemEvent.PlayerID)
		}
	})

	// Trigger some events
	fmt.Println("Triggering events...")

	// Create a new block to trigger block placed event
	grassBlock, err := world.CreateBlock("grass", 100.0, 200.0, 300.0)
	if err != nil {
		log.Printf("Failed to create grass block: %v", err)
	} else {
		fmt.Printf("Created grass block: %s\n", grassBlock.ID)
	}

	// Create a new item to trigger item crafted event
	armor, err := world.CreateItem("iron_chestplate", 1, "player1")
	if err != nil {
		log.Printf("Failed to create armor: %v", err)
	} else {
		fmt.Printf("Created iron chestplate: %s\n", armor.ID)
	}

	// Wait for event processing
	time.Sleep(100 * time.Millisecond)

	// Demonstrate plugin system
	fmt.Println("\n--- Plugin System ---")

	pluginManager := world.GetPluginManager()
	plugins := pluginManager.ListPlugins()
	fmt.Printf("Loaded plugins: %v\n", plugins)

	// Show available components
	fmt.Println("\nAvailable Components:")
	componentTypes := []string{"render", "physics", "inventory", "behavior", "crafting", "tool", "combat", "magic", "tech"}
	for _, compType := range componentTypes {
		if _, exists := entities.ComponentRegistry[compType]; exists {
			fmt.Printf("✅ %s\n", compType)
		} else {
			fmt.Printf("❌ %s\n", compType)
		}
	}

	// Demonstrate legacy compatibility
	fmt.Println("\n--- Legacy Compatibility ---")

	// Get block color using legacy method
	blockColor, err := world.GetBlockColor("stone")
	if err != nil {
		log.Printf("Failed to get block color: %v", err)
	} else {
		fmt.Printf("Stone block color: %v\n", blockColor)
	}

	// Get item properties using legacy method
	itemProps, err := world.GetItemProperties("wooden_pickaxe")
	if err != nil {
		log.Printf("Failed to get item properties: %v", err)
	} else {
		fmt.Printf("Pickaxe properties type: %T\n", itemProps)
		if propsMap, ok := itemProps.(map[string]interface{}); ok {
			fmt.Printf("Pickaxe has %d components\n", len(propsMap))
		}
	}

	// Get organism properties using legacy method
	organismProps, err := world.GetOrganismProperties("tree")
	if err != nil {
		log.Printf("Failed to get organism properties: %v", err)
	} else {
		fmt.Printf("Tree properties type: %T\n", organismProps)
		if propsMap, ok := organismProps.(map[string]interface{}); ok {
			fmt.Printf("Tree has %d components\n", len(propsMap))
		}
	}

	// Show world statistics
	fmt.Println("\n--- World Statistics ---")

	stats := world.GetStatistics()
	fmt.Printf("Initialized: %v\n", stats["initialized"])
	fmt.Printf("Total Entities: %d\n", stats["total_entities"])

	if typeCounts, ok := stats["entity_types"].(map[string]int); ok {
		fmt.Println("Entity Types:")
		for entityType, count := range typeCounts {
			fmt.Printf("  %s: %d\n", entityType, count)
		}
	}

	if componentCounts, ok := stats["component_counts"].(map[string]int); ok {
		fmt.Println("Component Counts:")
		for compType, count := range componentCounts {
			fmt.Printf("  %s: %d\n", compType, count)
		}
	}

	if tagCounts, ok := stats["tag_counts"].(map[string]int); ok {
		fmt.Println("Tag Counts:")
		for tag, count := range tagCounts {
			fmt.Printf("  %s: %d\n", tag, count)
		}
	}

	// Demonstrate flexibility with dynamic entity creation
	fmt.Println("\n--- Dynamic Entity Creation ---")

	// Create a custom entity using the entity manager directly

	// Create a custom "magic sword" entity by combining components
	magicSword := entities.NewEntity("magic_sword_001", "item")
	magicSword.Name = "Magic Sword of Fire"
	magicSword.Tags = []string{"item", "weapon", "magic", "fire"}

	// Add render component
	magicSword.AddComponent(&entities.RenderComponent{
		Type:     "render",
		Color:    entities.ColorFromSlice([]uint8{255, 100, 0, 255}), // Orange
		Pattern:  "glowing",
		Visible:  true,
		Scale:    1.2,
		Animated: true,
	})

	// Add combat component
	magicSword.AddComponent(&entities.CombatComponent{
		Type:       "combat",
		WeaponType: "sword",
		Damage:     15.0,
		Range:      1.5,
		Speed:      1.8,
		Health:     0,
		MaxHealth:  0,
		Effects:    []string{"fire", "magic"},
	})

	// Add magic component (from plugin)
	magicSword.AddComponent(&entities.MagicComponent{
		Type:       "magic",
		ManaCost:   5,
		SpellPower: 20,
		Spells:     []string{"fireball", "flame_strike"},
	})

	// Add inventory component
	magicSword.AddComponent(&entities.InventoryComponent{
		Type:              "inventory",
		StackSize:         1,
		MaxDurability:     500,
		CurrentDurability: 500,
		Container:         false,
		Weight:            2.5,
		Categories:        []string{"weapon", "magic", "legendary"},
	})

	// Register the custom entity
	world.GetSystemManager().AddEntity(magicSword)

	fmt.Printf("✅ Created custom magic sword: %s\n", magicSword.Name)
	fmt.Printf("   Components: %d\n", len(magicSword.Components))
	fmt.Printf("   Tags: %v\n", magicSword.Tags)

	// Show magic component
	if magicComp, has := magicSword.GetComponent("magic"); has {
		magic := magicComp.(*entities.MagicComponent)
		fmt.Printf("   Magic - Mana Cost: %d, Spell Power: %d, Spells: %v\n",
			magic.ManaCost, magic.SpellPower, magic.Spells)
	}

	// Final world statistics
	fmt.Println("\n--- Final World Statistics ---")
	world.DebugPrint()

	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("The entity system is working with:")
	fmt.Println("✅ Component-based architecture")
	fmt.Println("✅ Data-driven configuration")
	fmt.Println("✅ Event-driven communication")
	fmt.Println("✅ Plugin extensibility")
	fmt.Println("✅ Legacy compatibility")
	fmt.Println("✅ High performance")
}
