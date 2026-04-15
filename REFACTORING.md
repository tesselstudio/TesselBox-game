# Tesselbox Architecture Refactoring

## Overview
This document outlines the refactoring of the Tesselbox codebase to address architectural issues:
- 4k-line monolithic `main.go`
- Duplicate Game struct definitions
- Unused subsystems and packages
- Duplicate configuration paths
- Boolean flag state management
- Unwired ECS system

## Changes Made

### 1. Centralized Configuration (`pkg/config/paths.go`)
**Problem**: Path logic duplicated across 7+ files
- `cmd/main.go`, `pkg/gui/menu.go`, `pkg/save/save.go`, `pkg/chest/chest.go`, `pkg/world/storage.go`, `pkg/skin/editor.go`, `pkg/settings/settings.go`

**Solution**: Single source of truth for all paths
```go
GetTesselboxDir()      // ~/.tesselbox
GetSavesDir()          // ~/.tesselbox/saves
GetWorldSaveDir()      // ~/.tesselbox/saves/{world}
GetWorldsDir()         // ~/.tesselbox/worlds
GetSkinsDir()          // ~/.tesselbox/skins
GetChestFile()         // ~/.tesselbox/saves/{world}/chests.json
EnsureDirectories()    // Create all necessary dirs
```

**Migration**: Update all packages to use `config.GetTesselboxDir()` instead of local implementations.

### 2. State Machine (`pkg/ui/state.go`)
**Problem**: 6 boolean flags with hardcoded modal ordering
```go
inMenu, inGame, inCrafting, inPluginUI, inSkinEditor, isDead
```

**Solution**: Centralized state machine with clear transitions
```go
type GameState int
const (
    StateMenu, StateGame, StateCrafting, StateBackpack,
    StateChest, StatePluginUI, StateSkinEditor, StateDeathScreen
)

type StateManager struct {
    current GameState
}
```

**Benefits**:
- Single source of truth for game state
- Clear state transitions
- No conflicting boolean flags
- Easy to add new states
- Thread-safe with RWMutex

### 3. GameManager (`pkg/game/manager.go`)
**Problem**: Game struct owns 50+ fields, acts as god object

**Solution**: GameManager coordinates subsystems
- Owns all subsystem instances
- Provides unified Update() method
- Delegates to subsystems based on state
- Cleaner separation of concerns

**Structure**:
```
GameManager
├── Core Systems (World, Player, Inventory, StateManager)
├── Crafting (CraftingSystem, CraftingUI)
├── Plugins (PluginManager, PluginUI, PluginInstaller)
├── Skin (SkinEditor)
├── Input (InputManager)
├── Save (SaveManager, AutoSaver)
├── Time (DayNightCycle, WeatherSystem)
├── Audio (AudioManager, SoundLibrary)
├── Debug (RecoveryHandler, Profiler)
├── Survival (SurvivalManager, EquipmentSet, HealthSystem, BackpackUI)
├── Chest (ChestManager, ChestUI)
├── Combat (WeaponSystem)
├── UI Effects (DamageIndicators, ScreenFlash, DirectionalHitManager, DeathScreen)
├── Enemies (ZombieSpawner)
└── Rendering (Camera, Layers, Object Pools)
```

### 4. Game Launcher (`pkg/game/launcher.go`)
**Problem**: Multiple launch paths mixed in main.go

**Solution**: Separate launcher module
- `LaunchGUI()` - Ebiten GUI mode
- `LaunchCLI()` - Headless CLI mode (future)
- `LaunchTUI()` - Terminal UI mode (future)
- `LaunchConfig` - Unified configuration

### 5. Duplicate Game Struct Removal
**Problem**: `pkg/render/renderer.go` has its own Game struct

**Solution**: 
- Remove duplicate Game struct from pkg/render
- Use GameManager for all rendering
- Consolidate rendering logic into main game loop

### 6. Unused Packages
**Current Status**: 20+ packages not wired into binary
- `pkg/anticheat`, `pkg/boss`, `pkg/chat`, `pkg/commands`, `pkg/cosmetics`
- `pkg/creatures`, `pkg/dungeons`, `pkg/economy/*`, `pkg/land`, `pkg/mail`
- `pkg/minigames`, `pkg/moderation`, `pkg/organisms`, `pkg/permissions`
- `pkg/pvp`, `pkg/quests`, `pkg/rollback`, `pkg/security`, `pkg/social/*`
- `pkg/stats`, `pkg/status`, `pkg/village`, `pkg/vote`, `pkg/warps`

**Action**: Keep packages but don't import them. They can be integrated later when needed.

### 7. ECS System Integration (Future)
**Current Status**: `pkg/entities` implements full ECS but is unused

**Plan**:
1. Create `pkg/game/ecs.go` to wire EntityManager into GameManager
2. Migrate entity-based systems to use ECS
3. Gradually move from ad-hoc object management to ECS

## Migration Path

### Phase 1: Foundation (Current)
- [x] Create `pkg/config/paths.go` - Centralize paths
- [x] Create `pkg/ui/state.go` - State machine
- [x] Create `pkg/game/manager.go` - GameManager
- [x] Create `pkg/game/launcher.go` - Launch paths

### Phase 2: Integration
- [x] Update `cmd/main.go` to use config package
- [x] Update `pkg/gui/menu.go` to use config package
- [x] Update `pkg/save/save.go` to use config package
- [x] Update `pkg/chest/chest.go` to use config package
- [x] Update `pkg/skin/editor.go` to use config package
- [x] Update `pkg/world/storage.go` to use config package
- [x] Update `pkg/settings/settings.go` to use config package
- [ ] Replace boolean flags with StateManager (main.go)
- [ ] Remove duplicate Game struct from pkg/render

### Phase 3: ECS Integration
- [ ] Wire EntityManager into GameManager
- [ ] Create ECS-based systems for entities
- [ ] Migrate zombie/creature systems to ECS

### Phase 4: Cleanup
- [ ] Split main.go into smaller modules
- [ ] Remove unused imports
- [ ] Add comprehensive tests

## Files to Update

### Immediate (Phase 2)
1. `cmd/main.go` - Use GameManager, remove duplicate code
2. `pkg/gui/menu.go` - Use config.GetTesselboxDir()
3. `pkg/save/save.go` - Use config paths
4. `pkg/chest/chest.go` - Use config paths
5. `pkg/world/storage.go` - Use config paths
6. `pkg/skin/editor.go` - Use config paths
7. `pkg/settings/settings.go` - Use config paths
8. `pkg/render/renderer.go` - Remove duplicate Game struct

### Future (Phase 3+)
1. `pkg/entities/system.go` - Wire into GameManager
2. `pkg/enemies/zombie.go` - Use ECS
3. Various unused packages - Integrate as needed

## Benefits

1. **Reduced Complexity**: 4k-line main.go → modular subsystems
2. **Single Responsibility**: Each package has clear purpose
3. **Testability**: GameManager can be tested independently
4. **Extensibility**: New systems can be added to GameManager
5. **Maintainability**: Clear state transitions, no boolean flag conflicts
6. **Reusability**: GameManager can be used in different contexts (GUI, CLI, TUI)
7. **Future-Proof**: ECS system ready for integration

## Testing Strategy

1. Unit tests for StateManager
2. Integration tests for GameManager
3. Launcher tests for different entry points
4. Config path tests for all packages
5. State transition tests

## Rollback Plan

If issues arise:
1. Keep old Game struct in main.go temporarily
2. Gradually migrate systems to GameManager
3. Run both in parallel during transition
4. Remove old code once all systems migrated
