package warps

import (
	"fmt"
	"time"
)

// TPARequest represents a teleport request
type TPARequest struct {
	ID          string    `json:"id"`
	FromID      string    `json:"from_id"`
	FromName    string    `json:"from_name"`
	ToID        string    `json:"to_id"`
	ToName      string    `json:"to_name"`
	
	// Type: true = request to teleport to someone, false = request someone to teleport to you
	IsHere      bool      `json:"is_here"` // true = /tpahere, false = /tpa
	
	// Status
	Status      TPAStatus `json:"status"`
	
	// Timing
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	AcceptedAt  *time.Time `json:"accepted_at,omitempty"`
}

// TPAStatus represents request state
type TPAStatus int

const (
	TPAPending TPAStatus = iota
	TPAAccepted
	TPADenied
	TPAExpired
	TPACancelled
)

// NewTPARequest creates a new TPA request
func NewTPARequest(id, fromID, fromName, toID, toName string, isHere bool) *TPARequest {
	return &TPARequest{
		ID:        id,
		FromID:    fromID,
		FromName:  fromName,
		ToID:      toID,
		ToName:    toName,
		IsHere:    isHere,
		Status:    TPAPending,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(30 * time.Second),
	}
}

// Accept accepts the request
func (t *TPARequest) Accept() {
	if t.Status == TPAPending {
		t.Status = TPAAccepted
		now := time.Now()
		t.AcceptedAt = &now
	}
}

// Deny denies the request
func (t *TPARequest) Deny() {
	if t.Status == TPAPending {
		t.Status = TPADenied
	}
}

// Cancel cancels the request
func (t *TPARequest) Cancel() {
	if t.Status == TPAPending {
		t.Status = TPACancelled
	}
}

// IsExpired checks if request expired
func (t *TPARequest) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// GetTeleportPlayers returns who teleports and destination
func (t *TPARequest) GetTeleportPlayers() (teleporter, destination string) {
	if t.IsHere {
		// /tpahere: target teleports to sender
		return t.ToID, t.FromID
	}
	// /tpa: sender teleports to target
	return t.FromID, t.ToID
}

// TPAManager manages teleport requests
type TPAManager struct {
	requests    map[string]*TPARequest
	byPlayer    map[string]string // PlayerID -> RequestID (outgoing)
	
	requestCounter int
	cooldowns   map[string]time.Time // PlayerID -> Last request time
	cooldownDuration time.Duration
}

// NewTPAManager creates new TPA manager
func NewTPAManager() *TPAManager {
	return &TPAManager{
		requests:         make(map[string]*TPARequest),
		byPlayer:         make(map[string]string),
		requestCounter:   0,
		cooldowns:        make(map[string]time.Time),
		cooldownDuration: 3 * time.Second,
	}
}

// RequestTPA creates a TPA request
func (tm *TPAManager) RequestTPA(fromID, fromName, toID, toName string, isHere bool) (*TPARequest, error) {
	// Check cooldown
	if lastRequest, exists := tm.cooldowns[fromID]; exists {
		if time.Since(lastRequest) < tm.cooldownDuration {
			return nil, fmt.Errorf("request cooldown active")
		}
	}
	
	// Check if already has pending request
	if _, exists := tm.byPlayer[fromID]; exists {
		return nil, fmt.Errorf("already have pending request")
	}
	
	// Check if target has request from this player
	for _, req := range tm.requests {
		if req.FromID == fromID && req.ToID == toID && req.Status == TPAPending {
			return nil, fmt.Errorf("already requested to this player")
		}
	}
	
	tm.requestCounter++
	requestID := fmt.Sprintf("tpa_%d_%d", tm.requestCounter, time.Now().Unix())
	
	request := NewTPARequest(requestID, fromID, fromName, toID, toName, isHere)
	
	tm.requests[requestID] = request
	tm.byPlayer[fromID] = requestID
	tm.cooldowns[fromID] = time.Now()
	
	return request, nil
}

// GetRequest gets a request
func (tm *TPAManager) GetRequest(requestID string) (*TPARequest, bool) {
	request, exists := tm.requests[requestID]
	return request, exists
}

// GetRequestToPlayer gets pending request to a player
func (tm *TPAManager) GetRequestToPlayer(playerID string) *TPARequest {
	for _, req := range tm.requests {
		if req.ToID == playerID && req.Status == TPAPending {
			return req
		}
	}
	return nil
}

// AcceptRequest accepts a TPA
func (tm *TPAManager) AcceptRequest(requestID, playerID string) error {
	request, exists := tm.GetRequest(requestID)
	if !exists {
		return fmt.Errorf("request not found")
	}
	
	if request.ToID != playerID {
		return fmt.Errorf("not your request")
	}
	
	if request.Status != TPAPending {
		return fmt.Errorf("request not pending")
	}
	
	request.Accept()
	tm.cleanupRequest(requestID)
	
	return nil
}

// DenyRequest denies a TPA
func (tm *TPAManager) DenyRequest(requestID, playerID string) error {
	request, exists := tm.GetRequest(requestID)
	if !exists {
		return fmt.Errorf("request not found")
	}
	
	if request.ToID != playerID {
		return fmt.Errorf("not your request")
	}
	
	request.Deny()
	tm.cleanupRequest(requestID)
	
	return nil
}

// CancelRequest cancels outgoing TPA
func (tm *TPAManager) CancelRequest(playerID string) error {
	requestID, exists := tm.byPlayer[playerID]
	if !exists {
		return fmt.Errorf("no pending request")
	}
	
	request, exists := tm.GetRequest(requestID)
	if !exists {
		return fmt.Errorf("request not found")
	}
	
	request.Cancel()
	tm.cleanupRequest(requestID)
	
	return nil
}

// cleanupRequest removes a request
func (tm *TPAManager) cleanupRequest(requestID string) {
	request, exists := tm.requests[requestID]
	if !exists {
		return
	}
	
	delete(tm.byPlayer, request.FromID)
	delete(tm.requests, requestID)
}

// Update processes expired requests
func (tm *TPAManager) Update() {
	for id, request := range tm.requests {
		if request.Status == TPAPending && request.IsExpired() {
			request.Status = TPAExpired
			tm.cleanupRequest(id)
		}
	}
}

// HasPendingRequest checks if player has pending outgoing request
func (tm *TPAManager) HasPendingRequest(playerID string) bool {
	_, exists := tm.byPlayer[playerID]
	return exists
}

// ToggleTP blocks/unblocks TPA from player
func (tm *TPAManager) ToggleTP(playerID, targetID string, block bool) {
	// In real implementation, store in player preferences
	_ = playerID
	_ = targetID
	_ = block
}
