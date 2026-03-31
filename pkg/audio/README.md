# Audio System for TesselBox

This package provides a comprehensive audio system for the TesselBox game, including sound effects, background music, and ambient sounds.

## Features

- **Sound Effects**: Footsteps, block interactions, item pickups, UI interactions
- **Background Music**: Context-aware music that changes based on game state
- **Ambient Sounds**: Weather effects, environmental sounds
- **Volume Controls**: Separate controls for master, music, SFX, and ambient volumes
- **Audio Loading**: Embedded assets with fallback to placeholder sounds
- **Performance Optimized**: Object pooling and efficient audio management

## Components

### AudioManager
Core audio management system that handles:
- Loading and playing audio files
- Volume control and muting
- Player lifecycle management
- Audio context management

### SoundLibrary
Game-specific sound management that provides:
- Context-aware sound selection
- Surface-specific footstep sounds
- Game state-based music selection
- Environmental audio adaptation

### AudioLoader
Handles loading audio files from embedded assets:
- WAV file support (with placeholders for missing files)
- Automatic sound categorization
- Fallback placeholder generation

## Usage

### Basic Setup
```go
// Create audio manager
audioManager := audio.NewAudioManager()

// Create sound library
soundLibrary := audio.NewSoundLibrary(audioManager)

// Load audio files
loader := audio.NewAudioLoader(audioManager)
loader.LoadAllAudio()

// Initialize sound library
soundLibrary.InitializeDefaultSounds()
```

### Playing Sounds
```go
// Play a specific sound effect
soundLibrary.PlayFootstep("grass")
soundLibrary.PlayBlockSound("break", "stone")
soundLibrary.PlayItemSound("pickup")
soundLibrary.PlayUISound("click")

// Play background music
soundLibrary.PlayMusic("gameplay")

// Play ambient sounds
soundLibrary.PlayAmbientSound()
```

### Volume Control
```go
// Set master volume (0.0 to 1.0)
audioManager.SetMasterVolume(0.8)

// Set specific volume types
audioManager.SetMusicVolume(0.7)
audioManager.SetSFXVolume(0.9)
audioManager.SetAmbientVolume(0.6)

// Mute/unmute
audioManager.SetMuted(true)
```

## Audio File Structure

Audio files should be placed in the following structure:

```
pkg/audio/assets/
├── sfx/           # Sound effects
│   ├── footstep_grass.wav
│   ├── block_break.wav
│   ├── item_pickup.wav
│   └── ui_click.wav
├── music/         # Background music
│   ├── menu_music.wav
│   ├── gameplay_music.wav
│   └── creative_music.wav
└── ambient/       # Ambient sounds
    ├── wind.wav
    ├── rain.wav
    └── birds.wav
```

## Sound Effects

### Movement
- Footstep sounds for different surfaces (grass, stone, sand, water)
- Jump and landing sounds

### Block Interaction
- Block breaking sounds
- Block placement sounds
- Mining progress sounds

### Items
- Item pickup sounds
- Item drop sounds
- Inventory sounds
- Hotbar selection sounds

### UI
- Click sounds
- Hover sounds
- Menu navigation sounds

### Crafting
- Crafting start/complete sounds
- Smelting sounds

## Music Tracks

### Context-Aware Music
- Menu music
- Gameplay music (changes by biome/time)
- Creative mode music
- Underground music
- Night music
- Combat music
- Boss music

## Ambient Sounds

### Environmental
- Wind (varying by biome)
- Rain and thunder
- Bird sounds (forest biomes)
- Water flow sounds
- Fire sounds

### Weather
- Dynamic weather audio
- Storm effects
- Snow sounds

## Integration with Game

The audio system is integrated into the main game loop:

1. **Initialization**: Audio system is set up in `NewGame()`
2. **Updates**: Audio manager is updated each frame to clean up finished sounds
3. **Game Events**: Various game actions trigger appropriate sounds
4. **Context Updates**: Audio context changes based on player location and game state

## Performance Considerations

- **Object Pooling**: Reuses audio players to reduce garbage collection
- **Concurrent Playback**: Supports multiple simultaneous sounds
- **Memory Management**: Automatically cleans up finished audio players
- **Embedded Assets**: Audio files are embedded in the binary for distribution

## Placeholder Audio

If audio files are missing, the system generates placeholder sounds:
- Simple sine wave tones for testing
- Different frequencies for different sound types
- Ensures the game works without external audio files

## Future Enhancements

- **Audio Streaming**: For longer music tracks
- **3D Audio**: Positional audio for immersive gameplay
- **Dynamic Audio**: Music that responds to gameplay intensity
- **Audio Settings**: In-game audio configuration menu
- **Mod Support**: Allow mods to add custom audio files
