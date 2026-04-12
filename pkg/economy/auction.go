package economy

import (
	"fmt"
	"time"

	"tesselbox/pkg/items"
)

// AuctionStatus represents the status of an auction
type AuctionStatus int

const (
	AuctionActive AuctionStatus = iota
	AuctionEnded
	AuctionCancelled
	AuctionPending // Waiting for item delivery
)

// Bid represents a single bid
type Bid struct {
	BidderID string    `json:"bidder_id"`
	Amount   float64   `json:"amount"`
	Time     time.Time `json:"time"`
}

// Auction represents an auction listing
type Auction struct {
	ID       string `json:"id"`
	SellerID string `json:"seller_id"`
	WorldID  string `json:"world_id"`

	// Item
	Item     items.Item     `json:"item"`
	ItemType items.ItemType `json:"item_type"`
	Quantity int            `json:"quantity"`

	// Pricing
	StartPrice   float64 `json:"start_price"`
	BuyNowPrice  float64 `json:"buy_now_price"` // 0 = no buy now
	ReservePrice float64 `json:"reserve_price"` // 0 = no reserve
	CurrentBid   float64 `json:"current_bid"`
	HighBidder   string  `json:"high_bidder,omitempty"`

	// Bidding
	Bids         []Bid   `json:"bids,omitempty"`
	MinIncrement float64 `json:"min_increment"` // Minimum bid increment

	// Timing
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Extended  bool      `json:"extended"` // Extended due to last-minute bid

	// Status
	Status       AuctionStatus `json:"status"`
	CancelReason string        `json:"cancel_reason,omitempty"`

	// Meta
	Views    int      `json:"views"`
	Watchers []string `json:"watchers,omitempty"` // Players watching
}

// NewAuction creates a new auction
func NewAuction(id, sellerID, worldID string, item items.Item, quantity int, startPrice, buyNowPrice, reservePrice float64, duration time.Duration) *Auction {
	now := time.Now()
	return &Auction{
		ID:           id,
		SellerID:     sellerID,
		WorldID:      worldID,
		Item:         item,
		ItemType:     item.Type,
		Quantity:     quantity,
		StartPrice:   startPrice,
		BuyNowPrice:  buyNowPrice,
		ReservePrice: reservePrice,
		CurrentBid:   0,
		HighBidder:   "",
		Bids:         make([]Bid, 0),
		MinIncrement: calculateIncrement(startPrice),
		StartTime:    now,
		EndTime:      now.Add(duration),
		Extended:     false,
		Status:       AuctionActive,
		Views:        0,
		Watchers:     make([]string, 0),
	}
}

// calculateIncrement calculates minimum bid increment based on price
func calculateIncrement(price float64) float64 {
	switch {
	case price < 100:
		return 1
	case price < 1000:
		return 5
	case price < 10000:
		return 10
	default:
		return 100
	}
}

// IsActive checks if auction is still active
func (a *Auction) IsActive() bool {
	if a.Status != AuctionActive {
		return false
	}
	return time.Now().Before(a.EndTime)
}

// TimeRemaining returns time left
func (a *Auction) TimeRemaining() time.Duration {
	remaining := time.Until(a.EndTime)
	if remaining < 0 {
		return 0
	}
	return remaining
}

// PlaceBid places a new bid
func (a *Auction) PlaceBid(bidderID string, amount float64) error {
	if !a.IsActive() {
		return fmt.Errorf("auction is not active")
	}

	// Check minimum bid
	minBid := a.StartPrice
	if a.CurrentBid > 0 {
		minBid = a.CurrentBid + a.MinIncrement
	}

	if amount < minBid {
		return fmt.Errorf("bid must be at least %.2f", minBid)
	}

	// Anti-sniping: extend auction by 5 minutes if bid in last 5 minutes
	if a.TimeRemaining() < 5*time.Minute && !a.Extended {
		a.EndTime = a.EndTime.Add(5 * time.Minute)
		a.Extended = true
	}

	// Record bid
	bid := Bid{
		BidderID: bidderID,
		Amount:   amount,
		Time:     time.Now(),
	}
	a.Bids = append(a.Bids, bid)

	// Update current bid
	a.CurrentBid = amount
	a.HighBidder = bidderID

	return nil
}

// BuyNow executes buy now purchase
func (a *Auction) BuyNow(buyerID string) error {
	if !a.IsActive() {
		return fmt.Errorf("auction is not active")
	}

	if a.BuyNowPrice <= 0 {
		return fmt.Errorf("no buy now price set")
	}

	if a.CurrentBid > 0 {
		return fmt.Errorf("bids have been placed")
	}

	// End auction immediately
	a.CurrentBid = a.BuyNowPrice
	a.HighBidder = buyerID
	a.Status = AuctionEnded

	// Record the purchase as a bid
	bid := Bid{
		BidderID: buyerID,
		Amount:   a.BuyNowPrice,
		Time:     time.Now(),
	}
	a.Bids = append(a.Bids, bid)

	return nil
}

// End ends the auction
func (a *Auction) End() {
	if a.Status != AuctionActive {
		return
	}

	a.Status = AuctionEnded

	// Check if reserve was met
	if a.ReservePrice > 0 && a.CurrentBid < a.ReservePrice {
		// Reserve not met, cancel
		a.Status = AuctionCancelled
		a.CancelReason = "Reserve price not met"
		return
	}

	// Check if any bids were placed
	if len(a.Bids) == 0 {
		a.Status = AuctionCancelled
		a.CancelReason = "No bids placed"
		return
	}
}

// Cancel cancels the auction
func (a *Auction) Cancel(reason string) {
	a.Status = AuctionCancelled
	a.CancelReason = reason
}

// GetBidHistory returns bid history
func (a *Auction) GetBidHistory() []Bid {
	return a.Bids
}

// IsWinning returns if a player is the current high bidder
func (a *Auction) IsWinning(playerID string) bool {
	return a.HighBidder == playerID
}

// AddWatcher adds a player to watchers
func (a *Auction) AddWatcher(playerID string) {
	for _, id := range a.Watchers {
		if id == playerID {
			return // Already watching
		}
	}
	a.Watchers = append(a.Watchers, playerID)
}

// RemoveWatcher removes a player from watchers
func (a *Auction) RemoveWatcher(playerID string) {
	for i, id := range a.Watchers {
		if id == playerID {
			a.Watchers = append(a.Watchers[:i], a.Watchers[i+1:]...)
			return
		}
	}
}

// RecordView increments view count
func (a *Auction) RecordView() {
	a.Views++
}

// AuctionHouse manages all auctions
type AuctionHouse struct {
	auctions map[string]*Auction
	bySeller map[string][]string
	byStatus map[AuctionStatus][]string

	taxRate     float64
	minDuration time.Duration
	maxDuration time.Duration

	economy   *EconomyEngine
	walletMgr *WalletManager
}

// NewAuctionHouse creates a new auction house
func NewAuctionHouse(taxRate float64, economy *EconomyEngine, walletMgr *WalletManager) *AuctionHouse {
	return &AuctionHouse{
		auctions:    make(map[string]*Auction),
		bySeller:    make(map[string][]string),
		byStatus:    make(map[AuctionStatus][]string),
		taxRate:     taxRate,
		minDuration: 1 * time.Hour,
		maxDuration: 7 * 24 * time.Hour,
		economy:     economy,
		walletMgr:   walletMgr,
	}
}

// CreateAuction creates a new auction
func (ah *AuctionHouse) CreateAuction(id, sellerID, worldID string, item items.Item, quantity int, startPrice, buyNowPrice, reservePrice float64, duration time.Duration) (*Auction, error) {
	// Validate duration
	if duration < ah.minDuration {
		return nil, fmt.Errorf("duration too short (minimum %s)", ah.minDuration)
	}
	if duration > ah.maxDuration {
		return nil, fmt.Errorf("duration too long (maximum %s)", ah.maxDuration)
	}

	// Validate prices
	if startPrice <= 0 {
		return nil, fmt.Errorf("start price must be positive")
	}
	if buyNowPrice > 0 && buyNowPrice <= startPrice {
		return nil, fmt.Errorf("buy now price must be higher than start price")
	}
	if reservePrice > 0 && reservePrice < startPrice {
		return nil, fmt.Errorf("reserve price must be at least start price")
	}

	if _, exists := ah.auctions[id]; exists {
		return nil, fmt.Errorf("auction with ID '%s' already exists", id)
	}

	auction := NewAuction(id, sellerID, worldID, item, quantity, startPrice, buyNowPrice, reservePrice, duration)

	ah.auctions[id] = auction
	ah.bySeller[sellerID] = append(ah.bySeller[sellerID], id)
	ah.byStatus[AuctionActive] = append(ah.byStatus[AuctionActive], id)

	return auction, nil
}

// GetAuction gets an auction by ID
func (ah *AuctionHouse) GetAuction(auctionID string) (*Auction, bool) {
	auction, exists := ah.auctions[auctionID]
	return auction, exists
}

// GetActiveAuctions returns all active auctions
func (ah *AuctionHouse) GetActiveAuctions() []*Auction {
	result := make([]*Auction, 0)
	for _, auction := range ah.auctions {
		if auction.IsActive() {
			result = append(result, auction)
		}
	}
	return result
}

// GetAuctionsBySeller returns auctions by a seller
func (ah *AuctionHouse) GetAuctionsBySeller(sellerID string) []*Auction {
	auctionIDs := ah.bySeller[sellerID]
	auctions := make([]*Auction, 0, len(auctionIDs))

	for _, id := range auctionIDs {
		if auction, exists := ah.auctions[id]; exists {
			auctions = append(auctions, auction)
		}
	}

	return auctions
}

// SearchAuctions searches for auctions by item type
func (ah *AuctionHouse) SearchAuctions(itemType items.ItemType) []*Auction {
	result := make([]*Auction, 0)

	for _, auction := range ah.auctions {
		if auction.IsActive() && auction.ItemType == itemType {
			result = append(result, auction)
		}
	}

	return result
}

// PlaceBid places a bid on an auction
func (ah *AuctionHouse) PlaceBid(auctionID, bidderID string, amount float64) error {
	auction, exists := ah.GetAuction(auctionID)
	if !exists {
		return fmt.Errorf("auction not found")
	}

	// Check if bidder is seller
	if auction.SellerID == bidderID {
		return fmt.Errorf("cannot bid on your own auction")
	}

	// Check if bidder has funds
	bidderWallet := ah.walletMgr.GetWallet(bidderID)
	if bidderWallet == nil || !bidderWallet.CanAfford(amount) {
		return fmt.Errorf("insufficient funds")
	}

	// Place the bid
	if err := auction.PlaceBid(bidderID, amount); err != nil {
		return err
	}

	return nil
}

// BuyNow executes a buy now purchase
func (ah *AuctionHouse) BuyNow(auctionID, buyerID string) error {
	auction, exists := ah.GetAuction(auctionID)
	if !exists {
		return fmt.Errorf("auction not found")
	}

	// Check if buyer is seller
	if auction.SellerID == buyerID {
		return fmt.Errorf("cannot buy your own auction")
	}

	// Check if buyer has funds
	buyerWallet := ah.walletMgr.GetWallet(buyerID)
	if buyerWallet == nil || !buyerWallet.CanAfford(auction.BuyNowPrice) {
		return fmt.Errorf("insufficient funds")
	}

	// Execute buy now
	if err := auction.BuyNow(buyerID); err != nil {
		return err
	}

	// Process payment immediately
	return ah.processPayment(auction)
}

// EndAuction ends an auction and processes payment
func (ah *AuctionHouse) EndAuction(auctionID string) error {
	auction, exists := ah.GetAuction(auctionID)
	if !exists {
		return fmt.Errorf("auction not found")
	}

	// End the auction
	auction.End()

	// Process payment if successful
	if auction.Status == AuctionEnded && auction.HighBidder != "" {
		return ah.processPayment(auction)
	}

	return nil
}

// processPayment handles money transfer for completed auction
func (ah *AuctionHouse) processPayment(auction *Auction) error {
	if auction.Status != AuctionEnded {
		return fmt.Errorf("auction not ended")
	}

	if auction.HighBidder == "" || auction.CurrentBid == 0 {
		return fmt.Errorf("no winning bid")
	}

	// Calculate amounts
	tax := auction.CurrentBid * ah.taxRate * ah.economy.ShopTaxMod
	sellerReceives := auction.CurrentBid - tax

	// Transfer from buyer to seller
	buyerWallet := ah.walletMgr.GetWallet(auction.HighBidder)
	sellerWallet := ah.walletMgr.GetOrCreateWallet(auction.SellerID)

	if buyerWallet == nil {
		return fmt.Errorf("buyer wallet not found")
	}

	// Buyer pays
	_, success := buyerWallet.Remove(auction.CurrentBid, TransactionAuction, auction.SellerID, fmt.Sprintf("Auction win: %s", auction.ID))
	if !success {
		return fmt.Errorf("buyer payment failed")
	}

	// Seller receives
	sellerWallet.Add(sellerReceives, TransactionAuction, auction.HighBidder, fmt.Sprintf("Auction sale: %s (after tax)", auction.ID))

	return nil
}

// CancelAuction cancels an auction (only by seller or admin)
func (ah *AuctionHouse) CancelAuction(auctionID, cancellerID string, isAdmin bool) error {
	auction, exists := ah.GetAuction(auctionID)
	if !exists {
		return fmt.Errorf("auction not found")
	}

	// Check permission
	if auction.SellerID != cancellerID && !isAdmin {
		return fmt.Errorf("permission denied")
	}

	// Check if bids placed (only admin can cancel with bids)
	if len(auction.Bids) > 0 && !isAdmin {
		return fmt.Errorf("cannot cancel auction with bids")
	}

	auction.Cancel("Cancelled by " + cancellerID)

	return nil
}

// Update processes all active auctions, ending expired ones
func (ah *AuctionHouse) Update() {
	now := time.Now()

	for _, auction := range ah.auctions {
		if auction.IsActive() && now.After(auction.EndTime) {
			ah.EndAuction(auction.ID)
		}
	}
}

// GetWatchedAuctions returns auctions being watched by a player
func (ah *AuctionHouse) GetWatchedAuctions(playerID string) []*Auction {
	result := make([]*Auction, 0)

	for _, auction := range ah.auctions {
		for _, watcher := range auction.Watchers {
			if watcher == playerID {
				result = append(result, auction)
				break
			}
		}
	}

	return result
}

// GetEndingSoon returns auctions ending within duration
func (ah *AuctionHouse) GetEndingSoon(duration time.Duration) []*Auction {
	result := make([]*Auction, 0)
	now := time.Now()

	for _, auction := range ah.auctions {
		if auction.IsActive() && auction.EndTime.Before(now.Add(duration)) {
			result = append(result, auction)
		}
	}

	return result
}

// GetStats returns auction house statistics
func (ah *AuctionHouse) GetStats() (active, ended, cancelled int, totalVolume float64) {
	for _, auction := range ah.auctions {
		switch auction.Status {
		case AuctionActive:
			active++
		case AuctionEnded:
			ended++
			totalVolume += auction.CurrentBid
		case AuctionCancelled:
			cancelled++
		}
	}
	return
}
