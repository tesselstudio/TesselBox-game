package economy

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"
)

// CompanyType represents the type of company
type CompanyType int

const (
	CompanyShop CompanyType = iota
	CompanyGuild
	CompanyMinigame
)

// String returns company type name
func (c CompanyType) String() string {
	switch c {
	case CompanyShop:
		return "Shop"
	case CompanyGuild:
		return "Guild"
	case CompanyMinigame:
		return "Minigame"
	}
	return "Unknown"
}

// Company represents a company in the stock market
type Company struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	OwnerID     string      `json:"owner_id"`
	Type        CompanyType `json:"type"`
	
	// Shares
	TotalShares int         `json:"total_shares"`
	SharePrice  float64     `json:"share_price"`
	
	// Financials
	Revenue     float64     `json:"revenue"`
	Expenses    float64     `json:"expenses"`
	Profit      float64     `json:"profit"`
	Assets      float64     `json:"assets"`
	
	// History
	PriceHistory []StockPrice `json:"price_history,omitempty"`
	
	// Market
	Open        bool        `json:"open"`
	LastUpdated time.Time   `json:"last_updated"`
}

// StockPrice represents a price point in history
type StockPrice struct {
	Price     float64   `json:"price"`
	Volume    int       `json:"volume"`
	Timestamp time.Time `json:"timestamp"`
}

// NewCompany creates a new company
func NewCompany(id, name, ownerID string, companyType CompanyType, initialShares int, initialPrice float64) *Company {
	now := time.Now()
	return &Company{
		ID:           id,
		Name:         name,
		OwnerID:      ownerID,
		Type:         companyType,
		TotalShares:  initialShares,
		SharePrice:   initialPrice,
		Revenue:      0,
		Expenses:     0,
		Profit:       0,
		Assets:       float64(initialShares) * initialPrice,
		PriceHistory: []StockPrice{{Price: initialPrice, Volume: 0, Timestamp: now}},
		Open:         true,
		LastUpdated:  now,
	}
}

// CalculateMarketCap calculates market capitalization
func (c *Company) CalculateMarketCap() float64 {
	return float64(c.TotalShares) * c.SharePrice
}

// UpdatePrice updates share price based on supply/demand
func (c *Company) UpdatePrice(demand float64) {
	// Simple price adjustment based on demand (-1 to 1)
	change := c.SharePrice * demand * 0.05 // Max 5% change per update
	c.SharePrice += change
	
	if c.SharePrice < 1.0 {
		c.SharePrice = 1.0 // Minimum price
	}
	
	c.LastUpdated = time.Now()
	
	// Add to history
	c.PriceHistory = append(c.PriceHistory, StockPrice{
		Price:     c.SharePrice,
		Volume:    0,
		Timestamp: time.Now(),
	})
	
	// Keep last 100 price points
	if len(c.PriceHistory) > 100 {
		c.PriceHistory = c.PriceHistory[len(c.PriceHistory)-100:]
	}
}

// RecordRevenue records revenue
func (c *Company) RecordRevenue(amount float64) {
	c.Revenue += amount
	c.updateProfit()
}

// RecordExpenses records expenses
func (c *Company) RecordExpenses(amount float64) {
	c.Expenses += amount
	c.updateProfit()
}

// updateProfit updates profit
func (c *Company) updateProfit() {
	c.Profit = c.Revenue - c.Expenses
}

// GetPriceChange returns price change over period
func (c *Company) GetPriceChange(period time.Duration) float64 {
	if len(c.PriceHistory) < 2 {
		return 0
	}
	
	cutoff := time.Now().Add(-period)
	var oldPrice float64
	
	for _, price := range c.PriceHistory {
		if price.Timestamp.Before(cutoff) {
			oldPrice = price.Price
		}
	}
	
	if oldPrice == 0 {
		return 0
	}
	
	return ((c.SharePrice - oldPrice) / oldPrice) * 100
}

// Shareholding represents a player's share ownership
type Shareholding struct {
	PlayerID    string    `json:"player_id"`
	CompanyID   string    `json:"company_id"`
	Shares      int       `json:"shares"`
	AvgBuyPrice float64   `json:"avg_buy_price"`
	PurchasedAt time.Time `json:"purchased_at"`
}

// GetCurrentValue returns current value of holding
func (s *Shareholding) GetCurrentValue(currentPrice float64) float64 {
	return float64(s.Shares) * currentPrice
}

// GetProfitLoss returns profit/loss
func (s *Shareholding) GetProfitLoss(currentPrice float64) float64 {
	currentValue := s.GetCurrentValue(currentPrice)
	costBasis := float64(s.Shares) * s.AvgBuyPrice
	return currentValue - costBasis
}

// StockMarket manages the stock exchange
type StockMarket struct {
	companies     map[string]*Company
	holdings      map[string]map[string]*Shareholding // playerID -> companyID -> holding
	
	marketOpen    bool
	openTime      time.Time
	closeTime     time.Time
	
	// Transaction tax
	taxRate       float64
	
	walletMgr     *WalletManager
	
	storagePath   string
}

// NewStockMarket creates a new stock market
func NewStockMarket(walletMgr *WalletManager, storageDir string) *StockMarket {
	return &StockMarket{
		companies:     make(map[string]*Company),
		holdings:      make(map[string]map[string]*Shareholding),
		marketOpen:    true,
		openTime:      time.Now(),
		closeTime:     time.Now().Add(10 * time.Hour),
		taxRate:       0.01, // 1% transaction tax
		walletMgr:     walletMgr,
		storagePath:   filepath.Join(storageDir, "stockmarket.json"),
	}
}

// RegisterCompany registers a new company
func (sm *StockMarket) RegisterCompany(id, name, ownerID string, companyType CompanyType, initialShares int, initialPrice float64) (*Company, error) {
	if _, exists := sm.companies[id]; exists {
		return nil, fmt.Errorf("company with ID '%s' already exists", id)
	}
	
	company := NewCompany(id, name, ownerID, companyType, initialShares, initialPrice)
	sm.companies[id] = company
	
	return company, nil
}

// GetCompany gets a company
func (sm *StockMarket) GetCompany(companyID string) (*Company, bool) {
	company, exists := sm.companies[companyID]
	return company, exists
}

// GetAllCompanies returns all companies
func (sm *StockMarket) GetAllCompanies() []*Company {
	result := make([]*Company, 0, len(sm.companies))
	for _, company := range sm.companies {
		result = append(result, company)
	}
	return result
}

// GetOpenCompanies returns companies open for trading
func (sm *StockMarket) GetOpenCompanies() []*Company {
	result := make([]*Company, 0)
	for _, company := range sm.companies {
		if company.Open {
			result = append(result, company)
		}
	}
	return result
}

// IsMarketOpen checks if market is currently open
func (sm *StockMarket) IsMarketOpen() bool {
	now := time.Now()
	
	// Simple daily cycle: open 8am-6pm
	hour := now.Hour()
	return hour >= 8 && hour < 18
}

// BuyShares buys shares of a company
func (sm *StockMarket) BuyShares(playerID, companyID string, shares int) error {
	if !sm.IsMarketOpen() {
		return fmt.Errorf("market is closed")
	}
	
	company, exists := sm.GetCompany(companyID)
	if !exists {
		return fmt.Errorf("company not found")
	}
	
	if !company.Open {
		return fmt.Errorf("company is not open for trading")
	}
	
	// Calculate cost
	cost := float64(shares) * company.SharePrice
	tax := cost * sm.taxRate
	total := cost + tax
	
	// Check wallet
	wallet := sm.walletMgr.GetWallet(playerID)
	if wallet == nil || !wallet.CanAfford(total) {
		return fmt.Errorf("insufficient funds")
	}
	
	// Deduct money
	_, success := wallet.Remove(total, TransactionShop, "STOCK_MARKET", fmt.Sprintf("Bought %d shares of %s", shares, company.Name))
	if !success {
		return fmt.Errorf("payment failed")
	}
	
	// Create or update holding
	if _, exists := sm.holdings[playerID]; !exists {
		sm.holdings[playerID] = make(map[string]*Shareholding)
	}
	
	holding, exists := sm.holdings[playerID][companyID]
	if !exists {
		holding = &Shareholding{
			PlayerID:    playerID,
			CompanyID:   companyID,
			Shares:      0,
			AvgBuyPrice: 0,
			PurchasedAt: time.Now(),
		}
		sm.holdings[playerID][companyID] = holding
	}
	
	// Update average buy price
	totalCost := float64(holding.Shares)*holding.AvgBuyPrice + cost
	holding.Shares += shares
	holding.AvgBuyPrice = totalCost / float64(holding.Shares)
	
	// Increase demand (price goes up slightly)
	company.UpdatePrice(0.1)
	
	return nil
}

// SellShares sells shares of a company
func (sm *StockMarket) SellShares(playerID, companyID string, shares int) error {
	if !sm.IsMarketOpen() {
		return fmt.Errorf("market is closed")
	}
	
	company, exists := sm.GetCompany(companyID)
	if !exists {
		return fmt.Errorf("company not found")
	}
	
	// Check holding
	holding, exists := sm.holdings[playerID][companyID]
	if !exists || holding.Shares < shares {
		return fmt.Errorf("insufficient shares")
	}
	
	// Calculate proceeds
	proceeds := float64(shares) * company.SharePrice
	tax := proceeds * sm.taxRate
	net := proceeds - tax
	
	// Add money to wallet
	wallet := sm.walletMgr.GetOrCreateWallet(playerID)
	wallet.Add(net, TransactionShop, "STOCK_MARKET", fmt.Sprintf("Sold %d shares of %s", shares, company.Name))
	
	// Update holding
	holding.Shares -= shares
	if holding.Shares == 0 {
		delete(sm.holdings[playerID], companyID)
	}
	
	// Decrease demand (price goes down slightly)
	company.UpdatePrice(-0.1)
	
	return nil
}

// GetHolding gets a player's holding in a company
func (sm *StockMarket) GetHolding(playerID, companyID string) (*Shareholding, bool) {
	if playerHoldings, exists := sm.holdings[playerID]; exists {
		holding, exists := playerHoldings[companyID]
		return holding, exists
	}
	return nil, false
}

// GetPlayerHoldings gets all holdings for a player
func (sm *StockMarket) GetPlayerHoldings(playerID string) []*Shareholding {
	if playerHoldings, exists := sm.holdings[playerID]; exists {
		result := make([]*Shareholding, 0, len(playerHoldings))
		for _, holding := range playerHoldings {
			result = append(result, holding)
		}
		return result
	}
	return []*Shareholding{}
}

// GetPortfolioValue calculates total portfolio value for a player
func (sm *StockMarket) GetPortfolioValue(playerID string) float64 {
	total := 0.0
	
	for companyID, holding := range sm.holdings[playerID] {
		if company, exists := sm.GetCompany(companyID); exists {
			total += holding.GetCurrentValue(company.SharePrice)
		}
	}
	
	return total
}

// GetLeaderboard returns top investors by portfolio value
func (sm *StockMarket) GetLeaderboard(count int) []struct {
	PlayerID string
	Value      float64
} {
	// Calculate portfolio values
	portfolios := make(map[string]float64)
	
	for playerID := range sm.holdings {
		portfolios[playerID] = sm.GetPortfolioValue(playerID)
	}
	
	// Convert to slice
	result := make([]struct {
		PlayerID string
		Value      float64
	}, 0, len(portfolios))
	
	for id, value := range portfolios {
		result = append(result, struct {
			PlayerID string
			Value      float64
		}{id, value})
	}
	
	// Sort by value
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Value < result[j].Value {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	
	if count > len(result) {
		count = len(result)
	}
	
	return result[:count]
}

// Update runs daily updates on all companies
func (sm *StockMarket) Update() {
	for _, company := range sm.companies {
		// Update price based on company performance
		if company.Profit > 0 {
			// Profitable company = price goes up
			company.UpdatePrice(0.05)
		} else if company.Profit < 0 {
			// Unprofitable = price goes down
			company.UpdatePrice(-0.05)
		}
		
		// Random fluctuation
		fluctuation := (math.Sin(float64(time.Now().Unix())) * 0.02) // -2% to +2%
		company.UpdatePrice(fluctuation)
		
		// Reset daily metrics
		company.Revenue = 0
		company.Expenses = 0
		company.Profit = 0
	}
}

// PayDividends pays dividends to shareholders
func (sm *StockMarket) PayDividends(companyID string, dividendPerShare float64) error {
	company, exists := sm.GetCompany(companyID)
	if !exists {
		return fmt.Errorf("company not found")
	}
	
	totalDividend := 0.0
	
	// Pay each shareholder
	for playerID, playerHoldings := range sm.holdings {
		if holding, exists := playerHoldings[companyID]; exists {
			dividend := float64(holding.Shares) * dividendPerShare
			totalDividend += dividend
			
			// Add to wallet
			wallet := sm.walletMgr.GetOrCreateWallet(playerID)
			wallet.Add(dividend, TransactionInterest, companyID, fmt.Sprintf("Dividend from %s", company.Name))
		}
	}
	
	// Record as company expense
	company.RecordExpenses(totalDividend)
	
	return nil
}

// Save saves stock market data
func (sm *StockMarket) Save() error {
	data := struct {
		Companies map[string]*Company                     `json:"companies"`
		Holdings  map[string]map[string]*Shareholding     `json:"holdings"`
	}{
		Companies: sm.companies,
		Holdings:  sm.holdings,
	}
	
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}
	
	if err := os.WriteFile(sm.storagePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	
	return nil
}

// Load loads stock market data
func (sm *StockMarket) Load() error {
	data, err := os.ReadFile(sm.storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read: %w", err)
	}
	
	var loaded struct {
		Companies map[string]*Company                     `json:"companies"`
		Holdings  map[string]map[string]*Shareholding     `json:"holdings"`
	}
	
	if err := json.Unmarshal(data, &loaded); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}
	
	sm.companies = loaded.Companies
	if sm.companies == nil {
		sm.companies = make(map[string]*Company)
	}
	
	sm.holdings = loaded.Holdings
	if sm.holdings == nil {
		sm.holdings = make(map[string]map[string]*Shareholding)
	}
	
	return nil
}

// GetMarketStats returns market statistics
func (sm *StockMarket) GetMarketStats() (totalCompanies, openCompanies int, totalMarketCap float64) {
	totalCompanies = len(sm.companies)
	
	for _, company := range sm.companies {
		if company.Open {
			openCompanies++
		}
		totalMarketCap += company.CalculateMarketCap()
	}
	
	return
}
