package audio

import (
	"log"
	"sync"
	"time"
)

// BackgroundMusicManager handles looping background music for the game
type BackgroundMusicManager struct {
	audioManager *AudioManager
	currentTrack string
	isPlaying    bool
	mu           sync.RWMutex

	// Fade settings
	fadeDuration time.Duration
	fadeTimer    time.Time
	fadeStartVol float64
	fadeEndVol   float64
	isFading     bool
}

// NewBackgroundMusicManager creates a new background music manager
func NewBackgroundMusicManager(audioManager *AudioManager) *BackgroundMusicManager {
	return &BackgroundMusicManager{
		audioManager: audioManager,
		fadeDuration: 2 * time.Second,
	}
}

// PlayBackgroundMusic starts playing background music
// It will loop continuously until stopped
func (bmm *BackgroundMusicManager) PlayBackgroundMusic(trackName string) error {
	bmm.mu.Lock()
	defer bmm.mu.Unlock()

	// Stop current track if playing
	if bmm.isPlaying && bmm.currentTrack != "" {
		bmm.audioManager.StopSound(bmm.currentTrack)
	}

	// Play the new track
	if err := bmm.audioManager.PlaySound(trackName); err != nil {
		log.Printf("Failed to play background music '%s': %v", trackName, err)
		return err
	}

	bmm.currentTrack = trackName
	bmm.isPlaying = true
	log.Printf("Background music started: %s", trackName)

	return nil
}

// StopBackgroundMusic stops the current background music
func (bmm *BackgroundMusicManager) StopBackgroundMusic() {
	bmm.mu.Lock()
	defer bmm.mu.Unlock()

	if bmm.isPlaying && bmm.currentTrack != "" {
		bmm.audioManager.StopSound(bmm.currentTrack)
		bmm.isPlaying = false
		log.Printf("Background music stopped")
	}
}

// CrossfadeBackgroundMusic crossfades from current track to a new one
func (bmm *BackgroundMusicManager) CrossfadeBackgroundMusic(newTrackName string) error {
	bmm.mu.Lock()
	defer bmm.mu.Unlock()

	oldTrack := bmm.currentTrack

	// Start fade out of old track
	if bmm.isPlaying && oldTrack != "" {
		// Fade out old track
		go func() {
			time.Sleep(bmm.fadeDuration)
			bmm.audioManager.StopSound(oldTrack)
		}()
	}

	// Start fade in of new track
	if err := bmm.audioManager.PlaySound(newTrackName); err != nil {
		log.Printf("Failed to crossfade to '%s': %v", newTrackName, err)
		return err
	}

	bmm.currentTrack = newTrackName
	bmm.isPlaying = true
	log.Printf("Background music crossfaded to: %s", newTrackName)

	return nil
}

// IsPlaying returns whether background music is currently playing
func (bmm *BackgroundMusicManager) IsPlaying() bool {
	bmm.mu.RLock()
	defer bmm.mu.RUnlock()
	return bmm.isPlaying
}

// GetCurrentTrack returns the name of the currently playing track
func (bmm *BackgroundMusicManager) GetCurrentTrack() string {
	bmm.mu.RLock()
	defer bmm.mu.RUnlock()
	return bmm.currentTrack
}

// SetVolume sets the background music volume (0.0 to 1.0)
func (bmm *BackgroundMusicManager) SetVolume(volume float64) {
	bmm.audioManager.SetMusicVolume(volume)
}

// Update should be called regularly to handle fade effects
func (bmm *BackgroundMusicManager) Update() {
	// Fade effects can be implemented here if needed
}
