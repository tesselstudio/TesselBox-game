package economy

import (
	"fmt"
	"time"

	"tesselbox/pkg/items"
)

// ShopType represents the type of shop
type ShopType int

const (
	ShopTypePhysical ShopType = iota // Physical chest + sign
	ShopTypeVirtual                  // GUI-only shop
)

// ShopListing represents an item listing in a shop
type ShopListing struct {
	ItemType  items.ItemType `json:"item_type"`
	BuyPrice  float64        `json:"buy_price"`  // Shop buys from players at this price
	SellPrice float64        `json:"sell_price"` // Shop sells to players at this price
	Quantity  int            `json:"quantity"`   // Current stock
	MaxStock  int            `json:"max_stock"`  // Maximum stock capacity
	Dynamic   bool           `json:"dynamic"`    // Price adjusts with economy
}

// Sale represents a completed sale
type Sale struct {
	Timestamp time.Time      `json:"timestamp"`
	ItemType  items.ItemType `json:"item_type"`
	Quantity  int            `json:"quantity"`
	Price     float64        `json:"price"`
	BuyerID   string         `json:"buyer_id"`
	IsBuy     bool           `json:"is_buy"` // true = player bought from shop, false = player sold to shop
}

// Shop represents a player-owned shop
type Shop struct {
	ID      string   `json:"id"`
	OwnerID string   `json:"owner_id"`
	WorldID string   `json:"world_id"`
	Name    string   `json:"name"`
	Type    ShopType `json:"type"`

	// Location (for physical shops)
	X float64 `json:"x"`
	Y float64 `json:"y"`

	// Inventory
	Inventory map[string]ShopListing `json:"inventory"` // key: item type name

	// History
	Sales      []Sale  `json:"sales,omitempty"`
	TotalSales float64 `json:"total_sales"`
	TaxPaid    float64 `json:"tax_paid"`

	// State
	IsOpen      bool `json:"is_open"`
	AutoPricing bool `json:"auto_pricing"` // Adjust prices with economy

	// Meta
	CreatedAt   time.Time `json:"created_at"`
	LastRestock time.Time `json:"last_restock"`
}

// NewShop creates a new shop
func NewShop(id, ownerID, worldID, name string, shopType ShopType) *Shop {
	return &Shop{
		ID:          id,
		OwnerID:     ownerID,
		WorldID:     worldID,
		Name:        name,
		Type:        shopType,
		Inventory:   make(map[string]ShopListing),
		Sales:       make([]Sale, 0),
		IsOpen:      true,
		AutoPricing: false,
		CreatedAt:   time.Now(),
		LastRestock: time.Now(),
	}
}

// SetLocation sets the shop location
func (s *Shop) SetLocation(x, y float64) {
	s.X = x
	s.Y = y
}

// AddItem adds an item to the shop inventory
func (s *Shop) AddItem(itemType items.ItemType, buyPrice, sellPrice float64, quantity, maxStock int, dynamic bool) {
	key := fmt.Sprintf("%d", itemType)
	s.Inventory[key] = ShopListing{
		ItemType:  itemType,
		BuyPrice:  buyPrice,
		SellPrice: sellPrice,
		Quantity:  quantity,
		MaxStock:  maxStock,
		Dynamic:   dynamic,
	}
}

// RemoveItem removes an item from inventory
func (s *Shop) RemoveItem(itemType items.ItemType) {
	key := fmt.Sprintf("%d", itemType)
	delete(s.Inventory, key)
}

// GetItem gets an item listing
func (s *Shop) GetItem(itemType items.ItemType) (ShopListing, bool) {
	key := fmt.Sprintf("%d", itemType)
	listing, exists := s.Inventory[key]
	return listing, exists
}

// GetSellPrice gets the effective sell price (player buys from shop)
func (s *Shop) GetSellPrice(itemType items.ItemType, economyMultiplier float64) (float64, bool) {
	listing, exists := s.GetItem(itemType)
	if !exists {
		return 0, false
	}

	price := listing.SellPrice
	if listing.Dynamic {
		price *= economyMultiplier
	}
	return price, true
}

// GetBuyPrice gets the effective buy price (shop buys from player)
func (s *Shop) GetBuyPrice(itemType items.ItemType, economyMultiplier float64) (float64, bool) {
	listing, exists := s.GetItem(itemType)
	if !exists {
		return 0, false
	}

	price := listing.BuyPrice
	if listing.Dynamic {
		price *= economyMultiplier
	}
	return price, true
}

// CanSell checks if shop can sell quantity of item
func (s *Shop) CanSell(itemType items.ItemType, quantity int) bool {
	listing, exists := s.GetItem(itemType)
	if !exists {
		return false
	}
	return listing.Quantity >= quantity && s.IsOpen
}

// CanBuy checks if shop can buy quantity of item
func (s *Shop) CanBuy(itemType items.ItemType) bool {
	listing, exists := s.GetItem(itemType)
	if !exists {
		return false
	}
	return listing.Quantity < listing.MaxStock && s.IsOpen
}

// ExecuteSell sells items to a player
func (s *Shop) ExecuteSell(itemType items.ItemType, quantity int, buyerID string, price float64) bool {
	key := fmt.Sprintf("%d", itemType)
	listing, exists := s.Inventory[key]
	if !exists {
		return false
	}

	if listing.Quantity < quantity {
		return false
	}

	// Update inventory
	listing.Quantity -= quantity
	s.Inventory[key] = listing

	// Record sale
	sale := Sale{
		Timestamp: time.Now(),
		ItemType:  itemType,
		Quantity:  quantity,
		Price:     price,
		BuyerID:   buyerID,
		IsBuy:     true, // Player bought from shop
	}
	s.Sales = append(s.Sales, sale)
	s.TotalSales += price

	return true
}

// ExecuteBuy buys items from a player
func (s *Shop) ExecuteBuy(itemType items.ItemType, quantity int, sellerID string, price float64) bool {
	key := fmt.Sprintf("%d", itemType)
	listing, exists := s.Inventory[key]
	if !exists {
		return false
	}

	if listing.Quantity+quantity > listing.MaxStock {
		return false
	}

	// Update inventory
	listing.Quantity += quantity
	s.Inventory[key] = listing

	// Record sale
	sale := Sale{
		Timestamp: time.Now(),
		ItemType:  itemType,
		Quantity:  quantity,
		Price:     price,
		BuyerID:   sellerID,
		IsBuy:     false, // Player sold to shop
	}
	s.Sales = append(s.Sales, sale)

	return true
}

// Restock adds stock to all items (simulated supply)
func (s *Shop) Restock() {
	for key, listing := range s.Inventory {
		if listing.Quantity < listing.MaxStock {
			// Add 10% of max stock or 1, whichever is larger
			add := listing.MaxStock / 10
			if add < 1 {
				add = 1
			}
			listing.Quantity += add
			if listing.Quantity > listing.MaxStock {
				listing.Quantity = listing.MaxStock
			}
			s.Inventory[key] = listing
		}
	}
	s.LastRestock = time.Now()
}

// GetRecentSales returns recent sales
func (s *Shop) GetRecentSales(count int) []Sale {
	if count > len(s.Sales) {
		count = len(s.Sales)
	}

	start := len(s.Sales) - count
	if start < 0 {
		start = 0
	}

	result := make([]Sale, count)
	copy(result, s.Sales[start:])
	return result
}

// GetInventoryList returns all items in inventory
func (s *Shop) GetInventoryList() []ShopListing {
	result := make([]ShopListing, 0, len(s.Inventory))
	for _, listing := range s.Inventory {
		result = append(result, listing)
	}
	return result
}

// ShopManager manages all shops
type ShopManager struct {
	shops      map[string]*Shop
	byOwner    map[string][]string
	byLocation map[string]string // "world:x,y" -> shopID

	taxRate   float64
	economy   *EconomyEngine
	walletMgr *WalletManager
}

// NewShopManager creates a new shop manager
func NewShopManager(taxRate float64, economy *EconomyEngine, walletMgr *WalletManager) *ShopManager {
	return &ShopManager{
		shops:      make(map[string]*Shop),
		byOwner:    make(map[string][]string),
		byLocation: make(map[string]string),
		taxRate:    taxRate,
		economy:    economy,
		walletMgr:  walletMgr,
	}
}

// CreateShop creates a new shop
func (sm *ShopManager) CreateShop(id, ownerID, worldID, name string, shopType ShopType, x, y float64) (*Shop, error) {
	if _, exists := sm.shops[id]; exists {
		return nil, fmt.Errorf("shop with ID '%s' already exists", id)
	}

	shop := NewShop(id, ownerID, worldID, name, shopType)
	shop.SetLocation(x, y)

	sm.shops[id] = shop
	sm.byOwner[ownerID] = append(sm.byOwner[ownerID], id)

	// Register location for physical shops
	if shopType == ShopTypePhysical {
		locationKey := fmt.Sprintf("%s:%.0f,%.0f", worldID, x, y)
		sm.byLocation[locationKey] = id
	}

	return shop, nil
}

// GetShop gets a shop by ID
func (sm *ShopManager) GetShop(shopID string) (*Shop, bool) {
	shop, exists := sm.shops[shopID]
	return shop, exists
}

// GetShopAt gets shop at location
func (sm *ShopManager) GetShopAt(worldID string, x, y float64) (*Shop, bool) {
	locationKey := fmt.Sprintf("%s:%.0f,%.0f", worldID, x, y)
	shopID, exists := sm.byLocation[locationKey]
	if !exists {
		return nil, false
	}
	return sm.GetShop(shopID)
}

// GetShopsByOwner gets all shops for an owner
func (sm *ShopManager) GetShopsByOwner(ownerID string) []*Shop {
	shopIDs := sm.byOwner[ownerID]
	shops := make([]*Shop, 0, len(shopIDs))

	for _, id := range shopIDs {
		if shop, exists := sm.shops[id]; exists {
			shops = append(shops, shop)
		}
	}

	return shops
}

// GetAllShops returns all shops
func (sm *ShopManager) GetAllShops() []*Shop {
	shops := make([]*Shop, 0, len(sm.shops))
	for _, shop := range sm.shops {
		shops = append(shops, shop)
	}
	return shops
}

// GetOpenShops returns all open shops
func (sm *ShopManager) GetOpenShops() []*Shop {
	shops := make([]*Shop, 0)
	for _, shop := range sm.shops {
		if shop.IsOpen {
			shops = append(shops, shop)
		}
	}
	return shops
}

// DeleteShop deletes a shop
func (sm *ShopManager) DeleteShop(shopID string) error {
	shop, exists := sm.shops[shopID]
	if !exists {
		return fmt.Errorf("shop not found")
	}

	// Remove from owner list
	ownerShops := sm.byOwner[shop.OwnerID]
	for i, id := range ownerShops {
		if id == shopID {
			sm.byOwner[shop.OwnerID] = append(ownerShops[:i], ownerShops[i+1:]...)
			break
		}
	}

	// Remove from location
	if shop.Type == ShopTypePhysical {
		locationKey := fmt.Sprintf("%s:%.0f,%.0f", shop.WorldID, shop.X, shop.Y)
		delete(sm.byLocation, locationKey)
	}

	// Remove shop
	delete(sm.shops, shopID)

	return nil
}

// BuyFromShop handles a purchase from a shop
func (sm *ShopManager) BuyFromShop(shopID string, itemType items.ItemType, quantity int, buyerID string) (float64, error) {
	shop, exists := sm.GetShop(shopID)
	if !exists {
		return 0, fmt.Errorf("shop not found")
	}

	if !shop.IsOpen {
		return 0, fmt.Errorf("shop is closed")
	}

	// Get price with economy multiplier
	price, exists := shop.GetSellPrice(itemType, sm.economy.PriceMultiplier)
	if !exists {
		return 0, fmt.Errorf("item not available")
	}

	totalPrice := price * float64(quantity)

	// Check if shop has stock
	if !shop.CanSell(itemType, quantity) {
		return 0, fmt.Errorf("insufficient stock")
	}

	// Check if buyer has funds
	buyerWallet := sm.walletMgr.GetWallet(buyerID)
	if buyerWallet == nil || !buyerWallet.CanAfford(totalPrice) {
		return 0, fmt.Errorf("insufficient funds")
	}

	// Calculate tax
	tax := totalPrice * sm.taxRate * sm.economy.ShopTaxMod
	sellerReceives := totalPrice - tax

	// Execute transaction
	// Buyer pays
	_, success := buyerWallet.Remove(totalPrice, TransactionShop, shop.OwnerID, fmt.Sprintf("Bought %d items from %s", quantity, shop.Name))
	if !success {
		return 0, fmt.Errorf("payment failed")
	}

	// Seller receives (minus tax)
	sellerWallet := sm.walletMgr.GetOrCreateWallet(shop.OwnerID)
	sellerWallet.Add(sellerReceives, TransactionShop, buyerID, fmt.Sprintf("Sale to %s (after tax)", buyerID))

	// Shop executes sale
	shop.ExecuteSell(itemType, quantity, buyerID, totalPrice)
	shop.TaxPaid += tax

	return totalPrice, nil
}

// SellToShop handles selling to a shop
func (sm *ShopManager) SellToShop(shopID string, itemType items.ItemType, quantity int, sellerID string) (float64, error) {
	shop, exists := sm.GetShop(shopID)
	if !exists {
		return 0, fmt.Errorf("shop not found")
	}

	if !shop.IsOpen {
		return 0, fmt.Errorf("shop is closed")
	}

	// Get price with economy multiplier
	price, exists := shop.GetBuyPrice(itemType, sm.economy.PriceMultiplier)
	if !exists {
		return 0, fmt.Errorf("shop not buying this item")
	}

	totalPrice := price * float64(quantity)

	// Check if shop can buy (has space and owner has funds)
	if !shop.CanBuy(itemType) {
		return 0, fmt.Errorf("shop not accepting this item")
	}

	sellerWallet := sm.walletMgr.GetOrCreateWallet(sellerID)
	ownerWallet := sm.walletMgr.GetWallet(shop.OwnerID)

	// Check if owner has funds to buy
	if ownerWallet == nil || !ownerWallet.CanAfford(totalPrice) {
		return 0, fmt.Errorf("shop owner cannot afford purchase")
	}

	// Execute transaction
	// Owner pays
	_, success := ownerWallet.Remove(totalPrice, TransactionShop, sellerID, fmt.Sprintf("Bought %d items from %s", quantity, sellerID))
	if !success {
		return 0, fmt.Errorf("shop owner payment failed")
	}

	// Seller receives
	sellerWallet.Add(totalPrice, TransactionShop, shop.OwnerID, fmt.Sprintf("Sold to %s", shop.Name))

	// Shop executes buy
	shop.ExecuteBuy(itemType, quantity, sellerID, totalPrice)

	return totalPrice, nil
}

// SearchShops searches for shops selling an item
func (sm *ShopManager) SearchShops(itemType items.ItemType) []*Shop {
	result := make([]*Shop, 0)

	for _, shop := range sm.shops {
		if shop.IsOpen {
			if _, exists := shop.GetItem(itemType); exists {
				result = append(result, shop)
			}
		}
	}

	return result
}

// RestockAll restocks all shops
func (sm *ShopManager) RestockAll() {
	for _, shop := range sm.shops {
		shop.Restock()
	}
}

// GetTopShops returns top shops by total sales
func (sm *ShopManager) GetTopShops(count int) []*Shop {
	allShops := sm.GetAllShops()

	// Sort by total sales (bubble sort for simplicity)
	for i := 0; i < len(allShops); i++ {
		for j := i + 1; j < len(allShops); j++ {
			if allShops[i].TotalSales < allShops[j].TotalSales {
				allShops[i], allShops[j] = allShops[j], allShops[i]
			}
		}
	}

	if count > len(allShops) {
		count = len(allShops)
	}

	return allShops[:count]
}

// GetTotalTaxCollected returns total tax collected across all shops
func (sm *ShopManager) GetTotalTaxCollected() float64 {
	total := 0.0
	for _, shop := range sm.shops {
		total += shop.TaxPaid
	}
	return total
}
