package entities

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ============================================================================
// Plugin Configuration Manager
// ============================================================================

// PluginConfigManager manages plugin configurations
type PluginConfigManager struct {
	configPath     string
	globalConfig   *GlobalPluginConfig
	pluginConfigs  map[string]*PluginConfig
	configDefaults map[string]*PluginConfig
}

// GlobalPluginConfig contains global plugin settings
type GlobalPluginConfig struct {
	Enabled            bool           `yaml:"enabled"`
	PluginDirectory    string         `yaml:"pluginDirectory"`
	HotReload          bool           `yaml:"hotReload"`
	DefaultPermissions []string       `yaml:"defaultPermissions"`
	Security           SecurityConfig `yaml:"security"`
	Logging            LoggingConfig  `yaml:"logging"`
}

// SecurityConfig contains security settings for plugins
type SecurityConfig struct {
	SandboxEnabled bool     `yaml:"sandboxEnabled"`
	AllowedPaths   []string `yaml:"allowedPaths"`
	BlockedPaths   []string `yaml:"blockedPaths"`
	MaxMemory      int64    `yaml:"maxMemory"` // in bytes
	MaxCPU         int      `yaml:"maxCPU"`    // percentage
	Timeout        int      `yaml:"timeout"`   // in seconds
}

// LoggingConfig contains logging settings for plugins
type LoggingConfig struct {
	Level      string   `yaml:"level"` // debug, info, warn, error
	ToFile     bool     `yaml:"toFile"`
	LogDir     string   `yaml:"logDir"`
	MaxLogSize int64    `yaml:"maxLogSize"` // in bytes
	MaxLogs    int      `yaml:"maxLogs"`
	Plugins    []string `yaml:"plugins"` // specific plugins to log
}

// NewPluginConfigManager creates a new plugin configuration manager
func NewPluginConfigManager(configPath string) *PluginConfigManager {
	pcm := &PluginConfigManager{
		configPath:     configPath,
		pluginConfigs:  make(map[string]*PluginConfig),
		configDefaults: make(map[string]*PluginConfig),
	}

	// Set default global config
	pcm.globalConfig = &GlobalPluginConfig{
		Enabled:            true,
		PluginDirectory:    "plugins",
		HotReload:          false,
		DefaultPermissions: []string{"*"},
		Security: SecurityConfig{
			SandboxEnabled: false,
			AllowedPaths:   []string{},
			BlockedPaths:   []string{"/etc", "/sys", "/proc"},
			MaxMemory:      100 * 1024 * 1024, // 100MB
			MaxCPU:         50,                // 50%
			Timeout:        30,                // 30 seconds
		},
		Logging: LoggingConfig{
			Level:      "info",
			ToFile:     true,
			LogDir:     "logs/plugins",
			MaxLogSize: 10 * 1024 * 1024, // 10MB
			MaxLogs:    5,
			Plugins:    []string{},
		},
	}

	return pcm
}

// LoadGlobalConfig loads the global plugin configuration
func (pcm *PluginConfigManager) LoadGlobalConfig() error {
	configFile := filepath.Join(pcm.configPath, "plugins.yaml")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Create default config file
		return pcm.SaveGlobalConfig()
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read global config: %v", err)
	}

	if err := yaml.Unmarshal(data, pcm.globalConfig); err != nil {
		return fmt.Errorf("failed to parse global config: %v", err)
	}

	return nil
}

// SaveGlobalConfig saves the global plugin configuration
func (pcm *PluginConfigManager) SaveGlobalConfig() error {
	// Ensure config directory exists
	if err := os.MkdirAll(pcm.configPath, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	configFile := filepath.Join(pcm.configPath, "plugins.yaml")
	data, err := yaml.Marshal(pcm.globalConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal global config: %v", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write global config: %v", err)
	}

	return nil
}

// LoadPluginConfig loads configuration for a specific plugin
func (pcm *PluginConfigManager) LoadPluginConfig(pluginName string) (*PluginConfig, error) {
	configFile := filepath.Join(pcm.configPath, pluginName+".yaml")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Return default config
		return pcm.GetDefaultConfig(pluginName), nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin config for %s: %v", pluginName, err)
	}

	var config PluginConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse plugin config for %s: %v", pluginName, err)
	}

	// Store in cache
	pcm.pluginConfigs[pluginName] = &config
	return &config, nil
}

// SavePluginConfig saves configuration for a specific plugin
func (pcm *PluginConfigManager) SavePluginConfig(pluginName string, config *PluginConfig) error {
	// Ensure config directory exists
	if err := os.MkdirAll(pcm.configPath, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	configFile := filepath.Join(pcm.configPath, pluginName+".yaml")
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal plugin config for %s: %v", pluginName, err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write plugin config for %s: %v", pluginName, err)
	}

	// Update cache
	pcm.pluginConfigs[pluginName] = config
	return nil
}

// GetDefaultConfig returns the default configuration for a plugin
func (pcm *PluginConfigManager) GetDefaultConfig(pluginName string) *PluginConfig {
	// Check if we have a cached default
	if defaultConfig, exists := pcm.configDefaults[pluginName]; exists {
		return defaultConfig
	}

	// Create default config
	config := &PluginConfig{
		Enabled:     true,
		Priority:    100,
		Permissions: pcm.globalConfig.DefaultPermissions,
		Settings:    make(map[string]interface{}),
		AutoLoad:    true,
		AutoReload:  pcm.globalConfig.HotReload,
	}

	// Cache the default
	pcm.configDefaults[pluginName] = config
	return config
}

// SetDefaultConfig sets the default configuration for a plugin
func (pcm *PluginConfigManager) SetDefaultConfig(pluginName string, config *PluginConfig) {
	pcm.configDefaults[pluginName] = config
}

// GetGlobalConfig returns the global configuration
func (pcm *PluginConfigManager) GetGlobalConfig() *GlobalPluginConfig {
	return pcm.globalConfig
}

// UpdateGlobalConfig updates the global configuration
func (pcm *PluginConfigManager) UpdateGlobalConfig(config *GlobalPluginConfig) error {
	pcm.globalConfig = config
	return pcm.SaveGlobalConfig()
}

// GetPluginConfig returns the configuration for a plugin
func (pcm *PluginConfigManager) GetPluginConfig(pluginName string) (*PluginConfig, error) {
	if config, exists := pcm.pluginConfigs[pluginName]; exists {
		return config, nil
	}
	return pcm.LoadPluginConfig(pluginName)
}

// UpdatePluginConfig updates the configuration for a plugin
func (pcm *PluginConfigManager) UpdatePluginConfig(pluginName string, config *PluginConfig) error {
	return pcm.SavePluginConfig(pluginName, config)
}

// ListPluginConfigs returns a list of all plugin configurations
func (pcm *PluginConfigManager) ListPluginConfigs() map[string]*PluginConfig {
	result := make(map[string]*PluginConfig)
	for name, config := range pcm.pluginConfigs {
		result[name] = config
	}
	return result
}

// ValidatePluginConfig validates a plugin configuration
func (pcm *PluginConfigManager) ValidatePluginConfig(config *PluginConfig) error {
	// Validate permissions
	for _, permission := range config.Permissions {
		if !pcm.isValidPermission(permission) {
			return fmt.Errorf("invalid permission: %s", permission)
		}
	}

	// Validate priority
	if config.Priority < 0 || config.Priority > 1000 {
		return fmt.Errorf("priority must be between 0 and 1000")
	}

	// Validate settings
	if err := pcm.validateSettings(config.Settings); err != nil {
		return fmt.Errorf("invalid settings: %v", err)
	}

	return nil
}

// isValidPermission checks if a permission string is valid
func (pcm *PluginConfigManager) isValidPermission(permission string) bool {
	validPermissions := []string{
		"*",
		"entity.create",
		"entity.remove",
		"entity.modify",
		"entity.get",
		"entity.find",
		"component.create",
		"event.publish",
		"event.subscribe",
		"world.read",
		"world.modify",
		"template.get",
		"template.register",
		"system.register",
		"system.unregister",
	}

	for _, valid := range validPermissions {
		if permission == valid || strings.HasPrefix(permission, "custom.") {
			return true
		}
	}

	return false
}

// validateSettings validates plugin settings
func (pcm *PluginConfigManager) validateSettings(settings map[string]interface{}) error {
	// Basic validation - could be extended based on plugin requirements
	for key, value := range settings {
		if strings.Contains(key, "path") {
			if path, ok := value.(string); ok {
				if !pcm.isValidPath(path) {
					return fmt.Errorf("invalid path in setting %s: %s", key, path)
				}
			}
		}
	}
	return nil
}

// isValidPath checks if a path is valid and allowed
func (pcm *PluginConfigManager) isValidPath(path string) bool {
	// Check against blocked paths
	for _, blocked := range pcm.globalConfig.Security.BlockedPaths {
		if strings.HasPrefix(path, blocked) {
			return false
		}
	}

	// If allowed paths are specified, check against them
	if len(pcm.globalConfig.Security.AllowedPaths) > 0 {
		for _, allowed := range pcm.globalConfig.Security.AllowedPaths {
			if strings.HasPrefix(path, allowed) {
				return true
			}
		}
		return false
	}

	return true
}

// ============================================================================
// Plugin Configuration Templates
// ============================================================================

// CreateMagicPluginConfig creates a configuration for the magic plugin
func CreateMagicPluginConfig() *PluginConfig {
	return &PluginConfig{
		Enabled:     true,
		Priority:    50,
		Permissions: []string{"entity.create", "entity.modify", "event.publish", "world.modify"},
		Settings: map[string]interface{}{
			"manaRegenRate":         1.0,
			"maxMana":               100,
			"spellDamageMultiplier": 1.0,
			"enableParticles":       true,
			"particleDensity":       "medium",
		},
		AutoLoad:   true,
		AutoReload: false,
	}
}

// CreateTechPluginConfig creates a configuration for the tech plugin
func CreatePluginConfig() *PluginConfig {
	return &PluginConfig{
		Enabled:     true,
		Priority:    60,
		Permissions: []string{"entity.create", "entity.modify", "event.publish", "world.modify", "system.register"},
		Settings: map[string]interface{}{
			"powerConsumptionRate": 1.0,
			"maxPower":             1000,
			"efficiencyBonus":      0.1,
			"enableAutomation":     true,
			"researchSpeed":        1.0,
		},
		AutoLoad:   true,
		AutoReload: false,
	}
}

// CreateSecurityPluginConfig creates a configuration for a security-focused plugin
func CreateSecurityPluginConfig() *PluginConfig {
	return &PluginConfig{
		Enabled:     true,
		Priority:    10,                                                      // High priority for security
		Permissions: []string{"entity.get", "event.subscribe", "world.read"}, // Read-only permissions
		Settings: map[string]interface{}{
			"logLevel":            "info",
			"auditEvents":         true,
			"alertThreshold":      5,
			"blockedActions":      []string{"world.modify"},
			"requireConfirmation": true,
		},
		AutoLoad:   true,
		AutoReload: false,
	}
}

// ExportConfig exports all configurations to a single file
func (pcm *PluginConfigManager) ExportConfig(exportPath string) error {
	exportData := map[string]interface{}{
		"global":  pcm.globalConfig,
		"plugins": pcm.pluginConfigs,
	}

	data, err := yaml.Marshal(exportData)
	if err != nil {
		return fmt.Errorf("failed to marshal export data: %v", err)
	}

	if err := os.WriteFile(exportPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write export file: %v", err)
	}

	return nil
}

// ImportConfig imports configurations from a file
func (pcm *PluginConfigManager) ImportConfig(importPath string) error {
	data, err := os.ReadFile(importPath)
	if err != nil {
		return fmt.Errorf("failed to read import file: %v", err)
	}

	var importData map[string]interface{}
	if err := yaml.Unmarshal(data, &importData); err != nil {
		return fmt.Errorf("failed to parse import data: %v", err)
	}

	// Import global config
	if globalData, exists := importData["global"]; exists {
		globalConfig := &GlobalPluginConfig{}
		globalBytes, _ := yaml.Marshal(globalData)
		if err := yaml.Unmarshal(globalBytes, globalConfig); err == nil {
			pcm.globalConfig = globalConfig
			pcm.SaveGlobalConfig()
		}
	}

	// Import plugin configs
	if pluginsData, exists := importData["plugins"]; exists {
		if pluginsMap, ok := pluginsData.(map[string]interface{}); ok {
			for pluginName, configData := range pluginsMap {
				config := &PluginConfig{}
				configBytes, _ := yaml.Marshal(configData)
				if err := yaml.Unmarshal(configBytes, config); err == nil {
					pcm.SavePluginConfig(pluginName, config)
				}
			}
		}
	}

	return nil
}
