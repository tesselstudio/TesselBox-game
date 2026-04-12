package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// TwoFactorMethod represents 2FA type
type TwoFactorMethod int

const (
	TwoFactorTOTP TwoFactorMethod = iota // Time-based OTP
	TwoFactorEmail
	TwoFactorBackupCodes
)

// String returns method name
func (t TwoFactorMethod) String() string {
	switch t {
	case TwoFactorTOTP:
		return "TOTP"
	case TwoFactorEmail:
		return "Email"
	case TwoFactorBackupCodes:
		return "Backup"
	}
	return "Unknown"
}

// PlayerSecurity stores security settings
type PlayerSecurity struct {
	PlayerID string `json:"player_id"`

	// Password (hashed)
	PasswordHash string `json:"password_hash,omitempty"`
	Salt         string `json:"salt,omitempty"`

	// 2FA
	TwoFactorEnabled bool            `json:"two_factor_enabled"`
	TwoFactorMethod  TwoFactorMethod `json:"two_factor_method"`
	TwoFactorSecret  string          `json:"two_factor_secret,omitempty"`
	BackupCodes      []string        `json:"backup_codes,omitempty"`

	// IP Security
	LastIP      string   `json:"last_ip,omitempty"`
	TrustedIPs  []string `json:"trusted_ips,omitempty"`
	IPWhitelist bool     `json:"ip_whitelist"`

	// Session
	LastLogin     time.Time `json:"last_login"`
	LastLogout    time.Time `json:"last_logout"`
	SessionToken  string    `json:"session_token,omitempty"`
	SessionExpiry time.Time `json:"session_expiry"`

	// Security events
	FailedLogins int        `json:"failed_logins"`
	LockedUntil  *time.Time `json:"locked_until,omitempty"`

	// Settings
	AutoLogoutTime time.Duration `json:"auto_logout_time"`
}

// NewPlayerSecurity creates new security settings
func NewPlayerSecurity(playerID string) *PlayerSecurity {
	return &PlayerSecurity{
		PlayerID:         playerID,
		TwoFactorEnabled: false,
		TrustedIPs:       make([]string, 0),
		BackupCodes:      make([]string, 0),
		FailedLogins:     0,
		AutoLogoutTime:   30 * time.Minute,
	}
}

// EnableTwoFactor enables 2FA
func (ps *PlayerSecurity) EnableTwoFactor(method TwoFactorMethod) (string, error) {
	if ps.TwoFactorEnabled {
		return "", fmt.Errorf("2FA already enabled")
	}

	ps.TwoFactorMethod = method
	ps.TwoFactorEnabled = true

	switch method {
	case TwoFactorTOTP:
		// Generate TOTP secret
		secret := make([]byte, 20)
		if _, err := rand.Read(secret); err != nil {
			return "", err
		}
		ps.TwoFactorSecret = base32.StdEncoding.EncodeToString(secret)

		// Generate backup codes
		ps.BackupCodes = make([]string, 10)
		for i := range ps.BackupCodes {
			code := make([]byte, 8)
			if _, err := rand.Read(code); err != nil {
				return "", err
			}
			ps.BackupCodes[i] = base32.StdEncoding.EncodeToString(code)[:8]
		}

		return ps.TwoFactorSecret, nil

	case TwoFactorEmail:
		// Email 2FA would send codes via email
		return "", nil

	default:
		return "", fmt.Errorf("unknown 2FA method")
	}
}

// DisableTwoFactor disables 2FA
func (ps *PlayerSecurity) DisableTwoFactor() {
	ps.TwoFactorEnabled = false
	ps.TwoFactorSecret = ""
	ps.BackupCodes = make([]string, 0)
}

// VerifyTOTP verifies a TOTP code
func (ps *PlayerSecurity) VerifyTOTP(code string) bool {
	if !ps.TwoFactorEnabled || ps.TwoFactorMethod != TwoFactorTOTP {
		return false
	}

	// In real implementation, validate TOTP
	// For now, simple length check
	return len(code) == 6
}

// VerifyBackupCode verifies and consumes a backup code
func (ps *PlayerSecurity) VerifyBackupCode(code string) bool {
	for i, bc := range ps.BackupCodes {
		if subtle.ConstantTimeCompare([]byte(bc), []byte(code)) == 1 {
			// Remove used code
			ps.BackupCodes = append(ps.BackupCodes[:i], ps.BackupCodes[i+1:]...)
			return true
		}
	}
	return false
}

// IsIPTrusted checks if IP is trusted
func (ps *PlayerSecurity) IsIPTrusted(ip string) bool {
	if !ps.IPWhitelist {
		return true
	}

	for _, trusted := range ps.TrustedIPs {
		if trusted == ip {
			return true
		}
	}

	return false
}

// AddTrustedIP adds a trusted IP
func (ps *PlayerSecurity) AddTrustedIP(ip string) {
	for _, trusted := range ps.TrustedIPs {
		if trusted == ip {
			return // Already trusted
		}
	}
	ps.TrustedIPs = append(ps.TrustedIPs, ip)
}

// RemoveTrustedIP removes a trusted IP
func (ps *PlayerSecurity) RemoveTrustedIP(ip string) {
	for i, trusted := range ps.TrustedIPs {
		if trusted == ip {
			ps.TrustedIPs = append(ps.TrustedIPs[:i], ps.TrustedIPs[i+1:]...)
			return
		}
	}
}

// RecordFailedLogin records a failed login attempt
func (ps *PlayerSecurity) RecordFailedLogin() {
	ps.FailedLogins++

	// Lock account after 5 failed attempts
	if ps.FailedLogins >= 5 {
		lockUntil := time.Now().Add(30 * time.Minute)
		ps.LockedUntil = &lockUntil
	}
}

// ResetFailedLogins resets failed login count
func (ps *PlayerSecurity) ResetFailedLogins() {
	ps.FailedLogins = 0
	ps.LockedUntil = nil
}

// IsLocked checks if account is locked
func (ps *PlayerSecurity) IsLocked() bool {
	if ps.LockedUntil == nil {
		return false
	}

	if time.Now().After(*ps.LockedUntil) {
		ps.ResetFailedLogins()
		return false
	}

	return true
}

// GenerateSessionToken generates a new session token
func (ps *PlayerSecurity) GenerateSessionToken() string {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return ""
	}

	ps.SessionToken = base32.StdEncoding.EncodeToString(token)
	ps.SessionExpiry = time.Now().Add(ps.AutoLogoutTime)

	return ps.SessionToken
}

// ValidateSession validates a session token
func (ps *PlayerSecurity) ValidateSession(token string) bool {
	if ps.SessionToken == "" || ps.SessionToken != token {
		return false
	}

	if time.Now().After(ps.SessionExpiry) {
		return false
	}

	// Extend session
	ps.SessionExpiry = time.Now().Add(ps.AutoLogoutTime)

	return true
}

// ClearSession clears session
func (ps *PlayerSecurity) ClearSession() {
	ps.LastLogout = time.Now()
	ps.SessionToken = ""
}

// SecurityManager manages player security
type SecurityManager struct {
	players map[string]*PlayerSecurity

	storagePath string
}

// NewSecurityManager creates new manager
func NewSecurityManager(storageDir string) *SecurityManager {
	return &SecurityManager{
		players:     make(map[string]*PlayerSecurity),
		storagePath: filepath.Join(storageDir, "security.json"),
	}
}

// GetOrCreateSecurity gets or creates security settings
func (sm *SecurityManager) GetOrCreateSecurity(playerID string) *PlayerSecurity {
	if ps, exists := sm.players[playerID]; exists {
		return ps
	}

	ps := NewPlayerSecurity(playerID)
	sm.players[playerID] = ps
	return ps
}

// GetSecurity gets security settings
func (sm *SecurityManager) GetSecurity(playerID string) (*PlayerSecurity, bool) {
	ps, exists := sm.players[playerID]
	return ps, exists
}

// AuthenticatePlayer authenticates a player
func (sm *SecurityManager) AuthenticatePlayer(playerID, password, ip string, twoFactorCode string) (bool, []string) {
	ps := sm.GetOrCreateSecurity(playerID)

	// Check if locked
	if ps.IsLocked() {
		return false, []string{"Account locked due to too many failed attempts"}
	}

	// In real implementation, check password hash
	passwordValid := true // Placeholder

	if !passwordValid {
		ps.RecordFailedLogin()
		return false, []string{"Invalid credentials"}
	}

	// Check IP
	if ps.IPWhitelist && !ps.IsIPTrusted(ip) {
		return false, []string{"IP not whitelisted"}
	}

	// Check 2FA
	if ps.TwoFactorEnabled {
		var twoFactorValid bool

		switch ps.TwoFactorMethod {
		case TwoFactorTOTP:
			twoFactorValid = ps.VerifyTOTP(twoFactorCode)
			if !twoFactorValid {
				// Try backup codes
				twoFactorValid = ps.VerifyBackupCode(twoFactorCode)
			}
		default:
			twoFactorValid = true // Other methods handle separately
		}

		if !twoFactorValid {
			return false, []string{"Invalid 2FA code"}
		}
	}

	// Success
	ps.ResetFailedLogins()
	ps.LastIP = ip
	ps.LastLogin = time.Now()
	ps.AddTrustedIP(ip)
	ps.GenerateSessionToken()

	return true, nil
}

// ValidateSession validates a session
func (sm *SecurityManager) ValidateSession(playerID, token string) bool {
	ps, exists := sm.GetSecurity(playerID)
	if !exists {
		return false
	}

	return ps.ValidateSession(token)
}

// Logout logs out a player
func (sm *SecurityManager) Logout(playerID string) {
	if ps, exists := sm.GetSecurity(playerID); exists {
		ps.ClearSession()
	}
}

// EnableTwoFactor enables 2FA for player
func (sm *SecurityManager) EnableTwoFactor(playerID string, method TwoFactorMethod) (string, error) {
	ps := sm.GetOrCreateSecurity(playerID)
	return ps.EnableTwoFactor(method)
}

// DisableTwoFactor disables 2FA
func (sm *SecurityManager) DisableTwoFactor(playerID string) {
	if ps, exists := sm.GetSecurity(playerID); exists {
		ps.DisableTwoFactor()
	}
}

// GetTwoFactorStatus gets 2FA status
func (sm *SecurityManager) GetTwoFactorStatus(playerID string) struct {
	Enabled bool
	Method  string
} {
	ps, exists := sm.GetSecurity(playerID)
	if !exists {
		return struct {
			Enabled bool
			Method  string
		}{false, ""}
	}

	methodName := ""
	if ps.TwoFactorEnabled {
		methodName = ps.TwoFactorMethod.String()
	}

	return struct {
		Enabled bool
		Method  string
	}{ps.TwoFactorEnabled, methodName}
}

// Save saves security data
func (sm *SecurityManager) Save() error {
	data, err := json.MarshalIndent(sm.players, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	if err := os.WriteFile(sm.storagePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	return nil
}

// Load loads security data
func (sm *SecurityManager) Load() error {
	data, err := os.ReadFile(sm.storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read: %w", err)
	}

	var loaded map[string]*PlayerSecurity
	if err := json.Unmarshal(data, &loaded); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	sm.players = loaded
	if sm.players == nil {
		sm.players = make(map[string]*PlayerSecurity)
	}

	return nil
}

// CleanupSessions cleans up expired sessions
func (sm *SecurityManager) CleanupSessions() int {
	cleaned := 0

	for _, ps := range sm.players {
		if ps.SessionToken != "" && time.Now().After(ps.SessionExpiry) {
			ps.ClearSession()
			cleaned++
		}
	}

	return cleaned
}

// GetStats returns security statistics
func (sm *SecurityManager) GetStats() (total, with2FA, withWhitelist int) {
	total = len(sm.players)

	for _, ps := range sm.players {
		if ps.TwoFactorEnabled {
			with2FA++
		}
		if ps.IPWhitelist {
			withWhitelist++
		}
	}

	return
}
