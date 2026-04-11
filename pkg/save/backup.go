package save

import (
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	maxBackups      = 5
	backupRetention = 30 * 24 * time.Hour // 30 days
)

// BackupManager handles save file backups and recovery
type BackupManager struct {
	worldName  string
	playerName string
	backupDir  string
}

// NewBackupManager creates a new backup manager
func NewBackupManager(worldName, playerName, saveDir string) *BackupManager {
	backupDir := filepath.Join(saveDir, "backups")
	return &BackupManager{
		worldName:  worldName,
		playerName: playerName,
		backupDir:  backupDir,
	}
}

// CreateBackup creates a backup of the current save file
func (bm *BackupManager) CreateBackup(saveData *SaveData) error {
	// Ensure backup directory exists
	if err := os.MkdirAll(bm.backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Generate timestamped backup filename
	timestamp := time.Now().Format("20060102_150405")
	backupFilename := fmt.Sprintf("%s_%s_%s.json.gz", bm.worldName, bm.playerName, timestamp)
	backupPath := filepath.Join(bm.backupDir, backupFilename)

	// Marshal save data
	data, err := json.MarshalIndent(saveData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal backup data: %w", err)
	}

	// Compress and write backup
	file, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	if _, err := gzipWriter.Write(data); err != nil {
		gzipWriter.Close()
		return fmt.Errorf("failed to compress backup: %w", err)
	}

	if err := gzipWriter.Close(); err != nil {
		return fmt.Errorf("failed to close gzip writer: %w", err)
	}

	// Clean up old backups
	bm.cleanupOldBackups()

	return nil
}

// FindLatestBackup returns the most recent backup file path
func (bm *BackupManager) FindLatestBackup() (string, error) {
	entries, err := os.ReadDir(bm.backupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("no backups found")
		}
		return "", fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []os.DirEntry
	prefix := fmt.Sprintf("%s_%s_", bm.worldName, bm.playerName)

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), prefix) && strings.HasSuffix(entry.Name(), ".json.gz") {
			backups = append(backups, entry)
		}
	}

	if len(backups) == 0 {
		return "", fmt.Errorf("no backups found for %s/%s", bm.worldName, bm.playerName)
	}

	// Sort by modification time (newest first)
	sort.Slice(backups, func(i, j int) bool {
		infoI, _ := backups[i].Info()
		infoJ, _ := backups[j].Info()
		return infoI.ModTime().After(infoJ.ModTime())
	})

	return filepath.Join(bm.backupDir, backups[0].Name()), nil
}

// LoadBackup loads a backup file and returns the SaveData
func (bm *BackupManager) LoadBackup(backupPath string) (*SaveData, error) {
	file, err := os.Open(backupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open backup file: %w", err)
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzipReader.Close()

	data, err := io.ReadAll(gzipReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup data: %w", err)
	}

	var saveData SaveData
	if err := json.Unmarshal(data, &saveData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal backup data: %w", err)
	}

	return &saveData, nil
}

// cleanupOldBackups removes excess backups and old compressed backups
func (bm *BackupManager) cleanupOldBackups() {
	entries, err := os.ReadDir(bm.backupDir)
	if err != nil {
		return
	}

	prefix := fmt.Sprintf("%s_%s_", bm.worldName, bm.playerName)
	var backups []os.DirEntry

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), prefix) && strings.HasSuffix(entry.Name(), ".json.gz") {
			backups = append(backups, entry)
		}
	}

	// Sort by modification time (newest first)
	sort.Slice(backups, func(i, j int) bool {
		infoI, _ := backups[i].Info()
		infoJ, _ := backups[j].Info()
		return infoI.ModTime().After(infoJ.ModTime())
	})

	// Remove excess backups (keep only maxBackups)
	if len(backups) > maxBackups {
		for i := maxBackups; i < len(backups); i++ {
			path := filepath.Join(bm.backupDir, backups[i].Name())
			os.Remove(path)
		}
	}

	// Remove backups older than retention period
	cutoff := time.Now().Add(-backupRetention)
	for _, entry := range backups {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			path := filepath.Join(bm.backupDir, entry.Name())
			os.Remove(path)
		}
	}
}

// CalculateChecksum computes SHA256 checksum of save data
func CalculateChecksum(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
