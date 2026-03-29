package entities

import (
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"
)

// EventType represents the type of event
type EventType string

const (
	// Entity events
	EventEntityAdded    EventType = "entity_added"
	EventEntityRemoved  EventType = "entity_removed"
	EventEntityUpdated  EventType = "entity_updated"
	
	// Component events
	EventComponentAdded EventType = "component_added"
	EventComponentRemoved EventType = "component_removed"
	EventComponentUpdated EventType = "component_updated"
	
	// Interaction events
	EventBlockPlaced    EventType = "block_placed"
	EventBlockBroken    EventType = "block_broken"
	EventItemUsed       EventType = "item_used"
	EventItemCrafted    EventType = "item_crafted"
	
	// Combat events
	EventAttack         EventType = "attack"
	EventDamage         EventType = "damage"
	EventDeath          EventType = "death"
	EventHeal           EventType = "heal"
	
	// System events
	EventSystemStarted  EventType = "system_started"
	EventSystemStopped  EventType = "system_stopped"
	
	// World events
	EventWorldLoaded    EventType = "world_loaded"
	EventWorldSaved     EventType = "world_saved"
	EventChunkLoaded    EventType = "chunk_loaded"
	EventChunkUnloaded  EventType = "chunk_unloaded"
)

// Event represents a game event
type Event struct {
	Type      EventType       `json:"type"`
	Timestamp time.Time       `json:"timestamp"`
	Source    string          `json:"source"`
	Data      interface{}     `json:"data"`
	Priority  int             `json:"priority"`
	Cancelled bool            `json:"-"`
}

// EventHandler represents an event handler function
type EventHandler func(event Event)

// EventBus manages event publishing and subscription
type EventBus struct {
	subscribers map[EventType][]EventHandler
	handlers    map[string]EventHandler
	mutex       sync.RWMutex
	eventQueue  []Event
	maxQueue    int
	enabled     bool
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[EventType][]EventHandler),
		handlers:    make(map[string]EventHandler),
		eventQueue:  make([]Event, 0),
		maxQueue:    1000,
		enabled:     true,
	}
}

// Subscribe subscribes to an event type
func (eb *EventBus) Subscribe(eventType EventType, handler EventHandler) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()
	
	eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
	log.Printf("Subscribed to event: %s", eventType)
}

// Unsubscribe removes an event handler
func (eb *EventBus) Unsubscribe(eventType EventType, handler EventHandler) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()
	
	handlers := eb.subscribers[eventType]
	for i, h := range handlers {
		if reflect.ValueOf(h).Pointer() == reflect.ValueOf(handler).Pointer() {
			eb.subscribers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
	log.Printf("Unsubscribed from event: %s", eventType)
}

// Publish publishes an event to all subscribers
func (eb *EventBus) Publish(eventType EventType, data interface{}) {
	if !eb.enabled {
		return
	}
	
	event := Event{
		Type:      eventType,
		Timestamp: time.Now(),
		Source:    "system",
		Data:      data,
		Priority:  0,
		Cancelled: false,
	}
	
	eb.mutex.Lock()
	
	// Add to queue
	if len(eb.eventQueue) < eb.maxQueue {
		eb.eventQueue = append(eb.eventQueue, event)
	} else {
		log.Printf("Event queue full, dropping event: %s", eventType)
		eb.mutex.Unlock()
		return
	}
	
	eb.mutex.Unlock()
	
	// Process immediately for now (could be moved to batch processing)
	eb.processEvent(event)
}

// PublishWithSource publishes an event with a specific source
func (eb *EventBus) PublishWithSource(eventType EventType, source string, data interface{}) {
	if !eb.enabled {
		return
	}
	
	event := Event{
		Type:      eventType,
		Timestamp: time.Now(),
		Source:    source,
		Data:      data,
		Priority:  0,
		Cancelled: false,
	}
	
	eb.processEvent(event)
}

// PublishWithPriority publishes an event with a specific priority
func (eb *EventBus) PublishWithPriority(eventType EventType, source string, data interface{}, priority int) {
	if !eb.enabled {
		return
	}
	
	event := Event{
		Type:      eventType,
		Timestamp: time.Now(),
		Source:    source,
		Data:      data,
		Priority:  priority,
		Cancelled: false,
	}
	
	eb.processEvent(event)
}

// processEvent processes a single event
func (eb *EventBus) processEvent(event Event) {
	eb.mutex.RLock()
	handlers := eb.subscribers[event.Type]
	eb.mutex.RUnlock()
	
	for _, handler := range handlers {
		if event.Cancelled {
			break
		}
		
		// Run handler in goroutine to avoid blocking
		go func(h EventHandler) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Event handler panic: %v", r)
				}
			}()
			h(event)
		}(handler)
	}
}

// ProcessBatch processes events in the queue
func (eb *EventBus) ProcessBatch() {
	eb.mutex.Lock()
	events := make([]Event, len(eb.eventQueue))
	copy(events, eb.eventQueue)
	eb.eventQueue = eb.eventQueue[:0] // Clear queue
	eb.mutex.Unlock()
	
	// Sort by priority (higher priority first)
	for i := 0; i < len(events)-1; i++ {
		for j := i + 1; j < len(events); j++ {
			if events[i].Priority < events[j].Priority {
				events[i], events[j] = events[j], events[i]
			}
		}
	}
	
	// Process events
	for _, event := range events {
		eb.processEvent(event)
	}
}

// Enable enables the event bus
func (eb *EventBus) Enable() {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()
	eb.enabled = true
}

// Disable disables the event bus
func (eb *EventBus) Disable() {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()
	eb.enabled = false
}

// IsEnabled returns whether the event bus is enabled
func (eb *EventBus) IsEnabled() bool {
	eb.mutex.RLock()
	defer eb.mutex.RUnlock()
	return eb.enabled
}

// Clear clears all subscribers and the event queue
func (eb *EventBus) Clear() {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()
	
	eb.subscribers = make(map[EventType][]EventHandler)
	eb.eventQueue = eb.eventQueue[:0]
}

// GetQueueSize returns the current queue size
func (eb *EventBus) GetQueueSize() int {
	eb.mutex.RLock()
	defer eb.mutex.RUnlock()
	return len(eb.eventQueue)
}

// ============================================================================
// Event Data Types
// ============================================================================

// EntityEvent represents entity-related event data
type EntityEvent struct {
	Entity    *Entity `json:"entity"`
	Component string  `json:"component,omitempty"`
	OldValue  interface{} `json:"oldValue,omitempty"`
	NewValue  interface{} `json:"newValue,omitempty"`
}

// BlockEvent represents block-related event data
type BlockEvent struct {
	BlockType string  `json:"blockType"`
	Position  struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	} `json:"position"`
	PlayerID  string `json:"playerId,omitempty"`
	ToolUsed  string `json:"toolUsed,omitempty"`
}

// ItemEvent represents item-related event data
type ItemEvent struct {
	ItemType   string `json:"itemType"`
	Quantity   int    `json:"quantity"`
	PlayerID   string `json:"playerId,omitempty"`
	TargetID   string `json:"targetId,omitempty"`
	Success    bool   `json:"success"`
}

// CombatEvent represents combat-related event data
type CombatEvent struct {
	AttackerID string  `json:"attackerId"`
	TargetID   string  `json:"targetId"`
	Damage     float64 `json:"damage"`
	WeaponType string  `json:"weaponType"`
	Critical   bool    `json:"critical"`
	Killed     bool    `json:"killed"`
}

// SystemEvent represents system-related event data
type SystemEvent struct {
	SystemName string `json:"systemName"`
	Status     string `json:"status"` // "started", "stopped", "error"
	Message    string `json:"message,omitempty"`
}

// WorldEvent represents world-related event data
type WorldEvent struct {
	WorldName string `json:"worldName"`
	Position  struct {
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"position"`
	ChunkData interface{} `json:"chunkData,omitempty"`
}

// ============================================================================
// Utility Functions
// ============================================================================

// CreateEntityEvent creates an entity event
func CreateEntityEvent(entity *Entity, component string, oldValue, newValue interface{}) EntityEvent {
	return EntityEvent{
		Entity:    entity,
		Component: component,
		OldValue:  oldValue,
		NewValue:  newValue,
	}
}

// CreateBlockEvent creates a block event
func CreateBlockEvent(blockType string, x, y, z float64, playerID, toolUsed string) BlockEvent {
	event := BlockEvent{
		BlockType: blockType,
		PlayerID:  playerID,
		ToolUsed:  toolUsed,
	}
	event.Position.X = x
	event.Position.Y = y
	event.Position.Z = z
	return event
}

// CreateItemEvent creates an item event
func CreateItemEvent(itemType string, quantity int, playerID, targetID string, success bool) ItemEvent {
	return ItemEvent{
		ItemType: itemType,
		Quantity: quantity,
		PlayerID: playerID,
		TargetID: targetID,
		Success:  success,
	}
}

// CreateCombatEvent creates a combat event
func CreateCombatEvent(attackerID, targetID string, damage float64, weaponType string, critical, killed bool) CombatEvent {
	return CombatEvent{
		AttackerID: attackerID,
		TargetID:   targetID,
		Damage:     damage,
		WeaponType: weaponType,
		Critical:   critical,
		Killed:     killed,
	}
}

// CreateSystemEvent creates a system event
func CreateSystemEvent(systemName, status, message string) SystemEvent {
	return SystemEvent{
		SystemName: systemName,
		Status:     status,
		Message:    message,
	}
}

// CreateWorldEvent creates a world event
func CreateWorldEvent(worldName string, x, y int, chunkData interface{}) WorldEvent {
	event := WorldEvent{
		WorldName: worldName,
		ChunkData: chunkData,
	}
	event.Position.X = x
	event.Position.Y = y
	return event
}

// ============================================================================
// Event Listeners
// ============================================================================

// EventListener provides a convenient way to listen to multiple events
type EventListener struct {
	bus       *EventBus
	handlers  map[EventType]EventHandler
	active    bool
}

// NewEventListener creates a new event listener
func NewEventListener(bus *EventBus) *EventListener {
	return &EventListener{
		bus:      bus,
		handlers: make(map[EventType]EventHandler),
		active:   true,
	}
}

// Listen adds a handler for an event type
func (el *EventListener) Listen(eventType EventType, handler EventHandler) {
	if !el.active {
		return
	}
	
	el.handlers[eventType] = handler
	el.bus.Subscribe(eventType, handler)
}

// StopListening removes all handlers
func (el *EventListener) StopListening() {
	if !el.active {
		return
	}
	
	el.active = false
	for eventType, handler := range el.handlers {
		el.bus.Unsubscribe(eventType, handler)
	}
	el.handlers = make(map[EventType]EventHandler)
}

// IsActive returns whether the listener is active
func (el *EventListener) IsActive() bool {
	return el.active
}

// ============================================================================
// Event Filters
// ============================================================================

// EventFilter filters events based on criteria
type EventFilter struct {
	sourceFilter  map[string]bool
	typeFilter    map[EventType]bool
	priorityMin   int
	priorityMax   int
	customFilter  func(Event) bool
}

// NewEventFilter creates a new event filter
func NewEventFilter() *EventFilter {
	return &EventFilter{
		sourceFilter: make(map[string]bool),
		typeFilter:   make(map[EventType]bool),
		priorityMin:  0,
		priorityMax:  100,
	}
}

// BySource filters events by source
func (ef *EventFilter) BySource(source string) *EventFilter {
	ef.sourceFilter[source] = true
	return ef
}

// ByType filters events by type
func (ef *EventFilter) ByType(eventType EventType) *EventFilter {
	ef.typeFilter[eventType] = true
	return ef
}

// ByPriority filters events by priority range
func (ef *EventFilter) ByPriority(min, max int) *EventFilter {
	ef.priorityMin = min
	ef.priorityMax = max
	return ef
}

// Custom adds a custom filter function
func (ef *EventFilter) Custom(filter func(Event) bool) *EventFilter {
	ef.customFilter = filter
	return ef
}

// Matches checks if an event matches the filter criteria
func (ef *EventFilter) Matches(event Event) bool {
	// Check source filter
	if len(ef.sourceFilter) > 0 {
		if !ef.sourceFilter[event.Source] {
			return false
		}
	}
	
	// Check type filter
	if len(ef.typeFilter) > 0 {
		if !ef.typeFilter[event.Type] {
			return false
		}
	}
	
	// Check priority filter
	if event.Priority < ef.priorityMin || event.Priority > ef.priorityMax {
		return false
	}
	
	// Check custom filter
	if ef.customFilter != nil && !ef.customFilter(event) {
		return false
	}
	
	return true
}

// ============================================================================
// Event Logger
// ============================================================================

// EventLogger logs events to console or file
type EventLogger struct {
	bus        *EventBus
	filter     *EventFilter
	logLevel   int
	logToConsole bool
	logToFile   bool
	filePath    string
}

// NewEventLogger creates a new event logger
func NewEventLogger(bus *EventBus) *EventLogger {
	return &EventLogger{
		bus:          bus,
		filter:       NewEventFilter(),
		logLevel:     0,
		logToConsole: true,
		logToFile:    false,
	}
}

// Start begins logging events
func (el *EventLogger) Start() {
	handler := func(event Event) {
		if el.filter.Matches(event) {
			el.logEvent(event)
		}
	}
	
	// Subscribe to all event types
	eventTypes := []EventType{
		EventEntityAdded, EventEntityRemoved, EventEntityUpdated,
		EventComponentAdded, EventComponentRemoved, EventComponentUpdated,
		EventBlockPlaced, EventBlockBroken, EventItemUsed, EventItemCrafted,
		EventAttack, EventDamage, EventDeath, EventHeal,
		EventSystemStarted, EventSystemStopped,
		EventWorldLoaded, EventWorldSaved, EventChunkLoaded, EventChunkUnloaded,
	}
	
	for _, eventType := range eventTypes {
		el.bus.Subscribe(eventType, handler)
	}
}

// Stop stops logging events
func (el *EventLogger) Stop() {
	// Implementation would unsubscribe from all events
}

// SetFilter sets the event filter
func (el *EventLogger) SetFilter(filter *EventFilter) {
	el.filter = filter
}

// logEvent logs a single event
func (el *EventLogger) logEvent(event Event) {
	message := fmt.Sprintf("[%s] %s from %s", event.Timestamp.Format("15:04:05"), event.Type, event.Source)
	
	if el.logToConsole {
		log.Println(message)
	}
	
	// File logging would be implemented here
	if el.logToFile && el.filePath != "" {
		// Write to file
	}
}
