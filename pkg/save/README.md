# TesselBox Save System

## Overview
The TesselBox save system provides comprehensive game state persistence, including world data, player state, inventory, and game settings. The system is implemented in `pkg/save/save.go` and integrates seamlessly with the existing world storage system.

## Features

### Core Functionality
- **Complete Game State Saving**: Saves player position, health, inventory, camera position, and game state flags
- **World Integration**: Works with existing `pkg/world/storage.go` for chunk persistence
- **Player-Specific Saves**: Each player gets their own save file within a world
- **Multiple World Support**: Separate save directories for different worlds

### Auto-Save System
- **Configurable Intervals**: Auto-save frequency can be adjusted (default: 5 minutes)
- **Background Operation**: Runs in a separate goroutine without blocking gameplay
- **Force Save**: Manual save triggers available via API or keyboard shortcuts

### Save Management
- **Save Slots**: Multiple save files supported per world/player combination
- **Save Metadata**: Track save time, version, player stats, and world information
- **Delete Functionality**: Clean removal of unwanted saves
- **List Operations**: Browse available saves and their metadata

## File Structure
```
saves/
├── world_name/
│   ├── player_user_id.json      # Player-specific save data
│   └── metadata.json            # World metadata (from world storage)
worlds/
├── world_name/
│   ├── chunk_0_0.json          # World chunks (from world storage)
│   ├── chunk_0_1.json
│   └── metadata.json
```

## Usage

### Keyboard Shortcuts
- **F5**: Manual save
- **F9**: Manual load

### API Usage

#### Basic Save/Load
```go
// Create save manager
saveManager := save.NewSaveManager("world_name", "player_id")

// Save game
err := saveManager.SaveGame(gameState)

// Load game
saveData, err := saveManager.LoadGame()
err = saveManager.ApplySaveData(saveData, gameState)
```

#### Auto-Save
```go
// Create auto-saver with 5-minute interval
autoSaver := save.NewAutoSaver(saveManager, gameState, 5*time.Minute)

// Start auto-saving
autoSaver.Start()

// Stop auto-saving
autoSaver.Stop()

// Force immediate save
err := autoSaver.ForceSave()
```

#### Save Management
```go
// List all saves for this world
saves, err := saveManager.ListSaves()

// Get save metadata
saveInfo, err := saveManager.GetSaveInfo()

// Delete a save
err := saveManager.DeleteSave()
```

## Data Structures

### SaveData
Contains the complete serializable game state:
- Player position, velocity, health
- Inventory contents with item types and quantities
- Camera position
- Game state flags (menu, crafting, etc.)
- World metadata

### GameState
Runtime representation used for save/load operations:
- References to actual game objects
- Used to create SaveData for saving
- Target for applying loaded save data

### SaveInfo
Metadata about saved games:
- Save time and version
- Player statistics
- World information
- Chunk count

## Integration Points

### Main Game Loop
The save system is integrated into `cmd/client/main.go`:
- Save manager initialized in `NewGame()`
- Auto-saver started in `main()`
- Keyboard shortcuts handled in `handleGameInput()`

### World Storage
Leverages existing `pkg/world/storage.go`:
- Chunks saved/loaded via `WorldStorage`
- World metadata managed separately
- Seamless integration with world generation

### Authentication
Integrates with `pkg/auth`:
- Player ID used for save file naming
- Guest mode supported with default player name
- Multiple players per world supported

## Testing
Comprehensive test suite in `pkg/save/save_test.go`:
- Unit tests for all major functionality
- Auto-saver integration tests
- Save management operations
- Error handling verification

Run tests with:
```bash
go test ./pkg/save -v
```

## Configuration

### Auto-Save Interval
Default is 5 minutes, configurable via:
```go
autoSaver := save.NewAutoSaver(saveManager, gameState, interval)
```

### Save Directory Structure
Saves are stored in `saves/world_name/` by default, with worlds in `worlds/world_name/`.

## Error Handling
The save system provides comprehensive error handling:
- Graceful degradation when save directory is inaccessible
- Validation of save data integrity
- Recovery from corrupted save files
- Logging of save/load operations

## Performance Considerations
- Incremental saves: Only modified chunks are saved
- Background operation: Auto-save doesn't block gameplay
- Efficient serialization: JSON with proper formatting
- Memory management: Proper cleanup of temporary resources

## Future Enhancements
Potential improvements for Phase 2.1+:
- Compression for save files
- Save file versioning and migration
- Cloud save integration
- Save file encryption
- Multiple save slots per player
- Save file backup and recovery
