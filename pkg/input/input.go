package input

import (
	"encoding/json"
	"log"

	"tesselbox/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// KeyBinding represents a keyboard or mouse binding
type KeyBinding struct {
	Key      ebiten.Key `json:"key,omitempty"`
	Mouse    int        `json:"mouse,omitempty"` // 0=left, 1=right, 2=middle
	IsMouse  bool       `json:"is_mouse"`
	Action   string     `json:"action"`
	Modifier ebiten.Key `json:"modifier,omitempty"`
}

// InputConfig holds all key bindings
type InputConfig struct {
	Bindings map[string]KeyBinding `json:"bindings"`
}

// DefaultInputConfig returns the default key bindings
func DefaultInputConfig() *InputConfig {
	return &InputConfig{
		Bindings: map[string]KeyBinding{
			"move_left":  {Key: ebiten.KeyA, Action: "move_left"},
			"move_right": {Key: ebiten.KeyD, Action: "move_right"},
			"move_up":    {Key: ebiten.KeyW, Action: "move_up"},
			"move_down":  {Key: ebiten.KeyS, Action: "move_down"},
			"jump":       {Key: ebiten.KeySpace, Action: "jump"},
			"mine":       {Mouse: 0, IsMouse: true, Action: "mine"},
			"place":      {Mouse: 1, IsMouse: true, Action: "place"},
			"drop":       {Key: ebiten.KeyQ, Action: "drop"},
			"inventory":  {Key: ebiten.KeyE, Action: "inventory"},
			"crafting":   {Key: ebiten.KeyC, Action: "crafting"},
			"hotbar_1":   {Key: ebiten.Key1, Action: "hotbar_1"},
			"hotbar_2":   {Key: ebiten.Key2, Action: "hotbar_2"},
			"hotbar_3":   {Key: ebiten.Key3, Action: "hotbar_3"},
			"hotbar_4":   {Key: ebiten.Key4, Action: "hotbar_4"},
			"hotbar_5":   {Key: ebiten.Key5, Action: "hotbar_5"},
			"hotbar_6":   {Key: ebiten.Key6, Action: "hotbar_6"},
			"hotbar_7":   {Key: ebiten.Key7, Action: "hotbar_7"},
			"hotbar_8":   {Key: ebiten.Key8, Action: "hotbar_8"},
			"hotbar_9":   {Key: ebiten.Key9, Action: "hotbar_9"},
			"chat":       {Key: ebiten.KeyT, Action: "chat"},
			"command":    {Key: ebiten.KeySlash, Action: "command"},
			"menu":       {Key: ebiten.KeyEscape, Action: "menu"},
		},
	}
}

// InputManager handles input processing with configurable bindings
type InputManager struct {
	config *InputConfig
}

// NewInputManager creates a new input manager
func NewInputManager() *InputManager {
	im := &InputManager{
		config: DefaultInputConfig(),
	}

	// Try to load custom config
	im.LoadConfig()

	return im
}

// LoadConfig attempts to load input configuration from embedded assets
func (im *InputManager) LoadConfig() error {
	data, err := assets.GetConfigFile("input_config.json")
	if err != nil {
		log.Printf("No custom input config found, using defaults")
		return err
	}

	var config InputConfig
	if err := json.Unmarshal(data, &config); err != nil {
		log.Printf("Error parsing input config: %v", err)
		return err
	}

	im.config = &config
	log.Printf("Loaded custom input configuration")
	return nil
}

// IsActionPressed checks if an action is currently pressed
func (im *InputManager) IsActionPressed(action string) bool {
	binding, exists := im.config.Bindings[action]
	if !exists {
		return false
	}

	if binding.IsMouse {
		return ebiten.IsMouseButtonPressed(ebiten.MouseButton(binding.Mouse))
	} else {
		return ebiten.IsKeyPressed(binding.Key)
	}
}

// IsActionJustPressed checks if an action was just pressed
func (im *InputManager) IsActionJustPressed(action string) bool {
	binding, exists := im.config.Bindings[action]
	if !exists {
		return false
	}

	if binding.IsMouse {
		return inpututil.IsMouseButtonJustPressed(ebiten.MouseButton(binding.Mouse))
	} else {
		return inpututil.IsKeyJustPressed(binding.Key)
	}
}

// GetKeyForAction returns the key binding for an action (for UI display)
func (im *InputManager) GetKeyForAction(action string) string {
	binding, exists := im.config.Bindings[action]
	if !exists {
		return "Unbound"
	}

	if binding.IsMouse {
		switch binding.Mouse {
		case 0:
			return "Left Mouse"
		case 1:
			return "Right Mouse"
		case 2:
			return "Middle Mouse"
		default:
			return "Mouse"
		}
	} else {
		return binding.Key.String()
	}
}

// SetBinding updates a key binding
func (im *InputManager) SetBinding(action string, binding KeyBinding) {
	if im.config.Bindings == nil {
		im.config.Bindings = make(map[string]KeyBinding)
	}
	im.config.Bindings[action] = binding
}

// GetConfig returns the current input configuration
func (im *InputManager) GetConfig() *InputConfig {
	return im.config
}
