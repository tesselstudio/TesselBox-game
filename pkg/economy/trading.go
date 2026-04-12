package economy

import (
	"fmt"
	"time"

	"tesselbox/pkg/items"
)

// TradeStatus represents the status of a trade
type TradeStatus int

const (
	TradePending TradeStatus = iota
	TradeAccepted
	TradeConfirmed
	TradeCompleted
	TradeCancelled
	TradeExpired
)

// TradeOffer represents what a player is offering
type TradeOffer struct {
	PlayerID    string        `json:"player_id"`
	Money       float64       `json:"money"`
	Items       []items.Item  `json:"items"`
	Confirmed   bool          `json:"confirmed"`
}

// TradeSession represents a trade between two players
type TradeSession struct {
	ID            string      `json:"id"`
	InitiatorID   string      `json:"initiator_id"`
	PartnerID     string      `json:"partner_id"`
	
	// Offers
	InitiatorOffer TradeOffer  `json:"initiator_offer"`
	PartnerOffer   TradeOffer  `json:"partner_offer"`
	
	// Status
	Status        TradeStatus `json:"status"`
	
	// Timing
	CreatedAt     time.Time   `json:"created_at"`
	ExpiresAt     time.Time   `json:"expires_at"`
	CompletedAt   *time.Time  `json:"completed_at,omitempty"`
	
	// Security
	ConfirmTimer  int         `json:"confirm_timer"` // Seconds countdown
}

// NewTradeSession creates a new trade session
func NewTradeSession(id, initiatorID, partnerID string) *TradeSession {
	now := time.Now()
	return &TradeSession{
		ID:           id,
		InitiatorID:  initiatorID,
		PartnerID:    partnerID,
		InitiatorOffer: TradeOffer{
			PlayerID: initiatorID,
			Items:    make([]items.Item, 0),
		},
		PartnerOffer: TradeOffer{
			PlayerID: partnerID,
			Items:    make([]items.Item, 0),
		},
		Status:       TradePending,
		CreatedAt:     now,
		ExpiresAt:     now.Add(5 * time.Minute), // 5 minute timeout
		ConfirmTimer:  5,
	}
}

// AddItem adds an item to a player's offer
func (t *TradeSession) AddItem(playerID string, item items.Item) error {
	if t.Status != TradePending {
		return fmt.Errorf("trade is not pending")
	}
	
	if playerID == t.InitiatorID {
		t.InitiatorOffer.Items = append(t.InitiatorOffer.Items, item)
	} else if playerID == t.PartnerID {
		t.PartnerOffer.Items = append(t.PartnerOffer.Items, item)
	} else {
		return fmt.Errorf("not part of this trade")
	}
	
	return nil
}

// RemoveItem removes an item from a player's offer
func (t *TradeSession) RemoveItem(playerID string, index int) error {
	if t.Status != TradePending {
		return fmt.Errorf("trade is not pending")
	}
	
	var items *[]items.Item
	if playerID == t.InitiatorID {
		items = &t.InitiatorOffer.Items
	} else if playerID == t.PartnerID {
		items = &t.PartnerOffer.Items
	} else {
		return fmt.Errorf("not part of this trade")
	}
	
	if index < 0 || index >= len(*items) {
		return fmt.Errorf("invalid item index")
	}
	
	*items = append((*items)[:index], (*items)[index+1:]...)
	return nil
}

// SetMoney sets money in a player's offer
func (t *TradeSession) SetMoney(playerID string, amount float64) error {
	if t.Status != TradePending {
		return fmt.Errorf("trade is not pending")
	}
	
	if playerID == t.InitiatorID {
		t.InitiatorOffer.Money = amount
	} else if playerID == t.PartnerID {
		t.PartnerOffer.Money = amount
	} else {
		return fmt.Errorf("not part of this trade")
	}
	
	return nil
}

// Accept accepts the trade offer
func (t *TradeSession) Accept(playerID string) error {
	if t.Status != TradePending {
		return fmt.Errorf("trade is not pending")
	}
	
	if playerID == t.InitiatorID {
		t.InitiatorOffer.Confirmed = true
	} else if playerID == t.PartnerID {
		t.PartnerOffer.Confirmed = true
	} else {
		return fmt.Errorf("not part of this trade")
	}
	
	// Check if both accepted
	if t.InitiatorOffer.Confirmed && t.PartnerOffer.Confirmed {
		t.Status = TradeAccepted
		t.ConfirmTimer = 5 // Reset countdown
	}
	
	return nil
}

// Unaccept unaccepts the trade
func (t *TradeSession) Unaccept(playerID string) error {
	if t.Status != TradePending && t.Status != TradeAccepted {
		return fmt.Errorf("cannot unaccept at this stage")
	}
	
	if playerID == t.InitiatorID {
		t.InitiatorOffer.Confirmed = false
	} else if playerID == t.PartnerID {
		t.PartnerOffer.Confirmed = false
	} else {
		return fmt.Errorf("not part of this trade")
	}
	
	// If either unaccepts, go back to pending
	t.Status = TradePending
	
	return nil
}

// FinalConfirm final confirmation after countdown
func (t *TradeSession) FinalConfirm(playerID string) error {
	if t.Status != TradeAccepted {
		return fmt.Errorf("trade not ready for confirmation")
	}
	
	if playerID == t.InitiatorID {
		if t.InitiatorOffer.Confirmed {
			return fmt.Errorf("already confirmed")
		}
		t.InitiatorOffer.Confirmed = true
	} else if playerID == t.PartnerID {
		if t.PartnerOffer.Confirmed {
			return fmt.Errorf("already confirmed")
		}
		t.PartnerOffer.Confirmed = true
	} else {
		return fmt.Errorf("not part of this trade")
	}
	
	return nil
}

// Complete completes the trade
func (t *TradeSession) Complete() error {
	if t.Status == TradeCompleted {
		return fmt.Errorf("trade already completed")
	}
	
	if !t.InitiatorOffer.Confirmed || !t.PartnerOffer.Confirmed {
		return fmt.Errorf("both parties must confirm")
	}
	
	now := time.Now()
	t.Status = TradeCompleted
	t.CompletedAt = &now
	
	return nil
}

// Cancel cancels the trade
func (t *TradeSession) Cancel(playerID string) {
	if t.Status == TradeCompleted {
		return
	}
	
	t.Status = TradeCancelled
}

// IsExpired checks if trade has expired
func (t *TradeSession) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// CanModify checks if player can still modify offer
func (t *TradeSession) CanModify() bool {
	return t.Status == TradePending
}

// GetTotalItems gets total items in trade
func (t *TradeSession) GetTotalItems() int {
	return len(t.InitiatorOffer.Items) + len(t.PartnerOffer.Items)
}

// GetTotalMoney gets total money in trade
func (t *TradeSession) GetTotalMoney() float64 {
	return t.InitiatorOffer.Money + t.PartnerOffer.Money
}

// TradingSystem manages trade sessions
type TradingSystem struct {
	sessions      map[string]*TradeSession
	byPlayer      map[string]string // PlayerID -> SessionID
	
	sessionCounter int
	walletMgr      *WalletManager
}

// NewTradingSystem creates a new trading system
func NewTradingSystem(walletMgr *WalletManager) *TradingSystem {
	return &TradingSystem{
		sessions:       make(map[string]*TradeSession),
		byPlayer:       make(map[string]string),
		sessionCounter: 0,
		walletMgr:      walletMgr,
	}
}

// CreateTrade creates a new trade session
func (ts *TradingSystem) CreateTrade(initiatorID, partnerID string) (*TradeSession, error) {
	// Check if either is already trading
	if ts.IsTrading(initiatorID) {
		return nil, fmt.Errorf("already in a trade")
	}
	if ts.IsTrading(partnerID) {
		return nil, fmt.Errorf("partner is already in a trade")
	}
	
	ts.sessionCounter++
	sessionID := fmt.Sprintf("trade_%d_%d", ts.sessionCounter, time.Now().Unix())
	
	session := NewTradeSession(sessionID, initiatorID, partnerID)
	
	ts.sessions[sessionID] = session
	ts.byPlayer[initiatorID] = sessionID
	ts.byPlayer[partnerID] = sessionID
	
	return session, nil
}

// GetTrade gets a trade session
func (ts *TradingSystem) GetTrade(sessionID string) (*TradeSession, bool) {
	session, exists := ts.sessions[sessionID]
	return session, exists
}

// GetPlayerTrade gets a player's current trade
func (ts *TradingSystem) GetPlayerTrade(playerID string) (*TradeSession, bool) {
	sessionID, exists := ts.byPlayer[playerID]
	if !exists {
		return nil, false
	}
	return ts.GetTrade(sessionID)
}

// IsTrading checks if player is in a trade
func (ts *TradingSystem) IsTrading(playerID string) bool {
	_, exists := ts.byPlayer[playerID]
	return exists
}

// ExecuteTrade executes a completed trade
func (ts *TradingSystem) ExecuteTrade(sessionID string) error {
	session, exists := ts.GetTrade(sessionID)
	if !exists {
		return fmt.Errorf("trade not found")
	}
	
	if err := session.Complete(); err != nil {
		return err
	}
	
	// Transfer money
	if session.InitiatorOffer.Money > 0 || session.PartnerOffer.Money > 0 {
		initiatorWallet := ts.walletMgr.GetWallet(session.InitiatorID)
		partnerWallet := ts.walletMgr.GetWallet(session.PartnerID)
		
		if initiatorWallet == nil || partnerWallet == nil {
			return fmt.Errorf("wallet not found")
		}
		
		// Check funds
		if !initiatorWallet.CanAfford(session.InitiatorOffer.Money) {
			return fmt.Errorf("initiator cannot afford trade")
		}
		if !partnerWallet.CanAfford(session.PartnerOffer.Money) {
			return fmt.Errorf("partner cannot afford trade")
		}
		
		// Transfer from initiator to partner
		if session.InitiatorOffer.Money > 0 {
			_, success := initiatorWallet.Remove(session.InitiatorOffer.Money, TransactionTrade, session.PartnerID, "Trade payment")
			if !success {
				return fmt.Errorf("payment failed")
			}
			partnerWallet.Add(session.InitiatorOffer.Money, TransactionTrade, session.InitiatorID, "Trade received")
		}
		
		// Transfer from partner to initiator
		if session.PartnerOffer.Money > 0 {
			_, success := partnerWallet.Remove(session.PartnerOffer.Money, TransactionTrade, session.InitiatorID, "Trade payment")
			if !success {
				// Rollback first transfer (in real implementation)
				return fmt.Errorf("payment failed")
			}
			initiatorWallet.Add(session.PartnerOffer.Money, TransactionTrade, session.PartnerID, "Trade received")
		}
	}
	
	// Transfer items would happen here in real implementation
	
	// Cleanup
	ts.cleanupTrade(sessionID)
	
	return nil
}

// CancelTrade cancels a trade
func (ts *TradingSystem) CancelTrade(sessionID, playerID string) error {
	session, exists := ts.GetTrade(sessionID)
	if !exists {
		return fmt.Errorf("trade not found")
	}
	
	if session.InitiatorID != playerID && session.PartnerID != playerID {
		return fmt.Errorf("not part of this trade")
	}
	
	session.Cancel(playerID)
	ts.cleanupTrade(sessionID)
	
	return nil
}

// cleanupTrade removes a trade session
func (ts *TradingSystem) cleanupTrade(sessionID string) {
	session, exists := ts.sessions[sessionID]
	if !exists {
		return
	}
	
	delete(ts.byPlayer, session.InitiatorID)
	delete(ts.byPlayer, session.PartnerID)
	delete(ts.sessions, sessionID)
}

// Update updates all active trades
func (ts *TradingSystem) Update() {
	now := time.Now()
	
	for id, session := range ts.sessions {
		// Check expiration
		if session.IsExpired() {
			session.Status = TradeExpired
			ts.cleanupTrade(id)
			continue
		}
		
		// Handle confirmation countdown
		if session.Status == TradeAccepted && session.ConfirmTimer > 0 {
			// In real implementation, countdown every second
			// When timer hits 0, both need to re-confirm
			_ = now
		}
	}
}

// GetActiveTrades returns active trade count
func (ts *TradingSystem) GetActiveTrades() int {
	return len(ts.sessions)
}
