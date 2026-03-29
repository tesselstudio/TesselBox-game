package entities

import (
	"testing"
	"time"
)

// TestGameWorldInitialization tests the game world initialization
func TestGameWorldInitialization(t *testing.T) {
	world := NewGameWorld()
	
	// Test initial state
	if world.IsInitialized() {
		t.Error("World should not be initialized initially")
	}
	
	// Initialize world
	err := world.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize world: %v", err)
	}
	
	// Test initialized state
	if !world.IsInitialized() {
		t.Error("World should be initialized after Initialize()")
	}
	
	// Test statistics
	stats := world.GetStatistics()
	if initialized, ok := stats["initialized"].(bool); !ok || !initialized {
		t.Error("Statistics should show world as initialized")
	}
	
	// Clean up
	err = world.Shutdown()
	if err != nil {
		t.Errorf("Failed to shutdown world: %v", err)
	}
}

// TestEntityCreation tests entity creation from templates
func TestEntityCreation(t *testing.T) {
	world := NewGameWorld()
	err := world.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize world: %v", err)
	}
	defer world.Shutdown()
	
	// Test creating a block
	block, err := world.CreateBlock("stone", 10.0, 20.0, 30.0)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}
	
	if block.Type != "block" {
		t.Errorf("Expected block type 'block', got '%s'", block.Type)
	}
	
	if !block.HasComponent("render") {
		t.Error("Block should have render component")
	}
	
	if !block.HasComponent("physics") {
		t.Error("Block should have physics component")
	}
	
	// Test creating an item
	item, err := world.CreateItem("wooden_pickaxe", 1, "player1")
	if err != nil {
		t.Fatalf("Failed to create item: %v", err)
	}
	
	if item.Type != "item" {
		t.Errorf("Expected item type 'item', got '%s'", item.Type)
	}
	
	if !item.HasComponent("render") {
		t.Error("Item should have render component")
	}
	
	if !item.HasComponent("inventory") {
		t.Error("Item should have inventory component")
	}
	
	if !item.HasComponent("tool") {
		t.Error("Pickaxe should have tool component")
	}
	
	// Test creating an organism
	organism, err := world.CreateOrganism("tree", 50.0, 60.0, 70.0)
	if err != nil {
		t.Fatalf("Failed to create organism: %v", err)
	}
	
	if organism.Type != "organism" {
		t.Errorf("Expected organism type 'organism', got '%s'", organism.Type)
	}
	
	if !organism.HasComponent("render") {
		t.Error("Organism should have render component")
	}
	
	if !organism.HasComponent("behavior") {
		t.Error("Organism should have behavior component")
	}
	
	if !organism.HasComponent("combat") {
		t.Error("Organism should have combat component")
	}
}

// TestComponentSystem tests the component system
func TestComponentSystem(t *testing.T) {
	// Test component registration
	RegisterComponent("test", &TestComponent{})
	
	if _, exists := ComponentRegistry["test"]; !exists {
		t.Error("Test component should be registered")
	}
	
	// Test component creation
	component, err := CreateComponent("test")
	if err != nil {
		t.Fatalf("Failed to create test component: %v", err)
	}
	
	if component.GetType() != "test" {
		t.Errorf("Expected component type 'test', got '%s'", component.GetType())
	}
	
	// Test component cloning
	cloned := component.Clone()
	if cloned.GetType() != component.GetType() {
		t.Error("Cloned component should have same type")
	}
	
	// Test component validation
	err = component.Validate()
	if err != nil {
		t.Errorf("Component validation failed: %v", err)
	}
}

// TestSystemManager tests the system manager
func TestSystemManager(t *testing.T) {
	systemManager := NewSystemManager()
	
	// Create test entities
	entity1 := NewEntity("entity1", "test")
	entity1.AddComponent(&TestComponent{})
	
	entity2 := NewEntity("entity2", "test")
	entity2.AddComponent(&TestComponent{})
	
	// Register entities
	systemManager.AddEntity(entity1)
	systemManager.AddEntity(entity2)
	
	// Register test system
	testSystem := &TestSystem{}
	systemManager.RegisterSystem(testSystem)
	
	// Test system processing
	systemManager.Update(0.016) // 60 FPS frame time
	
	// Test entity queries
	entities := systemManager.FindEntitiesByComponent("test")
	if len(entities) != 2 {
		t.Errorf("Expected 2 entities with test component, got %d", len(entities))
	}
	
	entities = systemManager.FindEntitiesByTag("test")
	if len(entities) != 0 {
		t.Errorf("Expected 0 entities with test tag, got %d", len(entities))
	}
	
	// Test entity removal
	systemManager.RemoveEntity("entity1")
	entities = systemManager.FindEntitiesByComponent("test")
	if len(entities) != 1 {
		t.Errorf("Expected 1 entity after removal, got %d", len(entities))
	}
}

// TestEventBus tests the event system
func TestEventBus(t *testing.T) {
	eventBus := NewEventBus()
	
	// Test event subscription
	received := false
	eventBus.Subscribe(EventEntityAdded, func(event Event) {
		received = true
	})
	
	// Test event publishing
	eventBus.Publish(EventEntityAdded, EntityEvent{
		Entity: NewEntity("test", "test"),
	})
	
	// Wait a bit for async processing
	time.Sleep(10 * time.Millisecond)
	
	if !received {
		t.Error("Event should have been received")
	}
	
	// Test event filtering
	filter := NewEventFilter().ByType(EventEntityAdded)
	testEvent := Event{
		Type: EventEntityAdded,
		Source: "test",
		Data: nil,
		Priority: 0,
	}
	
	if !filter.Matches(testEvent) {
		t.Error("Event should match filter")
	}
	
	filter = NewEventFilter().ByType(EventEntityRemoved)
	if filter.Matches(testEvent) {
		t.Error("Event should not match different type filter")
	}
}

// TestPluginSystem tests the plugin system
func TestPluginSystem(t *testing.T) {
	world := NewGameWorld()
	err := world.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize world: %v", err)
	}
	defer world.Shutdown()
	
	// Test built-in plugin loading
	pluginManager := world.GetPluginManager()
	
	// Check if magic plugin components are available
	if _, exists := ComponentRegistry["magic"]; !exists {
		t.Error("Magic component should be registered from built-in plugin")
	}
	
	// Check if magic system is registered
	systems := world.GetSystemManager()
	// We can't directly access systems, but we can check if the plugin is loaded
	plugins := pluginManager.ListPlugins()
	
	// Built-in plugins should be loaded
	hasMagic := false
	hasTech := false
	for _, pluginName := range plugins {
		if pluginName == "magic" {
			hasMagic = true
		}
		if pluginName == "tech" {
			hasTech = true
		}
	}
	
	if !hasMagic {
		t.Error("Magic plugin should be loaded")
	}
	
	if !hasTech {
		t.Error("Tech plugin should be loaded")
	}
}

// TestLegacyCompatibility tests legacy compatibility methods
func TestLegacyCompatibility(t *testing.T) {
	world := NewGameWorld()
	err := world.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize world: %v", err)
	}
	defer world.Shutdown()
	
	// Test legacy block color method
	color, err := world.GetBlockColor("stone")
	if err != nil {
		t.Errorf("Failed to get block color: %v", err)
	}
	
	if color == nil {
		t.Error("Block color should not be nil")
	}
	
	// Test legacy item properties method
	properties, err := world.GetItemProperties("wooden_pickaxe")
	if err != nil {
		t.Errorf("Failed to get item properties: %v", err)
	}
	
	if properties == nil {
		t.Error("Item properties should not be nil")
	}
	
	// Test legacy organism properties method
	organismProps, err := world.GetOrganismProperties("tree")
	if err != nil {
		t.Errorf("Failed to get organism properties: %v", err)
	}
	
	if organismProps == nil {
		t.Error("Organism properties should not be nil")
	}
}

// TestPerformance tests performance with many entities
func TestPerformance(t *testing.T) {
	world := NewGameWorld()
	err := world.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize world: %v", err)
	}
	defer world.Shutdown()
	
	// Create many entities
	const entityCount = 1000
	start := time.Now()
	
	for i := 0; i < entityCount; i++ {
		_, err := world.CreateBlock("stone", float64(i), 0.0, 0.0)
		if err != nil {
			t.Errorf("Failed to create block %d: %v", i, err)
		}
	}
	
	creationTime := time.Since(start)
	t.Logf("Created %d entities in %v", entityCount, creationTime)
	
	// Test update performance
	start = time.Now()
	for i := 0; i < 100; i++ {
		world.Update(0.016)
	}
	updateTime := time.Since(start)
	t.Logf("Updated %d entities for 100 frames in %v", entityCount, updateTime)
	
	// Test query performance
	start = time.Now()
	entities := world.GetEntitiesByType("block")
	queryTime := time.Since(start)
	
	if len(entities) != entityCount {
		t.Errorf("Expected %d entities, got %d", entityCount, len(entities))
	}
	
	t.Logf("Queried %d entities in %v", entityCount, queryTime)
	
	// Performance assertions
	if creationTime > time.Second {
		t.Errorf("Entity creation too slow: %v", creationTime)
	}
	
	if updateTime > time.Second {
		t.Errorf("Entity update too slow: %v", updateTime)
	}
	
	if queryTime > 100*time.Millisecond {
		t.Errorf("Entity query too slow: %v", queryTime)
	}
}

// ============================================================================
// Test Components and Systems
// ============================================================================

// TestComponent is a simple test component
type TestComponent struct {
	Type     string `yaml:"type"`
	Value    int    `yaml:"value"`
	Enabled  bool   `yaml:"enabled"`
}

func (c *TestComponent) GetType() string { return "test" }
func (c *TestComponent) Clone() Component {
	clone := *c
	return &clone
}
func (c *TestComponent) Merge(other Component) {
	if tc, ok := other.(*TestComponent); ok {
		if tc.Value != 0 {
			c.Value = tc.Value
		}
		c.Enabled = tc.Enabled
	}
}
func (c *TestComponent) Validate() error {
	if c.Value < 0 {
		return fmt.Errorf("value cannot be negative")
	}
	return nil
}

// TestSystem is a simple test system
type TestSystem struct {
	name string
	requiredComponents []string
}

func NewTestSystem() *TestSystem {
	return &TestSystem{
		name: "test",
		requiredComponents: []string{"test"},
	}
}

func (ts *TestSystem) GetName() string { return ts.name }
func (ts *TestSystem) GetRequiredComponents() []string { return ts.requiredComponents }
func (ts *TestSystem) Matches(entity *Entity) bool {
	return entity.HasComponent("test")
}

func (ts *TestSystem) Process(deltaTime float64, entities []*Entity) {
	for _, entity := range entities {
		if !ts.Matches(entity) {
			continue
		}

		testComp, _ := entity.GetComponent("test")
		if test, ok := testComp.(*TestComponent); ok {
			// Simple processing: increment value if enabled
			if test.Enabled {
				test.Value++
			}
		}
	}
}

// ============================================================================
// Benchmarks
// ============================================================================

// BenchmarkEntityCreation benchmarks entity creation performance
func BenchmarkEntityCreation(b *testing.B) {
	world := NewGameWorld()
	err := world.Initialize()
	if err != nil {
		b.Fatalf("Failed to initialize world: %v", err)
	}
	defer world.Shutdown()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := world.CreateBlock("stone", float64(i), 0.0, 0.0)
		if err != nil {
			b.Fatalf("Failed to create block: %v", err)
		}
	}
}

// BenchmarkEntityUpdate benchmarks entity update performance
func BenchmarkEntityUpdate(b *testing.B) {
	world := NewGameWorld()
	err := world.Initialize()
	if err != nil {
		b.Fatalf("Failed to initialize world: %v", err)
	}
	defer world.Shutdown()
	
	// Create test entities
	for i := 0; i < 1000; i++ {
		_, err := world.CreateBlock("stone", float64(i), 0.0, 0.0)
		if err != nil {
			b.Fatalf("Failed to create block: %v", err)
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update(0.016)
	}
}

// BenchmarkEntityQuery benchmarks entity query performance
func BenchmarkEntityQuery(b *testing.B) {
	world := NewGameWorld()
	err := world.Initialize()
	if err != nil {
		b.Fatalf("Failed to initialize world: %v", err)
	}
	defer world.Shutdown()
	
	// Create test entities
	for i := 0; i < 1000; i++ {
		_, err := world.CreateBlock("stone", float64(i), 0.0, 0.0)
		if err != nil {
			b.Fatalf("Failed to create block: %v", err)
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = world.GetEntitiesByType("block")
	}
}
