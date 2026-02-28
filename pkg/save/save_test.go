package save

import (
	"os"
	"testing"
	"time"

	"tesselbox/pkg/items"
	"tesselbox/pkg/player"
	"tesselbox/pkg/world"
)

func TestSaveManager(t *testing.T) {
	// Create temporary test directory
	testDir := "test_saves"
	defer os.RemoveAll(testDir)

	// Create test world and player
	testWorld := world.NewWorld("test_world")
	testPlayer := player.NewPlayer(100, 200)
	testPlayer.Health = 75.5
	testPlayer.SelectedSlot = 3

	testInventory := items.NewInventory(10)
	testInventory.AddItem(items.DIRT_BLOCK, 5)
	testInventory.AddItem(items.STONE_BLOCK, 3)
	testInventory.SelectSlot(2)

	// Create game state
	gameState := &GameState{
		World:      testWorld,
		Player:     testPlayer,
		Inventory:  testInventory,
		CameraX:    50.5,
		CameraY:    100.5,
		InMenu:     false,
		InGame:     true,
		InCrafting: false,
	}

	// Test save manager
	saveManager := NewSaveManager(testDir, "test_player")

	// Test saving
	err := saveManager.SaveGame(gameState)
	if err != nil {
		t.Fatalf("Failed to save game: %v", err)
	}

	// Test loading
	saveData, err := saveManager.LoadGame()
	if err != nil {
		t.Fatalf("Failed to load game: %v", err)
	}

	// Verify save data
	if saveData.PlayerX != 100.0 {
		t.Errorf("Expected player X 100.0, got %f", saveData.PlayerX)
	}
	if saveData.PlayerY != 200.0 {
		t.Errorf("Expected player Y 200.0, got %f", saveData.PlayerY)
	}
	if saveData.PlayerHealth != 75.5 {
		t.Errorf("Expected player health 75.5, got %f", saveData.PlayerHealth)
	}
	if saveData.SelectedSlot != 3 {
		t.Errorf("Expected selected slot 3, got %d", saveData.SelectedSlot)
	}
	if len(saveData.InventorySlots) != 10 {
		t.Errorf("Expected 10 inventory slots, got %d", len(saveData.InventorySlots))
	}

	// Test applying save data
	newWorld := world.NewWorld("test_world")
	newPlayer := player.NewPlayer(0, 0)
	newInventory := items.NewInventory(10)

	newGameState := &GameState{
		World:      newWorld,
		Player:     newPlayer,
		Inventory:  newInventory,
		CameraX:    0,
		CameraY:    0,
		InMenu:     true,
		InGame:     false,
		InCrafting: true,
	}

	err = saveManager.ApplySaveData(saveData, newGameState)
	if err != nil {
		t.Fatalf("Failed to apply save data: %v", err)
	}

	// Verify applied data
	if newGameState.Player.X != 100.0 {
		t.Errorf("Expected applied player X 100.0, got %f", newGameState.Player.X)
	}
	if newGameState.Player.Y != 200.0 {
		t.Errorf("Expected applied player Y 200.0, got %f", newGameState.Player.Y)
	}
	if newGameState.Player.Health != 75.5 {
		t.Errorf("Expected applied player health 75.5, got %f", newGameState.Player.Health)
	}
	if newGameState.CameraX != 50.5 {
		t.Errorf("Expected applied camera X 50.5, got %f", newGameState.CameraX)
	}
	if newGameState.InGame != true {
		t.Errorf("Expected applied inGame true, got %v", newGameState.InGame)
	}
}

func TestAutoSaver(t *testing.T) {
	// Create temporary test directory
	testDir := "test_autosave"
	defer os.RemoveAll(testDir)

	// Create test game state
	testWorld := world.NewWorld("autosave_world")
	testPlayer := player.NewPlayer(0, 0)
	testInventory := items.NewInventory(5)

	gameState := &GameState{
		World:      testWorld,
		Player:     testPlayer,
		Inventory:  testInventory,
		CameraX:    0,
		CameraY:    0,
		InMenu:     false,
		InGame:     true,
		InCrafting: false,
	}

	saveManager := NewSaveManager(testDir, "autosave_player")
	autoSaver := NewAutoSaver(saveManager, gameState, 100*time.Millisecond)

	// Test force save
	err := autoSaver.ForceSave()
	if err != nil {
		t.Fatalf("Failed to force save: %v", err)
	}

	// Verify save exists
	saveData, err := saveManager.LoadGame()
	if err != nil {
		t.Fatalf("Failed to load autosaved game: %v", err)
	}

	if saveData.WorldName != "test_autosave" {
		t.Errorf("Expected world name 'test_autosave', got '%s'", saveData.WorldName)
	}

	// Test auto-saver start/stop
	autoSaver.Start()
	time.Sleep(150 * time.Millisecond) // Wait for at least one auto-save
	autoSaver.Stop()

	// Verify multiple saves don't cause errors
	for i := 0; i < 3; i++ {
		err = autoSaver.ForceSave()
		if err != nil {
			t.Errorf("Force save %d failed: %v", i, err)
		}
	}
}

func TestSaveManagerListAndDelete(t *testing.T) {
	// Create temporary test directory
	testDir := "test_list_delete"
	defer os.RemoveAll(testDir)

	saveManager := NewSaveManager(testDir, "list_test_player")

	// Create and save a game state
	testWorld := world.NewWorld("list_test_world")
	testPlayer := player.NewPlayer(0, 0)
	testInventory := items.NewInventory(5)

	gameState := &GameState{
		World:      testWorld,
		Player:     testPlayer,
		Inventory:  testInventory,
		CameraX:    0,
		CameraY:    0,
		InMenu:     false,
		InGame:     true,
		InCrafting: false,
	}

	err := saveManager.SaveGame(gameState)
	if err != nil {
		t.Fatalf("Failed to save game: %v", err)
	}

	// Test listing saves
	saves, err := saveManager.ListSaves()
	if err != nil {
		t.Fatalf("Failed to list saves: %v", err)
	}

	if len(saves) != 1 {
		t.Errorf("Expected 1 save, got %d", len(saves))
	}

	expectedSaveName := "player_list_test_player.json"
	if saves[0] != expectedSaveName {
		t.Errorf("Expected save name '%s', got '%s'", expectedSaveName, saves[0])
	}

	// Test getting save info
	saveInfo, err := saveManager.GetSaveInfo()
	if err != nil {
		t.Fatalf("Failed to get save info: %v", err)
	}

	if saveInfo.PlayerName != "list_test_player" {
		t.Errorf("Expected player name 'list_test_player', got '%s'", saveInfo.PlayerName)
	}
	if saveInfo.WorldName != "test_list_delete" {
		t.Errorf("Expected world name 'test_list_delete', got '%s'", saveInfo.WorldName)
	}

	// Test deleting save
	err = saveManager.DeleteSave()
	if err != nil {
		t.Fatalf("Failed to delete save: %v", err)
	}

	// Verify save is deleted
	saves, err = saveManager.ListSaves()
	if err != nil {
		t.Fatalf("Failed to list saves after deletion: %v", err)
	}

	if len(saves) != 0 {
		t.Errorf("Expected 0 saves after deletion, got %d", len(saves))
	}
}
