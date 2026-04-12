package economy

import (
	"fmt"
	"time"
)

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionEarn     TransactionType = "earn"
	TransactionSpend    TransactionType = "spend"
	TransactionTrade    TransactionType = "trade"
	TransactionTax      TransactionType = "tax"
	TransactionLoan     TransactionType = "loan"
	TransactionInterest TransactionType = "interest"
	TransactionPenalty  TransactionType = "penalty"
	TransactionGift     TransactionType = "gift"
	TransactionShop     TransactionType = "shop"
	TransactionAuction  TransactionType = "auction"
	TransactionJob      TransactionType = "job"
	TransactionMining   TransactionType = "mining"
	TransactionCombat   TransactionType = "combat"
	TransactionRefund   TransactionType = "refund"
)

// Transaction represents a single monetary transaction
type Transaction struct {
	ID          string          `json:"id"`
	Type        TransactionType `json:"type"`
	Amount      float64         `json:"amount"`
	From        string          `json:"from"` // Player ID or "SYSTEM"
	To          string          `json:"to"`   // Player ID or "SYSTEM"
	Description string          `json:"description"`
	Timestamp   time.Time       `json:"timestamp"`
	WorldID     string          `json:"world_id,omitempty"`
}

// NewTransaction creates a new transaction
func NewTransaction(txType TransactionType, amount float64, from, to, description string) *Transaction {
	return &Transaction{
		ID:          generateID(),
		Type:        txType,
		Amount:      amount,
		From:        from,
		To:          to,
		Description: description,
		Timestamp:   time.Now(),
	}
}

func generateID() string {
	return fmt.Sprintf("tx_%d", time.Now().UnixNano())
}

// Wallet represents a player's wallet
type Wallet struct {
	PlayerID    string  `json:"player_id"`
	Balance     float64 `json:"balance"`
	BankBalance float64 `json:"bank_balance"`

	// Transaction history
	Transactions []Transaction `json:"transactions,omitempty"`
	maxHistory   int

	// Credit
	CreditScore int     `json:"credit_score"` // 300-850
	LoanLimit   float64 `json:"loan_limit"`

	// Stats
	TotalEarned float64 `json:"total_earned"`
	TotalSpent  float64 `json:"total_spent"`
	TotalTraded float64 `json:"total_traded"`
}

// NewWallet creates a new wallet
func NewWallet(playerID string, startingBalance float64) *Wallet {
	return &Wallet{
		PlayerID:     playerID,
		Balance:      startingBalance,
		BankBalance:  0,
		Transactions: make([]Transaction, 0),
		maxHistory:   100,
		CreditScore:  500, // Average starting score
		LoanLimit:    1000,
		TotalEarned:  startingBalance,
		TotalSpent:   0,
		TotalTraded:  0,
	}
}

// SetMaxHistory sets the maximum transaction history to keep
func (w *Wallet) SetMaxHistory(max int) {
	w.maxHistory = max
	w.trimHistory()
}

// trimHistory removes old transactions if over limit
func (w *Wallet) trimHistory() {
	if len(w.Transactions) > w.maxHistory {
		w.Transactions = w.Transactions[len(w.Transactions)-w.maxHistory:]
	}
}

// Add adds money to wallet
func (w *Wallet) Add(amount float64, txType TransactionType, from, description string) *Transaction {
	if amount <= 0 {
		return nil
	}

	w.Balance += amount
	w.TotalEarned += amount

	tx := NewTransaction(txType, amount, from, w.PlayerID, description)
	w.Transactions = append(w.Transactions, *tx)
	w.trimHistory()

	return tx
}

// Remove removes money from wallet (returns false if insufficient)
func (w *Wallet) Remove(amount float64, txType TransactionType, to, description string) (*Transaction, bool) {
	if amount <= 0 {
		return nil, true
	}

	if w.Balance < amount {
		return nil, false
	}

	w.Balance -= amount
	w.TotalSpent += amount

	tx := NewTransaction(txType, amount, w.PlayerID, to, description)
	w.Transactions = append(w.Transactions, *tx)
	w.trimHistory()

	return tx, true
}

// Transfer transfers money to another wallet (returns false if insufficient)
func (w *Wallet) Transfer(amount float64, target *Wallet, description string) bool {
	if amount <= 0 || w.Balance < amount {
		return false
	}

	// Remove from source
	w.Balance -= amount
	w.TotalTraded += amount

	txOut := NewTransaction(TransactionTrade, amount, w.PlayerID, target.PlayerID, description)
	w.Transactions = append(w.Transactions, *txOut)
	w.trimHistory()

	// Add to target
	target.Balance += amount
	target.TotalTraded += amount

	txIn := NewTransaction(TransactionTrade, amount, w.PlayerID, target.PlayerID, description)
	target.Transactions = append(target.Transactions, *txIn)
	target.trimHistory()

	return true
}

// CanAfford checks if wallet has enough balance
func (w *Wallet) CanAfford(amount float64) bool {
	return w.Balance >= amount
}

// GetBalance returns the total balance (wallet + bank)
func (w *Wallet) GetBalance() float64 {
	return w.Balance + w.BankBalance
}

// DepositToBank moves money to bank
func (w *Wallet) DepositToBank(amount float64) bool {
	if amount <= 0 || w.Balance < amount {
		return false
	}

	w.Balance -= amount
	w.BankBalance += amount

	tx := NewTransaction(TransactionInterest, amount, w.PlayerID, "BANK", "Bank deposit")
	w.Transactions = append(w.Transactions, *tx)
	w.trimHistory()

	return true
}

// WithdrawFromBank moves money from bank to wallet
func (w *Wallet) WithdrawFromBank(amount float64) bool {
	if amount <= 0 || w.BankBalance < amount {
		return false
	}

	w.BankBalance -= amount
	w.Balance += amount

	tx := NewTransaction(TransactionInterest, amount, "BANK", w.PlayerID, "Bank withdrawal")
	w.Transactions = append(w.Transactions, *tx)
	w.trimHistory()

	return true
}

// GetRecentTransactions returns recent transactions
func (w *Wallet) GetRecentTransactions(count int) []Transaction {
	if count > len(w.Transactions) {
		count = len(w.Transactions)
	}

	start := len(w.Transactions) - count
	if start < 0 {
		start = 0
	}

	result := make([]Transaction, count)
	copy(result, w.Transactions[start:])

	return result
}

// GetTransactionsByType returns transactions of a specific type
func (w *Wallet) GetTransactionsByType(txType TransactionType) []Transaction {
	result := make([]Transaction, 0)

	for _, tx := range w.Transactions {
		if tx.Type == txType {
			result = append(result, tx)
		}
	}

	return result
}

// CalculateNetWorth calculates the player's total worth
func (w *Wallet) CalculateNetWorth(inventoryValue float64) float64 {
	return w.Balance + w.BankBalance + inventoryValue
}

// UpdateCreditScore updates the credit score based on activity
func (w *Wallet) UpdateCreditScore() {
	// Base score
	score := 500

	// Increase for consistent play
	if len(w.Transactions) > 50 {
		score += 50
	}

	// Increase for high balance
	if w.BankBalance > 10000 {
		score += 50
	}

	// Decrease for negative behavior (would track separately)

	// Clamp to 300-850 range
	if score < 300 {
		score = 300
	}
	if score > 850 {
		score = 850
	}

	w.CreditScore = score

	// Update loan limit based on credit
	w.LoanLimit = float64(score) * 10 // $3000 to $8500
}

// WalletManager manages all player wallets
type WalletManager struct {
	wallets         map[string]*Wallet
	startingBalance float64

	// Callbacks
	OnTransaction func(tx *Transaction)
}

// NewWalletManager creates a new wallet manager
func NewWalletManager(startingBalance float64) *WalletManager {
	return &WalletManager{
		wallets:         make(map[string]*Wallet),
		startingBalance: startingBalance,
	}
}

// GetOrCreateWallet gets or creates a wallet for a player
func (wm *WalletManager) GetOrCreateWallet(playerID string) *Wallet {
	if wallet, exists := wm.wallets[playerID]; exists {
		return wallet
	}

	wallet := NewWallet(playerID, wm.startingBalance)
	wm.wallets[playerID] = wallet

	return wallet
}

// GetWallet gets a wallet (returns nil if not exists)
func (wm *WalletManager) GetWallet(playerID string) *Wallet {
	return wm.wallets[playerID]
}

// HasWallet checks if a player has a wallet
func (wm *WalletManager) HasWallet(playerID string) bool {
	_, exists := wm.wallets[playerID]
	return exists
}

// Transfer transfers between two players
func (wm *WalletManager) Transfer(fromID, toID string, amount float64, description string) bool {
	fromWallet := wm.GetWallet(fromID)
	toWallet := wm.GetWallet(toID)

	if fromWallet == nil || toWallet == nil {
		return false
	}

	success := fromWallet.Transfer(amount, toWallet, description)

	if success && wm.OnTransaction != nil {
		tx := NewTransaction(TransactionTrade, amount, fromID, toID, description)
		wm.OnTransaction(tx)
	}

	return success
}

// SystemAdd adds money from the system to a player
func (wm *WalletManager) SystemAdd(playerID string, amount float64, txType TransactionType, description string) *Transaction {
	wallet := wm.GetOrCreateWallet(playerID)
	tx := wallet.Add(amount, txType, "SYSTEM", description)

	if tx != nil && wm.OnTransaction != nil {
		wm.OnTransaction(tx)
	}

	return tx
}

// SystemRemove removes money from a player to the system
func (wm *WalletManager) SystemRemove(playerID string, amount float64, txType TransactionType, description string) (*Transaction, bool) {
	wallet := wm.GetWallet(playerID)
	if wallet == nil {
		return nil, false
	}

	tx, success := wallet.Remove(amount, txType, "SYSTEM", description)

	if success && wm.OnTransaction != nil {
		wm.OnTransaction(tx)
	}

	return tx, success
}

// GetTotalMoneyInCirculation returns total money in all wallets
func (wm *WalletManager) GetTotalMoneyInCirculation() float64 {
	total := 0.0
	for _, wallet := range wm.wallets {
		total += wallet.GetBalance()
	}
	return total
}

// GetRichestPlayers returns top N players by balance
func (wm *WalletManager) GetRichestPlayers(count int) []*Wallet {
	// Convert map to slice
	allWallets := make([]*Wallet, 0, len(wm.wallets))
	for _, wallet := range wm.wallets {
		allWallets = append(allWallets, wallet)
	}

	// Simple bubble sort by balance
	for i := 0; i < len(allWallets); i++ {
		for j := i + 1; j < len(allWallets); j++ {
			if allWallets[i].GetBalance() < allWallets[j].GetBalance() {
				allWallets[i], allWallets[j] = allWallets[j], allWallets[i]
			}
		}
	}

	// Return top N
	if count > len(allWallets) {
		count = len(allWallets)
	}
	return allWallets[:count]
}

// GetAverageBalance returns average player balance
func (wm *WalletManager) GetAverageBalance() float64 {
	if len(wm.wallets) == 0 {
		return 0
	}
	return wm.GetTotalMoneyInCirculation() / float64(len(wm.wallets))
}

// MiningReward calculates reward for mining a block
func (wm *WalletManager) MiningReward(blockType string) float64 {
	// Base rewards
	rewards := map[string]float64{
		"coal_ore":    1.0,
		"iron_ore":    2.0,
		"gold_ore":    5.0,
		"diamond_ore": 10.0,
		"emerald_ore": 20.0,
		"stone":       0.1,
		"dirt":        0.05,
	}

	if reward, exists := rewards[blockType]; exists {
		return reward
	}
	return 0
}

// CombatReward calculates reward for killing a mob
func (wm *WalletManager) CombatReward(mobType string) float64 {
	// Base bounties
	bounties := map[string]float64{
		"zombie":   5.0,
		"skeleton": 8.0,
		"creeper":  10.0,
		"spider":   6.0,
		"enderman": 15.0,
		"witch":    20.0,
		"boss":     100.0,
		"player":   50.0,
	}

	if bounty, exists := bounties[mobType]; exists {
		return bounty
	}
	return 0
}
