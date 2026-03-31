package audio

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

// Audio files will be embedded when they exist
//go:embed assets/sfx/*.wav
//go:embed assets/music/*.wav
//go:embed assets/ambient/*.wav
var audioFS embed.FS

// AudioLoader handles loading audio files from embedded assets
type AudioLoader struct {
	manager *AudioManager
}

// NewAudioLoader creates a new audio loader
func NewAudioLoader(manager *AudioManager) *AudioLoader {
	return &AudioLoader{
		manager: manager,
	}
}

// LoadAllAudio loads all audio files from embedded assets with robust fallback
func (al *AudioLoader) LoadAllAudio() error {
	log.Printf("Loading audio files from embedded assets...")
	
	loadedCount := 0
	errorCount := 0
	
	// Try to load from embedded directories
	directories := []struct {
		path      string
		audioType AudioType
		volume    float64
		loop      bool
	}{
		{"assets/sfx", AudioTypeSFX, 0.7, false},
		{"assets/music", AudioTypeMusic, 0.4, true},
		{"assets/ambient", AudioTypeAmbient, 0.3, true},
	}
	
	for _, dir := range directories {
		if err := al.loadAudioFromDir(dir.path, dir.audioType, dir.volume, dir.loop); err != nil {
			log.Printf("Warning: Failed to load audio from %s: %v", dir.path, err)
			errorCount++
		} else {
			loadedCount++
		}
	}
	
	// Always load placeholder sounds to ensure basic audio functionality
	al.LoadPlaceholderSounds()
	
	if loadedCount > 0 {
		log.Printf("Audio loading completed: %d directories loaded, %d directories used placeholders", loadedCount, errorCount)
	} else {
		log.Printf("Audio loading completed with placeholder sounds only")
	}
	
	return nil
}

// loadAudioFromDir loads all audio files from a specific directory
func (al *AudioLoader) loadAudioFromDir(dir string, audioType AudioType, volume float64, loop bool) error {
	entries, err := fs.ReadDir(audioFS, dir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dir, err)
	}
	
	loadedCount := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		name := entry.Name()
		if !al.isAudioFile(name) {
			continue
		}
		
		// Construct full path
		fullPath := filepath.Join(dir, name)
		
		// Load audio data
		if err := al.loadSingleAudio(fullPath, audioType, volume, loop); err != nil {
			log.Printf("Failed to load audio %s: %v", fullPath, err)
			continue
		}
		
		loadedCount++
	}
	
	log.Printf("Loaded %d audio files from %s", loadedCount, dir)
	return nil
}

// loadSingleAudio loads a single audio file
func (al *AudioLoader) loadSingleAudio(path string, audioType AudioType, volume float64, loop bool) error {
	data, err := fs.ReadFile(audioFS, path)
	if err != nil {
		return fmt.Errorf("failed to read audio file %s: %w", path, err)
	}
	
	// Extract sound name from path (remove directory and extension)
	name := al.extractSoundName(path)
	
	return al.manager.LoadSound(name, data, audioType, volume, loop)
}

// isAudioFile checks if a file is an audio file based on extension
func (al *AudioLoader) isAudioFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".wav" || ext == ".mp3" || ext == ".ogg"
}

// extractSoundName extracts a clean sound name from file path
func (al *AudioLoader) extractSoundName(path string) string {
	// Get just the filename
	filename := filepath.Base(path)
	
	// Remove extension
	name := filename[:len(filename)-len(filepath.Ext(filename))]
	
	// Convert to lowercase and replace spaces/hyphens with underscores
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	
	return name
}

// LoadSpecificAudio loads a specific audio file if it exists
func (al *AudioLoader) LoadSpecificAudio(soundName string, audioType AudioType, volume float64, loop bool) error {
	// Try different possible paths for the sound
	possiblePaths := []string{
		fmt.Sprintf("assets/sfx/%s.wav", soundName),
		fmt.Sprintf("assets/music/%s.wav", soundName),
		fmt.Sprintf("assets/ambient/%s.wav", soundName),
		fmt.Sprintf("assets/%s.wav", soundName),
	}
	
	for _, path := range possiblePaths {
		data, err := fs.ReadFile(audioFS, path)
		if err == nil {
			// File found, load it
			return al.manager.LoadSound(soundName, data, audioType, volume, loop)
		}
	}
	
	return fmt.Errorf("audio file not found for sound: %s", soundName)
}

// ListAvailableAudio returns a list of all available audio files
func (al *AudioLoader) ListAvailableAudio() (sfx, music, ambient []string) {
	al.listAudioInDir("assets/sfx", &sfx)
	al.listAudioInDir("assets/music", &music)
	al.listAudioInDir("assets/ambient", &ambient)
	al.listAudioInDir("assets", &sfx) // Root level SFX
	
	return sfx, music, ambient
}

// listAudioInDir lists audio files in a specific directory
func (al *AudioLoader) listAudioInDir(dir string, audioList *[]string) {
	entries, err := fs.ReadDir(audioFS, dir)
	if err != nil {
		return
	}
	
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		name := entry.Name()
		if al.isAudioFile(name) {
			soundName := al.extractSoundName(filepath.Join(dir, name))
			*audioList = append(*audioList, soundName)
		}
	}
}

// CreatePlaceholderAudio creates simple placeholder audio data for testing
// This generates a simple sine wave tone for missing sounds
func (al *AudioLoader) CreatePlaceholderAudio(frequency float64, duration float64, sampleRate int) []byte {
	samples := int(duration * float64(sampleRate))
	data := make([]byte, samples*2) // 16-bit audio
	
	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)
		value := int16(32767 * 0.1 * sin(2*3.14159*frequency*t)) // 10% volume
		
		// Little-endian 16-bit PCM
		data[i*2] = byte(value & 0xFF)
		data[i*2+1] = byte(value >> 8)
	}
	
	return data
}

// sin is a simple sine function for placeholder audio generation
func sin(x float64) float64 {
	// Simple Taylor series approximation for sine
	// sin(x) ≈ x - x³/6 + x⁵/120 - x⁷/5040
	x = modFloat(x, 2*3.14159)
	
	x2 := x * x
	x3 := x2 * x
	x5 := x3 * x2
	x7 := x5 * x2
	
	return x - x3/6 + x5/120 - x7/5040
}

// modFloat is float modulo for placeholder audio generation
func modFloat(a, b float64) float64 {
	for a >= b {
		a -= b
	}
	for a < 0 {
		a += b
	}
	return a
}

// LoadPlaceholderSounds creates placeholder sounds for missing audio files with improved reliability
func (al *AudioLoader) LoadPlaceholderSounds() {
	log.Printf("Creating placeholder audio sounds...")
	
	// Create placeholder for common sound effects
	placeholderSounds := map[string]struct {
		frequency float64
		duration  float64
		audioType AudioType
		volume    float64
		loop      bool
	}{
		"ui_click":           {440, 0.1, AudioTypeSFX, 0.5, false},
		"ui_hover":           {880, 0.05, AudioTypeSFX, 0.3, false},
		"ui_open":            {660, 0.15, AudioTypeSFX, 0.6, false},
		"ui_close":           {330, 0.15, AudioTypeSFX, 0.6, false},
		"block_place":        {660, 0.15, AudioTypeSFX, 0.6, false},
		"block_break":        {220, 0.2, AudioTypeSFX, 0.7, false},
		"item_pickup":        {880, 0.1, AudioTypeSFX, 0.6, false},
		"item_drop":          {440, 0.1, AudioTypeSFX, 0.5, false},
		"hotbar_select":      {550, 0.08, AudioTypeSFX, 0.4, false},
		"footstep_grass":     {330, 0.1, AudioTypeSFX, 0.4, false},
		"footstep_stone":     {440, 0.1, AudioTypeSFX, 0.5, false},
		"footstep_dirt":      {380, 0.1, AudioTypeSFX, 0.45, false},
		"footstep_sand":      {290, 0.1, AudioTypeSFX, 0.4, false},
		"jump":               {550, 0.15, AudioTypeSFX, 0.6, false},
		"land":               {220, 0.1, AudioTypeSFX, 0.5, false},
		"craft_complete":     {880, 0.2, AudioTypeSFX, 0.6, false},
		"menu_music":         {220, 10.0, AudioTypeMusic, 0.3, true},
		"gameplay_music":     {110, 15.0, AudioTypeMusic, 0.4, true},
		"wind":               {80, 20.0, AudioTypeAmbient, 0.2, true},
		"rain":               {200, 25.0, AudioTypeAmbient, 0.3, true},
		"night":              {60, 30.0, AudioTypeAmbient, 0.15, true},
	}
	
	createdCount := 0
	skippedCount := 0
	
	for name, config := range placeholderSounds {
		// Check if sound already exists to avoid duplicates
		if al.manager.HasSound(name) {
			skippedCount++
			continue
		}
		
		data := al.CreatePlaceholderAudio(config.frequency, config.duration, SampleRate)
		if err := al.manager.LoadSound(name, data, config.audioType, config.volume, config.loop); err != nil {
			log.Printf("Failed to create placeholder sound %s: %v", name, err)
		} else {
			createdCount++
		}
	}
	
	log.Printf("Placeholder audio sounds created: %d new, %d skipped (already exist)", createdCount, skippedCount)
}
