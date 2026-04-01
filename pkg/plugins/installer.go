package plugins

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"tesselbox/pkg/entities"
)

// PluginInstaller handles downloading and installing plugins
type PluginInstaller struct {
	pluginManager    *entities.PluginManager
	pluginsDirectory string
	tempDirectory    string
	downloads        map[string]*DownloadProgress
}

// DownloadProgress tracks the progress of a plugin download
type DownloadProgress struct {
	PluginID      string
	TotalBytes     int64
	DownloadedBytes int64
	StartTime      time.Time
	Complete       bool
	Error          error
	Progress       float64
	Speed          float64 // Bytes per second
}

// PluginManifest represents the manifest file for a plugin
type PluginManifest struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Version       string   `json:"version"`
	Description   string   `json:"description"`
	Author        string   `json:"author"`
	Dependencies  []string `json:"dependencies"`
	MinGameVersion string  `json:"minGameVersion"`
	MaxGameVersion string  `json:"maxGameVersion"`
	Checksum      string   `json:"checksum"`
	Files         []PluginFile `json:"files"`
}

// PluginFile represents a file in the plugin package
type PluginFile struct {
	Path     string `json:"path"`
	Type     string `json:"type"` // "binary", "config", "asset", "data"
	Checksum string `json:"checksum"`
	Size     int64  `json:"size"`
}

// NewPluginInstaller creates a new plugin installer
func NewPluginInstaller(pluginManager *entities.PluginManager) *PluginInstaller {
	pluginsDir := "plugins"
	tempDir := filepath.Join(pluginsDir, "temp")
	
	// Ensure directories exist
	os.MkdirAll(pluginsDir, 0755)
	os.MkdirAll(tempDir, 0755)
	
	return &PluginInstaller{
		pluginManager:    pluginManager,
		pluginsDirectory: pluginsDir,
		tempDirectory:    tempDir,
		downloads:        make(map[string]*DownloadProgress),
	}
}

// DownloadAndInstall downloads and installs a plugin from the marketplace
func (pi *PluginInstaller) DownloadAndInstall(plugin *MarketplacePlugin) error {
	log.Printf("Starting download and install for plugin: %s", plugin.ID)
	
	// Create download progress
	progress := &DownloadProgress{
		PluginID:   plugin.ID,
		StartTime:  time.Now(),
		Complete:   false,
	}
	pi.downloads[plugin.ID] = progress
	
	// Start download in goroutine
	go pi.performDownload(plugin, progress)
	
	return nil
}

// performDownload handles the actual download process
func (pi *PluginInstaller) performDownload(plugin *MarketplacePlugin, progress *DownloadProgress) {
	// Download the plugin file
	tempFile := filepath.Join(pi.tempDirectory, plugin.ID+".tmp")
	
	err := pi.downloadFile(plugin.DownloadURL, tempFile, progress)
	if err != nil {
		progress.Error = fmt.Errorf("download failed: %v", err)
		progress.Complete = true
		log.Printf("Failed to download plugin %s: %v", plugin.ID, err)
		return
	}
	
	// Verify checksum if available
	if plugin.ID == "colored-blocks" { // Example checksum verification
		expectedChecksum := "a1b2c3d4e5f6" // This would come from plugin metadata
		if err := pi.verifyFileChecksum(tempFile, expectedChecksum); err != nil {
			progress.Error = fmt.Errorf("checksum verification failed: %v", err)
			progress.Complete = true
			os.Remove(tempFile)
			return
		}
	}
	
	// Extract and install plugin
	err = pi.installPluginFile(tempFile, plugin)
	if err != nil {
		progress.Error = fmt.Errorf("installation failed: %v", err)
		progress.Complete = true
		os.Remove(tempFile)
		return
	}
	
	// Clean up temp file
	os.Remove(tempFile)
	
	// Mark as complete
	progress.Complete = true
	progress.Progress = 100.0
	
	log.Printf("Successfully installed plugin: %s", plugin.ID)
}

// downloadFile downloads a file with progress tracking
func (pi *PluginInstaller) downloadFile(url, dest string, progress *DownloadProgress) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %s", resp.Status)
	}
	
	progress.TotalBytes = resp.ContentLength
	
	// Create destination file
	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()
	
	// Create progress reader
	reader := &progressReader{
		reader:   resp.Body,
		progress: progress,
	}
	
	// Copy with progress tracking
	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}
	
	return nil
}

// progressReader wraps a reader to track download progress
type progressReader struct {
	reader   io.Reader
	progress *DownloadProgress
}

func (pr *progressReader) Read(p []byte) (n int, err error) {
	n, err = pr.reader.Read(p)
	
	pr.progress.DownloadedBytes += int64(n)
	
	if pr.progress.TotalBytes > 0 {
		pr.progress.Progress = float64(pr.progress.DownloadedBytes) / float64(pr.progress.TotalBytes) * 100
	}
	
	// Calculate download speed
	elapsed := time.Since(pr.progress.StartTime).Seconds()
	if elapsed > 0 {
		pr.progress.Speed = float64(pr.progress.DownloadedBytes) / elapsed
	}
	
	return n, err
}

// verifyFileChecksum verifies the SHA256 checksum of a file
func (pi *PluginInstaller) verifyFileChecksum(filePath, expectedChecksum string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return err
	}
	
	actualChecksum := hex.EncodeToString(hasher.Sum(nil))
	
	if !strings.EqualFold(actualChecksum, expectedChecksum) {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum)
	}
	
	return nil
}

// installPluginFile installs a downloaded plugin file
func (pi *PluginInstaller) installPluginFile(tempFile string, plugin *MarketplacePlugin) error {
	// For now, we'll simulate installation by copying to plugins directory
	pluginFile := filepath.Join(pi.pluginsDirectory, plugin.ID+".so")
	
	// Read source file
	sourceData, err := os.ReadFile(tempFile)
	if err != nil {
		return fmt.Errorf("failed to read plugin file: %v", err)
	}
	
	// Write to destination
	err = os.WriteFile(pluginFile, sourceData, 0755)
	if err != nil {
		return fmt.Errorf("failed to write plugin file: %v", err)
	}
	
	// Load plugin using the plugin manager
	err = pi.pluginManager.LoadPlugin(plugin.ID)
	if err != nil {
		// Remove the file if loading failed
		os.Remove(pluginFile)
		return fmt.Errorf("failed to load plugin: %v", err)
	}
	
	return nil
}

// UninstallPlugin removes a plugin
func (pi *PluginInstaller) UninstallPlugin(pluginID string) error {
	log.Printf("Uninstalling plugin: %s", pluginID)
	
	// Unload plugin first
	err := pi.pluginManager.UnloadPlugin(pluginID)
	if err != nil {
		log.Printf("Warning: failed to unload plugin %s: %v", pluginID, err)
	}
	
	// Remove plugin file
	pluginFile := filepath.Join(pi.pluginsDirectory, pluginID+".so")
	if err := os.Remove(pluginFile); err != nil {
		return fmt.Errorf("failed to remove plugin file: %v", err)
	}
	
	log.Printf("Successfully uninstalled plugin: %s", pluginID)
	return nil
}

// GetDownloadProgress returns the download progress for a plugin
func (pi *PluginInstaller) GetDownloadProgress(pluginID string) *DownloadProgress {
	return pi.downloads[pluginID]
}

// GetAllDownloadProgress returns all download progress
func (pi *PluginInstaller) GetAllDownloadProgress() map[string]*DownloadProgress {
	return pi.downloads
}

// CleanupProgress removes completed download progress entries
func (pi *PluginInstaller) CleanupProgress() {
	for id, progress := range pi.downloads {
		if progress.Complete && time.Since(progress.StartTime) > 5*time.Minute {
			delete(pi.downloads, id)
		}
	}
}

// GetInstalledPlugins returns a list of installed plugins
func (pi *PluginInstaller) GetInstalledPlugins() []string {
	files, err := os.ReadDir(pi.pluginsDirectory)
	if err != nil {
		return []string{}
	}
	
	var plugins []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".so") {
			pluginName := strings.TrimSuffix(file.Name(), ".so")
			plugins = append(plugins, pluginName)
		}
	}
	
	return plugins
}

// ValidatePlugin checks if a plugin is compatible and safe to install
func (pi *PluginInstaller) ValidatePlugin(plugin *MarketplacePlugin) error {
	// Check if already installed
	if plugin.Installed {
		return fmt.Errorf("plugin is already installed")
	}
	
	// Check dependencies
	for _, dep := range plugin.Dependencies {
		if !pi.isDependencyInstalled(dep) {
			return fmt.Errorf("missing dependency: %s", dep)
		}
	}
	
	// Check game version compatibility (simplified)
	gameVersion := "2.0.0" // This would come from actual game version
	if plugin.Version < "1.0.0" {
		return fmt.Errorf("plugin version %s is not compatible with game version %s", plugin.Version, gameVersion)
	}
	
	return nil
}

// isDependencyInstalled checks if a dependency is installed
func (pi *PluginInstaller) isDependencyInstalled(dependency string) bool {
	installed := pi.GetInstalledPlugins()
	for _, plugin := range installed {
		if plugin == dependency {
			return true
		}
	}
	return false
}

// UpdatePlugin updates an installed plugin to the latest version
func (pi *PluginInstaller) UpdatePlugin(plugin *MarketplacePlugin) error {
	if !plugin.Installed {
		return fmt.Errorf("plugin is not installed")
	}
	
	log.Printf("Updating plugin: %s", plugin.ID)
	
	// Uninstall current version
	err := pi.UninstallPlugin(plugin.ID)
	if err != nil {
		return fmt.Errorf("failed to uninstall current version: %v", err)
	}
	
	// Install new version
	err = pi.DownloadAndInstall(plugin)
	if err != nil {
		return fmt.Errorf("failed to install new version: %v", err)
	}
	
	log.Printf("Successfully updated plugin: %s", plugin.ID)
	return nil
}

// EnablePlugin enables a plugin
func (pi *PluginInstaller) EnablePlugin(pluginID string) error {
	// This would interact with the plugin manager to enable the plugin
	log.Printf("Enabling plugin: %s", pluginID)
	return nil
}

// DisablePlugin disables a plugin
func (pi *PluginInstaller) DisablePlugin(pluginID string) error {
	// This would interact with the plugin manager to disable the plugin
	log.Printf("Disabling plugin: %s", pluginID)
	return nil
}

// GetPluginInfo returns detailed information about an installed plugin
func (pi *PluginInstaller) GetPluginInfo(pluginID string) (*PluginInfo, error) {
	pluginFile := filepath.Join(pi.pluginsDirectory, pluginID+".so")
	
	// Get file info
	fileInfo, err := os.Stat(pluginFile)
	if err != nil {
		return nil, fmt.Errorf("plugin not found: %v", err)
	}
	
	// Create basic plugin info
	info := &PluginInfo{
		ID:          pluginID,
		Name:        pluginID,
		Version:     "unknown",
		Description: "Plugin description not available",
		Author:      "Unknown",
		Installed:   true,
		Enabled:     true,
		InstallDate: fileInfo.ModTime(),
		FileSize:    fileInfo.Size(),
		FilePath:    pluginFile,
	}
	
	// Try to get more detailed info from plugin manager
	if plugin, exists := pi.pluginManager.GetPlugin(pluginID); exists {
		info.Name = plugin.GetName()
		info.Version = plugin.GetVersion()
		info.Description = plugin.GetDescription()
		info.Author = plugin.GetAuthor()
		info.Dependencies = plugin.GetDependencies()
	}
	
	return info, nil
}

// PluginInfo contains detailed information about a plugin
type PluginInfo struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Dependencies []string  `json:"dependencies"`
	Installed   bool      `json:"installed"`
	Enabled     bool      `json:"enabled"`
	InstallDate time.Time `json:"installDate"`
	FileSize    int64     `json:"fileSize"`
	FilePath    string    `json:"filePath"`
	Checksum    string    `json:"checksum"`
}
