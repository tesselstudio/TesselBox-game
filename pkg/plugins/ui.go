package plugins

import (
	"fmt"
	"image/color"
	"log"
	"strings"
	"time"

	"tesselbox/pkg/entities"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// PluginUI represents the plugin manager interface
type PluginUI struct {
	pluginManager *entities.PluginManager
	
	// UI state
	currentView   PluginView
	selectedIndex int
	scrollOffset  int
	searchQuery   string
	searchActive  bool
	
	// Marketplace data
	availablePlugins []*MarketplacePlugin
	categories       []string
	selectedCategory string
	
	// Installation state
	installingPlugins map[string]bool
	installProgress   map[string]float64
	
	// Visual properties
	backgroundColor color.RGBA
	accentColor     color.RGBA
	textColor       color.RGBA
	successColor    color.RGBA
	errorColor      color.RGBA
	
	// Animations
	animTimer       float64
	fadeAlpha       float64
	
	// For solid color drawing
	whiteImage *ebiten.Image
}

// PluginView represents different views in the plugin manager
type PluginView int

const (
	ViewInstalled PluginView = iota
	ViewMarketplace
	ViewCategories
	ViewSearch
	ViewDetails
)

// MarketplacePlugin represents a plugin available in the marketplace
type MarketplacePlugin struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	Category    string            `json:"category"`
	Tags        []string          `json:"tags"`
	Downloads   int               `json:"downloads"`
	Rating      float64           `json:"rating"`
	Dependencies []string         `json:"dependencies"`
	DownloadURL string            `json:"downloadUrl"`
	Installed   bool              `json:"installed"`
	Enabled     bool              `json:"enabled"`
	LastUpdate  time.Time         `json:"lastUpdate"`
	Screenshots []string          `json:"screenshots"`
	Features    []string          `json:"features"`
	Compatibility map[string]string `json:"compatibility"`
}

// NewPluginUI creates a new plugin manager UI
func NewPluginUI(pluginManager *entities.PluginManager) *PluginUI {
	// Create a 1x1 white image for solid color drawing
	whiteImage := ebiten.NewImage(1, 1)
	whiteImage.Fill(color.RGBA{255, 255, 255, 255})

	ui := &PluginUI{
		pluginManager:     pluginManager,
		currentView:       ViewMarketplace,
		selectedIndex:     0,
		scrollOffset:      0,
		searchQuery:       "",
		searchActive:      false,
		categories:        []string{"All", "Blocks", "Items", "Creatures", "UI", "Tools", "Gameplay", "Visual"},
		selectedCategory:  "All",
		installingPlugins: make(map[string]bool),
		installProgress:   make(map[string]float64),
		backgroundColor:   color.RGBA{15, 20, 35, 255},
		accentColor:       color.RGBA{120, 180, 255, 255},
		textColor:         color.RGBA{255, 255, 255, 255},
		successColor:      color.RGBA{100, 255, 100, 255},
		errorColor:        color.RGBA{255, 100, 100, 255},
		animTimer:         0.0,
		fadeAlpha:         1.0,
		whiteImage:        whiteImage,
	}

	// Initialize marketplace data
	ui.initializeMarketplace()
	
	return ui
}

// initializeMarketplace sets up the marketplace with sample plugins
func (ui *PluginUI) initializeMarketplace() {
	ui.availablePlugins = []*MarketplacePlugin{
		{
			ID:          "magic-system",
			Name:        "Magic System",
			Version:     "1.2.0",
			Description: "Adds magic spells, wands, and enchanting system to the game",
			Author:      "TesselBox Team",
			Category:    "Gameplay",
			Tags:        []string{"magic", "spells", "enchanting"},
			Downloads:   15420,
			Rating:      4.8,
			Dependencies: []string{},
			DownloadURL: "https://plugins.tesselbox.com/magic-system.so",
			Installed:   false,
			Enabled:     false,
			LastUpdate:  time.Now().AddDate(0, -1, -15),
			Features:    []string{"Spell casting", "Mana system", "Enchanting table", "Magic items"},
		},
		{
			ID:          "colored-blocks",
			Name:        "Colored Blocks Pack",
			Version:     "2.1.0",
			Description: "Adds blocks with multiple color variations and patterns",
			Author:      "Artist Collective",
			Category:    "Blocks",
			Tags:        []string{"blocks", "colors", "decoration"},
			Downloads:   32100,
			Rating:      4.9,
			Dependencies: []string{},
			DownloadURL: "https://plugins.tesselbox.com/colored-blocks.so",
			Installed:   false,
			Enabled:     false,
			LastUpdate:  time.Now().AddDate(0, -0, -5),
			Features:    []string{"Multi-color blocks", "Pattern blocks", "Gradient blocks", "Custom textures"},
		},
		{
			ID:          "advanced-crafting",
			Name:        "Advanced Crafting",
			Version:     "1.5.2",
			Description: "Enhanced crafting system with new recipes and workstations",
			Author:      "Crafting Masters",
			Category:    "Gameplay",
			Tags:        []string{"crafting", "recipes", "workstations"},
			Downloads:   28900,
			Rating:      4.7,
			Dependencies: []string{},
			DownloadURL: "https://plugins.tesselbox.com/advanced-crafting.so",
			Installed:   false,
			Enabled:     false,
			LastUpdate:  time.Now().AddDate(0, -2, -10),
			Features:    []string{"New recipes", "Advanced workstations", "Quality system", "Tool upgrades"},
		},
		{
			ID:          "creature-mobs",
			Name:        "Creature Mobs",
			Version:     "1.0.0",
			Description: "Adds various creatures and mobs to the world",
			Author:      "Wildlife Team",
			Category:    "Creatures",
			Tags:        []string{"creatures", "mobs", "animals"},
			Downloads:   19800,
			Rating:      4.6,
			Dependencies: []string{},
			DownloadURL: "https://plugins.tesselbox.com/creature-mobs.so",
			Installed:   false,
			Enabled:     false,
			LastUpdate:  time.Now().AddDate(0, -1, -20),
			Features:    []string{"Friendly animals", "Hostile mobs", "Boss creatures", "AI behaviors"},
		},
		{
			ID:          "ui-enhancement",
			Name:        "UI Enhancement Pack",
			Version:     "1.3.1",
			Description: "Improves the user interface with better layouts and visual effects",
			Author:      "UI Designers",
			Category:    "UI",
			Tags:        []string{"ui", "interface", "visual"},
			Downloads:   41200,
			Rating:      4.9,
			Dependencies: []string{},
			DownloadURL: "https://plugins.tesselbox.com/ui-enhancement.so",
			Installed:   false,
			Enabled:     false,
			LastUpdate:  time.Now().AddDate(0, -0, -15),
			Features:    []string{"Better inventory", "Enhanced tooltips", "Animations", "Themes"},
		},
		{
			ID:          "tech-modpack",
			Name:        "Technology Modpack",
			Version:     "2.0.0",
			Description: "Adds technology, machines, and automation to the game",
			Author:      "Tech Engineers",
			Category:    "Tools",
			Tags:        []string{"technology", "machines", "automation"},
			Downloads:   35600,
			Rating:      4.8,
			Dependencies: []string{"advanced-crafting"},
			DownloadURL: "https://plugins.tesselbox.com/tech-modpack.so",
			Installed:   false,
			Enabled:     false,
			LastUpdate:  time.Now().AddDate(0, -0, -8),
			Features:    []string{"Machines", "Automation", "Power systems", "Tech items"},
		},
	}
}

// Update handles UI updates and input
func (ui *PluginUI) Update() error {
	// Update animations
	ui.animTimer += 0.016

	// Handle keyboard input
	ui.handleKeyboardInput()

	// Handle mouse input
	ui.handleMouseInput()

	// Update installation progress
	ui.updateInstallProgress()

	return nil
}

// handleKeyboardInput processes keyboard input
func (ui *PluginUI) handleKeyboardInput() {
	// ESC to go back
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if ui.currentView == ViewDetails {
			ui.currentView = ViewMarketplace
		} else if ui.currentView != ViewMarketplace {
			ui.currentView = ViewMarketplace
		}
	}

	// Tab to switch views
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		ui.switchView()
	}

	// Arrow keys for navigation
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		ui.selectedIndex--
		if ui.selectedIndex < 0 {
			ui.selectedIndex = len(ui.getVisiblePlugins()) - 1
		}
		ui.updateScrollOffset()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		ui.selectedIndex++
		maxIndex := len(ui.getVisiblePlugins()) - 1
		if ui.selectedIndex > maxIndex {
			ui.selectedIndex = 0
		}
		ui.updateScrollOffset()
	}

	// Enter to select/install
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		ui.handleSelection()
	}

	// Search input
	if ui.searchActive {
		ui.handleSearchInput()
	}
}

// handleMouseInput processes mouse input
func (ui *PluginUI) handleMouseInput() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		ui.handleMouseClick(mx, my)
	}

	// Mouse wheel for scrolling
	_, scrollY := ebiten.Wheel()
	if scrollY != 0 {
		ui.handleScrolling(scrollY)
	}
}

// handleSearchInput processes search input
func (ui *PluginUI) handleSearchInput() {
	// Handle typing for search
	for key := ebiten.KeyA; key <= ebiten.KeyZ; key++ {
		if inpututil.IsKeyJustPressed(key) {
			char := string(rune('a' + int(key-ebiten.KeyA)))
			ui.searchQuery += char
		}
	}

	// Space
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		ui.searchQuery += " "
	}

	// Backspace
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		if len(ui.searchQuery) > 0 {
			ui.searchQuery = ui.searchQuery[:len(ui.searchQuery)-1]
		}
	}
}

// handleSelection handles the current selection
func (ui *PluginUI) handleSelection() {
	plugins := ui.getVisiblePlugins()
	if ui.selectedIndex >= 0 && ui.selectedIndex < len(plugins) {
		plugin := plugins[ui.selectedIndex]
		
		switch ui.currentView {
		case ViewMarketplace, ViewSearch:
			ui.currentView = ViewDetails
		case ViewInstalled:
			// Toggle enable/disable
			ui.togglePlugin(plugin.ID)
		case ViewDetails:
			// Install/uninstall
			if plugin.Installed {
				ui.uninstallPlugin(plugin.ID)
			} else {
				ui.installPlugin(plugin.ID)
			}
		}
	}
}

// switchView cycles through available views
func (ui *PluginUI) switchView() {
	switch ui.currentView {
	case ViewMarketplace:
		ui.currentView = ViewInstalled
	case ViewInstalled:
		ui.currentView = ViewCategories
	case ViewCategories:
		ui.currentView = ViewSearch
	case ViewSearch:
		ui.currentView = ViewMarketplace
	case ViewDetails:
		ui.currentView = ViewMarketplace
	}
}

// getVisiblePlugins returns plugins filtered by current view
func (ui *PluginUI) getVisiblePlugins() []*MarketplacePlugin {
	switch ui.currentView {
	case ViewMarketplace:
		return ui.getPluginsByCategory(ui.selectedCategory)
	case ViewInstalled:
		return ui.getInstalledPlugins()
	case ViewSearch:
		return ui.searchPlugins(ui.searchQuery)
	case ViewDetails:
		return ui.getVisiblePlugins() // Return current selection
	default:
		return ui.availablePlugins
	}
}

// getPluginsByCategory filters plugins by category
func (ui *PluginUI) getPluginsByCategory(category string) []*MarketplacePlugin {
	if category == "All" {
		return ui.availablePlugins
	}

	var filtered []*MarketplacePlugin
	for _, plugin := range ui.availablePlugins {
		if plugin.Category == category {
			filtered = append(filtered, plugin)
		}
	}
	return filtered
}

// getInstalledPlugins returns installed plugins
func (ui *PluginUI) getInstalledPlugins() []*MarketplacePlugin {
	var installed []*MarketplacePlugin
	for _, plugin := range ui.availablePlugins {
		if plugin.Installed {
			installed = append(installed, plugin)
		}
	}
	return installed
}

// searchPlugins searches plugins by query
func (ui *PluginUI) searchPlugins(query string) []*MarketplacePlugin {
	if query == "" {
		return ui.availablePlugins
	}

	query = strings.ToLower(query)
	var results []*MarketplacePlugin
	for _, plugin := range ui.availablePlugins {
		if strings.Contains(strings.ToLower(plugin.Name), query) ||
			strings.Contains(strings.ToLower(plugin.Description), query) ||
			strings.Contains(strings.ToLower(plugin.Author), query) {
			results = append(results, plugin)
		}
	}
	return results
}

// installPlugin starts plugin installation
func (ui *PluginUI) installPlugin(pluginID string) {
	ui.installingPlugins[pluginID] = true
	ui.installProgress[pluginID] = 0.0
	
	// Simulate installation process
	go func() {
		for i := 0; i <= 100; i += 5 {
			// TODO: Add mutex protection for thread safety
			ui.installProgress[pluginID] = float64(i)
			time.Sleep(100 * time.Millisecond)
		}
		
		// Mark as installed (thread safety issue - needs mutex)
		for _, plugin := range ui.availablePlugins {
			if plugin.ID == pluginID {
				plugin.Installed = true
				plugin.Enabled = true
				break
			}
		}
		
		delete(ui.installingPlugins, pluginID)
		delete(ui.installProgress, pluginID)
		
		log.Printf("Plugin %s installed successfully", pluginID)
	}()
}

// uninstallPlugin uninstalls a plugin
func (ui *PluginUI) uninstallPlugin(pluginID string) {
	for _, plugin := range ui.availablePlugins {
		if plugin.ID == pluginID {
			plugin.Installed = false
			plugin.Enabled = false
			break
		}
	}
	
	log.Printf("Plugin %s uninstalled", pluginID)
}

// togglePlugin enables/disables a plugin
func (ui *PluginUI) togglePlugin(pluginID string) {
	for _, plugin := range ui.availablePlugins {
		if plugin.ID == pluginID {
			plugin.Enabled = !plugin.Enabled
			log.Printf("Plugin %s %s", pluginID, map[bool]string{true: "enabled", false: "disabled"}[plugin.Enabled])
			break
		}
	}
}

// updateScrollOffset adjusts scroll offset based on selection
func (ui *PluginUI) updateScrollOffset() {
	maxVisible := 8
	plugins := ui.getVisiblePlugins()
	
	if ui.selectedIndex < ui.scrollOffset {
		ui.scrollOffset = ui.selectedIndex
	} else if ui.selectedIndex >= ui.scrollOffset+maxVisible {
		ui.scrollOffset = ui.selectedIndex - maxVisible + 1
	}
	
	// Ensure scroll offset is valid
	if ui.scrollOffset < 0 {
		ui.scrollOffset = 0
	} else if ui.scrollOffset > len(plugins)-maxVisible {
		ui.scrollOffset = max(0, len(plugins)-maxVisible)
	}
}

// handleScrolling handles mouse wheel scrolling
func (ui *PluginUI) handleScrolling(scrollY float64) {
	maxVisible := 8
	plugins := ui.getVisiblePlugins()
	
	if scrollY > 0 {
		// Scroll up
		ui.scrollOffset--
		if ui.scrollOffset < 0 {
			ui.scrollOffset = 0
		}
	} else if scrollY < 0 {
		// Scroll down
		maxScroll := max(0, len(plugins)-maxVisible)
		if ui.scrollOffset < maxScroll {
			ui.scrollOffset++
		}
	}
	
	// Adjust selection if needed
	if ui.selectedIndex < ui.scrollOffset {
		ui.selectedIndex = ui.scrollOffset
	} else if ui.selectedIndex >= ui.scrollOffset+maxVisible {
		ui.selectedIndex = ui.scrollOffset + maxVisible - 1
	}
}

// handleMouseClick handles mouse clicks on UI elements
func (ui *PluginUI) handleMouseClick(mx, my int) {
	// Implementation for mouse click handling
	// This would check if mouse is over specific UI elements
	// and take appropriate actions
}

// updateInstallProgress updates installation progress
func (ui *PluginUI) updateInstallProgress() {
	// Progress is updated in the installation goroutines
}

// Draw renders the plugin UI
func (ui *PluginUI) Draw(screen *ebiten.Image) {
	// Draw background
	screen.Fill(ui.backgroundColor)

	// Draw header
	ui.drawHeader(screen)

	// Draw current view
	switch ui.currentView {
	case ViewMarketplace:
		ui.drawMarketplace(screen)
	case ViewInstalled:
		ui.drawInstalled(screen)
	case ViewCategories:
		ui.drawCategories(screen)
	case ViewSearch:
		ui.drawSearch(screen)
	case ViewDetails:
		ui.drawDetails(screen)
	}

	// Draw footer
	ui.drawFooter(screen)
}

// drawHeader draws the UI header
func (ui *PluginUI) drawHeader(screen *ebiten.Image) {
	// Draw title
	title := "PLUGIN MANAGER"
	ebitenutil.DebugPrintAt(screen, title, 10, 10)

	// Draw view tabs
	tabs := []string{"Marketplace", "Installed", "Categories", "Search"}
	tabX := 300
	for i, tab := range tabs {
		if PluginView(i) == ui.currentView {
			ebitenutil.DebugPrintAt(screen, tab, tabX+i*100, 10)
		} else {
			ebitenutil.DebugPrintAt(screen, tab, tabX+i*100, 10)
		}
	}
}

// drawMarketplace draws the marketplace view
func (ui *PluginUI) drawMarketplace(screen *ebiten.Image) {
	plugins := ui.getVisiblePlugins()
	ui.drawPluginList(screen, plugins)
}

// drawInstalled draws the installed plugins view
func (ui *PluginUI) drawInstalled(screen *ebiten.Image) {
	plugins := ui.getInstalledPlugins()
	ui.drawPluginList(screen, plugins)
}

// drawCategories draws the categories view
func (ui *PluginUI) drawCategories(screen *ebiten.Image) {
	y := 80
	for i, category := range ui.categories {
		if category == ui.selectedCategory {
			// Highlight selected category
			ebitenutil.DebugPrintAt(screen, "> "+category, 25, y)
		} else {
			ebitenutil.DebugPrintAt(screen, "  "+category, 25, y)
		}
		
		if i == ui.selectedIndex {
			// Draw selection indicator
			ebitenutil.DrawRect(screen, 20, float64(y-5), 200, 25, ui.accentColor)
		}
		
		y += 30
	}
}

// drawSearch draws the search view
func (ui *PluginUI) drawSearch(screen *ebiten.Image) {
	// Draw search box
	ebitenutil.DrawRect(screen, 20, 80, 400, 30, color.RGBA{40, 40, 50, 255})
	ebitenutil.DebugPrintAt(screen, "Search: "+ui.searchQuery, 25, 85)
	
	// Draw search results
	plugins := ui.searchPlugins(ui.searchQuery)
	ui.drawPluginList(screen, plugins)
}

// drawDetails draws the plugin details view
func (ui *PluginUI) drawDetails(screen *ebiten.Image) {
	plugins := ui.getVisiblePlugins()
	if ui.selectedIndex >= 0 && ui.selectedIndex < len(plugins) {
		plugin := plugins[ui.selectedIndex]
		ui.drawPluginDetails(screen, plugin)
	}
}

// drawPluginList draws a list of plugins
func (ui *PluginUI) drawPluginList(screen *ebiten.Image, plugins []*MarketplacePlugin) {
	y := 80
	maxVisible := 8
	
	visibleStart := ui.scrollOffset
	visibleEnd := min(ui.scrollOffset+maxVisible, len(plugins))
	
	for i := visibleStart; i < visibleEnd; i++ {
		plugin := plugins[i]
		yPos := y + (i-visibleStart)*60
		
		// Draw selection highlight
		if i == ui.selectedIndex {
			ebitenutil.DrawRect(screen, 20, float64(yPos-5), 760, 50, color.RGBA{80, 120, 180, 100})
		}
		
		// Draw plugin info
		ebitenutil.DebugPrintAt(screen, plugin.Name, 30, yPos)
		ebitenutil.DebugPrintAt(screen, plugin.Description, 30, yPos+20)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("v%s by %s", plugin.Version, plugin.Author), 500, yPos)
		
		// Draw status
		statusText := "Not Installed"
		if plugin.Installed {
			if plugin.Enabled {
				statusText = "Enabled"
			} else {
				statusText = "Disabled"
			}
		}
		ebitenutil.DebugPrintAt(screen, statusText, 650, yPos)
		
		// Draw installation progress
		if ui.installingPlugins[plugin.ID] {
			progress := ui.installProgress[plugin.ID]
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Installing: %.0f%%", progress), 650, yPos+20)
		}
	}
	
	// Draw scroll indicators
	if len(plugins) > maxVisible {
		if ui.scrollOffset > 0 {
			ebitenutil.DebugPrintAt(screen, "▲", 390, 60)
		}
		if ui.scrollOffset < len(plugins)-maxVisible {
			ebitenutil.DebugPrintAt(screen, "▼", 390, 560)
		}
	}
}

// drawPluginDetails draws detailed information about a plugin
func (ui *PluginUI) drawPluginDetails(screen *ebiten.Image, plugin *MarketplacePlugin) {
	y := 80
	
	// Plugin name and version
	ebitenutil.DebugPrintAt(screen, plugin.Name, 30, y)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Version %s", plugin.Version), 30, y+20)
	
	// Author and rating
	y += 50
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("By %s", plugin.Author), 30, y)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Rating: %.1f/5.0 (%d downloads)", plugin.Rating, plugin.Downloads), 30, y+20)
	
	// Description
	y += 50
	ebitenutil.DebugPrintAt(screen, "Description:", 30, y)
	y += 20
	lines := ui.wrapText(plugin.Description, 70)
	for _, line := range lines {
		ebitenutil.DebugPrintAt(screen, line, 30, y)
		y += 20
	}
	
	// Features
	y += 20
	ebitenutil.DebugPrintAt(screen, "Features:", 30, y)
	y += 20
	for _, feature := range plugin.Features {
		ebitenutil.DebugPrintAt(screen, "• "+feature, 50, y)
		y += 20
	}
	
	// Dependencies
	if len(plugin.Dependencies) > 0 {
		y += 20
		ebitenutil.DebugPrintAt(screen, "Dependencies:", 30, y)
		y += 20
		for _, dep := range plugin.Dependencies {
			ebitenutil.DebugPrintAt(screen, "• "+dep, 50, y)
			y += 20
		}
	}
	
	// Install button
	y += 30
	buttonText := "INSTALL"
	if plugin.Installed {
		buttonText = "UNINSTALL"
	}
	buttonColor := ui.successColor
	if plugin.Installed {
		buttonColor = ui.errorColor
	}
	ebitenutil.DrawRect(screen, 30, float64(y), 150, 40, buttonColor)
	ebitenutil.DebugPrintAt(screen, buttonText, 70, y+15)
}

// drawFooter draws the UI footer
func (ui *PluginUI) drawFooter(screen *ebiten.Image) {
	// Draw help text
	helpText := "↑↓ Navigate | Enter Select | Tab Switch Views | ESC Back | Search: Type"
	ebitenutil.DebugPrintAt(screen, helpText, 10, 680)
}

// wrapText wraps text to fit within a maximum width
func (ui *PluginUI) wrapText(text string, maxWidth int) []string {
	words := strings.Split(text, " ")
	lines := []string{}
	currentLine := ""
	
	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word
		
		if len(testLine) > maxWidth {
			if currentLine != "" {
				lines = append(lines, currentLine)
				currentLine = word
			} else {
				lines = append(lines, word)
			}
		} else {
			currentLine = testLine
		}
	}
	
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	
	return lines
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
