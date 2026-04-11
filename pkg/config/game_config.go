package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// GameConfig holds all game configuration settings
type GameConfig struct {
	mu sync.RWMutex

	// Key bindings
	KeyBindings KeyBindingConfig `json:"keyBindings"`

	// Video settings
	Video VideoConfig `json:"video"`

	// Audio settings
	Audio AudioConfig `json:"audio"`

	// Accessibility settings
	Accessibility AccessibilityConfig `json:"accessibility"`

	// Gameplay settings
	Gameplay GameplayConfig `json:"gameplay"`

	// Internal
	path string
}

// KeyBindingConfig holds all key bindings
type KeyBindingConfig struct {
	MoveUp      string `json:"moveUp"`      // Default: "W"
	MoveDown    string `json:"moveDown"`    // Default: "S"
	MoveLeft    string `json:"moveLeft"`    // Default: "A"
	MoveRight   string `json:"moveRight"`   // Default: "D"
	Jump        string `json:"jump"`        // Default: "Space"
	Sprint      string `json:"sprint"`     // Default: "Shift"
	Inventory   string `json:"inventory"`   // Default: "I"
	Crafting    string `json:"crafting"`    // Default: "C"
	BlockLib    string `json:"blockLib"`    // Default: "B"
	DropItem    string `json:"dropItem"`    // Default: "Q"
	QuickSave   string `json:"quickSave"`   // Default: "F5"
	QuickLoad   string `json:"quickLoad"`   // Default: "F9"
	DebugToggle string `json:"debugToggle"` // Default: "F3"
}

// VideoConfig holds video/display settings
type VideoConfig struct {
	ResolutionX     int    `json:"resolutionX"`     // Default: 1280
	ResolutionY     int    `json:"resolutionY"`     // Default: 720
	Fullscreen      bool   `json:"fullscreen"`      // Default: false
	VSync           bool   `json:"vsync"`           // Default: true
	RenderDistance  int    `json:"renderDistance"`  // Default: 8
	UIScale         float64 `json:"uiScale"`        // Default: 1.0
	ShowFPS         bool   `json:"showFps"`         // Default: false
}

// AudioConfig holds audio settings
type AudioConfig struct {
	MasterVolume float64 `json:"masterVolume"` // 0.0-1.0, default: 1.0
	MusicVolume  float64 `json:"musicVolume"`  // 0.0-1.0, default: 0.7
	SFXVolume    float64 `json:"sfxVolume"`    // 0.0-1.0, default: 1.0
	UISounds     bool    `json:"uiSounds"`     // default: true
}

// AccessibilityConfig holds accessibility settings
type AccessibilityConfig struct {
	ColorblindMode string  `json:"colorblindMode"` // "none", "protanopia", "deuteranopia", "tritanopia"
	HighContrast   bool    `json:"highContrast"`
	LargeFonts     bool    `json:"largeFonts"`
	ReduceMotion   bool    `json:"reduceMotion"`
	ScreenReader   bool    `json:"screenReader"`
}

// GameplayConfig holds gameplay settings
type GameplayConfig struct {
	AutoSave       bool `json:"autoSave"`       // default: true
	AutoSaveInterval int  `json:"autoSaveInterval"` // minutes, default: 5
	CreativeMode   bool `json:"creativeMode"`   // default: false
	ShowTutorials  bool `json:"showTutorials"`    // default: true
}

// NewGameConfig creates a new config with default values
func NewGameConfig(configPath string) *GameConfig {
	return &GameConfig{
		path: configPath,
		KeyBindings: KeyBindingConfig{
			MoveUp:      "W",
			MoveDown:    "S",
			MoveLeft:    "A",
			MoveRight:   "D",
			Jump:        "Space",
			Sprint:      "Shift",
			Inventory:   "I",
			Crafting:    "C",
			BlockLib:    "B",
			DropItem:    "Q",
			QuickSave:   "F5",
			QuickLoad:   "F9",
			DebugToggle: "F3",
		},
		Video: VideoConfig{
			ResolutionX:    1280,
			ResolutionY:    720,
			Fullscreen:     false,
			VSync:          true,
			RenderDistance: 8,
			UIScale:        1.0,
			ShowFPS:        false,
		},
		Audio: AudioConfig{
			MasterVolume: 1.0,
			MusicVolume:  0.7,
			SFXVolume:    1.0,
			UISounds:     true,
		},
		Accessibility: AccessibilityConfig{
			ColorblindMode: "none",
			HighContrast:   false,
			LargeFonts:     false,
			ReduceMotion:   false,
			ScreenReader:   false,
		},
		Gameplay: GameplayConfig{
			AutoSave:         true,
			AutoSaveInterval: 5,
			CreativeMode:     false,
			ShowTutorials:    true,
		},
	}
}

// Load loads config from disk
func (gc *GameConfig) Load() error {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	data, err := os.ReadFile(gc.path)
	if err != nil {
		if os.IsNotExist(err) {
			// No config exists, save defaults
			return gc.Save()
		}
		return fmt.Errorf("failed to read config: %w", err)
	}

	if err := json.Unmarshal(data, gc); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	return nil
}

// Save saves config to disk
func (gc *GameConfig) Save() error {
	gc.mu.RLock()
	defer gc.mu.RUnlock()

	// Ensure directory exists
	dir := filepath.Dir(gc.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(gc, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(gc.path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// GetKeyBinding returns a key binding
func (gc *GameConfig) GetKeyBinding(action string) string {
	gc.mu.RLock()
	defer gc.mu.RUnlock()

	switch action {
	case "moveUp":
		return gc.KeyBindings.MoveUp
	case "moveDown":
		return gc.KeyBindings.MoveDown
	case "moveLeft":
		return gc.KeyBindings.MoveLeft
	case "moveRight":
		return gc.KeyBindings.MoveRight
	case "jump":
		return gc.KeyBindings.Jump
	case "sprint":
		return gc.KeyBindings.Sprint
	case "inventory":
		return gc.KeyBindings.Inventory
	case "crafting":
		return gc.KeyBindings.Crafting
	case "blockLib":
		return gc.KeyBindings.BlockLib
	case "dropItem":
		return gc.KeyBindings.DropItem
	case "quickSave":
		return gc.KeyBindings.QuickSave
	case "quickLoad":
		return gc.KeyBindings.QuickLoad
	case "debugToggle":
		return gc.KeyBindings.DebugToggle
	default:
		return ""
	}
}

// SetKeyBinding sets a key binding
func (gc *GameConfig) SetKeyBinding(action, key string) {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	switch action {
	case "moveUp":
		gc.KeyBindings.MoveUp = key
	case "moveDown":
		gc.KeyBindings.MoveDown = key
	case "moveLeft":
		gc.KeyBindings.MoveLeft = key
	case "moveRight":
		gc.KeyBindings.MoveRight = key
	case "jump":
		gc.KeyBindings.Jump = key
	case "sprint":
		gc.KeyBindings.Sprint = key
	case "inventory":
		gc.KeyBindings.Inventory = key
	case "crafting":
		gc.KeyBindings.Crafting = key
	case "blockLib":
		gc.KeyBindings.BlockLib = key
	case "dropItem":
		gc.KeyBindings.DropItem = key
	case "quickSave":
		gc.KeyBindings.QuickSave = key
	case "quickLoad":
		gc.KeyBindings.QuickLoad = key
	case "debugToggle":
		gc.KeyBindings.DebugToggle = key
	}
}

// GetVideo returns video config (read-only copy)
func (gc *GameConfig) GetVideo() VideoConfig {
	gc.mu.RLock()
	defer gc.mu.RUnlock()
	return gc.Video
}

// SetVideo updates video config
func (gc *GameConfig) SetVideo(v VideoConfig) {
	gc.mu.Lock()
	defer gc.mu.Unlock()
	gc.Video = v
}

// GetAudio returns audio config
func (gc *GameConfig) GetAudio() AudioConfig {
	gc.mu.RLock()
	defer gc.mu.RUnlock()
	return gc.Audio
}

// SetAudio updates audio config
func (gc *GameConfig) SetAudio(a AudioConfig) {
	gc.mu.Lock()
	defer gc.mu.Unlock()
	gc.Audio = a
}
