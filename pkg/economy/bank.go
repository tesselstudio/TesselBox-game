package economy

import (
	"fmt"
	"time"
)

// LoanStatus represents the status of a loan
type LoanStatus int

const (
	LoanActive LoanStatus = iota
	LoanPaid
	LoanDefaulted
	LoanRepossessing
)

// Loan represents a player loan
type Loan struct {
	ID           string  `json:"id"`
	PlayerID     string  `json:"player_id"`
	Principal    float64 `json:"principal"`
	InterestRate float64 `json:"interest_rate"` // Daily rate
	TotalOwed    float64 `json:"total_owed"`
	Paid         float64 `json:"paid"`
	Balance      float64 `json:"balance"` // Remaining to pay

	IssuedAt    time.Time `json:"issued_at"`
	DueDate     time.Time `json:"due_date"`
	LastPayment time.Time `json:"last_payment"`

	Status     LoanStatus `json:"status"`
	Collateral []string   `json:"collateral,omitempty"` // Item IDs held as collateral

	// Penalties
	LateFees    float64    `json:"late_fees"`
	DefaultedAt *time.Time `json:"defaulted_at,omitempty"`
}

// NewLoan creates a new loan
func NewLoan(id, playerID string, principal, interestRate float64, duration time.Duration, collateral []string) *Loan {
	now := time.Now()

	// Calculate total owed with compound interest
	days := int(duration.Hours() / 24)
	totalOwed := principal
	for i := 0; i < days; i++ {
		totalOwed *= (1 + interestRate)
	}

	return &Loan{
		ID:           id,
		PlayerID:     playerID,
		Principal:    principal,
		InterestRate: interestRate,
		TotalOwed:    totalOwed,
		Paid:         0,
		Balance:      totalOwed,
		IssuedAt:     now,
		DueDate:      now.Add(duration),
		LastPayment:  now,
		Status:       LoanActive,
		Collateral:   collateral,
		LateFees:     0,
	}
}

// CalculateDailyInterest calculates interest for today
func (l *Loan) CalculateDailyInterest() float64 {
	return l.Balance * l.InterestRate
}

// ApplyDailyInterest applies daily interest (call once per day)
func (l *Loan) ApplyDailyInterest() {
	if l.Status != LoanActive {
		return
	}

	interest := l.CalculateDailyInterest()
	l.TotalOwed += interest
	l.Balance += interest
}

// MakePayment makes a payment on the loan
func (l *Loan) MakePayment(amount float64) (paidOff bool, remaining float64, err error) {
	if l.Status != LoanActive {
		return false, 0, fmt.Errorf("loan is not active")
	}

	if amount <= 0 {
		return false, l.Balance, fmt.Errorf("payment must be positive")
	}

	l.Paid += amount
	l.Balance -= amount
	l.LastPayment = time.Now()

	if l.Balance <= 0 {
		l.Status = LoanPaid
		l.Balance = 0
		return true, 0, nil
	}

	return false, l.Balance, nil
}

// IsOverdue checks if loan is past due
func (l *Loan) IsOverdue() bool {
	if l.Status != LoanActive {
		return false
	}
	return time.Now().After(l.DueDate)
}

// DaysOverdue returns days overdue (0 if not overdue)
func (l *Loan) DaysOverdue() int {
	if !l.IsOverdue() {
		return 0
	}
	return int(time.Since(l.DueDate).Hours() / 24)
}

// ApplyLateFees applies late fees (call daily when overdue)
func (l *Loan) ApplyLateFees() {
	if !l.IsOverdue() {
		return
	}

	// 1% late fee per day overdue
	lateFee := l.Balance * 0.01
	l.LateFees += lateFee
	l.TotalOwed += lateFee
	l.Balance += lateFee
}

// MarkDefaulted marks loan as defaulted
func (l *Loan) MarkDefaulted() {
	if l.Status == LoanActive && l.IsOverdue() {
		now := time.Now()
		l.Status = LoanDefaulted
		l.DefaultedAt = &now
	}
}

// TimeRemaining returns time until due
func (l *Loan) TimeRemaining() time.Duration {
	return time.Until(l.DueDate)
}

// BankAccount represents a player's bank account
type BankAccount struct {
	PlayerID     string    `json:"player_id"`
	Balance      float64   `json:"balance"`
	SavingsRate  float64   `json:"savings_rate"` // Daily interest rate
	LastInterest time.Time `json:"last_interest"`

	// Safety deposit
	SafetyDeposit   []SafetyDepositItem `json:"safety_deposit,omitempty"`
	MaxDepositSlots int                 `json:"max_deposit_slots"`

	// Credit
	CreditScore   int     `json:"credit_score"` // 300-850
	CreditLimit   float64 `json:"credit_limit"`
	TotalBorrowed float64 `json:"total_borrowed"`
	TotalRepaid   float64 `json:"total_repaid"`

	// Transaction history
	History []BankTransaction `json:"history,omitempty"`
}

// SafetyDepositItem represents an item in safety deposit
type SafetyDepositItem struct {
	Slot     int       `json:"slot"`
	ItemType string    `json:"item_type"`
	Quantity int       `json:"quantity"`
	StoredAt time.Time `json:"stored_at"`
}

// BankTransaction represents a bank transaction
type BankTransaction struct {
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"` // "deposit", "withdraw", "interest", "fee"
	Amount    float64   `json:"amount"`
	Balance   float64   `json:"balance"`
}

// NewBankAccount creates a new bank account
func NewBankAccount(playerID string) *BankAccount {
	return &BankAccount{
		PlayerID:        playerID,
		Balance:         0,
		SavingsRate:     0.02 / 365, // 2% APR daily
		LastInterest:    time.Now(),
		SafetyDeposit:   make([]SafetyDepositItem, 0),
		MaxDepositSlots: 8,
		CreditScore:     500, // Average
		CreditLimit:     1000,
		History:         make([]BankTransaction, 0),
	}
}

// Deposit deposits money into account
func (ba *BankAccount) Deposit(amount float64, wallet *Wallet) bool {
	if wallet == nil {
		return false
	}

	_, success := wallet.Remove(amount, TransactionSpend, "BANK", "Bank deposit")
	if !success {
		return false
	}

	ba.Balance += amount
	ba.recordTransaction("deposit", amount)
	return true
}

// Withdraw withdraws money to wallet
func (ba *BankAccount) Withdraw(amount float64, wallet *Wallet) bool {
	if ba.Balance < amount {
		return false
	}

	ba.Balance -= amount
	wallet.Add(amount, TransactionEarn, "BANK", "Bank withdrawal")
	ba.recordTransaction("withdraw", -amount)
	return true
}

// ApplyInterest applies daily interest to savings
func (ba *BankAccount) ApplyInterest() {
	now := time.Now()

	// Only apply once per day
	if now.Sub(ba.LastInterest) < 24*time.Hour {
		return
	}

	interest := ba.Balance * ba.SavingsRate
	if interest > 0 {
		ba.Balance += interest
		ba.recordTransaction("interest", interest)
	}

	ba.LastInterest = now
}

// CanAffordLoan checks if player can afford loan payments
func (ba *BankAccount) CanAffordLoan(amount float64, duration time.Duration) bool {
	// Simple check: can they pay back within duration based on credit score
	minCreditScore := 500
	if ba.CreditScore < minCreditScore {
		return false
	}

	// Check if amount is within credit limit
	if amount > ba.CreditLimit {
		return false
	}

	return true
}

// recordTransaction records a transaction
func (ba *BankAccount) recordTransaction(txType string, amount float64) {
	ba.History = append(ba.History, BankTransaction{
		Timestamp: time.Now(),
		Type:      txType,
		Amount:    amount,
		Balance:   ba.Balance,
	})

	// Keep last 50 transactions
	if len(ba.History) > 50 {
		ba.History = ba.History[len(ba.History)-50:]
	}
}

// UpdateCreditScore updates credit score based on loan history
func (ba *BankAccount) UpdateCreditScore(loans []Loan) {
	baseScore := 500

	// Increase for good repayment
	goodLoans := 0
	totalPaid := 0.0
	for _, loan := range loans {
		totalPaid += loan.Paid
		if loan.Status == LoanPaid {
			goodLoans++
		}
	}

	baseScore += goodLoans * 20
	baseScore += int(totalPaid / 1000) // +1 per 1000 repaid

	// Decrease for defaults
	defaults := 0
	for _, loan := range loans {
		if loan.Status == LoanDefaulted {
			defaults++
		}
	}
	baseScore -= defaults * 100

	// Clamp to 300-850
	if baseScore < 300 {
		baseScore = 300
	}
	if baseScore > 850 {
		baseScore = 850
	}

	ba.CreditScore = baseScore
	ba.CreditLimit = float64(baseScore) * 10 // $3000 to $8500
}

// Bank manages all bank accounts and loans
type Bank struct {
	accounts map[string]*BankAccount
	loans    map[string]*Loan
	byPlayer map[string][]string // PlayerID -> Loan IDs

	loanCounter int

	// Settings
	MinCreditScore   int
	BaseInterestRate float64
	MaxLoanDuration  time.Duration

	// Economy ref for integration
	walletMgr *WalletManager
}

// NewBank creates a new bank
func NewBank(walletMgr *WalletManager) *Bank {
	return &Bank{
		accounts:         make(map[string]*BankAccount),
		loans:            make(map[string]*Loan),
		byPlayer:         make(map[string][]string),
		loanCounter:      0,
		MinCreditScore:   500,
		BaseInterestRate: 0.05 / 365,          // 5% APR
		MaxLoanDuration:  30 * 24 * time.Hour, // 30 days
		walletMgr:        walletMgr,
	}
}

// GetOrCreateAccount gets or creates a bank account
func (b *Bank) GetOrCreateAccount(playerID string) *BankAccount {
	if account, exists := b.accounts[playerID]; exists {
		return account
	}

	account := NewBankAccount(playerID)
	b.accounts[playerID] = account
	return account
}

// GetAccount gets an account (nil if not exists)
func (b *Bank) GetAccount(playerID string) *BankAccount {
	return b.accounts[playerID]
}

// ApplyDailyProcesses applies interest to all accounts and loans
func (b *Bank) ApplyDailyProcesses() {
	// Apply savings interest
	for _, account := range b.accounts {
		account.ApplyInterest()
	}
	_ = time.Now() // Satisfy linter

	// Apply loan interest and check defaults
	for _, loan := range b.loans {
		if loan.Status == LoanActive {
			// Apply daily interest
			loan.ApplyDailyInterest()

			// Check for default (30 days overdue)
			if loan.IsOverdue() && loan.DaysOverdue() > 30 {
				loan.MarkDefaulted()
			} else if loan.IsOverdue() {
				// Apply late fees
				loan.ApplyLateFees()
			}
		}
	}
}

// IssueLoan issues a new loan
func (b *Bank) IssueLoan(playerID string, amount float64, duration time.Duration, collateral []string) (*Loan, error) {
	// Validate
	if duration > b.MaxLoanDuration {
		return nil, fmt.Errorf("loan duration too long (max %s)", b.MaxLoanDuration)
	}

	account := b.GetOrCreateAccount(playerID)

	// Check credit
	if account.CreditScore < b.MinCreditScore {
		return nil, fmt.Errorf("credit score too low (need %d, have %d)", b.MinCreditScore, account.CreditScore)
	}

	// Check existing loans
	existingLoans := b.GetPlayerLoans(playerID)
	activeLoanAmount := 0.0
	for _, loan := range existingLoans {
		if loan.Status == LoanActive {
			activeLoanAmount += loan.Balance
		}
	}

	if activeLoanAmount+amount > account.CreditLimit {
		return nil, fmt.Errorf("exceeds credit limit")
	}

	// Create loan
	b.loanCounter++
	loanID := fmt.Sprintf("loan_%d_%d", b.loanCounter, time.Now().Unix())

	// Adjust interest rate based on credit score
	interestRate := b.BaseInterestRate
	if account.CreditScore > 700 {
		interestRate *= 0.8 // Good credit = lower rate
	} else if account.CreditScore < 600 {
		interestRate *= 1.5 // Poor credit = higher rate
	}

	loan := NewLoan(loanID, playerID, amount, interestRate, duration, collateral)

	// Disburse funds
	wallet := b.walletMgr.GetOrCreateWallet(playerID)
	wallet.Add(amount, TransactionLoan, "BANK", "Loan disbursement")

	// Register loan
	b.loans[loanID] = loan
	b.byPlayer[playerID] = append(b.byPlayer[playerID], loanID)

	account.TotalBorrowed += amount

	return loan, nil
}

// GetLoan gets a loan by ID
func (b *Bank) GetLoan(loanID string) (*Loan, bool) {
	loan, exists := b.loans[loanID]
	return loan, exists
}

// GetPlayerLoans returns all loans for a player
func (b *Bank) GetPlayerLoans(playerID string) []Loan {
	loanIDs := b.byPlayer[playerID]
	loans := make([]Loan, 0, len(loanIDs))

	for _, id := range loanIDs {
		if loan, exists := b.loans[id]; exists {
			loans = append(loans, *loan)
		}
	}

	return loans
}

// GetActiveLoans returns active loans for a player
func (b *Bank) GetActiveLoans(playerID string) []Loan {
	allLoans := b.GetPlayerLoans(playerID)
	active := make([]Loan, 0)

	for _, loan := range allLoans {
		if loan.Status == LoanActive {
			active = append(active, loan)
		}
	}

	return active
}

// MakePayment makes a payment on a loan
func (b *Bank) MakePayment(loanID string, amount float64) (bool, float64, error) {
	loan, exists := b.GetLoan(loanID)
	if !exists {
		return false, 0, fmt.Errorf("loan not found")
	}

	wallet := b.walletMgr.GetWallet(loan.PlayerID)
	if wallet == nil || !wallet.CanAfford(amount) {
		return false, 0, fmt.Errorf("insufficient funds")
	}

	// Take payment from wallet
	_, success := wallet.Remove(amount, TransactionSpend, "BANK", fmt.Sprintf("Loan payment %s", loanID))
	if !success {
		return false, 0, fmt.Errorf("payment failed")
	}

	// Apply to loan
	paidOff, remaining, err := loan.MakePayment(amount)
	if err != nil {
		// Refund if loan error
		wallet.Add(amount, TransactionRefund, "BANK", "Loan payment refund")
		return false, 0, err
	}

	// Update account stats
	account := b.GetOrCreateAccount(loan.PlayerID)
	account.TotalRepaid += amount

	return paidOff, remaining, nil
}

// GetTotalDeposits returns total bank deposits
func (b *Bank) GetTotalDeposits() float64 {
	total := 0.0
	for _, account := range b.accounts {
		total += account.Balance
	}
	return total
}

// GetTotalOutstandingLoans returns total outstanding loan amounts
func (b *Bank) GetTotalOutstandingLoans() float64 {
	total := 0.0
	for _, loan := range b.loans {
		if loan.Status == LoanActive {
			total += loan.Balance
		}
	}
	return total
}
