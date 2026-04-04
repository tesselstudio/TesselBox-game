package settings

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

// Settings holds all game configuration
// Auto-saved to ~/.tesselbox/settings.yaml
type Settings struct {
	// Audio settings
	Volume float64 `yaml:"volume"`
	Mute   bool    `yaml:"mute"`

	// Display settings
	Fullscreen bool `yaml:"fullscreen"`
	VSync      bool `yaml:"vsync"`
	ShowDebug  bool `yaml:"show_debug"`

	// Gameplay settings
	AutoSaveInterval int  `yaml:"autosave_interval"`
	CreativeMode     bool `yaml:"creative_mode_default"`
	ChunkRenderDist  int  `yaml:"chunk_render_distance"`

	// Input settings
	Keybinds map[string]string `yaml:"keybinds"`

	// Window settings
	WindowWidth  int `yaml:"window_width"`
	WindowHeight int `yaml:"window_height"`

	mu sync.RWMutex
}

// DefaultSettings returns settings with sensible defaults
func DefaultSettings() *Settings {
	return &Settings{
		Volume:           0.7,
		Mute:             false,
		Fullscreen:       false,
		VSync:            true,
		ShowDebug:        false,
		AutoSaveInterval: 5,
		CreativeMode:     false,
		ChunkRenderDist:  8,
		Keybinds: map[string]string{
			"move_left":  "A",
			"move_right": "D",
			"jump":       "Space",
			"crouch":     "S",
			"inventory":  "E",
			"crafting":   "C",
			"menu":       "Escape",
			"drop":       "Q",
		},
		WindowWidth:  1280,
		WindowHeight: 720,
	}
}

// Manager handles loading and saving settings
type Manager struct {
	settings *Settings
	path     string
	mu       sync.RWMutex
}

// NewManager creates a settings manager with the given config directory
func NewManager(configDir string) *Manager {
	if configDir == "" {
		home, _ := os.UserHomeDir()
		configDir = filepath.Join(home, ".tesselbox")
	}

	return &Manager{
		settings: DefaultSettings(),
		path:     filepath.Join(configDir, "settings.yaml"),
	}
}

// Load reads settings from disk or creates defaults
func (m *Manager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Ensure directory exists
	dir := filepath.Dir(m.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	// Try to load existing settings
	data, err := os.ReadFile(m.path)
	if err != nil {
		if os.IsNotExist(err) {
			// No settings file yet - save defaults
			return m.saveLocked()
		}
		return fmt.Errorf("failed to read settings: %w", err)
	}

	// Parse settings
	if err := yaml.Unmarshal(data, m.settings); err != nil {
		return fmt.Errorf("failed to parse settings: %w", err)
	}

	return nil
}

// Save writes settings to disk
func (m *Manager) Save() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.saveLocked()
}

func (m *Manager) saveLocked() error {
	data, err := yaml.Marshal(m.settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := os.WriteFile(m.path, data, 0644); err != nil {
		return fmt.Errorf("failed to write settings: %w", err)
	}

	return nil
}

// Get returns the current settings (copy)
func (m *Manager) Get() *Settings {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy to prevent external modification
	// Create a copy without the mutex
	return &Settings{
		Volume:           m.settings.Volume,
		Mute:             m.settings.Mute,
		Fullscreen:       m.settings.Fullscreen,
		VSync:            m.settings.VSync,
		ShowDebug:        m.settings.ShowDebug,
		AutoSaveInterval: m.settings.AutoSaveInterval,
		CreativeMode:     m.settings.CreativeMode,
		ChunkRenderDist:  m.settings.ChunkRenderDist,
		WindowWidth:      m.settings.WindowWidth,
		WindowHeight:     m.settings.WindowHeight,
		Keybinds:         m.settings.Keybinds,
	}
}

// GetVolume returns the volume level
func (m *Manager) GetVolume() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.settings.Mute {
		return 0
	}
	return m.settings.Volume
}

// SetVolume updates the volume
func (m *Manager) SetVolume(v float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.settings.Volume = clamp(v, 0, 1)
	return m.saveLocked()
}

// SetMute toggles mute
func (m *Manager) SetMute(muted bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.settings.Mute = muted
	return m.saveLocked()
}

// GetKeybind returns a keybind or the default
func (m *Manager) GetKeybind(action string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if key, ok := m.settings.Keybinds[action]; ok {
		return key
	}
	return DefaultSettings().Keybinds[action]
}

// SetKeybind updates a keybind
func (m *Manager) SetKeybind(action, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.settings.Keybinds[action] = key
	return m.saveLocked()
}

// GetBool returns a boolean setting
func (m *Manager) GetBool(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	switch name {
	case "fullscreen":
		return m.settings.Fullscreen
	case "vsync":
		return m.settings.VSync
	case "show_debug":
		return m.settings.ShowDebug
	case "creative_mode_default":
		return m.settings.CreativeMode
	default:
		return false
	}
}

// SetBool updates a boolean setting
func (m *Manager) SetBool(name string, value bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	switch name {
	case "fullscreen":
		m.settings.Fullscreen = value
	case "vsync":
		m.settings.VSync = value
	case "show_debug":
		m.settings.ShowDebug = value
	case "creative_mode_default":
		m.settings.CreativeMode = value
	default:
		return fmt.Errorf("unknown setting: %s", name)
	}
	return m.saveLocked()
}

// GetInt returns an integer setting
func (m *Manager) GetInt(name string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	switch name {
	case "autosave_interval":
		return m.settings.AutoSaveInterval
	case "chunk_render_distance":
		return m.settings.ChunkRenderDist
	case "window_width":
		return m.settings.WindowWidth
	case "window_height":
		return m.settings.WindowHeight
	default:
		return 0
	}
}

// SetInt updates an integer setting
func (m *Manager) SetInt(name string, value int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	switch name {
	case "autosave_interval":
		m.settings.AutoSaveInterval = value
	case "chunk_render_distance":
		m.settings.ChunkRenderDist = value
	case "window_width":
		m.settings.WindowWidth = value
	case "window_height":
		m.settings.WindowHeight = value
	default:
		return fmt.Errorf("unknown setting: %s", name)
	}
	return m.saveLocked()
}

// ResetToDefaults restores all settings to defaults
func (m *Manager) ResetToDefaults() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.settings = DefaultSettings()
	return m.saveLocked()
}

// GetSettingsPath returns the path to the settings file
func (m *Manager) GetSettingsPath() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.path
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
