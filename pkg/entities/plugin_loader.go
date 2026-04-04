package entities

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

// ============================================================================
// Plugin Configuration
// ============================================================================

// PluginConfig represents the configuration for a plugin
type PluginConfig struct {
	Enabled     bool                   `yaml:"enabled"`
	Priority    int                    `yaml:"priority"`
	Permissions []string               `yaml:"permissions"`
	Settings    map[string]interface{} `yaml:"settings"`
	AutoLoad    bool                   `yaml:"autoLoad"`
	AutoReload  bool                   `yaml:"autoReload"`
}

// DefaultPluginConfig returns the default plugin configuration
func DefaultPluginConfig() *PluginConfig {
	return &PluginConfig{
		Enabled:     true,
		Priority:    100,
		Permissions: []string{"*"}, // Grant all permissions by default
		Settings:    make(map[string]interface{}),
		AutoLoad:    true,
		AutoReload:  false,
	}
}

// ============================================================================
// Plugin Discovery
// ============================================================================

// PluginDiscovery handles discovering plugins in directories
type PluginDiscovery struct {
	directories []string
	watcher     *fsnotify.Watcher
	mutex       sync.RWMutex
}

// NewPluginDiscovery creates a new plugin discovery
func NewPluginDiscovery() *PluginDiscovery {
	watcher, _ := fsnotify.NewWatcher()
	return &PluginDiscovery{
		directories: make([]string, 0),
		watcher:     watcher,
	}
}

// AddDirectory adds a directory to search for plugins
func (pd *PluginDiscovery) AddDirectory(path string) error {
	pd.mutex.Lock()
	defer pd.mutex.Unlock()

	// Check if directory exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("plugin directory %s does not exist", path)
	}

	pd.directories = append(pd.directories, path)

	// Add to file watcher if available
	if pd.watcher != nil {
		if err := pd.watcher.Add(path); err != nil {
			log.Printf("Failed to watch plugin directory %s: %v", path, err)
		}
	}

	log.Printf("Added plugin directory: %s", path)
	return nil
}

// DiscoverPlugins discovers all plugins in the configured directories
func (pd *PluginDiscovery) DiscoverPlugins() ([]*PluginMetadata, error) {
	pd.mutex.RLock()
	defer pd.mutex.RUnlock()

	var plugins []*PluginMetadata

	for _, dir := range pd.directories {
		dirPlugins, err := pd.discoverInDirectory(dir)
		if err != nil {
			log.Printf("Failed to discover plugins in %s: %v", dir, err)
			continue
		}
		plugins = append(plugins, dirPlugins...)
	}

	return plugins, nil
}

// discoverInDirectory discovers plugins in a specific directory
func (pd *PluginDiscovery) discoverInDirectory(dir string) ([]*PluginMetadata, error) {
	var plugins []*PluginMetadata

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Check for .so files (compiled plugins)
		if strings.HasSuffix(path, ".so") {
			metadata, err := pd.analyzePlugin(path, "compiled")
			if err != nil {
				log.Printf("Failed to analyze plugin %s: %v", path, err)
				return nil
			}
			plugins = append(plugins, metadata)
			return nil
		}

		// Check for plugin directories (source plugins)
		if d.IsDir() || strings.HasSuffix(path, ".go") {
			// Check if this is a plugin directory
			if pd.isPluginDirectory(path) {
				metadata, err := pd.analyzePlugin(path, "source")
				if err != nil {
					log.Printf("Failed to analyze plugin directory %s: %v", path, err)
					return nil
				}
				plugins = append(plugins, metadata)
			}
		}

		return nil
	})

	return plugins, err
}

// isPluginDirectory checks if a directory contains a plugin
func (pd *PluginDiscovery) isPluginDirectory(path string) bool {
	// Check for plugin.yaml or plugin.json
	configPath := filepath.Join(path, "plugin.yaml")
	if _, err := os.Stat(configPath); err == nil {
		return true
	}

	configPath = filepath.Join(path, "plugin.json")
	if _, err := os.Stat(configPath); err == nil {
		return true
	}

	// Check for main.go with plugin package
	mainPath := filepath.Join(path, "main.go")
	if _, err := os.Stat(mainPath); err == nil {
		// Could check file contents here for plugin package
		return true
	}

	return false
}

// analyzePlugin analyzes a plugin and returns metadata
func (pd *PluginDiscovery) analyzePlugin(path string, pluginType string) (*PluginMetadata, error) {
	metadata := &PluginMetadata{
		Path:       path,
		Type:       pluginType,
		Discovered: time.Now(),
	}

	// Try to read plugin configuration
	configPath := strings.TrimSuffix(path, filepath.Ext(path)) + ".yaml"
	if _, err := os.Stat(configPath); err == nil {
		if err := pd.loadPluginConfig(configPath, metadata); err != nil {
			log.Printf("Failed to load plugin config %s: %v", configPath, err)
		}
	} else {
		// Try JSON config
		configPath = strings.TrimSuffix(path, filepath.Ext(path)) + ".json"
		if _, err := os.Stat(configPath); err == nil {
			if err := pd.loadPluginConfigJSON(configPath, metadata); err != nil {
				log.Printf("Failed to load plugin config %s: %v", configPath, err)
			}
		}
	}

	// Extract basic info from filename if no config found
	if metadata.Name == "" {
		metadata.Name = filepath.Base(strings.TrimSuffix(path, filepath.Ext(path)))
	}

	return metadata, nil
}

// loadPluginConfig loads plugin configuration from YAML
func (pd *PluginDiscovery) loadPluginConfig(path string, metadata *PluginMetadata) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var config struct {
		PluginInfo `yaml:",inline"`
		Config     PluginConfig `yaml:"config"`
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return err
	}

	metadata.PluginInfo = config.PluginInfo
	metadata.Config = config.Config
	return nil
}

// loadPluginConfigJSON loads plugin configuration from JSON
func (pd *PluginDiscovery) loadPluginConfigJSON(path string, metadata *PluginMetadata) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var config struct {
		PluginInfo `json:",inline"`
		Config     PluginConfig `json:"config"`
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	metadata.PluginInfo = config.PluginInfo
	metadata.Config = config.Config
	return nil
}

// ============================================================================
// Plugin Metadata
// ============================================================================

// PluginMetadata contains metadata about a discovered plugin
type PluginMetadata struct {
	PluginInfo
	Path       string       `yaml:"path"`
	Type       string       `yaml:"type"` // "compiled" or "source"
	Config     PluginConfig `yaml:"config"`
	Discovered time.Time    `yaml:"discovered"`
	Loaded     bool         `yaml:"loaded"`
	LoadTime   time.Time    `yaml:"loadTime,omitempty"`
	Error      string       `yaml:"error,omitempty"`
}

// ============================================================================
// Enhanced Plugin Manager
// ============================================================================

// EnhancedPluginManager extends the basic PluginManager with advanced features
type EnhancedPluginManager struct {
	*PluginManager
	discovery   *PluginDiscovery
	configs     map[string]*PluginConfig
	hotReload   bool
	fileWatcher *fsnotify.Watcher
	reloadMutex sync.Mutex
	world       interface{} // Will be *world.World
}

// NewEnhancedPluginManager creates a new enhanced plugin manager
func NewEnhancedPluginManager(entityManager *EntityManager, systemManager *SystemManager, eventBus *EventBus) *EnhancedPluginManager {
	base := NewPluginManager(entityManager, systemManager, eventBus)

	epm := &EnhancedPluginManager{
		PluginManager: base,
		discovery:     NewPluginDiscovery(),
		configs:       make(map[string]*PluginConfig),
		hotReload:     false,
		world:         nil,
	}

	// Initialize file watcher for hot reload
	if watcher, err := fsnotify.NewWatcher(); err == nil {
		epm.fileWatcher = watcher
		go epm.watchFiles()
	}

	return epm
}

// SetWorld sets the world reference for plugin API
func (epm *EnhancedPluginManager) SetWorld(world interface{}) {
	epm.world = world
	// Note: EventBus doesn't have world field anymore
}

// AddPluginDirectory adds a directory to search for plugins
func (epm *EnhancedPluginManager) AddPluginDirectory(path string) error {
	return epm.discovery.AddDirectory(path)
}

// DiscoverAndLoad discovers and loads all plugins
func (epm *EnhancedPluginManager) DiscoverAndLoad() error {
	plugins, err := epm.discovery.DiscoverPlugins()
	if err != nil {
		return fmt.Errorf("failed to discover plugins: %v", err)
	}

	// Sort plugins by priority
	sortPluginsByPriority(plugins)

	// Load plugins in order
	for _, metadata := range plugins {
		if metadata.Config.AutoLoad && metadata.Config.Enabled {
			if err := epm.LoadPluginFromMetadata(metadata); err != nil {
				log.Printf("Failed to load plugin %s: %v", metadata.Name, err)
				metadata.Error = err.Error()
			}
		}
	}

	return nil
}

// LoadPluginFromMetadata loads a plugin from metadata
func (epm *EnhancedPluginManager) LoadPluginFromMetadata(metadata *PluginMetadata) error {
	epm.reloadMutex.Lock()
	defer epm.reloadMutex.Unlock()

	// Check if already loaded
	if epm.IsLoaded(metadata.Name) {
		return fmt.Errorf("plugin %s is already loaded", metadata.Name)
	}

	// Store config
	epm.configs[metadata.Name] = &metadata.Config

	var plugin Plugin
	var err error

	switch metadata.Type {
	case "compiled":
		plugin, err = epm.loadCompiledPlugin(metadata)
	case "source":
		plugin, err = epm.loadSourcePlugin(metadata)
	default:
		return fmt.Errorf("unsupported plugin type: %s", metadata.Type)
	}

	if err != nil {
		return fmt.Errorf("failed to load plugin %s: %v", metadata.Name, err)
	}

	// Create plugin API
	api := NewPluginAPI(epm.PluginManager, metadata.Name)

	// Grant permissions
	for _, permission := range metadata.Config.Permissions {
		if permission == "*" {
			// Grant all permissions
			api.allowedActions["*"] = true
		} else {
			api.allowedActions[permission] = true
		}
	}

	// Initialize plugin
	if enhancedPlugin, ok := plugin.(EnhancedPlugin); ok {
		if err := enhancedPlugin.OnLoad(api); err != nil {
			return fmt.Errorf("plugin %s OnLoad failed: %v", metadata.Name, err)
		}
	} else {
		if err := plugin.Initialize(epm.PluginManager); err != nil {
			return fmt.Errorf("plugin %s Initialize failed: %v", metadata.Name, err)
		}
	}

	// Store plugin
	epm.plugins[metadata.Name] = plugin
	metadata.Loaded = true
	metadata.LoadTime = time.Now()

	log.Printf("Loaded plugin: %s v%s (%s)", plugin.GetName(), plugin.GetVersion(), metadata.Type)
	return nil
}

// loadCompiledPlugin loads a compiled plugin (.so file)
func (epm *EnhancedPluginManager) loadCompiledPlugin(metadata *PluginMetadata) (Plugin, error) {
	pluginObj, err := plugin.Open(metadata.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin %s: %v", metadata.Path, err)
	}

	// Look for the plugin symbol
	sym, err := pluginObj.Lookup("Plugin")
	if err != nil {
		return nil, fmt.Errorf("plugin %s does not export a Plugin symbol: %v", metadata.Name, err)
	}

	// Convert symbol to Plugin interface
	pluginInstance, ok := sym.(Plugin)
	if !ok {
		return nil, fmt.Errorf("plugin %s does not implement the Plugin interface", metadata.Name)
	}

	return pluginInstance, nil
}

// loadSourcePlugin loads a source plugin (for development)
func (epm *EnhancedPluginManager) loadSourcePlugin(metadata *PluginMetadata) (Plugin, error) {
	// Check if this is a Go plugin directory
	mainGoPath := filepath.Join(metadata.Path, "main.go")
	if _, err := os.Stat(mainGoPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("no main.go found in plugin directory %s", metadata.Path)
	}

	// Create a temporary build directory
	tempDir, err := os.MkdirTemp("", "tesselbox-plugin-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Build the plugin as a shared library (.so file)
	soPath := filepath.Join(tempDir, "plugin.so")

	// Use go build to create the plugin
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", soPath, mainGoPath)
	cmd.Dir = metadata.Path

	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("failed to build plugin: %v\nOutput: %s", err, string(output))
	}

	// Load the compiled plugin
	pluginObj, err := plugin.Open(soPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open compiled plugin: %v", err)
	}

	// Look for the plugin symbol
	sym, err := pluginObj.Lookup("Plugin")
	if err != nil {
		return nil, fmt.Errorf("plugin does not export a Plugin symbol: %v", err)
	}

	// Convert symbol to Plugin interface
	pluginInstance, ok := sym.(Plugin)
	if !ok {
		return nil, fmt.Errorf("plugin does not implement the Plugin interface")
	}

	return pluginInstance, nil
}

// UnloadPluginWithMetadata unloads a plugin and updates metadata
func (epm *EnhancedPluginManager) UnloadPluginWithMetadata(metadata *PluginMetadata) error {
	epm.reloadMutex.Lock()
	defer epm.reloadMutex.Unlock()

	plugin, exists := epm.plugins[metadata.Name]
	if !exists {
		return fmt.Errorf("plugin %s is not loaded", metadata.Name)
	}

	// Create plugin API
	api := NewPluginAPI(epm.PluginManager, metadata.Name)

	// Shutdown plugin
	if enhancedPlugin, ok := plugin.(EnhancedPlugin); ok {
		if err := enhancedPlugin.OnUnload(api); err != nil {
			log.Printf("Warning: plugin %s OnUnload failed: %v", metadata.Name, err)
		}
	} else {
		if err := plugin.Shutdown(); err != nil {
			log.Printf("Warning: plugin %s shutdown failed: %v", metadata.Name, err)
		}
	}

	// Remove plugin
	delete(epm.plugins, metadata.Name)
	delete(epm.configs, metadata.Name)

	metadata.Loaded = false
	metadata.LoadTime = time.Time{}

	log.Printf("Unloaded plugin: %s", metadata.Name)
	return nil
}

// ReloadPluginWithMetadata reloads a plugin
func (epm *EnhancedPluginManager) ReloadPluginWithMetadata(metadata *PluginMetadata) error {
	if err := epm.UnloadPluginWithMetadata(metadata); err != nil {
		return err
	}
	return epm.LoadPluginFromMetadata(metadata)
}

// EnableHotReload enables hot reloading of plugins
func (epm *EnhancedPluginManager) EnableHotReload(enable bool) {
	epm.hotReload = enable
	if enable && epm.fileWatcher != nil {
		// Add plugin directories to watcher
		for _, dir := range epm.discovery.directories {
			epm.fileWatcher.Add(dir)
		}
	}
}

// watchFiles watches for file changes and reloads plugins
func (epm *EnhancedPluginManager) watchFiles() {
	if epm.fileWatcher == nil {
		return
	}

	for {
		select {
		case event, ok := <-epm.fileWatcher.Events:
			if !ok {
				return
			}

			if !epm.hotReload {
				continue
			}

			// Check if this is a plugin file change
			if epm.shouldHandleEvent(event) {
				epm.handleFileChange(event)
			}

		case err, ok := <-epm.fileWatcher.Errors:
			if !ok {
				return
			}
			log.Printf("File watcher error: %v", err)
		}
	}
}

// shouldHandleEvent checks if we should handle a file system event
func (epm *EnhancedPluginManager) shouldHandleEvent(event fsnotify.Event) bool {
	// Only handle write events for plugin files
	if event.Op&fsnotify.Write != fsnotify.Write {
		return false
	}

	// Check if it's a plugin file
	return strings.HasSuffix(event.Name, ".so") ||
		strings.HasSuffix(event.Name, ".yaml") ||
		strings.HasSuffix(event.Name, ".json")
}

// handleFileChange handles a file change event
func (epm *EnhancedPluginManager) handleFileChange(event fsnotify.Event) {
	// Find which plugin this affects
	plugins, err := epm.discovery.DiscoverPlugins()
	if err != nil {
		log.Printf("Failed to rediscover plugins after file change: %v", err)
		return
	}

	for _, metadata := range plugins {
		if metadata.Path == event.Name ||
			strings.HasPrefix(event.Name, strings.TrimSuffix(metadata.Path, filepath.Ext(metadata.Path))) {

			// Check if auto-reload is enabled for this plugin
			if metadata.Config.AutoReload && metadata.Loaded {
				log.Printf("Auto-reloading plugin %s due to file change", metadata.Name)
				if err := epm.ReloadPluginWithMetadata(metadata); err != nil {
					log.Printf("Failed to auto-reload plugin %s: %v", metadata.Name, err)
				}
			}
			break
		}
	}
}

// GetPluginConfig gets the configuration for a plugin
func (epm *EnhancedPluginManager) GetPluginConfig(pluginName string) (*PluginConfig, error) {
	config, exists := epm.configs[pluginName]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", pluginName)
	}
	return config, nil
}

// EnablePlugin enables a plugin
func (epm *EnhancedPluginManager) EnablePlugin(name string) error {
	epm.mutex.Lock()
	defer epm.mutex.Unlock()

	_, exists := epm.plugins[name]
	if !exists {
		return fmt.Errorf("plugin not found: %s", name)
	}

	// Enable the plugin in info
	if epm.pluginInfo[name] != nil {
		epm.pluginInfo[name].Enabled = true
	}

	return nil
}

// DisablePlugin disables a plugin
func (epm *EnhancedPluginManager) DisablePlugin(name string) error {
	epm.mutex.Lock()
	defer epm.mutex.Unlock()

	_, exists := epm.plugins[name]
	if !exists {
		return fmt.Errorf("plugin not found: %s", name)
	}

	// Disable the plugin in info
	if epm.pluginInfo[name] != nil {
		epm.pluginInfo[name].Enabled = false
	}

	return nil
}

// SetPluginConfig sets the configuration for a plugin
func (epm *EnhancedPluginManager) SetPluginConfig(pluginName string, config *PluginConfig) error {
	epm.configs[pluginName] = config

	// Notify plugin of config change if it's loaded
	if plugin, exists := epm.plugins[pluginName]; exists {
		if enhancedPlugin, ok := plugin.(EnhancedPlugin); ok {
			api := NewPluginAPI(epm.PluginManager, pluginName)
			if err := enhancedPlugin.OnConfigChange(api, config.Settings); err != nil {
				log.Printf("Plugin %s config change failed: %v", pluginName, err)
			}
		}
	}

	return nil
}

// GetPluginMetadata gets metadata for all plugins
func (epm *EnhancedPluginManager) GetPluginMetadata() ([]*PluginMetadata, error) {
	return epm.discovery.DiscoverPlugins()
}

// Shutdown shuts down the enhanced plugin manager
func (epm *EnhancedPluginManager) Shutdown() error {
	// Unload all plugins
	if err := epm.UnloadAllPlugins(); err != nil {
		log.Printf("Failed to unload all plugins: %v", err)
	}

	// Close file watcher
	if epm.fileWatcher != nil {
		epm.fileWatcher.Close()
	}

	return nil
}

// ============================================================================
// Utility Functions
// ============================================================================

// sortPluginsByPriority sorts plugins by priority (lower priority = loaded first)
func sortPluginsByPriority(plugins []*PluginMetadata) {
	for i := 0; i < len(plugins)-1; i++ {
		for j := i + 1; j < len(plugins); j++ {
			if plugins[i].Config.Priority > plugins[j].Config.Priority {
				plugins[i], plugins[j] = plugins[j], plugins[i]
			}
		}
	}
}
