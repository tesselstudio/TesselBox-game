package audio

import (
	"log"
)

// SoundEffect represents different game sound effects
type SoundEffect string

const (
	// Movement sounds
	SFXFootstepGrass SoundEffect = "footstep_grass"
	SFXFootstepStone SoundEffect = "footstep_stone"
	SFXFootstepSand  SoundEffect = "footstep_sand"
	SFXFootstepWater SoundEffect = "footstep_water"
	SFXJump          SoundEffect = "jump"
	SFXLand          SoundEffect = "land"

	// Block interaction sounds
	SFXBlockBreak     SoundEffect = "block_break"
	SFXBlockPlace     SoundEffect = "block_place"
	SFXBlockHit       SoundEffect = "block_hit"
	SFXMiningComplete SoundEffect = "mining_complete"

	// Item sounds
	SFXItemPickup     SoundEffect = "item_pickup"
	SFXItemDrop       SoundEffect = "item_drop"
	SFXInventoryOpen  SoundEffect = "inventory_open"
	SFXInventoryClose SoundEffect = "inventory_close"
	SFXHotbarSelect   SoundEffect = "hotbar_select"

	// Crafting sounds
	SFXCraftingStart    SoundEffect = "crafting_start"
	SFXCraftingComplete SoundEffect = "crafting_complete"
	SFXSmeltingStart    SoundEffect = "smelting_start"
	SFXSmeltingComplete SoundEffect = "smelting_complete"

	// UI sounds
	SFXUIClick      SoundEffect = "ui_click"
	SFXUIHover      SoundEffect = "ui_hover"
	SFXUIOpen       SoundEffect = "ui_open"
	SFXUIClose      SoundEffect = "ui_close"
	SFXMenuNavigate SoundEffect = "menu_navigate"
	SFXMenuSelect   SoundEffect = "menu_select"

	// Ambient sounds
	SFXWind      SoundEffect = "wind"
	SFXRain      SoundEffect = "rain"
	SFXThunder   SoundEffect = "thunder"
	SFXBirds     SoundEffect = "birds"
	SFXWaterFlow SoundEffect = "water_flow"
	SFXFire      SoundEffect = "fire"

	// Combat sounds (if implemented)
	SFXHit      SoundEffect = "hit"
	SFXMiss     SoundEffect = "miss"
	SFXCritical SoundEffect = "critical"

	// Weather sounds
	SFXWeatherClear SoundEffect = "weather_clear"
	SFXWeatherRain  SoundEffect = "weather_rain"
	SFXWeatherStorm SoundEffect = "weather_storm"
	SFXWeatherSnow  SoundEffect = "weather_snow"
)

// MusicTrack represents different background music tracks
type MusicTrack string

const (
	MusicMenu        MusicTrack = "menu_music"
	MusicGameplay    MusicTrack = "gameplay_music"
	MusicCreative    MusicTrack = "creative_music"
	MusicNight       MusicTrack = "night_music"
	MusicUnderground MusicTrack = "underground_music"
	MusicCombat      MusicTrack = "combat_music"
	MusicBoss        MusicTrack = "boss_music"
	MusicVictory     MusicTrack = "victory_music"
	MusicDanger      MusicTrack = "danger_music"
)

// SoundLibrary manages game-specific sounds and their variations
type SoundLibrary struct {
	manager *AudioManager

	// Sound variations for more variety
	footstepVariations map[SoundEffect][]string
	miningVariations   map[SoundEffect][]string

	// Current context
	currentBiome   string
	currentWeather string
	isUnderground  bool
	timeOfDay      string
}

// NewSoundLibrary creates a new sound library
func NewSoundLibrary(manager *AudioManager) *SoundLibrary {
	return &SoundLibrary{
		manager:            manager,
		footstepVariations: make(map[SoundEffect][]string),
		miningVariations:   make(map[SoundEffect][]string),
	}
}

// InitializeDefaultSounds sets up default sound mappings
func (sl *SoundLibrary) InitializeDefaultSounds() {
	// Footstep variations
	sl.footstepVariations[SFXFootstepGrass] = []string{
		"footstep_grass_1",
		"footstep_grass_2",
		"footstep_grass_3",
		"footstep_grass_4",
	}
	sl.footstepVariations[SFXFootstepStone] = []string{
		"footstep_stone_1",
		"footstep_stone_2",
		"footstep_stone_3",
	}
	sl.footstepVariations[SFXFootstepSand] = []string{
		"footstep_sand_1",
		"footstep_sand_2",
	}
	sl.footstepVariations[SFXFootstepWater] = []string{
		"footstep_water_1",
		"footstep_water_2",
	}

	// Mining variations
	sl.miningVariations[SFXBlockHit] = []string{
		"block_hit_1",
		"block_hit_2",
		"block_hit_3",
		"block_hit_4",
	}
}

// PlayFootstep plays a footstep sound based on the surface type
func (sl *SoundLibrary) PlayFootstep(surfaceType string) {
	var soundEffect SoundEffect

	switch surfaceType {
	case "grass", "dirt", "farmland":
		soundEffect = SFXFootstepGrass
	case "stone", "cobblestone", "bedrock":
		soundEffect = SFXFootstepStone
	case "sand", "gravel":
		soundEffect = SFXFootstepSand
	case "water", "ice":
		soundEffect = SFXFootstepWater
	default:
		soundEffect = SFXFootstepStone // Default to stone
	}

	variations, exists := sl.footstepVariations[soundEffect]
	if exists && len(variations) > 0 {
		sl.manager.PlayRandomSFX(variations)
	} else {
		sl.manager.PlaySound(string(soundEffect))
	}
}

// PlayMiningSound plays a mining-related sound
func (sl *SoundLibrary) PlayMiningSound(miningProgress float64) {
	if miningProgress >= 100.0 {
		sl.manager.PlaySound(string(SFXMiningComplete))
	} else {
		variations, exists := sl.miningVariations[SFXBlockHit]
		if exists && len(variations) > 0 {
			sl.manager.PlayRandomSFX(variations)
		} else {
			sl.manager.PlaySound(string(SFXBlockHit))
		}
	}
}

// PlayBlockSound plays a sound for block interaction
func (sl *SoundLibrary) PlayBlockSound(action string, blockType string) {
	switch action {
	case "break":
		sl.manager.PlaySound(string(SFXBlockBreak))
	case "place":
		sl.manager.PlaySound(string(SFXBlockPlace))
	default:
		sl.manager.PlaySound(string(SFXBlockHit))
	}
}

// PlayItemSound plays a sound for item interaction
func (sl *SoundLibrary) PlayItemSound(action string) {
	switch action {
	case "pickup":
		sl.manager.PlaySound(string(SFXItemPickup))
	case "drop":
		sl.manager.PlaySound(string(SFXItemDrop))
	case "inventory_open":
		sl.manager.PlaySound(string(SFXInventoryOpen))
	case "inventory_close":
		sl.manager.PlaySound(string(SFXInventoryClose))
	case "hotbar_select":
		sl.manager.PlaySound(string(SFXHotbarSelect))
	default:
		sl.manager.PlaySound(string(SFXUIClick))
	}
}

// PlayUISound plays a UI interaction sound
func (sl *SoundLibrary) PlayUISound(action string) {
	switch action {
	case "click":
		sl.manager.PlaySound(string(SFXUIClick))
	case "hover":
		sl.manager.PlaySound(string(SFXUIHover))
	case "open":
		sl.manager.PlaySound(string(SFXUIOpen))
	case "close":
		sl.manager.PlaySound(string(SFXUIClose))
	case "navigate":
		sl.manager.PlaySound(string(SFXMenuNavigate))
	case "select":
		sl.manager.PlaySound(string(SFXMenuSelect))
	default:
		sl.manager.PlaySound(string(SFXUIClick))
	}
}

// PlayCraftingSound plays crafting-related sounds
func (sl *SoundLibrary) PlayCraftingSound(action string) {
	switch action {
	case "start":
		sl.manager.PlaySound(string(SFXCraftingStart))
	case "complete":
		sl.manager.PlaySound(string(SFXCraftingComplete))
	case "smelting_start":
		sl.manager.PlaySound(string(SFXSmeltingStart))
	case "smelting_complete":
		sl.manager.PlaySound(string(SFXSmeltingComplete))
	default:
		sl.manager.PlaySound(string(SFXUIClick))
	}
}

// PlayMusic plays background music based on context
func (sl *SoundLibrary) PlayMusic(context string) {
	var track MusicTrack

	switch context {
	case "menu":
		track = MusicMenu
	case "gameplay":
		if sl.isUnderground {
			track = MusicUnderground
		} else if sl.timeOfDay == "night" {
			track = MusicNight
		} else {
			track = MusicGameplay
		}
	case "creative":
		track = MusicCreative
	case "combat":
		track = MusicCombat
	case "boss":
		track = MusicBoss
	case "victory":
		track = MusicVictory
	case "danger":
		track = MusicDanger
	default:
		track = MusicGameplay
	}

	sl.manager.PlaySound(string(track))
}

// PlayAmbientSound plays ambient sounds based on environment
func (sl *SoundLibrary) PlayAmbientSound() {
	// Stop current ambient sounds first
	sl.manager.StopSound(string(SFXWind))
	sl.manager.StopSound(string(SFXRain))
	sl.manager.StopSound(string(SFXBirds))
	sl.manager.StopSound(string(SFXWaterFlow))

	// Play appropriate ambient sounds based on context
	if sl.currentWeather == "rain" || sl.currentWeather == "storm" {
		sl.manager.PlaySound(string(SFXRain))
		if sl.currentWeather == "storm" {
			// Thunder would be played randomly
		}
	} else if sl.currentBiome == "forest" && sl.timeOfDay == "day" {
		sl.manager.PlaySound(string(SFXBirds))
	}

	// Always play wind in certain conditions
	if sl.currentBiome == "mountains" || sl.currentBiome == "tundra" {
		sl.manager.PlaySound(string(SFXWind))
	}
}

// UpdateContext updates the current game context for audio
func (sl *SoundLibrary) UpdateContext(biome, weather string, isUnderground bool, timeOfDay string) {
	sl.currentBiome = biome
	sl.currentWeather = weather
	sl.isUnderground = isUnderground
	sl.timeOfDay = timeOfDay

	// Update ambient sounds based on new context
	sl.PlayAmbientSound()
}

// PlayWeatherSound plays weather-specific sounds
func (sl *SoundLibrary) PlayWeatherSound(weatherEvent string) {
	switch weatherEvent {
	case "thunder":
		sl.manager.PlaySound(string(SFXThunder))
	case "rain_start":
		sl.manager.PlaySound(string(SFXWeatherRain))
	case "rain_stop":
		sl.manager.StopSound(string(SFXWeatherRain))
	case "storm_start":
		sl.manager.PlaySound(string(SFXWeatherStorm))
	case "storm_stop":
		sl.manager.StopSound(string(SFXWeatherStorm))
	case "snow_start":
		sl.manager.PlaySound(string(SFXWeatherSnow))
	case "snow_stop":
		sl.manager.StopSound(string(SFXWeatherSnow))
	}
}

// GetSoundEffectNames returns all available sound effect names
func (sl *SoundLibrary) GetSoundEffectNames() []string {
	effects := []SoundEffect{
		SFXFootstepGrass, SFXFootstepStone, SFXFootstepSand, SFXFootstepWater,
		SFXJump, SFXLand, SFXBlockBreak, SFXBlockPlace, SFXBlockHit,
		SFXMiningComplete, SFXItemPickup, SFXItemDrop, SFXInventoryOpen,
		SFXInventoryClose, SFXHotbarSelect, SFXCraftingStart, SFXCraftingComplete,
		SFXSmeltingStart, SFXSmeltingComplete, SFXUIClick, SFXUIHover,
		SFXUIOpen, SFXUIClose, SFXMenuNavigate, SFXMenuSelect,
		SFXWind, SFXRain, SFXThunder, SFXBirds, SFXWaterFlow,
		SFXFire, SFXHit, SFXMiss, SFXCritical, SFXWeatherClear,
		SFXWeatherRain, SFXWeatherStorm, SFXWeatherSnow,
	}

	names := make([]string, len(effects))
	for i, effect := range effects {
		names[i] = string(effect)
	}

	return names
}

// GetMusicTrackNames returns all available music track names
func (sl *SoundLibrary) GetMusicTrackNames() []string {
	tracks := []MusicTrack{
		MusicMenu, MusicGameplay, MusicCreative, MusicNight,
		MusicUnderground, MusicCombat, MusicBoss, MusicVictory, MusicDanger,
	}

	names := make([]string, len(tracks))
	for i, track := range tracks {
		names[i] = string(track)
	}

	return names
}

// LoadDefaultSoundEffects attempts to load common sound effects
// This would typically load from embedded assets or files
func (sl *SoundLibrary) LoadDefaultSoundEffects() {
	log.Printf("Loading default sound effects...")

	// This would be implemented to load actual audio files
	// For now, we'll just log what would be loaded

	soundEffects := sl.GetSoundEffectNames()
	for _, effect := range soundEffects {
		log.Printf("Would load sound effect: %s", effect)
		// sl.manager.LoadSound(effect, audioData, AudioTypeSFX, 1.0, false)
	}

	musicTracks := sl.GetMusicTrackNames()
	for _, track := range musicTracks {
		log.Printf("Would load music track: %s", track)
		// sl.manager.LoadSound(track, audioData, AudioTypeMusic, 0.7, true)
	}

	log.Printf("Sound library initialized with %d sound effects and %d music tracks",
		len(soundEffects), len(musicTracks))
}
