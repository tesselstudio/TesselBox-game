package save

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ValidationResult contains validation status and any errors
type ValidationResult struct {
	Valid      bool
	Errors     []string
	Warnings   []string
	Checksum   string
	CanRecover bool
	BackupPath string
}

// SaveValidator handles save file validation and corruption detection
type SaveValidator struct {
	backupManager *BackupManager
}

// NewSaveValidator creates a new save validator
func NewSaveValidator(backupManager *BackupManager) *SaveValidator {
	return &SaveValidator{
		backupManager: backupManager,
	}
}

// ValidateSaveData checks if save data is valid and complete
func (sv *SaveValidator) ValidateSaveData(data *SaveData) ValidationResult {
	result := ValidationResult{
		Valid:      true,
		Errors:     []string{},
		Warnings:   []string{},
		CanRecover: false,
	}

	// Check required fields
	if data.Version == "" {
		result.Errors = append(result.Errors, "Missing version field")
		result.Valid = false
	}

	if data.WorldName == "" {
		result.Errors = append(result.Errors, "Missing world name")
		result.Valid = false
	}

	if data.PlayerName == "" {
		result.Errors = append(result.Errors, "Missing player name")
		result.Valid = false
	}

	// Validate timestamps
	if data.SaveTime.IsZero() {
		result.Warnings = append(result.Warnings, "Missing save timestamp")
	}

	// Validate player state
	if data.PlayerHealth < 0 || data.PlayerHealth > data.PlayerMaxHealth {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Invalid player health: %.2f", data.PlayerHealth))
	}

	if data.PlayerX == 0 && data.PlayerY == 0 {
		result.Warnings = append(result.Warnings, "Player at origin (0,0) - possible spawn issue")
	}

	// Validate inventory
	for i, slot := range data.InventorySlots {
		if slot.Quantity < 0 {
			result.Errors = append(result.Errors, fmt.Sprintf("Invalid quantity in inventory slot %d: %d", i, slot.Quantity))
			result.Valid = false
		}
		if slot.Durability < 0 {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Negative durability in slot %d", i))
		}
	}

	// Validate world seed
	if data.Seed == 0 && data.GameMode == "survival" {
		result.Warnings = append(result.Warnings, "Zero seed in survival mode - world may be unplayable")
	}

	// Check if we can recover from backup
	if !result.Valid {
		backupPath, err := sv.backupManager.FindLatestBackup()
		if err == nil {
			result.CanRecover = true
			result.BackupPath = backupPath
		}
	}

	return result
}

// ValidateSaveFile validates a save file on disk
func (sv *SaveValidator) ValidateSaveFile(savePath string) ValidationResult {
	result := ValidationResult{
		Valid:      true,
		Errors:     []string{},
		Warnings:   []string{},
		CanRecover: false,
	}

	// Check if file exists
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		result.Errors = append(result.Errors, "Save file does not exist")
		result.Valid = false
		return result
	}

	// Read file
	data, err := os.ReadFile(savePath)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to read save file: %v", err))
		result.Valid = false
		return result
	}

	// Calculate checksum
	result.Checksum = CalculateChecksum(data)

	// Try to unmarshal
	var saveData SaveData
	if err := json.Unmarshal(data, &saveData); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Corrupted save file (invalid JSON): %v", err))
		result.Valid = false

		// Check for backup
		backupPath, backupErr := sv.backupManager.FindLatestBackup()
		if backupErr == nil {
			result.CanRecover = true
			result.BackupPath = backupPath
		}

		return result
	}

	// Validate the data structure
	dataResult := sv.ValidateSaveData(&saveData)
	result.Errors = append(result.Errors, dataResult.Errors...)
	result.Warnings = append(result.Warnings, dataResult.Warnings...)
	result.Valid = dataResult.Valid
	result.CanRecover = dataResult.CanRecover
	result.BackupPath = dataResult.BackupPath

	return result
}

// TryRecover attempts to load from backup if save is corrupted
func (sv *SaveValidator) TryRecover() (*SaveData, error) {
	backupPath, err := sv.backupManager.FindLatestBackup()
	if err != nil {
		return nil, fmt.Errorf("no backup available for recovery: %w", err)
	}

	saveData, err := sv.backupManager.LoadBackup(backupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load backup: %w", err)
	}

	return saveData, nil
}

// ValidateWorldName checks if a world name is valid
func ValidateWorldName(name string) error {
	// Max length
	if len(name) > 32 {
		return fmt.Errorf("world name too long (max 32 characters)")
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
			return fmt.Errorf("'%s' is a reserved name", name)
		}
	}

	// Check for path traversal
	if strings.Contains(name, "..") || strings.Contains(name, "/") || strings.Contains(name, "\\") {
		return fmt.Errorf("world name contains invalid characters")
	}

	// Check for valid characters
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9_\-\s]+$`)
	if !validPattern.MatchString(name) {
		return fmt.Errorf("world name contains invalid characters (use only letters, numbers, spaces, hyphens, underscores)")
	}

	return nil
}

// ValidatePlayerName checks if a player name is valid
func ValidatePlayerName(name string) error {
	// Max length
	if len(name) > 16 {
		return fmt.Errorf("player name too long (max 16 characters)")
	}

	// Min length
	if len(name) < 1 {
		return fmt.Errorf("player name cannot be empty")
	}

	// Check for valid characters
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !validPattern.MatchString(name) {
		return fmt.Errorf("player name contains invalid characters (use only letters, numbers, underscores)")
	}

	return nil
}

// ValidateSeed validates and sanitizes a world seed
func ValidateSeed(seed int64) error {
	// Check for reasonable range
	if seed < -9223372036854775808 || seed > 9223372036854775807 {
		return fmt.Errorf("seed out of valid range")
	}
	return nil
}

// VerifyChecksum verifies a file against its expected checksum
func VerifyChecksum(filePath, expectedChecksum string) (bool, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}

	actualChecksum := CalculateChecksum(data)
	return actualChecksum == expectedChecksum, nil
}
