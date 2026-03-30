package audio

import (
	"log"
	"sync"

	"tesselbox/pkg/settings"
)

// Manager handles all game audio (stub implementation for v0.3-alpha)
// Full audio with beep library coming in v0.4
type Manager struct {
	enabled    bool
	volume     float64
	settings   *settings.Manager
	mu         sync.RWMutex
}

// NewManager creates an audio manager
func NewManager(settingsMgr *settings.Manager) *Manager {
	return &Manager{
		enabled:  true,
		volume:   0.7,
		settings: settingsMgr,
	}
}

// Initialize sets up the audio system
func (m *Manager) Initialize() error {
	// Load settings
	if m.settings != nil {
		m.mu.Lock()
		m.volume = m.settings.GetVolume()
		m.enabled = !m.settings.GetBool("mute")
		m.mu.Unlock()
	}

	log.Println("Audio system initialized (stub - no sound output)")
	return nil
}

// PlayMineSound plays a mining sound effect
func (m *Manager) PlayMineSound() {
	if !m.isEnabled() {
		return
	}
	log.Println("[AUDIO] Mine sound")
}

// PlayPlaceSound plays a block placement sound
func (m *Manager) PlayPlaceSound() {
	if !m.isEnabled() {
		return
	}
	log.Println("[AUDIO] Place sound")
}

// PlayJumpSound plays a jump sound effect
func (m *Manager) PlayJumpSound() {
	if !m.isEnabled() {
		return
	}
	log.Println("[AUDIO] Jump sound")
}

// PlayStepSound plays a footstep sound
func (m *Manager) PlayStepSound() {
	if !m.isEnabled() {
		return
	}
	// Too frequent to log
}

// PlayInventorySound plays inventory interaction sound
func (m *Manager) PlayInventorySound() {
	if !m.isEnabled() {
		return
	}
	log.Println("[AUDIO] Inventory sound")
}

// PlayMenuSound plays menu navigation sound
func (m *Manager) PlayMenuSound() {
	if !m.isEnabled() {
		return
	}
	// Too frequent to log
}

// SetVolume updates the audio volume
func (m *Manager) SetVolume(v float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.volume = v
}

// SetEnabled toggles audio on/off
func (m *Manager) SetEnabled(enabled bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.enabled = enabled
}

// isEnabled returns whether audio is enabled
func (m *Manager) isEnabled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.enabled
}

// getVolume returns the current volume level
func (m *Manager) getVolume() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.volume
}

// UpdateSettings reloads audio settings from the settings manager
func (m *Manager) UpdateSettings() {
	if m.settings == nil {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.volume = m.settings.GetVolume()
	m.enabled = !m.settings.GetBool("mute")
}

// Close shuts down the audio system
func (m *Manager) Close() error {
	return nil
}
