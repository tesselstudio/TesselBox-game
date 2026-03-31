package audio

import (
	"bytes"
	"io"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const (
	// Audio sample rate
	SampleRate = 44100
	
	// Audio buffer size
	BufferSize = 1024
)

// AudioType represents different types of audio in the game
type AudioType int

const (
	AudioTypeMusic AudioType = iota
	AudioTypeSFX
	AudioTypeAmbient
)

// Sound represents a loaded audio file
type Sound struct {
	Name     string
	Data     []byte
	Type     AudioType
	Volume   float64
	Loop     bool
	LoadedAt time.Time
}

// AudioPlayer represents an audio player instance
type AudioPlayer struct {
	Player     *audio.Player
	Sound      *Sound
	Volume     float64
	IsPlaying  bool
	IsPaused   bool
	StartedAt  time.Time
	Duration   time.Duration
}

// AudioManager manages all audio in the game
type AudioManager struct {
	// Audio context
	audioContext *audio.Context
	
	// Sound storage
	sounds map[string]*Sound
	mu     sync.RWMutex
	
	// Active players
	players map[string]*AudioPlayer
	playerMu sync.RWMutex
	
	// Volume controls
	masterVolume float64
	musicVolume  float64
	sfxVolume    float64
	ambientVolume float64
	
	// Settings
	muted     bool
	audioEnabled bool
	
	// Player pools for performance
	musicPlayers  []*AudioPlayer
	sfxPlayers    []*AudioPlayer
	ambientPlayers []*AudioPlayer
}

// NewAudioManager creates a new audio manager
func NewAudioManager() *AudioManager {
	ctx := audio.NewContext(SampleRate)
	
	return &AudioManager{
		audioContext:   ctx,
		sounds:         make(map[string]*Sound),
		players:        make(map[string]*AudioPlayer),
		masterVolume:   1.0,
		musicVolume:    0.7,
		sfxVolume:      0.8,
		ambientVolume:  0.6,
		audioEnabled:   true,
		muted:          false,
		musicPlayers:   make([]*AudioPlayer, 0, 5),   // Max 5 concurrent music tracks
		sfxPlayers:     make([]*AudioPlayer, 0, 20),  // Max 20 concurrent SFX
		ambientPlayers: make([]*AudioPlayer, 0, 10),  // Max 10 concurrent ambient sounds
	}
}

// LoadSound loads an audio file from memory
func (am *AudioManager) LoadSound(name string, data []byte, audioType AudioType, volume float64, loop bool) error {
	am.mu.Lock()
	defer am.mu.Unlock()
	
	sound := &Sound{
		Name:     name,
		Data:     data,
		Type:     audioType,
		Volume:   volume,
		Loop:     loop,
		LoadedAt: time.Now(),
	}
	
	am.sounds[name] = sound
	log.Printf("Loaded audio: %s (type: %v, volume: %.2f)", name, audioType, volume)
	
	return nil
}

// PlaySound plays a sound by name
func (am *AudioManager) PlaySound(name string) error {
	if !am.audioEnabled || am.muted {
		return nil
	}
	
	am.mu.RLock()
	sound, exists := am.sounds[name]
	am.mu.RUnlock()
	
	if !exists {
		log.Printf("Sound not found: %s", name)
		return nil
	}
	
	// Calculate final volume based on type and settings
	finalVolume := am.calculateVolume(sound)
	
	// Create audio player from sound data
	stream, err := wav.DecodeWithSampleRate(SampleRate, bytes.NewReader(sound.Data))
	if err != nil {
		return err
	}
	
	// Create infinite loop if needed
	var audioStream io.Reader = stream
	if sound.Loop {
		audioStream = audio.NewInfiniteLoop(stream, stream.Length())
	}
	
	player, err := am.audioContext.NewPlayer(audioStream)
	if err != nil {
		return err
	}
	player.SetVolume(finalVolume)
	
	// Create player record
	audioPlayer := &AudioPlayer{
		Player:    player,
		Sound:     sound,
		Volume:    finalVolume,
		IsPlaying: true,
		StartedAt: time.Now(),
	}
	
	// Add to appropriate player pool
	am.addPlayerToPool(audioPlayer)
	
	// Store in active players
	am.playerMu.Lock()
	playerID := am.generatePlayerID()
	am.players[playerID] = audioPlayer
	am.playerMu.Unlock()
	
	// Start playing
	player.Play()
	
	return nil
}

// StopSound stops a playing sound by name
func (am *AudioManager) StopSound(name string) {
	am.playerMu.Lock()
	defer am.playerMu.Unlock()
	
	for id, player := range am.players {
		if player.Sound.Name == name {
			player.Player.Close()
			delete(am.players, id)
			am.removeFromPlayerPool(player)
			break
		}
	}
}

// StopAllSounds stops all currently playing sounds
func (am *AudioManager) StopAllSounds() {
	am.playerMu.Lock()
	defer am.playerMu.Unlock()
	
	for id, player := range am.players {
		player.Player.Close()
		delete(am.players, id)
	}
	
	// Clear all player pools
	am.musicPlayers = am.musicPlayers[:0]
	am.sfxPlayers = am.sfxPlayers[:0]
	am.ambientPlayers = am.ambientPlayers[:0]
}

// SetMasterVolume sets the master volume (0.0 to 1.0)
func (am *AudioManager) SetMasterVolume(volume float64) {
	if volume < 0.0 {
		volume = 0.0
	} else if volume > 1.0 {
		volume = 1.0
	}
	
	am.masterVolume = volume
	am.updateAllPlayerVolumes()
}

// SetMusicVolume sets the music volume (0.0 to 1.0)
func (am *AudioManager) SetMusicVolume(volume float64) {
	if volume < 0.0 {
		volume = 0.0
	} else if volume > 1.0 {
		volume = 1.0
	}
	
	am.musicVolume = volume
	am.updatePlayerVolumesByType(AudioTypeMusic)
}

// SetSFXVolume sets the sound effects volume (0.0 to 1.0)
func (am *AudioManager) SetSFXVolume(volume float64) {
	if volume < 0.0 {
		volume = 0.0
	} else if volume > 1.0 {
		volume = 1.0
	}
	
	am.sfxVolume = volume
	am.updatePlayerVolumesByType(AudioTypeSFX)
}

// SetAmbientVolume sets the ambient volume (0.0 to 1.0)
func (am *AudioManager) SetAmbientVolume(volume float64) {
	if volume < 0.0 {
		volume = 0.0
	} else if volume > 1.0 {
		volume = 1.0
	}
	
	am.ambientVolume = volume
	am.updatePlayerVolumesByType(AudioTypeAmbient)
}

// SetMuted mutes or unmutes all audio
func (am *AudioManager) SetMuted(muted bool) {
	am.muted = muted
	
	if muted {
		am.StopAllSounds()
	}
}

// SetEnabled enables or disables audio system
func (am *AudioManager) SetEnabled(enabled bool) {
	am.audioEnabled = enabled
	
	if !enabled {
		am.StopAllSounds()
	}
}

// Update should be called regularly to clean up finished players
func (am *AudioManager) Update() {
	am.playerMu.Lock()
	defer am.playerMu.Unlock()
	
	for id, player := range am.players {
		if player.Player.IsPlaying() {
			continue
		}
		
		// Player finished, clean it up
		player.Player.Close()
		delete(am.players, id)
		am.removeFromPlayerPool(player)
	}
}

// GetLoadedSounds returns a list of all loaded sounds
func (am *AudioManager) GetLoadedSounds() []string {
	am.mu.RLock()
	defer am.mu.RUnlock()
	
	sounds := make([]string, 0, len(am.sounds))
	for name := range am.sounds {
		sounds = append(sounds, name)
	}
	
	return sounds
}

// HasSound checks if a sound is already loaded
func (am *AudioManager) HasSound(name string) bool {
	am.mu.RLock()
	defer am.mu.RUnlock()
	
	_, exists := am.sounds[name]
	return exists
}

// GetPlayingSounds returns a list of currently playing sounds
func (am *AudioManager) GetPlayingSounds() []string {
	am.playerMu.RLock()
	defer am.playerMu.RUnlock()
	
	sounds := make([]string, 0, len(am.players))
	for _, player := range am.players {
		if player.IsPlaying {
			sounds = append(sounds, player.Sound.Name)
		}
	}
	
	return sounds
}

// Helper methods

func (am *AudioManager) calculateVolume(sound *Sound) float64 {
	if am.muted || !am.audioEnabled {
		return 0.0
	}
	
	typeVolume := 1.0
	switch sound.Type {
	case AudioTypeMusic:
		typeVolume = am.musicVolume
	case AudioTypeSFX:
		typeVolume = am.sfxVolume
	case AudioTypeAmbient:
		typeVolume = am.ambientVolume
	}
	
	return am.masterVolume * typeVolume * sound.Volume
}

func (am *AudioManager) updateAllPlayerVolumes() {
	am.playerMu.RLock()
	defer am.playerMu.RUnlock()
	
	for _, player := range am.players {
		newVolume := am.calculateVolume(player.Sound)
		player.Player.SetVolume(newVolume)
		player.Volume = newVolume
	}
}

func (am *AudioManager) updatePlayerVolumesByType(audioType AudioType) {
	am.playerMu.RLock()
	defer am.playerMu.RUnlock()
	
	for _, player := range am.players {
		if player.Sound.Type == audioType {
			newVolume := am.calculateVolume(player.Sound)
			player.Player.SetVolume(newVolume)
			player.Volume = newVolume
		}
	}
}

func (am *AudioManager) addPlayerToPool(player *AudioPlayer) {
	switch player.Sound.Type {
	case AudioTypeMusic:
		am.musicPlayers = append(am.musicPlayers, player)
	case AudioTypeSFX:
		am.sfxPlayers = append(am.sfxPlayers, player)
	case AudioTypeAmbient:
		am.ambientPlayers = append(am.ambientPlayers, player)
	}
}

func (am *AudioManager) removeFromPlayerPool(player *AudioPlayer) {
	switch player.Sound.Type {
	case AudioTypeMusic:
		for i, p := range am.musicPlayers {
			if p == player {
				am.musicPlayers = append(am.musicPlayers[:i], am.musicPlayers[i+1:]...)
				break
			}
		}
	case AudioTypeSFX:
		for i, p := range am.sfxPlayers {
			if p == player {
				am.sfxPlayers = append(am.sfxPlayers[:i], am.sfxPlayers[i+1:]...)
				break
			}
		}
	case AudioTypeAmbient:
		for i, p := range am.ambientPlayers {
			if p == player {
				am.ambientPlayers = append(am.ambientPlayers[:i], am.ambientPlayers[i+1:]...)
				break
			}
		}
	}
}

func (am *AudioManager) generatePlayerID() string {
	return "player_" + time.Now().Format("20060102150405") + "_" + 
		   string(rune(rand.Intn(26)+'A')) + string(rune(rand.Intn(26)+'a'))
}

// PlayRandomSFX plays a random sound effect from a list
func (am *AudioManager) PlayRandomSFX(sounds []string) {
	if len(sounds) == 0 {
		return
	}
	
	randomIndex := rand.Intn(len(sounds))
	am.PlaySound(sounds[randomIndex])
}

// PlayMusicWithFade plays music with a fade-in effect
func (am *AudioManager) PlayMusicWithFade(name string, fadeDuration time.Duration) {
	// For now, just play normally. Fade effects can be implemented later
	am.PlaySound(name)
}

// CrossfadeMusic crossfades from one music track to another
func (am *AudioManager) CrossfadeMusic(fromName, toName string, duration time.Duration) {
	// Stop current music and start new one
	if fromName != "" {
		am.StopSound(fromName)
	}
	if toName != "" {
		am.PlaySound(toName)
	}
}

// Cleanup should be called when shutting down the game
func (am *AudioManager) Cleanup() {
	am.StopAllSounds()
	am.audioContext = nil
}
