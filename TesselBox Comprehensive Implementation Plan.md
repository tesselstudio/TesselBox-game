# TesselBox Implementation Plan
## Overview
This plan addresses critical security issues, implements missing core functionality, and adds new features to the TesselBox hexagonal voxel game.
## Phase 1: Critical Security & Startup Fixes (Immediate Priority)
### 1.1 done

### 1.2 done

### 1.3 Add CSRF Protection to OAuth
**Problem:** OAuth state parameter is hardcoded to "state"
**Solution:** Generate and validate cryptographically secure state tokens
## Phase 2: Core Game Systems (Week 1-2)
### 2.1 Complete Save/Load System
**Current State:** Storage code exists but player state not persisted
**Implementation:**
* Create `SaveData` struct with world, player, inventory state
* Add `pkg/save/save.go` for unified save management
* Implement auto-save with configurable interval
* Add save slot management (multiple worlds)
### 2.2 Complete Inventory System Backend
**Current State:** `pkg/items/items.go` has basic inventory but not integrated
**Implementation:**
* Connect inventory to player in Game struct
* Add item pickup from mining
* Implement hotbar rendering with item icons
* Add item drop and throw mechanics
### 2.3 Implement Crafting System
**Current State:** `pkg/crafting/crafting.go` has recipe logic
**Implementation:**
* Create recipes.json with basic recipes
* Add crafting UI overlay
* Connect crafting to inventory
* Add crafting stations (workbench, furnace)
## Phase 3: Rendering & Performance (Week 2-3)
### 3.1 Fix Chunk Memory Management
**Problem:** Chunks never unloaded, memory leak
**Solution:**
* Implement actual chunk unloading in `UnloadDistantChunks()`
* Save modified chunks before unloading
* Add chunk loading/unloading events
### 3.2 Implement Spatial Partitioning
**Problem:** Linear searches for hexagon lookups
**Solution:**
* Add spatial hash grid for O(1) hexagon lookups
* Optimize collision detection
* Batch draw calls
### 3.3 Optimize Particle System
**Problem:** Particles create unnecessary allocations
**Solution:**
* Implement object pooling for particles
* Pre-allocate particle arrays
* Add particle lifecycle management
## Phase 4: Block Placement & Mining Enhancement (Week 3-4)
### 4.1 Block Placement System
* Right-click to place blocks
* Ghost preview before placement
* Placement validation (collision check)
### 4.2 Mining Refinement
* Tool-based mining speeds
* Block hardness affects mining time
* Visual feedback (cracks, progress bar)
* Drop items when block destroyed
## Phase 5: Day/Night & Weather (Week 4-5)
### 5.1 Day/Night Cycle
* Time system with configurable day length
* Lighting changes based on time
* Ambient color shifts
### 5.2 Weather System
* Weather states (clear, rain, storm)
* Particle effects for rain/snow
* Weather affects gameplay (visibility, movement)
## Phase 6: Combat & Creatures (Week 5-6)
### 6.1 Basic Combat System
* Health/damage calculations
* Attack animations
* Knockback physics
### 6.2 Creature System
* Basic AI state machine (idle, wander, chase)
* Pathfinding (A* algorithm)
* Spawn system based on biome
## Implementation Order (Priority)
3. **HIGH:** Complete save/load (Phase 2.1)
4. **HIGH:** Complete inventory integration (Phase 2.2)
5. **HIGH:** Implement crafting UI (Phase 2.3)
6. **MEDIUM:** Fix chunk memory (Phase 3.1)
7. **MEDIUM:** Block placement (Phase 4.1)
8. **MEDIUM:** Mining refinement (Phase 4.2)
9. **LOW:** Day/night cycle (Phase 5.1)
10. **LOW:** Combat system (Phase 6.1)
