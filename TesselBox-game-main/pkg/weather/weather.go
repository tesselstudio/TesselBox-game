package weather

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// WeatherType represents different types of weather
type WeatherType int

const (
	Clear WeatherType = iota
	Rain
	Storm
	Snow
)

// WeatherState represents the current weather condition
type WeatherState struct {
	Type             WeatherType
	Intensity        float64 // 0-1, how strong the weather is
	Duration         time.Duration
	StartTime        time.Time
	ParticleCount    int
	ParticleLifetime float64
	WindSpeed        float64
	WindDirection    float64 // radians
}

// WeatherSystem manages weather effects and transitions
type WeatherSystem struct {
	CurrentWeather     *WeatherState
	NextWeather        *WeatherState
	TransitionTime     time.Time
	TransitionDuration time.Duration
	IsTransitioning    bool

	// Weather probabilities (per minute)
	ClearToRainProb  float64
	RainToStormProb  float64
	StormToClearProb float64
	RainToClearProb  float64
	ClearToSnowProb  float64 // In cold biomes
	SnowToClearProb  float64

	// Particle effects
	rainParticles []*WeatherParticle
	snowParticles []*WeatherParticle
}

// WeatherParticle represents a weather effect particle
type WeatherParticle struct {
	X, Y    float64
	VX, VY  float64
	Life    float64
	MaxLife float64
	Active  bool
}

// NewWeatherSystem creates a new weather system
func NewWeatherSystem() *WeatherSystem {
	ws := &WeatherSystem{
		CurrentWeather: &WeatherState{
			Type:          Clear,
			Intensity:     0.0,
			Duration:      5 * time.Minute,
			StartTime:     time.Now(),
			ParticleCount: 0,
		},
		ClearToRainProb:  0.02,  // 2% chance per minute to start raining
		RainToStormProb:  0.01,  // 1% chance per minute to become storm
		StormToClearProb: 0.05,  // 5% chance per minute to clear up
		RainToClearProb:  0.03,  // 3% chance per minute to stop raining
		ClearToSnowProb:  0.005, // 0.5% chance per minute for snow (cold areas)
		SnowToClearProb:  0.02,  // 2% chance per minute to stop snowing
	}

	// Pre-allocate particle pools
	ws.rainParticles = make([]*WeatherParticle, 500)
	ws.snowParticles = make([]*WeatherParticle, 200)

	for i := range ws.rainParticles {
		ws.rainParticles[i] = &WeatherParticle{}
	}
	for i := range ws.snowParticles {
		ws.snowParticles[i] = &WeatherParticle{}
	}

	return ws
}

// Update updates the weather system
func (ws *WeatherSystem) Update(deltaTime float64, screenWidth, screenHeight int) {
	now := time.Now()

	// Handle weather transitions
	if ws.IsTransitioning {
		if now.After(ws.TransitionTime.Add(ws.TransitionDuration)) {
			// Transition complete
			ws.CurrentWeather = ws.NextWeather
			ws.NextWeather = nil
			ws.IsTransitioning = false
		}
	} else {
		// Check if current weather should end or change
		if now.After(ws.CurrentWeather.StartTime.Add(ws.CurrentWeather.Duration)) {
			ws.transitionToNewWeather()
		}
	}

	// Update particles
	ws.updateParticles(deltaTime, screenWidth, screenHeight)
}

// transitionToNewWeather transitions to a new weather state
func (ws *WeatherSystem) transitionToNewWeather() {
	currentType := ws.CurrentWeather.Type

	var newType WeatherType
	var duration time.Duration

	// Determine next weather based on current weather and probabilities
	randVal := rand.Float64()

	switch currentType {
	case Clear:
		if randVal < ws.ClearToRainProb {
			newType = Rain
			duration = time.Duration(300+rand.Intn(600)) * time.Second // 5-15 minutes
		} else if randVal < ws.ClearToRainProb+ws.ClearToSnowProb {
			newType = Snow
			duration = time.Duration(180+rand.Intn(420)) * time.Second // 3-10 minutes
		} else {
			newType = Clear
			duration = time.Duration(60+rand.Intn(240)) * time.Second // 1-5 minutes
		}
	case Rain:
		if randVal < ws.RainToStormProb {
			newType = Storm
			duration = time.Duration(120+rand.Intn(300)) * time.Second // 2-7 minutes
		} else if randVal < ws.RainToStormProb+ws.RainToClearProb {
			newType = Clear
			duration = time.Duration(60+rand.Intn(180)) * time.Second // 1-4 minutes
		} else {
			newType = Rain
			duration = time.Duration(180+rand.Intn(420)) * time.Second // 3-10 minutes
		}
	case Storm:
		if randVal < ws.StormToClearProb {
			newType = Clear
			duration = time.Duration(30+rand.Intn(90)) * time.Second // 0.5-2.5 minutes
		} else {
			newType = Storm
			duration = time.Duration(60+rand.Intn(180)) * time.Second // 1-4 minutes
		}
	case Snow:
		if randVal < ws.SnowToClearProb {
			newType = Clear
			duration = time.Duration(60+rand.Intn(120)) * time.Second // 1-3 minutes
		} else {
			newType = Snow
			duration = time.Duration(120+rand.Intn(300)) * time.Second // 2-7 minutes
		}
	}

	ws.NextWeather = &WeatherState{
		Type:             newType,
		Intensity:        0.0, // Will be set during transition
		Duration:         duration,
		StartTime:        time.Now(),
		ParticleCount:    ws.getParticleCountForWeather(newType),
		ParticleLifetime: ws.getParticleLifetimeForWeather(newType),
		WindSpeed:        ws.getWindSpeedForWeather(newType),
		WindDirection:    rand.Float64() * 2 * 3.14159, // Random direction
	}

	ws.TransitionTime = time.Now()
	ws.TransitionDuration = 30 * time.Second // 30 second transition
	ws.IsTransitioning = true
}

// getParticleCountForWeather returns the number of particles for a weather type
func (ws *WeatherSystem) getParticleCountForWeather(weatherType WeatherType) int {
	switch weatherType {
	case Rain:
		return 200
	case Storm:
		return 400
	case Snow:
		return 100
	default:
		return 0
	}
}

// getParticleLifetimeForWeather returns particle lifetime for a weather type
func (ws *WeatherSystem) getParticleLifetimeForWeather(weatherType WeatherType) float64 {
	switch weatherType {
	case Rain, Storm:
		return 2.0 // 2 seconds
	case Snow:
		return 8.0 // 8 seconds
	default:
		return 0
	}
}

// getWindSpeedForWeather returns wind speed for a weather type
func (ws *WeatherSystem) getWindSpeedForWeather(weatherType WeatherType) float64 {
	switch weatherType {
	case Clear:
		return 10.0
	case Rain:
		return 50.0
	case Storm:
		return 150.0
	case Snow:
		return 20.0
	default:
		return 10.0
	}
}

// updateParticles updates weather particles
func (ws *WeatherSystem) updateParticles(deltaTime float64, screenWidth, screenHeight int) {
	weather := ws.CurrentWeather
	if ws.IsTransitioning {
		weather = ws.NextWeather
	}

	// Update intensity during transitions
	if ws.IsTransitioning {
		elapsed := time.Since(ws.TransitionTime).Seconds()
		progress := elapsed / ws.TransitionDuration.Seconds()
		if progress < 0.5 {
			// Fading out old weather
			weather.Intensity = 1.0 - (progress * 2.0)
		} else {
			// Fading in new weather
			weather.Intensity = (progress - 0.5) * 2.0
		}
	} else {
		weather.Intensity = 1.0
	}

	// Update particles based on weather type
	switch weather.Type {
	case Rain, Storm:
		ws.updateRainParticles(deltaTime, screenWidth, screenHeight, weather)
	case Snow:
		ws.updateSnowParticles(deltaTime, screenWidth, screenHeight, weather)
	default:
		// Clear weather - no particles
	}
}

// updateRainParticles updates rain particles
func (ws *WeatherSystem) updateRainParticles(deltaTime float64, screenWidth, screenHeight int, weather *WeatherState) {
	targetCount := int(float64(weather.ParticleCount) * weather.Intensity)

	// Spawn new particles
	for i := 0; i < targetCount-len(ws.getActiveRainParticles()); i++ {
		particle := ws.getInactiveRainParticle()
		if particle != nil {
			particle.X = rand.Float64()*float64(screenWidth+100) - 50
			particle.Y = -10
			particle.VX = weather.WindSpeed * 0.3
			particle.VY = 300 + rand.Float64()*100 // Fast falling
			particle.Life = weather.ParticleLifetime
			particle.MaxLife = weather.ParticleLifetime
			particle.Active = true
		}
	}

	// Update existing particles
	for _, particle := range ws.rainParticles {
		if particle.Active {
			particle.X += particle.VX * deltaTime
			particle.Y += particle.VY * deltaTime
			particle.Life -= deltaTime

			// Deactivate when off screen or dead
			if particle.Y > float64(screenHeight+50) || particle.Life <= 0 {
				particle.Active = false
			}
		}
	}
}

// updateSnowParticles updates snow particles
func (ws *WeatherSystem) updateSnowParticles(deltaTime float64, screenWidth, screenHeight int, weather *WeatherState) {
	targetCount := int(float64(weather.ParticleCount) * weather.Intensity)

	// Spawn new particles
	for i := 0; i < targetCount-len(ws.getActiveSnowParticles()); i++ {
		particle := ws.getInactiveSnowParticle()
		if particle != nil {
			particle.X = rand.Float64()*float64(screenWidth+100) - 50
			particle.Y = -10
			particle.VX = weather.WindSpeed * 0.1 * (rand.Float64()*2 - 1) // Gentle horizontal drift
			particle.VY = 30 + rand.Float64()*20                           // Slow falling
			particle.Life = weather.ParticleLifetime
			particle.MaxLife = weather.ParticleLifetime
			particle.Active = true
		}
	}

	// Update existing particles
	for _, particle := range ws.snowParticles {
		if particle.Active {
			particle.X += particle.VX * deltaTime
			particle.Y += particle.VY * deltaTime
			particle.Life -= deltaTime

			// Deactivate when off screen or dead
			if particle.Y > float64(screenHeight+50) || particle.Life <= 0 {
				particle.Active = false
			}
		}
	}
}

// getInactiveRainParticle returns an inactive rain particle
func (ws *WeatherSystem) getInactiveRainParticle() *WeatherParticle {
	for _, p := range ws.rainParticles {
		if !p.Active {
			return p
		}
	}
	return nil
}

// getInactiveSnowParticle returns an inactive snow particle
func (ws *WeatherSystem) getInactiveSnowParticle() *WeatherParticle {
	for _, p := range ws.snowParticles {
		if !p.Active {
			return p
		}
	}
	return nil
}

// getActiveRainParticles returns active rain particles
func (ws *WeatherSystem) getActiveRainParticles() []*WeatherParticle {
	active := []*WeatherParticle{}
	for _, p := range ws.rainParticles {
		if p.Active {
			active = append(active, p)
		}
	}
	return active
}

// getActiveSnowParticles returns active snow particles
func (ws *WeatherSystem) getActiveSnowParticles() []*WeatherParticle {
	active := []*WeatherParticle{}
	for _, p := range ws.snowParticles {
		if p.Active {
			active = append(active, p)
		}
	}
	return active
}

// GetWeatherInfo returns current weather information
func (ws *WeatherSystem) GetWeatherInfo() (WeatherType, float64, string) {
	weather := ws.CurrentWeather
	if ws.IsTransitioning {
		weather = ws.NextWeather
	}

	var weatherName string
	switch weather.Type {
	case Clear:
		weatherName = "Clear"
	case Rain:
		weatherName = "Rain"
	case Storm:
		weatherName = "Storm"
	case Snow:
		weatherName = "Snow"
	}

	return weather.Type, weather.Intensity, weatherName
}

// AffectsVisibility returns true if current weather reduces visibility
func (ws *WeatherSystem) AffectsVisibility() bool {
	weatherType, intensity, _ := ws.GetWeatherInfo()
	return (weatherType == Rain || weatherType == Storm || weatherType == Snow) && intensity > 0.3
}

// AffectsMovement returns movement speed multiplier due to weather
func (ws *WeatherSystem) AffectsMovement() float64 {
	weatherType, intensity, _ := ws.GetWeatherInfo()

	switch weatherType {
	case Rain:
		return 1.0 - intensity*0.1 // Up to 10% slower in rain
	case Storm:
		return 1.0 - intensity*0.3 // Up to 30% slower in storms
	case Snow:
		return 1.0 - intensity*0.2 // Up to 20% slower in snow
	default:
		return 1.0 // Normal speed
	}
}

// Draw renders weather particles
func (ws *WeatherSystem) Draw(screen *ebiten.Image, cameraX, cameraY float64) {
	weather := ws.CurrentWeather
	if ws.IsTransitioning {
		weather = ws.NextWeather
	}

	switch weather.Type {
	case Rain, Storm:
		ws.drawRainParticles(screen, cameraX, cameraY, weather.Intensity)
	case Snow:
		ws.drawSnowParticles(screen, cameraX, cameraY, weather.Intensity)
	}
}

// drawRainParticles draws rain particles
func (ws *WeatherSystem) drawRainParticles(screen *ebiten.Image, cameraX, cameraY, intensity float64) {
	for _, particle := range ws.rainParticles {
		if particle.Active {
			screenX := particle.X - cameraX
			screenY := particle.Y - cameraY

			// Draw rain drop as a line
			length := 8.0 * intensity
			alpha := uint8(200 * intensity * (particle.Life / particle.MaxLife))

			// Draw as a vertical line
			ebitenutil.DrawLine(screen, screenX, screenY, screenX, screenY+length, color.RGBA{200, 220, 255, alpha})
		}
	}
}

// drawSnowParticles draws snow particles
func (ws *WeatherSystem) drawSnowParticles(screen *ebiten.Image, cameraX, cameraY, intensity float64) {
	for _, particle := range ws.snowParticles {
		if particle.Active {
			screenX := particle.X - cameraX
			screenY := particle.Y - cameraY

			// Draw snowflake as small circles
			size := 2.0 + rand.Float64()*2.0 // Random size
			alpha := uint8(180 * intensity * (particle.Life / particle.MaxLife))

			// Draw multiple small dots for snowflake effect
			for i := 0; i < 3; i++ {
				offsetX := (rand.Float64() - 0.5) * 4
				offsetY := (rand.Float64() - 0.5) * 4
				ebitenutil.DrawRect(screen, screenX+offsetX-size/2, screenY+offsetY-size/2, size, size, color.RGBA{255, 255, 255, alpha})
			}
		}
	}
}
