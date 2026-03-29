package internal

import (
	"testing"
	"time"

	"tesselbox/pkg/entities"
	"tesselbox/pkg/world"
)

// MockEntityManager provides a mock entity manager for testing
type MockEntityManager struct {
	entities map[uint64]*entities.Entity
	nextID  uint64
}

func NewMockEntityManager() *MockEntityManager {
	return &MockEntityManager{
		entities: make(map[uint64]*entities.Entity),
		nextID:   1,
	}
}

func (m *MockEntityManager) CreateEntity(entityType string) *entities.Entity {
	id := m.nextID
	m.nextID++
	
	entity := &entities.Entity{
		ID:   id,
		Type: entityType,
		Tags: []string{},
	}
	
	m.entities[id] = entity
	return entity
}

func (m *MockEntityManager) GetEntity(id uint64) (*entities.Entity, bool) {
	entity, exists := m.entities[id]
	return entity, exists
}

func (m *MockEntityManager) RemoveEntity(id uint64) {
	delete(m.entities, id)
}

func (m *MockEntityManager) GetAllEntities() []*entities.Entity {
	result := make([]*entities.Entity, 0, len(m.entities))
	for _, entity := range m.entities {
		result = append(result, entity)
	}
	return result
}

// MockEventBus provides a mock event bus for testing
type MockEventBus struct {
	events []entities.Event
}

func NewMockEventBus() *MockEventBus {
	return &MockEventBus{
		events: make([]entities.Event, 0),
	}
}

func (m *MockEventBus) Publish(event entities.Event) {
	m.events = append(m.events, event)
}

func (m *MockEventBus) GetEvents() []entities.Event {
	return m.events
}

func (m *MockEventBus) Clear() {
	m.events = make([]entities.Event, 0)
}

// MockWorld provides a mock world for testing
type MockWorld struct {
	blocks map[string]*world.Hexagon
}

func NewMockWorld() *MockWorld {
	return &MockWorld{
		blocks: make(map[string]*world.Hexagon),
	}
}

func (m *MockWorld) AddBlock(x, y float64, blockType string) {
	key := fmt.Sprintf("%.1f,%.1f", x, y)
	m.blocks[key] = &world.Hexagon{
		X:         x,
		Y:         y,
		BlockType:  world.STONE, // Simplified for testing
	}
}

func (m *MockWorld) GetBlock(x, y float64) *world.Hexagon {
	key := fmt.Sprintf("%.1f,%.1f", x, y)
	return m.blocks[key]
}

// TestHelper provides common test utilities
type TestHelper struct {
	t *testing.T
}

func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{t: t}
}

// AssertNoError fails the test if err is not nil
func (th *TestHelper) AssertNoError(err error, message string) {
	if err != nil {
		th.t.Fatalf("%s: %v", message, err)
	}
}

// AssertError fails the test if err is nil
func (th *TestHelper) AssertError(err error, message string) {
	if err == nil {
		th.t.Fatalf("%s: expected error but got nil", message)
	}
}

// AssertEqual fails the test if expected and actual don't match
func (th *TestHelper) AssertEqual(expected, actual interface{}, message string) {
	if expected != actual {
		th.t.Fatalf("%s: expected %v, got %v", message, expected, actual)
	}
}

// AssertNotEqual fails the test if expected and actual do match
func (th *TestHelper) AssertNotEqual(expected, actual interface{}, message string) {
	if expected == actual {
		th.t.Fatalf("%s: expected %v to not equal %v", message, expected, actual)
	}
}

// AssertTrue fails the test if condition is false
func (th *TestHelper) AssertTrue(condition bool, message string) {
	if !condition {
		th.t.Fatalf("%s: expected true, got false", message)
	}
}

// AssertFalse fails the test if condition is true
func (th *TestHelper) AssertFalse(condition bool, message string) {
	if condition {
		th.t.Fatalf("%s: expected false, got true", message)
	}
}

// CreateTestEntity creates a test entity with basic components
func CreateTestEntity(id uint64, entityType string) *entities.Entity {
	entity := &entities.Entity{
		ID:   id,
		Type: entityType,
		Tags: []string{"test"},
	}
	
	// Add basic components for testing
	entity.AddComponent(&entities.RenderComponent{
		Type:  "render",
		Color: []uint8{255, 0, 0, 255},
		Scale: 1.0,
	})
	
	return entity
}

// CreateTestComponent creates a test component
func CreateTestComponent(componentType string) entities.Component {
	switch componentType {
	case "render":
		return &entities.RenderComponent{
			Type:  "render",
			Color: []uint8{255, 255, 255, 255},
			Scale: 1.0,
		}
	case "physics":
		return &entities.PhysicsComponent{
			Type:     "physics",
			VelocityX: 0,
			VelocityY: 0,
			Mass:     1.0,
		}
	default:
		return &entities.BaseComponent{
			Type: componentType,
		}
	}
}

// WaitForCondition waits for a condition to be true or timeout
func WaitForCondition(condition func() bool, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}
	return false
}

// CleanupTestData cleans up test data
func CleanupTestData() {
	// Add any cleanup logic here
	// This could include removing test files, clearing caches, etc.
}
