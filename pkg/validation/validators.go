package validation

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// Validator provides input validation functions
type Validator struct{}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateWorldName validates a world name
func (v *Validator) ValidateWorldName(name string) error {
	// Check empty
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("world name cannot be empty")
	}

	// Max length
	if len(name) > 32 {
		return fmt.Errorf("world name too long (max 32 characters, got %d)", len(name))
	}

	// Min length
	if len(name) < 1 {
		return fmt.Errorf("world name cannot be empty")
	}

	// Check for reserved names (Windows reserved)
	reserved := []string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}
	upperName := strings.ToUpper(name)
	for _, reservedName := range reserved {
		if upperName == reservedName {
			return fmt.Errorf("'%s' is a reserved system name", name)
		}
	}

	// Check for path traversal
	if strings.Contains(name, "..") || strings.Contains(name, "/") || strings.Contains(name, "\\") {
		return fmt.Errorf("world name cannot contain path separators")
	}

	// Check for invalid characters
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9_\-\s]+$`)
	if !validPattern.MatchString(name) {
		return fmt.Errorf("world name contains invalid characters (use only letters, numbers, spaces, hyphens, underscores)")
	}

	// Check for leading/trailing spaces
	if name != strings.TrimSpace(name) {
		return fmt.Errorf("world name cannot have leading or trailing spaces")
	}

	return nil
}

// SanitizeWorldName cleans and validates a world name
func (v *Validator) SanitizeWorldName(name string) (string, error) {
	// Trim spaces
	name = strings.TrimSpace(name)

	// Replace multiple spaces with single space
	name = regexp.MustCompile(`\s+`).ReplaceAllString(name, " ")

	// Validate
	if err := v.ValidateWorldName(name); err != nil {
		return "", err
	}

	return name, nil
}

// ValidatePlayerName validates a player name
func (v *Validator) ValidatePlayerName(name string) error {
	// Check empty
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("player name cannot be empty")
	}

	// Max length
	if len(name) > 16 {
		return fmt.Errorf("player name too long (max 16 characters, got %d)", len(name))
	}

	// Check for valid characters (alphanumeric + underscore only)
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !validPattern.MatchString(name) {
		return fmt.Errorf("player name contains invalid characters (use only letters, numbers, underscores)")
	}

	// Check for reserved words (common offensive/banned terms could be added here)
	reserved := []string{"admin", "moderator", "system", "server", "null", "undefined"}
	lowerName := strings.ToLower(name)
	for _, word := range reserved {
		if strings.Contains(lowerName, word) {
			return fmt.Errorf("player name contains reserved word")
		}
	}

	return nil
}

// SanitizePlayerName cleans a player name
func (v *Validator) SanitizePlayerName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if err := v.ValidatePlayerName(name); err != nil {
		return "", err
	}
	return name, nil
}

// ValidateSeed validates a world seed
func (v *Validator) ValidateSeed(seed int64) error {
	// int64 is always valid range, but we could add sanity checks
	// e.g., seed == 0 might indicate an error in some contexts
	return nil
}

// ValidateSeedString validates and converts a seed string to int64
func (v *Validator) ValidateSeedString(seedStr string) (int64, error) {
	if strings.TrimSpace(seedStr) == "" {
		return 0, fmt.Errorf("seed cannot be empty")
	}

	// Check for valid characters
	validPattern := regexp.MustCompile(`^-?\d+$`)
	if !validPattern.MatchString(seedStr) {
		return 0, fmt.Errorf("seed must be a valid number")
	}

	// Try to parse
	var seed int64
	_, err := fmt.Sscanf(seedStr, "%d", &seed)
	if err != nil {
		return 0, fmt.Errorf("invalid seed number: %w", err)
	}

	return seed, nil
}

// SanitizePath ensures a path stays within the allowed directory
func (v *Validator) SanitizePath(baseDir, userPath string) (string, error) {
	// Clean the path
	cleanPath := filepath.Clean(userPath)

	// Ensure it's not absolute
	if filepath.IsAbs(cleanPath) {
		return "", fmt.Errorf("absolute paths are not allowed")
	}

	// Join with base directory
	fullPath := filepath.Join(baseDir, cleanPath)

	// Verify the resolved path is still within baseDir
	resolvedPath, err := filepath.EvalSymlinks(fullPath)
	if err != nil {
		// Path doesn't exist, check if it's within bounds anyway
		resolvedPath = fullPath
	}

	// Check if resolved path is within baseDir
	absBase, _ := filepath.Abs(baseDir)
	absResolved, _ := filepath.Abs(resolvedPath)

	if !strings.HasPrefix(absResolved, absBase) {
		return "", fmt.Errorf("path escapes base directory")
	}

	return fullPath, nil
}

// ValidateFileName validates a filename
func (v *Validator) ValidateFileName(filename string) error {
	// Check empty
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	// Max length
	if len(filename) > 255 {
		return fmt.Errorf("filename too long")
	}

	// Check for invalid characters
	invalidChars := []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*"}
	for _, char := range invalidChars {
		if strings.Contains(filename, char) {
			return fmt.Errorf("filename contains invalid character: %s", char)
		}
	}

	return nil
}

// ValidateChatMessage validates a chat/command message
func (v *Validator) ValidateChatMessage(message string) error {
	// Max length
	if len(message) > 500 {
		return fmt.Errorf("message too long (max 500 characters)")
	}

	// Check for null bytes
	if strings.Contains(message, "\x00") {
		return fmt.Errorf("message contains invalid characters")
	}

	return nil
}

// SanitizeChatMessage cleans a chat message
func (v *Validator) SanitizeChatMessage(message string) string {
	// Trim whitespace
	message = strings.TrimSpace(message)

	// Replace multiple spaces
	message = regexp.MustCompile(`\s+`).ReplaceAllString(message, " ")

	// Limit length
	if len(message) > 500 {
		message = message[:500]
	}

	return message
}

// IsValidKey validates a key binding string
func (v *Validator) IsValidKey(key string) bool {
	validKeys := []string{
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
		"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9", "F10", "F11", "F12",
		"Space", "Enter", "Escape", "Tab", "Backspace", "Delete",
		"Shift", "Ctrl", "Alt",
		"Up", "Down", "Left", "Right",
	}

	key = strings.ToUpper(key)
	for _, valid := range validKeys {
		if key == valid {
			return true
		}
	}
	return false
}

// Global validator instance
var DefaultValidator = NewValidator()

// Convenience functions using default validator
func ValidateWorldName(name string) error {
	return DefaultValidator.ValidateWorldName(name)
}

func SanitizeWorldName(name string) (string, error) {
	return DefaultValidator.SanitizeWorldName(name)
}

func ValidatePlayerName(name string) error {
	return DefaultValidator.ValidatePlayerName(name)
}

func ValidateSeedString(seedStr string) (int64, error) {
	return DefaultValidator.ValidateSeedString(seedStr)
}
