package mail

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"tesselbox/pkg/items"
)

// MailStatus represents the status of mail
type MailStatus int

const (
	MailUnread MailStatus = iota
	MailRead
	MailDeleted
)

// MailAttachment represents an attached item or money
type MailAttachment struct {
	Money    float64          `json:"money,omitempty"`
	Items    []items.Item     `json:"items,omitempty"`
	COD      float64          `json:"cod,omitempty"` // Cash on delivery amount
}

// MailMessage represents a mail message
type MailMessage struct {
	ID          string         `json:"id"`
	FromID      string         `json:"from_id"`
	FromName    string         `json:"from_name"`
	ToID        string         `json:"to_id"`
	ToName      string         `json:"to_name"`
	Subject     string         `json:"subject"`
	Body        string         `json:"body"`
	
	// Attachments
	Attachments MailAttachment `json:"attachments,omitempty"`
	
	// Status
	Status      MailStatus     `json:"status"`
	IsReturned  bool           `json:"is_returned"`
	
	// Timing
	SentAt      time.Time      `json:"sent_at"`
	ExpiresAt   time.Time      `json:"expires_at"`
	ReadAt      *time.Time     `json:"read_at,omitempty"`
}

// NewMailMessage creates a new mail message
func NewMailMessage(fromID, fromName, toID, toName, subject, body string) *MailMessage {
	return &MailMessage{
		ID:          generateMailID(),
		FromID:      fromID,
		FromName:    fromName,
		ToID:        toID,
		ToName:      toName,
		Subject:     subject,
		Body:        body,
		Attachments: MailAttachment{},
		Status:      MailUnread,
		IsReturned:  false,
		SentAt:      time.Now(),
		ExpiresAt:   time.Now().Add(30 * 24 * time.Hour), // 30 days
	}
}

func generateMailID() string {
	return fmt.Sprintf("mail_%d", time.Now().UnixNano())
}

// AttachMoney attaches money to the mail
func (m *MailMessage) AttachMoney(amount float64) {
	m.Attachments.Money = amount
}

// AttachItem attaches an item to the mail
func (m *MailMessage) AttachItem(item items.Item) {
	m.Attachments.Items = append(m.Attachments.Items, item)
}

// SetCOD sets cash on delivery
func (m *MailMessage) SetCOD(amount float64) {
	m.Attachments.COD = amount
}

// MarkRead marks the mail as read
func (m *MailMessage) MarkRead() {
	if m.Status == MailUnread {
		m.Status = MailRead
		now := time.Now()
		m.ReadAt = &now
	}
}

// MarkDeleted marks the mail as deleted
func (m *MailMessage) MarkDeleted() {
	m.Status = MailDeleted
}

// HasAttachments checks if mail has attachments
func (m *MailMessage) HasAttachments() bool {
	return m.Attachments.Money > 0 || len(m.Attachments.Items) > 0
}

// IsExpired checks if mail has expired
func (m *MailMessage) IsExpired() bool {
	return time.Now().After(m.ExpiresAt)
}

// IsCOD checks if this is a COD mail
func (m *MailMessage) IsCOD() bool {
	return m.Attachments.COD > 0
}

// Mailbox represents a player's mailbox
type Mailbox struct {
	PlayerID     string          `json:"player_id"`
	Messages     []MailMessage   `json:"messages"`
	Capacity     int             `json:"capacity"`
	TotalReceived int            `json:"total_received"`
	TotalSent    int             `json:"total_sent"`
}

// NewMailbox creates a new mailbox
func NewMailbox(playerID string) *Mailbox {
	return &Mailbox{
		PlayerID:     playerID,
		Messages:     make([]MailMessage, 0),
		Capacity:     100,
		TotalReceived: 0,
		TotalSent:    0,
	}
}

// CanReceive checks if mailbox can receive more mail
func (mb *Mailbox) CanReceive() bool {
	return len(mb.Messages) < mb.Capacity
}

// AddMessage adds a message to the mailbox
func (mb *Mailbox) AddMessage(msg MailMessage) error {
	if !mb.CanReceive() {
		return fmt.Errorf("mailbox is full")
	}
	
	mb.Messages = append(mb.Messages, msg)
	mb.TotalReceived++
	
	return nil
}

// GetUnreadCount returns number of unread messages
func (mb *Mailbox) GetUnreadCount() int {
	count := 0
	for _, msg := range mb.Messages {
		if msg.Status == MailUnread {
			count++
		}
	}
	return count
}

// GetMessage gets a message by ID
func (mb *Mailbox) GetMessage(messageID string) *MailMessage {
	for i := range mb.Messages {
		if mb.Messages[i].ID == messageID {
			return &mb.Messages[i]
		}
	}
	return nil
}

// GetUnreadMessages returns unread messages
func (mb *Mailbox) GetUnreadMessages() []MailMessage {
	result := make([]MailMessage, 0)
	for _, msg := range mb.Messages {
		if msg.Status == MailUnread {
			result = append(result, msg)
		}
	}
	return result
}

// DeleteMessage deletes a message
func (mb *Mailbox) DeleteMessage(messageID string) bool {
	for i, msg := range mb.Messages {
		if msg.ID == messageID {
			mb.Messages = append(mb.Messages[:i], mb.Messages[i+1:]...)
			return true
		}
	}
	return false
}

// CleanupExpired removes expired messages
func (mb *Mailbox) CleanupExpired() int {
	removed := 0
	newMessages := make([]MailMessage, 0)
	
	for _, msg := range mb.Messages {
		if !msg.IsExpired() {
			newMessages = append(newMessages, msg)
		} else {
			removed++
		}
	}
	
	mb.Messages = newMessages
	return removed
}

// MailSystem manages all mailboxes
type MailSystem struct {
	mailboxes   map[string]*Mailbox
	
	// Settings
	maxAttachmentItems int
	maxAttachmentMoney float64
	basePostage       float64
	
	storagePath       string
}

// NewMailSystem creates a new mail system
func NewMailSystem(storageDir string) *MailSystem {
	return &MailSystem{
		mailboxes:          make(map[string]*Mailbox),
		maxAttachmentItems: 10,
		maxAttachmentMoney: 1000000,
		basePostage:        5.0,
		storagePath:        filepath.Join(storageDir, "mail.json"),
	}
}

// GetOrCreateMailbox gets or creates a mailbox
func (ms *MailSystem) GetOrCreateMailbox(playerID string) *Mailbox {
	if mailbox, exists := ms.mailboxes[playerID]; exists {
		return mailbox
	}
	
	mailbox := NewMailbox(playerID)
	ms.mailboxes[playerID] = mailbox
	return mailbox
}

// GetMailbox gets a mailbox (nil if not exists)
func (ms *MailSystem) GetMailbox(playerID string) *Mailbox {
	return ms.mailboxes[playerID]
}

// HasMailbox checks if player has a mailbox
func (ms *MailSystem) HasMailbox(playerID string) bool {
	_, exists := ms.mailboxes[playerID]
	return exists
}

// SendMail sends mail from one player to another
func (ms *MailSystem) SendMail(fromID, fromName, toID, toName, subject, body string, money float64, items []items.Item, cod float64) (*MailMessage, float64, error) {
	// Check recipient has mailbox
	mailbox := ms.GetOrCreateMailbox(toID)
	
	if !mailbox.CanReceive() {
		return nil, 0, fmt.Errorf("recipient's mailbox is full")
	}
	
	// Validate attachments
	if money > ms.maxAttachmentMoney {
		return nil, 0, fmt.Errorf("money amount exceeds limit")
	}
	
	if len(items) > ms.maxAttachmentItems {
		return nil, 0, fmt.Errorf("too many items attached")
	}
	
	// Create message
	msg := NewMailMessage(fromID, fromName, toID, toName, subject, body)
	
	if money > 0 {
		msg.AttachMoney(money)
	}
	
	for _, item := range items {
		msg.AttachItem(item)
	}
	
	if cod > 0 {
		msg.SetCOD(cod)
	}
	
	// Calculate postage
	postage := ms.basePostage
	if len(items) > 0 {
		postage += float64(len(items)) * 2.0 // 2 per item
	}
	if money > 0 {
		postage += money * 0.01 // 1% fee
	}
	
	// Add to recipient's mailbox
	if err := mailbox.AddMessage(*msg); err != nil {
		return nil, 0, err
	}
	
	return msg, postage, nil
}

// SendSystemMail sends mail from the system
func (ms *MailSystem) SendSystemMail(toID, toName, subject, body string, money float64, items []items.Item) (*MailMessage, error) {
	mailbox := ms.GetOrCreateMailbox(toID)
	
	if !mailbox.CanReceive() {
		return nil, fmt.Errorf("mailbox is full")
	}
	
	msg := NewMailMessage("SYSTEM", "System", toID, toName, subject, body)
	
	if money > 0 {
		msg.AttachMoney(money)
	}
	
	for _, item := range items {
		msg.AttachItem(item)
	}
	
	if err := mailbox.AddMessage(*msg); err != nil {
		return nil, err
	}
	
	return msg, nil
}

// ReturnMail returns mail to sender
func (ms *MailSystem) ReturnMail(playerID, messageID string) error {
	mailbox := ms.GetMailbox(playerID)
	if mailbox == nil {
		return fmt.Errorf("mailbox not found")
	}
	
	msg := mailbox.GetMessage(messageID)
	if msg == nil {
		return fmt.Errorf("message not found")
	}
	
	// Can't return system mail
	if msg.FromID == "SYSTEM" {
		return fmt.Errorf("cannot return system mail")
	}
	
	// Mark as returned
	msg.IsReturned = true
	msg.MarkDeleted()
	
	// Create return message
	returnMsg := NewMailMessage(
		playerID,
		"You",
		msg.FromID,
		msg.FromName,
		"Returned: "+msg.Subject,
		"This message was returned to you.",
	)
	returnMsg.Attachments = msg.Attachments
	
	// Send to original sender
	senderBox := ms.GetOrCreateMailbox(msg.FromID)
	if !senderBox.CanReceive() {
		return fmt.Errorf("sender's mailbox is full")
	}
	
	return senderBox.AddMessage(*returnMsg)
}

// ClaimAttachments claims attachments from mail
func (ms *MailSystem) ClaimAttachments(playerID, messageID string) (*MailAttachment, error) {
	mailbox := ms.GetMailbox(playerID)
	if mailbox == nil {
		return nil, fmt.Errorf("mailbox not found")
	}
	
	msg := mailbox.GetMessage(messageID)
	if msg == nil {
		return nil, fmt.Errorf("message not found")
	}
	
	if !msg.HasAttachments() {
		return nil, fmt.Errorf("no attachments to claim")
	}
	
	// Check COD
	if msg.IsCOD() {
		return nil, fmt.Errorf("COD mail - payment required")
	}
	
	// Get attachments
	attachments := msg.Attachments
	
	// Clear attachments
	msg.Attachments = MailAttachment{}
	
	// Mark as read
	msg.MarkRead()
	
	return &attachments, nil
}

// PayCOD pays for COD mail and claims attachments
func (ms *MailSystem) PayCOD(playerID, messageID string) (*MailAttachment, error) {
	mailbox := ms.GetMailbox(playerID)
	if mailbox == nil {
		return nil, fmt.Errorf("mailbox not found")
	}
	
	msg := mailbox.GetMessage(messageID)
	if msg == nil {
		return nil, fmt.Errorf("message not found")
	}
	
	if !msg.IsCOD() {
		return nil, fmt.Errorf("not a COD mail")
	}
	
	// In real implementation, deduct money from player
	// For now, just claim
	attachments := msg.Attachments
	msg.Attachments = MailAttachment{}
	msg.MarkRead()
	
	return &attachments, nil
}

// MarkRead marks mail as read
func (ms *MailSystem) MarkRead(playerID, messageID string) error {
	mailbox := ms.GetMailbox(playerID)
	if mailbox == nil {
		return fmt.Errorf("mailbox not found")
	}
	
	msg := mailbox.GetMessage(messageID)
	if msg == nil {
		return fmt.Errorf("message not found")
	}
	
	msg.MarkRead()
	return nil
}

// DeleteMail deletes mail
func (ms *MailSystem) DeleteMail(playerID, messageID string) error {
	mailbox := ms.GetMailbox(playerID)
	if mailbox == nil {
		return fmt.Errorf("mailbox not found")
	}
	
	if !mailbox.DeleteMessage(messageID) {
		return fmt.Errorf("message not found")
	}
	
	return nil
}

// GetUnreadCount gets unread mail count for a player
func (ms *MailSystem) GetUnreadCount(playerID string) int {
	mailbox := ms.GetMailbox(playerID)
	if mailbox == nil {
		return 0
	}
	return mailbox.GetUnreadCount()
}

// GetTotalUnread gets total unread across all mailboxes
func (ms *MailSystem) GetTotalUnread() int {
	total := 0
	for _, mailbox := range ms.mailboxes {
		total += mailbox.GetUnreadCount()
	}
	return total
}

// CleanupExpired cleans up expired mail in all mailboxes
func (ms *MailSystem) CleanupExpired() int {
	totalRemoved := 0
	for _, mailbox := range ms.mailboxes {
		totalRemoved += mailbox.CleanupExpired()
	}
	return totalRemoved
}

// Save saves all mailboxes
func (ms *MailSystem) Save() error {
	data, err := json.MarshalIndent(ms.mailboxes, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}
	
	if err := os.WriteFile(ms.storagePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	
	return nil
}

// Load loads all mailboxes
func (ms *MailSystem) Load() error {
	data, err := os.ReadFile(ms.storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read: %w", err)
	}
	
	var loaded map[string]*Mailbox
	if err := json.Unmarshal(data, &loaded); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}
	
	ms.mailboxes = loaded
	if ms.mailboxes == nil {
		ms.mailboxes = make(map[string]*Mailbox)
	}
	
	return nil
}

// GetStats returns mail system statistics
func (ms *MailSystem) GetStats() (totalMailboxes, totalMessages, unreadMessages int) {
	totalMailboxes = len(ms.mailboxes)
	
	for _, mailbox := range ms.mailboxes {
		totalMessages += len(mailbox.Messages)
		unreadMessages += mailbox.GetUnreadCount()
	}
	
	return
}
