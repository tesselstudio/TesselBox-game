package gametime

import (
	"fmt"
	"math"
	"time"
)

// TimeOfDay represents different times of day
type TimeOfDay int

const (
	Dawn TimeOfDay = iota
	Morning
	Noon
	Afternoon
	Dusk
	Night
	Midnight
)

// DayNightCycle manages the day/night cycle and lighting
type DayNightCycle struct {
	// Time configuration
	DayLengthSeconds float64   // Length of a full day in real seconds
	StartTime        time.Time // When the cycle started
	Paused           bool

	// Current time values
	GameTime float64 // Current game time (0-1, where 1 = full day)
	DayCount int     // Number of days passed

	// Lighting values (0-1, where 1 is full brightness)
	AmbientLight float64
	SkyLight     float64
	BlockLight   float64

	// Colors for different times
	SkyColors map[TimeOfDay][3]float64 // RGB values
}

// NewDayNightCycle creates a new day/night cycle
func NewDayNightCycle(dayLengthSeconds float64) *DayNightCycle {
	cycle := &DayNightCycle{
		DayLengthSeconds: dayLengthSeconds,
		StartTime:        time.Now(),
		GameTime:         0.3, // Start at morning
		DayCount:         0,
		SkyColors:        make(map[TimeOfDay][3]float64),
	}

	// Initialize sky colors for different times
	cycle.SkyColors[Dawn] = [3]float64{0.8, 0.6, 0.4}      // Orange
	cycle.SkyColors[Morning] = [3]float64{0.6, 0.8, 1.0}   // Light blue
	cycle.SkyColors[Noon] = [3]float64{0.4, 0.7, 1.0}      // Blue
	cycle.SkyColors[Afternoon] = [3]float64{0.5, 0.8, 1.0} // Light blue
	cycle.SkyColors[Dusk] = [3]float64{0.9, 0.5, 0.3}      // Red-orange
	cycle.SkyColors[Night] = [3]float64{0.1, 0.1, 0.3}     // Dark blue
	cycle.SkyColors[Midnight] = [3]float64{0.0, 0.0, 0.1}  // Very dark

	// Update lighting initially
	cycle.updateLighting()

	return cycle
}

// Update updates the day/night cycle
func (dnc *DayNightCycle) Update() {
	if dnc.Paused {
		return
	}

	elapsed := time.Since(dnc.StartTime).Seconds()
	dnc.GameTime = math.Mod(elapsed/dnc.DayLengthSeconds, 1.0)

	// Check if a new day has started
	newDayCount := int(elapsed / dnc.DayLengthSeconds)
	if newDayCount > dnc.DayCount {
		dnc.DayCount = newDayCount
		// Could trigger day change events here
	}

	dnc.updateLighting()
}

// updateLighting calculates lighting values based on current time
func (dnc *DayNightCycle) updateLighting() {
	// Calculate lighting intensity (0-1)
	// Daytime: GameTime 0.2-0.8 (dawn to dusk)
	// Nighttime: GameTime 0.8-0.2 (dusk to dawn)

	var lightIntensity float64

	if dnc.GameTime >= 0.2 && dnc.GameTime <= 0.8 {
		// Daytime - full brightness
		lightIntensity = 1.0
	} else if dnc.GameTime > 0.8 && dnc.GameTime < 0.9 {
		// Dusk - fading
		fadeProgress := (dnc.GameTime - 0.8) / 0.1
		lightIntensity = 1.0 - fadeProgress*0.7
	} else if dnc.GameTime >= 0.9 || dnc.GameTime <= 0.1 {
		// Night - very dark
		lightIntensity = 0.1
	} else if dnc.GameTime > 0.1 && dnc.GameTime < 0.2 {
		// Dawn - brightening
		fadeProgress := (dnc.GameTime - 0.1) / 0.1
		lightIntensity = 0.1 + fadeProgress*0.9
	}

	// Clamp values
	if lightIntensity < 0.1 {
		lightIntensity = 0.1
	}
	if lightIntensity > 1.0 {
		lightIntensity = 1.0
	}

	dnc.AmbientLight = lightIntensity
	dnc.SkyLight = lightIntensity
	dnc.BlockLight = lightIntensity * 0.8 // Blocks are slightly darker
}

// GetCurrentTimeOfDay returns the current time of day
func (dnc *DayNightCycle) GetCurrentTimeOfDay() TimeOfDay {
	if dnc.GameTime < 0.1 {
		return Midnight
	} else if dnc.GameTime < 0.2 {
		return Dawn
	} else if dnc.GameTime < 0.4 {
		return Morning
	} else if dnc.GameTime < 0.6 {
		return Noon
	} else if dnc.GameTime < 0.8 {
		return Afternoon
	} else if dnc.GameTime < 0.9 {
		return Dusk
	} else {
		return Night
	}
}

// GetSkyColor returns the current sky color as RGB values
func (dnc *DayNightCycle) GetSkyColor() (float64, float64, float64) {
	timeOfDay := dnc.GetCurrentTimeOfDay()
	color := dnc.SkyColors[timeOfDay]

	// Interpolate between time periods for smoother transitions
	var nextTime TimeOfDay
	var interpolationFactor float64

	switch timeOfDay {
	case Midnight:
		nextTime = Dawn
		if dnc.GameTime < 0.1 {
			interpolationFactor = dnc.GameTime / 0.1
		} else {
			interpolationFactor = (dnc.GameTime - 0.9) / 0.1
		}
	case Dawn:
		nextTime = Morning
		interpolationFactor = (dnc.GameTime - 0.1) / 0.1
	case Morning:
		nextTime = Noon
		interpolationFactor = (dnc.GameTime - 0.2) / 0.2
	case Noon:
		nextTime = Afternoon
		interpolationFactor = (dnc.GameTime - 0.4) / 0.2
	case Afternoon:
		nextTime = Dusk
		interpolationFactor = (dnc.GameTime - 0.6) / 0.2
	case Dusk:
		nextTime = Night
		interpolationFactor = (dnc.GameTime - 0.8) / 0.1
	case Night:
		nextTime = Midnight
		interpolationFactor = (dnc.GameTime - 0.9) / 0.1
	}

	if interpolationFactor > 0 && interpolationFactor < 1 {
		nextColor := dnc.SkyColors[nextTime]
		r := color[0] + (nextColor[0]-color[0])*interpolationFactor
		g := color[1] + (nextColor[1]-color[1])*interpolationFactor
		b := color[2] + (nextColor[2]-color[2])*interpolationFactor
		return r, g, b
	}

	return color[0], color[1], color[2]
}

// GetTimeString returns a human-readable time string
func (dnc *DayNightCycle) GetTimeString() string {
	timeOfDay := dnc.GetCurrentTimeOfDay()
	var timeName string

	switch timeOfDay {
	case Midnight:
		timeName = "Midnight"
	case Dawn:
		timeName = "Dawn"
	case Morning:
		timeName = "Morning"
	case Noon:
		timeName = "Noon"
	case Afternoon:
		timeName = "Afternoon"
	case Dusk:
		timeName = "Dusk"
	case Night:
		timeName = "Night"
	}

	hours := int(dnc.GameTime * 24)
	minutes := int((dnc.GameTime*24 - float64(hours)) * 60)

	return fmt.Sprintf("Day %d, %s (%02d:%02d)", dnc.DayCount+1, timeName, hours, minutes)
}

// SetTime sets the game time (0-1)
func (dnc *DayNightCycle) SetTime(gameTime float64) {
	dnc.GameTime = math.Mod(gameTime, 1.0)
	if dnc.GameTime < 0 {
		dnc.GameTime += 1.0
	}
	dnc.StartTime = time.Now().Add(-time.Duration(dnc.GameTime*dnc.DayLengthSeconds) * time.Second)
	dnc.updateLighting()
}

// SetDayLength changes the length of a day
func (dnc *DayNightCycle) SetDayLength(seconds float64) {
	if seconds > 0 {
		dnc.DayLengthSeconds = seconds
	}
}

// Pause pauses the day/night cycle
func (dnc *DayNightCycle) Pause() {
	dnc.Paused = true
}

// Resume resumes the day/night cycle
func (dnc *DayNightCycle) Resume() {
	dnc.Paused = false
	dnc.StartTime = time.Now().Add(-time.Duration(dnc.GameTime*dnc.DayLengthSeconds) * time.Second)
}
