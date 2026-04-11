package debug

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"
)

const crashLogFilename = "crash.log"

// PanicInfo contains information about a recovered panic
type PanicInfo struct {
	Message     string
	Stack       string
	Time        time.Time
	Recoverable bool
}

// RecoveryHandler handles panic recovery and logging
type RecoveryHandler struct {
	crashLogPath string
	onRecover    func(PanicInfo)
}

// NewRecoveryHandler creates a new recovery handler
func NewRecoveryHandler(tesselboxDir string, onRecover func(PanicInfo)) *RecoveryHandler {
	return &RecoveryHandler{
		crashLogPath: filepath.Join(tesselboxDir, crashLogFilename),
		onRecover:    onRecover,
	}
}

// Recover catches panics, logs them, and calls the recovery callback
func (rh *RecoveryHandler) Recover() {
	if r := recover(); r != nil {
		panicInfo := PanicInfo{
			Message:     fmt.Sprintf("%v", r),
			Stack:       string(debug.Stack()),
			Time:        time.Now(),
			Recoverable: true,
		}

		// Log the crash
		rh.logCrash(panicInfo)

		// Call recovery callback
		if rh.onRecover != nil {
			rh.onRecover(panicInfo)
		}
	}
}

// logCrash writes crash information to the crash log file
func (rh *RecoveryHandler) logCrash(info PanicInfo) {
	// Ensure directory exists
	dir := filepath.Dir(rh.crashLogPath)
	os.MkdirAll(dir, 0755)

	// Open log file (append mode)
	file, err := os.OpenFile(rh.crashLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to write crash log: %v", err)
		return
	}
	defer file.Close()

	// Write crash info
	crashEntry := fmt.Sprintf(
		"=== CRASH at %s ===\nMessage: %s\nStack:\n%s\n===================\n\n",
		info.Time.Format("2006-01-02 15:04:05"),
		info.Message,
		info.Stack,
	)

	if _, err := file.WriteString(crashEntry); err != nil {
		log.Printf("Failed to write crash entry: %v", err)
	}
}

// TryEmergencySave attempts to save game state during a panic
func (rh *RecoveryHandler) TryEmergencySave(saveFunc func() error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Emergency save failed: %v", r)
		}
	}()

	if err := saveFunc(); err != nil {
		log.Printf("Emergency save error: %v", err)
	} else {
		log.Printf("Emergency save successful")
	}
}

// GetCrashLogPath returns the path to the crash log file
func (rh *RecoveryHandler) GetCrashLogPath() string {
	return rh.crashLogPath
}

// ReadCrashLog reads the contents of the crash log
func (rh *RecoveryHandler) ReadCrashLog() (string, error) {
	data, err := os.ReadFile(rh.crashLogPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return string(data), nil
}

// ClearCrashLog clears the crash log file
func (rh *RecoveryHandler) ClearCrashLog() error {
	return os.Remove(rh.crashLogPath)
}

// HasCrashLog checks if a crash log exists
func (rh *RecoveryHandler) HasCrashLog() bool {
	_, err := os.Stat(rh.crashLogPath)
	return !os.IsNotExist(err)
}
