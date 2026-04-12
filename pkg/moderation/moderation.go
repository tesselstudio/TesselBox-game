package moderation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// PunishmentType represents the type of punishment
type PunishmentType int

const (
	PunishmentWarn PunishmentType = iota
	PunishmentKick
	PunishmentTempBan
	PunishmentBan
	PunishmentMute
)

// String returns human-readable punishment name
func (p PunishmentType) String() string {
	switch p {
	case PunishmentWarn:
		return "Warning"
	case PunishmentKick:
		return "Kick"
	case PunishmentTempBan:
		return "Temp Ban"
	case PunishmentBan:
		return "Ban"
	case PunishmentMute:
		return "Mute"
	}
	return "Unknown"
}

// Punishment represents a punishment record
type Punishment struct {
	ID          string         `json:"id"`
	Type        PunishmentType `json:"type"`
	PlayerID    string         `json:"player_id"`
	PlayerName  string         `json:"player_name"`
	IssuedBy    string         `json:"issued_by"`    // Staff member ID
	IssuerName  string         `json:"issuer_name"`
	Reason      string         `json:"reason"`
	IssuedAt    time.Time      `json:"issued_at"`
	ExpiresAt   *time.Time     `json:"expires_at,omitempty"`
	Duration    time.Duration  `json:"duration"`
	Active      bool           `json:"active"`
	Evidence    string         `json:"evidence,omitempty"`
}

// IsExpired checks if punishment has expired
func (p *Punishment) IsExpired() bool {
	if p.ExpiresAt == nil {
		return false // Permanent
	}
	return time.Now().After(*p.ExpiresAt)
}

// TimeRemaining returns time until expiry
func (p *Punishment) TimeRemaining() time.Duration {
	if p.ExpiresAt == nil {
		return 0
	}
	remaining := time.Until(*p.ExpiresAt)
	if remaining < 0 {
		return 0
	}
	return remaining
}

// PlayerReport represents a player report
type PlayerReport struct {
	ID          string         `json:"id"`
	ReporterID  string         `json:"reporter_id"`
	ReporterName string        `json:"reporter_name"`
	TargetID    string         `json:"target_id"`
	TargetName  string         `json:"target_name"`
	Type        ReportType     `json:"type"`
	Description string         `json:"description"`
	Evidence    []string       `json:"evidence,omitempty"` // Screenshot paths, etc
	Status      ReportStatus   `json:"status"`
	AssignedTo  string         `json:"assigned_to,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	ResolvedAt  *time.Time     `json:"resolved_at,omitempty"`
	Resolution  string         `json:"resolution,omitempty"`
}

// ReportType represents the type of report
type ReportType int

const (
	ReportChat ReportType = iota
	ReportGriefing
	ReportCheating
	ReportHarassment
	ReportExploiting
	ReportOther
)

// String returns report type name
func (r ReportType) String() string {
	switch r {
	case ReportChat:
		return "Chat Violation"
	case ReportGriefing:
		return "Griefing"
	case ReportCheating:
		return "Cheating"
	case ReportHarassment:
		return "Harassment"
	case ReportExploiting:
		return "Exploiting"
	case ReportOther:
		return "Other"
	}
	return "Unknown"
}

// ReportStatus represents report status
type ReportStatus int

const (
	ReportOpen ReportStatus = iota
	ReportInvestigating
	ReportResolved
	ReportDismissed
)

// String returns status name
func (r ReportStatus) String() string {
	switch r {
	case ReportOpen:
		return "Open"
	case ReportInvestigating:
		return "Investigating"
	case ReportResolved:
		return "Resolved"
	case ReportDismissed:
		return "Dismissed"
	}
	return "Unknown"
}

// ModerationManager manages punishments and reports
type ModerationManager struct {
	punishments []Punishment
	reports     []PlayerReport
	
	storagePath string
}

// NewModerationManager creates a new moderation manager
func NewModerationManager(storageDir string) *ModerationManager {
	return &ModerationManager{
		punishments: make([]Punishment, 0),
		reports:     make([]PlayerReport, 0),
		storagePath: filepath.Join(storageDir, "moderation"),
	}
}

// IssuePunishment issues a new punishment
func (mm *ModerationManager) IssuePunishment(pType PunishmentType, playerID, playerName, issuerID, issuerName, reason string, duration time.Duration, evidence string) (*Punishment, error) {
	punishment := Punishment{
		ID:         generateID(),
		Type:       pType,
		PlayerID:   playerID,
		PlayerName: playerName,
		IssuedBy:   issuerID,
		IssuerName: issuerName,
		Reason:     reason,
		IssuedAt:   time.Now(),
		Duration:   duration,
		Active:     true,
		Evidence:   evidence,
	}
	
	// Set expiry if not permanent
	if duration > 0 {
		expires := punishment.IssuedAt.Add(duration)
		punishment.ExpiresAt = &expires
	}
	
	mm.punishments = append(mm.punishments, punishment)
	
	return &punishment, nil
}

// Warn issues a warning
func (mm *ModerationManager) Warn(playerID, playerName, issuerID, issuerName, reason string) *Punishment {
	p, _ := mm.IssuePunishment(PunishmentWarn, playerID, playerName, issuerID, issuerName, reason, 0, "")
	return p
}

// Kick issues a kick
func (mm *ModerationManager) Kick(playerID, playerName, issuerID, issuerName, reason string) *Punishment {
	p, _ := mm.IssuePunishment(PunishmentKick, playerID, playerName, issuerID, issuerName, reason, 0, "")
	return p
}

// TempBan issues a temporary ban
func (mm *ModerationManager) TempBan(playerID, playerName, issuerID, issuerName, reason string, duration time.Duration) *Punishment {
	p, _ := mm.IssuePunishment(PunishmentTempBan, playerID, playerName, issuerID, issuerName, reason, duration, "")
	return p
}

// Ban issues a permanent ban
func (mm *ModerationManager) Ban(playerID, playerName, issuerID, issuerName, reason string) *Punishment {
	p, _ := mm.IssuePunishment(PunishmentBan, playerID, playerName, issuerID, issuerName, reason, 0, "")
	return p
}

// Mute issues a mute
func (mm *ModerationManager) Mute(playerID, playerName, issuerID, issuerName, reason string, duration time.Duration) *Punishment {
	p, _ := mm.IssuePunishment(PunishmentMute, playerID, playerName, issuerID, issuerName, reason, duration, "")
	return p
}

// RevokePunishment revokes an active punishment
func (mm *ModerationManager) RevokePunishment(punishmentID, revokedBy, reason string) error {
	for i, p := range mm.punishments {
		if p.ID == punishmentID && p.Active {
			mm.punishments[i].Active = false
			return nil
		}
	}
	return fmt.Errorf("punishment not found or already expired")
}

// GetActivePunishments returns active punishments for a player
func (mm *ModerationManager) GetActivePunishments(playerID string) []Punishment {
	result := make([]Punishment, 0)
	now := time.Now()
	
	for _, p := range mm.punishments {
		if p.PlayerID != playerID {
			continue
		}
		
		// Check if expired
		if p.ExpiresAt != nil && now.After(*p.ExpiresAt) {
			continue
		}
		
		if p.Active {
			result = append(result, p)
		}
	}
	
	return result
}

// IsBanned checks if player is currently banned
func (mm *ModerationManager) IsBanned(playerID string) bool {
	punishments := mm.GetActivePunishments(playerID)
	for _, p := range punishments {
		if p.Type == PunishmentBan || p.Type == PunishmentTempBan {
			return true
		}
	}
	return false
}

// IsMuted checks if player is currently muted
func (mm *ModerationManager) IsMuted(playerID string) bool {
	punishments := mm.GetActivePunishments(playerID)
	for _, p := range punishments {
		if p.Type == PunishmentMute {
			return true
		}
	}
	return false
}

// GetBanReason returns the ban reason if player is banned
func (mm *ModerationManager) GetBanReason(playerID string) string {
	punishments := mm.GetActivePunishments(playerID)
	for _, p := range punishments {
		if p.Type == PunishmentBan || p.Type == PunishmentTempBan {
			return p.Reason
		}
	}
	return ""
}

// GetMuteTimeRemaining returns remaining mute time
func (mm *ModerationManager) GetMuteTimeRemaining(playerID string) time.Duration {
	punishments := mm.GetActivePunishments(playerID)
	for _, p := range punishments {
		if p.Type == PunishmentMute {
			return p.TimeRemaining()
		}
	}
	return 0
}

// SubmitReport submits a player report
func (mm *ModerationManager) SubmitReport(reporterID, reporterName, targetID, targetName string, rType ReportType, description string, evidence []string) *PlayerReport {
	report := PlayerReport{
		ID:           generateID(),
		ReporterID:   reporterID,
		ReporterName: reporterName,
		TargetID:     targetID,
		TargetName:   targetName,
		Type:         rType,
		Description:  description,
		Evidence:     evidence,
		Status:       ReportOpen,
		CreatedAt:    time.Now(),
	}
	
	mm.reports = append(mm.reports, report)
	
	return &report
}

// GetReport gets a report by ID
func (mm *ModerationManager) GetReport(reportID string) (*PlayerReport, bool) {
	for i, r := range mm.reports {
		if r.ID == reportID {
			return &mm.reports[i], true
		}
	}
	return nil, false
}

// GetOpenReports returns all open reports
func (mm *ModerationManager) GetOpenReports() []PlayerReport {
	result := make([]PlayerReport, 0)
	for _, r := range mm.reports {
		if r.Status == ReportOpen || r.Status == ReportInvestigating {
			result = append(result, r)
		}
	}
	return result
}

// GetReportsByTarget returns reports against a player
func (mm *ModerationManager) GetReportsByTarget(targetID string) []PlayerReport {
	result := make([]PlayerReport, 0)
	for _, r := range mm.reports {
		if r.TargetID == targetID {
			result = append(result, r)
		}
	}
	return result
}

// AssignReport assigns a report to a staff member
func (mm *ModerationManager) AssignReport(reportID, staffID string) error {
	for i, r := range mm.reports {
		if r.ID == reportID {
			mm.reports[i].AssignedTo = staffID
			mm.reports[i].Status = ReportInvestigating
			return nil
		}
	}
	return fmt.Errorf("report not found")
}

// ResolveReport resolves a report
func (mm *ModerationManager) ResolveReport(reportID, resolution string) error {
	for i, r := range mm.reports {
		if r.ID == reportID {
			now := time.Now()
			mm.reports[i].Status = ReportResolved
			mm.reports[i].ResolvedAt = &now
			mm.reports[i].Resolution = resolution
			return nil
		}
	}
	return fmt.Errorf("report not found")
}

// DismissReport dismisses a report
func (mm *ModerationManager) DismissReport(reportID, reason string) error {
	for i, r := range mm.reports {
		if r.ID == reportID {
			now := time.Now()
			mm.reports[i].Status = ReportDismissed
			mm.reports[i].ResolvedAt = &now
			mm.reports[i].Resolution = reason
			return nil
		}
	}
	return fmt.Errorf("report not found")
}

// GetPlayerHistory returns punishment history for a player
func (mm *ModerationManager) GetPlayerHistory(playerID string) []Punishment {
	result := make([]Punishment, 0)
	for _, p := range mm.punishments {
		if p.PlayerID == playerID {
			result = append(result, p)
		}
	}
	return result
}

// GetStats returns moderation statistics
func (mm *ModerationManager) GetStats() (totalPunishments, activePunishments, openReports, totalReports int) {
	totalPunishments = len(mm.punishments)
	
	for _, p := range mm.punishments {
		if p.Active && !p.IsExpired() {
			activePunishments++
		}
	}
	
	totalReports = len(mm.reports)
	for _, r := range mm.reports {
		if r.Status == ReportOpen || r.Status == ReportInvestigating {
			openReports++
		}
	}
	
	return
}

// Save saves moderation data
func (mm *ModerationManager) Save() error {
	// Ensure directory exists
	if err := os.MkdirAll(mm.storagePath, 0755); err != nil {
		return fmt.Errorf("failed to create moderation directory: %w", err)
	}
	
	// Save punishments
	punishmentsPath := filepath.Join(mm.storagePath, "punishments.json")
	punishmentsData, err := json.MarshalIndent(mm.punishments, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal punishments: %w", err)
	}
	if err := os.WriteFile(punishmentsPath, punishmentsData, 0644); err != nil {
		return fmt.Errorf("failed to write punishments: %w", err)
	}
	
	// Save reports
	reportsPath := filepath.Join(mm.storagePath, "reports.json")
	reportsData, err := json.MarshalIndent(mm.reports, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal reports: %w", err)
	}
	if err := os.WriteFile(reportsPath, reportsData, 0644); err != nil {
		return fmt.Errorf("failed to write reports: %w", err)
	}
	
	return nil
}

// Load loads moderation data
func (mm *ModerationManager) Load() error {
	// Load punishments
	punishmentsPath := filepath.Join(mm.storagePath, "punishments.json")
	if data, err := os.ReadFile(punishmentsPath); err == nil {
		if err := json.Unmarshal(data, &mm.punishments); err != nil {
			return fmt.Errorf("failed to unmarshal punishments: %w", err)
		}
	}
	
	// Load reports
	reportsPath := filepath.Join(mm.storagePath, "reports.json")
	if data, err := os.ReadFile(reportsPath); err == nil {
		if err := json.Unmarshal(data, &mm.reports); err != nil {
			return fmt.Errorf("failed to unmarshal reports: %w", err)
		}
	}
	
	return nil
}

// CleanupExpired removes expired punishments from active list (keeps history)
func (mm *ModerationManager) CleanupExpired() {
	now := time.Now()
	for i := range mm.punishments {
		if mm.punishments[i].Active && mm.punishments[i].ExpiresAt != nil && now.After(*mm.punishments[i].ExpiresAt) {
			mm.punishments[i].Active = false
		}
	}
}

func generateID() string {
	return fmt.Sprintf("mod_%d", time.Now().UnixNano())
}
