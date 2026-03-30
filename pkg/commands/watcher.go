package commands

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"tesselbox/pkg/entities"
)

// PluginWatcher monitors the plugins directory for changes
type PluginWatcher struct {
	pluginManager *entities.PluginManager
	pluginsDir    string
	watchedFiles  map[string]time.Time
	interval      time.Duration
	stopChan      chan bool
	running       bool
}

// NewPluginWatcher creates a new plugin watcher
func NewPluginWatcher(pm *entities.PluginManager, pluginsDir string) *PluginWatcher {
	return &PluginWatcher{
		pluginManager: pm,
		pluginsDir:    pluginsDir,
		watchedFiles:  make(map[string]time.Time),
		interval:      2 * time.Second,
		stopChan:      make(chan bool),
		running:       false,
	}
}

// Start begins watching the plugins directory
func (pw *PluginWatcher) Start() {
	if pw.running {
		return
	}

	pw.running = true
	go pw.watchLoop()

	log.Println("Plugin watcher started")
}

// Stop stops watching the plugins directory
func (pw *PluginWatcher) Stop() {
	if !pw.running {
		return
	}

	pw.running = false
	close(pw.stopChan)

	log.Println("Plugin watcher stopped")
}

// watchLoop continuously checks for plugin changes
func (pw *PluginWatcher) watchLoop() {
	ticker := time.NewTicker(pw.interval)
	defer ticker.Stop()

	// Initial scan
	pw.scanPlugins()

	for {
		select {
		case <-ticker.C:
			pw.scanPlugins()
		case <-pw.stopChan:
			return
		}
	}
}

// scanPlugins checks for new, modified, or deleted plugins
func (pw *PluginWatcher) scanPlugins() {
	// Ensure plugins directory exists
	if _, err := os.Stat(pw.pluginsDir); os.IsNotExist(err) {
		return
	}

	// Read directory
	entries, err := os.ReadDir(pw.pluginsDir)
	if err != nil {
		log.Printf("Failed to read plugins directory: %v", err)
		return
	}

	currentFiles := make(map[string]time.Time)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		// Only watch .so and .dll files
		if !isPluginFile(name) {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		path := filepath.Join(pw.pluginsDir, name)
		currentFiles[path] = info.ModTime()

		// Check if file is new or modified
		if lastMod, exists := pw.watchedFiles[path]; !exists {
			// New plugin detected
			pw.handleNewPlugin(path)
		} else if info.ModTime().After(lastMod) {
			// Modified plugin detected
			pw.handleModifiedPlugin(path)
		}
	}

	// Check for deleted plugins
	for path := range pw.watchedFiles {
		if _, exists := currentFiles[path]; !exists {
			pw.handleDeletedPlugin(path)
		}
	}

	// Update watched files
	pw.watchedFiles = currentFiles
}

// handleNewPlugin handles a newly detected plugin
func (pw *PluginWatcher) handleNewPlugin(path string) {
	log.Printf("New plugin detected: %s", path)

	if pw.pluginManager == nil {
		return
	}

	// Try to load the plugin
	if err := pw.pluginManager.LoadPlugin(path); err != nil {
		log.Printf("Auto-load failed for %s: %v", path, err)
	} else {
		log.Printf("Auto-loaded plugin: %s", path)
	}
}

// handleModifiedPlugin handles a modified plugin
func (pw *PluginWatcher) handleModifiedPlugin(path string) {
	log.Printf("Plugin modified: %s", path)

	if pw.pluginManager == nil {
		return
	}

	// Get plugin name from path
	name := filepath.Base(path)

	// Check if plugin is loaded
	plugins := pw.pluginManager.GetLoadedPlugins()
	isLoaded := false
	for _, p := range plugins {
		if p == name {
			isLoaded = true
			break
		}
	}

	if isLoaded {
		// Reload the plugin
		if err := pw.pluginManager.UnloadPlugin(name); err != nil {
			log.Printf("Failed to unload plugin %s for reload: %v", name, err)
			return
		}

		if err := pw.pluginManager.LoadPlugin(path); err != nil {
			log.Printf("Failed to reload plugin %s: %v", name, err)
		} else {
			log.Printf("Auto-reloaded plugin: %s", name)
		}
	}
}

// handleDeletedPlugin handles a deleted plugin
func (pw *PluginWatcher) handleDeletedPlugin(path string) {
	log.Printf("Plugin deleted: %s", path)

	if pw.pluginManager == nil {
		return
	}

	// Get plugin name from path
	name := filepath.Base(path)

	// Unload the plugin
	if err := pw.pluginManager.UnloadPlugin(name); err != nil {
		log.Printf("Failed to unload deleted plugin %s: %v", name, err)
	} else {
		log.Printf("Auto-unloaded deleted plugin: %s", name)
	}
}

// isPluginFile checks if a file is a plugin file
func isPluginFile(name string) bool {
	return filepath.Ext(name) == ".so" || filepath.Ext(name) == ".dll"
}
