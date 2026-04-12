package chat

import (
	"fmt"
	"strings"
	"time"
)

// MessageType represents the type of chat message
type MessageType int

const (
	MsgGlobal MessageType = iota
	MsgWhisper
	MsgParty
	MsgGuild
	MsgSystem
	MsgAnnouncement
	MsgLocal
)

// String returns the message type name
func (m MessageType) String() string {
	switch m {
	case MsgGlobal:
		return "global"
	case MsgWhisper:
		return "whisper"
	case MsgParty:
		return "party"
	case MsgGuild:
		return "guild"
	case MsgSystem:
		return "system"
	case MsgAnnouncement:
		return "announcement"
	case MsgLocal:
		return "local"
	}
	return "unknown"
}

// ChatMessage represents a chat message
type ChatMessage struct {
	ID         string      `json:"id"`
	Type       MessageType `json:"type"`
	SenderID   string      `json:"sender_id"`
	SenderName string      `json:"sender_name"`
	Content    string      `json:"content"`
	Timestamp  time.Time   `json:"timestamp"`

	// Target info (for whispers, etc)
	TargetID   string `json:"target_id,omitempty"`
	TargetName string `json:"target_name,omitempty"`

	// Channel info
	Channel string `json:"channel,omitempty"` // Party ID, Guild ID, etc

	// Metadata
	IsDeleted bool `json:"is_deleted"`
}

// NewChatMessage creates a new chat message
func NewChatMessage(msgType MessageType, senderID, senderName, content string) *ChatMessage {
	return &ChatMessage{
		ID:         generateMessageID(),
		Type:       msgType,
		SenderID:   senderID,
		SenderName: senderName,
		Content:    content,
		Timestamp:  time.Now(),
		IsDeleted:  false,
	}
}

func generateMessageID() string {
	return fmt.Sprintf("msg_%d", time.Now().UnixNano())
}

// IsFiltered checks if message contains filtered words
func (m *ChatMessage) IsFiltered(filterWords []string) bool {
	lowerContent := strings.ToLower(m.Content)
	for _, word := range filterWords {
		if strings.Contains(lowerContent, strings.ToLower(word)) {
			return true
		}
	}
	return false
}

// IsSpam checks for spam patterns
func (m *ChatMessage) IsSpam(previousMessages []ChatMessage) bool {
	// Check repeat messages
	recentCount := 0
	for _, msg := range previousMessages {
		if time.Since(msg.Timestamp) < 5*time.Second {
			if msg.SenderID == m.SenderID && msg.Content == m.Content {
				recentCount++
			}
		}
	}
	if recentCount > 2 {
		return true
	}

	// Check caps percentage
	if len(m.Content) > 5 {
		capsCount := 0
		for _, c := range m.Content {
			if c >= 'A' && c <= 'Z' {
				capsCount++
			}
		}
		capsPercent := float64(capsCount) / float64(len(m.Content))
		if capsPercent > 0.7 {
			return true
		}
	}

	return false
}

// Whisper creates a whisper message
func (m *ChatMessage) Whisper(targetID, targetName string) *ChatMessage {
	m.Type = MsgWhisper
	m.TargetID = targetID
	m.TargetName = targetName
	return m
}

// ChatManager manages chat history and channels
type ChatManager struct {
	messages   []ChatMessage
	maxHistory int

	// Muted players
	mutedPlayers map[string]*time.Time // PlayerID -> Mute expiry

	// Filter
	filterWords   []string
	filterEnabled bool

	// Callbacks
	OnMessage func(msg *ChatMessage)
	OnWhisper func(fromID, toID, content string)

	// Command prefix
	cmdPrefix string
}

// NewChatManager creates a new chat manager
func NewChatManager() *ChatManager {
	return &ChatManager{
		messages:      make([]ChatMessage, 0),
		maxHistory:    1000,
		mutedPlayers:  make(map[string]*time.Time),
		filterWords:   make([]string, 0),
		filterEnabled: true,
		cmdPrefix:     "/",
	}
}

// SetMaxHistory sets max message history
func (cm *ChatManager) SetMaxHistory(max int) {
	cm.maxHistory = max
	cm.trimHistory()
}

// trimHistory removes old messages
func (cm *ChatManager) trimHistory() {
	if len(cm.messages) > cm.maxHistory {
		cm.messages = cm.messages[len(cm.messages)-cm.maxHistory:]
	}
}

// IsMuted checks if player is muted
func (cm *ChatManager) IsMuted(playerID string) bool {
	if expiry, exists := cm.mutedPlayers[playerID]; exists {
		if time.Now().After(*expiry) {
			delete(cm.mutedPlayers, playerID)
			return false
		}
		return true
	}
	return false
}

// Mute mutes a player
func (cm *ChatManager) Mute(playerID string, duration time.Duration) {
	expiry := time.Now().Add(duration)
	cm.mutedPlayers[playerID] = &expiry
}

// Unmute unmutes a player
func (cm *ChatManager) Unmute(playerID string) {
	delete(cm.mutedPlayers, playerID)
}

// GetMuteTimeRemaining returns remaining mute time
func (cm *ChatManager) GetMuteTimeRemaining(playerID string) time.Duration {
	if expiry, exists := cm.mutedPlayers[playerID]; exists {
		remaining := time.Until(*expiry)
		if remaining > 0 {
			return remaining
		}
	}
	return 0
}

// AddFilterWord adds a word to the filter
func (cm *ChatManager) AddFilterWord(word string) {
	word = strings.ToLower(word)
	for _, w := range cm.filterWords {
		if w == word {
			return // Already exists
		}
	}
	cm.filterWords = append(cm.filterWords, word)
}

// RemoveFilterWord removes a word from the filter
func (cm *ChatManager) RemoveFilterWord(word string) {
	word = strings.ToLower(word)
	for i, w := range cm.filterWords {
		if w == word {
			cm.filterWords = append(cm.filterWords[:i], cm.filterWords[i+1:]...)
			return
		}
	}
}

// EnableFilter enables chat filter
func (cm *ChatManager) EnableFilter() {
	cm.filterEnabled = true
}

// DisableFilter disables chat filter
func (cm *ChatManager) DisableFilter() {
	cm.filterEnabled = false
}

// SendGlobal sends a global message
func (cm *ChatManager) SendGlobal(senderID, senderName, content string) *ChatMessage {
	// Check mute
	if cm.IsMuted(senderID) {
		return nil
	}

	msg := NewChatMessage(MsgGlobal, senderID, senderName, content)

	// Check filter
	if cm.filterEnabled && msg.IsFiltered(cm.filterWords) {
		msg.Content = cm.censorMessage(content)
	}

	// Check spam
	recent := cm.getRecentMessages(senderID, 10)
	if msg.IsSpam(recent) {
		return nil // Block spam
	}

	cm.addMessage(*msg)

	if cm.OnMessage != nil {
		cm.OnMessage(msg)
	}

	return msg
}

// SendWhisper sends a whisper
func (cm *ChatManager) SendWhisper(senderID, senderName, targetID, targetName, content string) *ChatMessage {
	// Check mute
	if cm.IsMuted(senderID) {
		return nil
	}

	msg := NewChatMessage(MsgWhisper, senderID, senderName, content)
	msg.Whisper(targetID, targetName)

	cm.addMessage(*msg)

	if cm.OnWhisper != nil {
		cm.OnWhisper(senderID, targetID, content)
	}

	return msg
}

// SendSystem sends a system message
func (cm *ChatManager) SendSystem(content string) *ChatMessage {
	msg := NewChatMessage(MsgSystem, "SYSTEM", "System", content)
	cm.addMessage(*msg)

	if cm.OnMessage != nil {
		cm.OnMessage(msg)
	}

	return msg
}

// SendAnnouncement sends an announcement
func (cm *ChatManager) SendAnnouncement(content string) *ChatMessage {
	msg := NewChatMessage(MsgAnnouncement, "SYSTEM", "Announcement", content)
	cm.addMessage(*msg)

	if cm.OnMessage != nil {
		cm.OnMessage(msg)
	}

	return msg
}

// SendParty sends a party message
func (cm *ChatManager) SendParty(senderID, senderName, partyID, content string) *ChatMessage {
	if cm.IsMuted(senderID) {
		return nil
	}

	msg := NewChatMessage(MsgParty, senderID, senderName, content)
	msg.Channel = partyID

	cm.addMessage(*msg)

	if cm.OnMessage != nil {
		cm.OnMessage(msg)
	}

	return msg
}

// SendGuild sends a guild message
func (cm *ChatManager) SendGuild(senderID, senderName, guildID, content string) *ChatMessage {
	if cm.IsMuted(senderID) {
		return nil
	}

	msg := NewChatMessage(MsgGuild, senderID, senderName, content)
	msg.Channel = guildID

	cm.addMessage(*msg)

	if cm.OnMessage != nil {
		cm.OnMessage(msg)
	}

	return msg
}

// addMessage adds a message to history
func (cm *ChatManager) addMessage(msg ChatMessage) {
	cm.messages = append(cm.messages, msg)
	cm.trimHistory()
}

// censorMessage replaces filtered words with ***
func (cm *ChatManager) censorMessage(content string) string {
	censored := content
	for _, word := range cm.filterWords {
		replacement := strings.Repeat("*", len(word))
		censored = strings.ReplaceAll(strings.ToLower(censored), word, replacement)
	}
	return censored
}

// getRecentMessages gets recent messages from a player
func (cm *ChatManager) getRecentMessages(playerID string, count int) []ChatMessage {
	result := make([]ChatMessage, 0)
	for i := len(cm.messages) - 1; i >= 0 && len(result) < count; i-- {
		if cm.messages[i].SenderID == playerID {
			result = append(result, cm.messages[i])
		}
	}
	return result
}

// GetHistory returns chat history
func (cm *ChatManager) GetHistory(count int) []ChatMessage {
	if count > len(cm.messages) {
		count = len(cm.messages)
	}
	start := len(cm.messages) - count
	if start < 0 {
		start = 0
	}

	result := make([]ChatMessage, count)
	copy(result, cm.messages[start:])
	return result
}

// GetHistoryByType returns history filtered by type
func (cm *ChatManager) GetHistoryByType(msgType MessageType, count int) []ChatMessage {
	result := make([]ChatMessage, 0)
	for i := len(cm.messages) - 1; i >= 0 && len(result) < count; i-- {
		if cm.messages[i].Type == msgType {
			result = append(result, cm.messages[i])
		}
	}
	return result
}

// DeleteMessage marks a message as deleted
func (cm *ChatManager) DeleteMessage(messageID string) bool {
	for i, msg := range cm.messages {
		if msg.ID == messageID {
			cm.messages[i].IsDeleted = true
			return true
		}
	}
	return false
}

// GetPlayerMessages returns all messages from a player
func (cm *ChatManager) GetPlayerMessages(playerID string, count int) []ChatMessage {
	result := make([]ChatMessage, 0)
	for i := len(cm.messages) - 1; i >= 0 && len(result) < count; i-- {
		if cm.messages[i].SenderID == playerID {
			result = append(result, cm.messages[i])
		}
	}
	return result
}

// ParseCommand parses a chat command
func (cm *ChatManager) ParseCommand(content string) (cmd string, args []string, isCmd bool) {
	if !strings.HasPrefix(content, cm.cmdPrefix) {
		return "", nil, false
	}

	// Remove prefix and split
	withoutPrefix := strings.TrimPrefix(content, cm.cmdPrefix)
	parts := strings.Fields(withoutPrefix)

	if len(parts) == 0 {
		return "", nil, false
	}

	return strings.ToLower(parts[0]), parts[1:], true
}

// ChatCommands returns available chat commands
func (cm *ChatManager) ChatCommands() map[string]string {
	return map[string]string{
		"msg":      "/msg <player> <message> - Send private message",
		"reply":    "/reply <message> - Reply to last whisper",
		"party":    "/party <message> - Send to party",
		"guild":    "/guild <message> - Send to guild",
		"ignore":   "/ignore <player> - Ignore player",
		"unignore": "/unignore <player> - Stop ignoring",
	}
}
