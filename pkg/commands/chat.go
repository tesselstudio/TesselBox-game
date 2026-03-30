package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"tesselbox/pkg/entities"
)

// ChatHandler processes chat messages and commands
type ChatHandler struct {
	pluginManager *entities.PluginManager
	showChat      bool
	chatHistory   []string
	currentInput  string
	maxHistory    int
}

// NewChatHandler creates a new chat handler
func NewChatHandler(pm *entities.PluginManager) *ChatHandler {
	return &ChatHandler{
		pluginManager: pm,
		showChat:      false,
		chatHistory:   make([]string, 0),
		currentInput:  "",
		maxHistory:    100,
	}
}

// ToggleChat toggles the chat display
func (ch *ChatHandler) ToggleChat() {
	ch.showChat = !ch.showChat
}

// IsChatOpen returns whether the chat is currently open
func (ch *ChatHandler) IsChatOpen() bool {
	return ch.showChat
}

// AddMessage adds a message to the chat history
func (ch *ChatHandler) AddMessage(msg string) {
	ch.chatHistory = append(ch.chatHistory, msg)
	if len(ch.chatHistory) > ch.maxHistory {
		ch.chatHistory = ch.chatHistory[1:]
	}
}

// GetHistory returns the chat history
func (ch *ChatHandler) GetHistory() []string {
	return ch.chatHistory
}

// ProcessInput processes a chat message or command
func (ch *ChatHandler) ProcessInput(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}

	// Add to history
	ch.AddMessage("> " + input)

	// Check if it's a command
	if strings.HasPrefix(input, "/") {
		return ch.handleCommand(input)
	}

	// Regular chat message
	return fmt.Sprintf("[Chat] %s", input)
}

// handleCommand processes slash commands
func (ch *ChatHandler) handleCommand(cmd string) string {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return "Unknown command"
	}

	command := strings.ToLower(parts[0])
	args := parts[1:]

	switch command {
	case "/plugin":
		return ch.handlePluginCommand(args)
	case "/help":
		return ch.handleHelpCommand()
	default:
		return fmt.Sprintf("Unknown command: %s", command)
	}
}

// handlePluginCommand handles /plugin subcommands
func (ch *ChatHandler) handlePluginCommand(args []string) string {
	if len(args) == 0 {
		return "Usage: /plugin <load|unload|list|enable|disable> [name]"
	}

	subcommand := strings.ToLower(args[0])
	pluginName := ""
	if len(args) > 1 {
		pluginName = args[1]
	}

	switch subcommand {
	case "load":
		if pluginName == "" {
			return "Usage: /plugin load <filename>"
		}
		return ch.loadPlugin(pluginName)

	case "unload":
		if pluginName == "" {
			return "Usage: /plugin unload <name>"
		}
		return ch.unloadPlugin(pluginName)

	case "list":
		return ch.listPlugins()

	case "enable":
		if pluginName == "" {
			return "Usage: /plugin enable <name>"
		}
		return ch.enablePlugin(pluginName)

	case "disable":
		if pluginName == "" {
			return "Usage: /plugin disable <name>"
		}
		return ch.disablePlugin(pluginName)

	case "reload":
		if pluginName == "" {
			return "Usage: /plugin reload <name>"
		}
		return ch.reloadPlugin(pluginName)

	default:
		return fmt.Sprintf("Unknown plugin subcommand: %s", subcommand)
	}
}

// loadPlugin loads a plugin from file
func (ch *ChatHandler) loadPlugin(filename string) string {
	// Ensure .so extension
	if !strings.HasSuffix(filename, ".so") && !strings.HasSuffix(filename, ".dll") {
		filename += ".so" // Default to .so for Linux
	}

	pluginPath := filepath.Join("plugins", filename)

	// Check if file exists
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		return fmt.Sprintf("Plugin file not found: %s", pluginPath)
	}

	// Load the plugin
	if ch.pluginManager == nil {
		return "Plugin manager not initialized"
	}

	err := ch.pluginManager.LoadPlugin(pluginPath)
	if err != nil {
		log.Printf("Failed to load plugin %s: %v", filename, err)
		return fmt.Sprintf("Failed to load plugin: %v", err)
	}

	log.Printf("Plugin loaded: %s", filename)
	return fmt.Sprintf("Plugin loaded: %s", filename)
}

// unloadPlugin unloads a plugin
func (ch *ChatHandler) unloadPlugin(name string) string {
	if ch.pluginManager == nil {
		return "Plugin manager not initialized"
	}

	err := ch.pluginManager.UnloadPlugin(name)
	if err != nil {
		return fmt.Sprintf("Failed to unload plugin: %v", err)
	}

	return fmt.Sprintf("Plugin unloaded: %s", name)
}

// listPlugins lists all loaded plugins
func (ch *ChatHandler) listPlugins() string {
	if ch.pluginManager == nil {
		return "Plugin manager not initialized"
	}

	plugins := ch.pluginManager.GetLoadedPlugins()
	if len(plugins) == 0 {
		return "No plugins loaded"
	}

	var sb strings.Builder
	sb.WriteString("Loaded plugins:\n")

	// Sort plugins by name for consistent output
	sort.Strings(plugins)

	for _, name := range plugins {
		info := ch.pluginManager.GetPluginInfo(name)
		status := "enabled"
		if !info.Enabled {
			status = "disabled"
		}
		sb.WriteString(fmt.Sprintf("  - %s [%s] v%s\n", name, status, info.Version))
	}

	return strings.TrimSuffix(sb.String(), "\n")
}

// enablePlugin enables a plugin
func (ch *ChatHandler) enablePlugin(name string) string {
	if ch.pluginManager == nil {
		return "Plugin manager not initialized"
	}

	if err := ch.pluginManager.EnablePlugin(name); err != nil {
		return fmt.Sprintf("Failed to enable plugin: %v", err)
	}

	return fmt.Sprintf("Plugin enabled: %s", name)
}

// disablePlugin disables a plugin
func (ch *ChatHandler) disablePlugin(name string) string {
	if ch.pluginManager == nil {
		return "Plugin manager not initialized"
	}

	if err := ch.pluginManager.DisablePlugin(name); err != nil {
		return fmt.Sprintf("Failed to disable plugin: %v", err)
	}

	return fmt.Sprintf("Plugin disabled: %s", name)
}

// reloadPlugin reloads a plugin
func (ch *ChatHandler) reloadPlugin(name string) string {
	if ch.pluginManager == nil {
		return "Plugin manager not initialized"
	}

	// Get plugin path first
	info := ch.pluginManager.GetPluginInfo(name)
	if info.Path == "" {
		return fmt.Sprintf("Plugin not found: %s", name)
	}

	// Unload
	if err := ch.pluginManager.UnloadPlugin(name); err != nil {
		return fmt.Sprintf("Failed to unload plugin: %v", err)
	}

	// Reload
	if err := ch.pluginManager.LoadPlugin(info.Path); err != nil {
		return fmt.Sprintf("Failed to reload plugin: %v", err)
	}

	return fmt.Sprintf("Plugin reloaded: %s", name)
}

// handleHelpCommand shows available commands
func (ch *ChatHandler) handleHelpCommand() string {
	return `Available commands:
  /plugin load <file>   - Load a plugin from file
  /plugin unload <name> - Unload a plugin
  /plugin list          - List loaded plugins
  /plugin enable <name> - Enable a plugin
  /plugin disable <name>- Disable a plugin
  /plugin reload <name> - Reload a plugin
  /help                 - Show this help message`
}
