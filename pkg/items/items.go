package items

import "image/color"

// ItemType represents the type of an item
type ItemType int

const (
	NONE ItemType = iota
	DIRT_BLOCK
	GRASS_BLOCK
	STONE_BLOCK
	SAND_BLOCK
	LOG_BLOCK
	COAL
	IRON_INGOT
	GOLD_INGOT
	DIAMOND
	IRON_PICKAXE
	STONE_PICKAXE
	WOODEN_PICKAXE
	PLANKS
	STICK
	WORKBENCH
	FURNACE
	GEL
	STRING
	ROTTEN_FLESH
	// Weapons
	WOODEN_SWORD
	STONE_SWORD
	IRON_SWORD
	DIAMOND_SWORD
	BOW
	MAGIC_WAND
	// Armor
	LEATHER_HELMET
	LEATHER_CHESTPLATE
	LEATHER_LEGGINGS
	LEATHER_BOOTS
	IRON_HELMET
	IRON_CHESTPLATE
	IRON_LEGGINGS
	IRON_BOOTS
	DIAMOND_HELMET
	DIAMOND_CHESTPLATE
	DIAMOND_LEGGINGS
	DIAMOND_BOOTS
	ANVIL
)

// ItemProperties defines the properties of an item type
type ItemProperties struct {
	ID          ItemType
	Name        string
	IconColor   color.RGBA
	Description string
	StackSize   int
	Durability  int // For tools, -1 for indestructible
	IsTool      bool
	ToolPower   float64 // Mining speed multiplier
	IsPlaceable bool    // Can be placed as a block
	BlockType   string  // Corresponding block type if placeable
	// Weapon properties
	IsWeapon     bool
	WeaponDamage float64
	WeaponRange  float64
	WeaponSpeed  float64 // Attacks per second
	WeaponType   string  // "melee", "ranged", "magic"
	// Armor properties
	IsArmor      bool
	ArmorType    string // "helmet", "chestplate", "leggings", "boots"
	ArmorDefense float64
}

var ItemDefinitions = map[ItemType]*ItemProperties{
	NONE: {
		ID:         NONE,
		Name:       "None",
		IconColor:  color.RGBA{0, 0, 0, 0},
		StackSize:  64,
		Durability: -1,
		IsTool:     false,
	},
	DIRT_BLOCK: {
		ID:          DIRT_BLOCK,
		Name:        "Dirt",
		IconColor:   color.RGBA{139, 90, 43, 255},
		Description: "A block of dirt",
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
		IsPlaceable: true,
		BlockType:   "dirt",
	},
	GRASS_BLOCK: {
		ID:          GRASS_BLOCK,
		Name:        "Grass Block",
		IconColor:   color.RGBA{100, 200, 100, 255},
		Description: "A block of grass",
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
		IsPlaceable: true,
		BlockType:   "grass",
	},
	STONE_BLOCK: {
		ID:          STONE_BLOCK,
		Name:        "Stone",
		IconColor:   color.RGBA{169, 169, 169, 255},
		Description: "A block of stone",
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
		IsPlaceable: true,
		BlockType:   "stone",
	},
	SAND_BLOCK: {
		ID:          SAND_BLOCK,
		Name:        "Sand",
		IconColor:   color.RGBA{238, 214, 175, 255},
		Description: "A block of sand",
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
		IsPlaceable: true,
		BlockType:   "sand",
	},
	LOG_BLOCK: {
		ID:          LOG_BLOCK,
		Name:        "Log",
		IconColor:   color.RGBA{139, 69, 19, 255},
		Description: "A wooden log",
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
		IsPlaceable: true,
		BlockType:   "log",
	},
	COAL: {
		ID:          COAL,
		Name:        "Coal",
		IconColor:   color.RGBA{45, 45, 45, 255},
		Description: "A piece of coal",
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
	},
	IRON_INGOT: {
		ID:          IRON_INGOT,
		Name:        "Iron Ingot",
		IconColor:   color.RGBA{169, 166, 150, 255},
		Description: "An iron ingot",
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
	},
	GOLD_INGOT: {
		ID:          GOLD_INGOT,
		Name:        "Gold Ingot",
		IconColor:   color.RGBA{255, 215, 0, 255},
		Description: "A gold ingot",
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
	},
	DIAMOND: {
		ID:          DIAMOND,
		Name:        "Diamond",
		IconColor:   color.RGBA{0, 255, 255, 255},
		Description: "A diamond",
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
	},
	PLANKS: {
		ID:          PLANKS,
		Name:        "Wooden Planks",
		IconColor:   color.RGBA{205, 133, 63, 255},
		Description: "Crafted wooden planks",
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
		IsPlaceable: true,
		BlockType:   "plank",
	},
	STICK: {
		ID:          STICK,
		Name:        "Stick",
		IconColor:   color.RGBA{139, 69, 19, 255},
		Description: "A wooden stick",
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
	},
	WORKBENCH: {
		ID:          WORKBENCH,
		Name:        "Workbench",
		IconColor:   color.RGBA{139, 69, 19, 255},
		Description: "A crafting station",
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
		IsPlaceable: true,
		BlockType:   "workbench",
	},
	FURNACE: {
		ID:          FURNACE,
		Name:        "Furnace",
		IconColor:   color.RGBA{169, 169, 169, 255},
		Description: "Used for smelting",
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
		IsPlaceable: true,
		BlockType:   "furnace",
	},
	ANVIL: {
		ID:          ANVIL,
		Name:        "Anvil",
		IconColor:   color.RGBA{169, 169, 169, 255},
		Description: "Used for forging weapons and armor",
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
		IsPlaceable: true,
		BlockType:   "anvil",
	},
	WOODEN_PICKAXE: {
		ID:          WOODEN_PICKAXE,
		Name:        "Wooden Pickaxe",
		IconColor:   color.RGBA{139, 69, 19, 255},
		Description: "A basic wooden pickaxe",
		StackSize:   1,
		Durability:  60,
		IsTool:      true,
		ToolPower:   2.0,
	},
	STONE_PICKAXE: {
		ID:          STONE_PICKAXE,
		Name:        "Stone Pickaxe",
		IconColor:   color.RGBA{169, 169, 169, 255},
		Description: "A sturdy stone pickaxe",
		StackSize:   1,
		Durability:  132,
		IsTool:      true,
		ToolPower:   4.0,
	},
	IRON_PICKAXE: {
		ID:          IRON_PICKAXE,
		Name:        "Iron Pickaxe",
		IconColor:   color.RGBA{169, 166, 150, 255},
		Description: "A durable iron pickaxe",
		StackSize:   1,
		Durability:  251,
		IsTool:      true,
		ToolPower:   6.0,
	},
	GEL: {
		ID:          GEL,
		Name:        "Gel",
		IconColor:   color.RGBA{0, 255, 0, 255},
		Description: "A sticky substance from slimes",
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
	},
	STRING: {
		ID:          STRING,
		Name:        "String",
		IconColor:   color.RGBA{255, 255, 255, 255},
		Description: "A piece of string from spiders",
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
	},
	ROTTEN_FLESH: {
		ID:          ROTTEN_FLESH,
		Name:        "Rotten Flesh",
		IconColor:   color.RGBA{139, 69, 19, 255},
		Description: "Decaying flesh from zombies",
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
	},
	WOODEN_SWORD: {
		ID:           WOODEN_SWORD,
		Name:         "Wooden Sword",
		IconColor:    color.RGBA{139, 69, 19, 255},
		Description:  "A basic wooden sword",
		StackSize:    1,
		Durability:   60,
		IsWeapon:     true,
		WeaponDamage: 4.0,
		WeaponRange:  80.0,
		WeaponSpeed:  1.2,
		WeaponType:   "melee",
	},
	STONE_SWORD: {
		ID:           STONE_SWORD,
		Name:         "Stone Sword",
		IconColor:    color.RGBA{169, 169, 169, 255},
		Description:  "A stone sword",
		StackSize:    1,
		Durability:   132,
		IsWeapon:     true,
		WeaponDamage: 6.0,
		WeaponRange:  80.0,
		WeaponSpeed:  1.3,
		WeaponType:   "melee",
	},
	IRON_SWORD: {
		ID:           IRON_SWORD,
		Name:         "Iron Sword",
		IconColor:    color.RGBA{169, 166, 150, 255},
		Description:  "An iron sword",
		StackSize:    1,
		Durability:   251,
		IsWeapon:     true,
		WeaponDamage: 8.0,
		WeaponRange:  80.0,
		WeaponSpeed:  1.4,
		WeaponType:   "melee",
	},
	DIAMOND_SWORD: {
		ID:           DIAMOND_SWORD,
		Name:         "Diamond Sword",
		IconColor:    color.RGBA{0, 255, 255, 255},
		Description:  "A diamond sword",
		StackSize:    1,
		Durability:   1562,
		IsWeapon:     true,
		WeaponDamage: 10.0,
		WeaponRange:  80.0,
		WeaponSpeed:  1.5,
		WeaponType:   "melee",
	},
	BOW: {
		ID:           BOW,
		Name:         "Bow",
		IconColor:    color.RGBA{139, 69, 19, 255},
		Description:  "A wooden bow",
		StackSize:    1,
		Durability:   385,
		IsWeapon:     true,
		WeaponDamage: 6.0,
		WeaponRange:  200.0,
		WeaponSpeed:  1.0,
		WeaponType:   "ranged",
	},
	MAGIC_WAND: {
		ID:           MAGIC_WAND,
		Name:         "Magic Wand",
		IconColor:    color.RGBA{255, 0, 255, 255},
		Description:  "A magical wand",
		StackSize:    1,
		Durability:   100,
		IsWeapon:     true,
		WeaponDamage: 8.0,
		WeaponRange:  150.0,
		WeaponSpeed:  0.8,
		WeaponType:   "magic",
	},
	LEATHER_HELMET: {
		ID:           LEATHER_HELMET,
		Name:         "Leather Helmet",
		IconColor:    color.RGBA{139, 69, 19, 255},
		Description:  "Basic leather helmet",
		StackSize:    1,
		Durability:   56,
		IsArmor:      true,
		ArmorType:    "helmet",
		ArmorDefense: 1.0,
	},
	LEATHER_CHESTPLATE: {
		ID:           LEATHER_CHESTPLATE,
		Name:         "Leather Chestplate",
		IconColor:    color.RGBA{139, 69, 19, 255},
		Description:  "Basic leather chestplate",
		StackSize:    1,
		Durability:   81,
		IsArmor:      true,
		ArmorType:    "chestplate",
		ArmorDefense: 3.0,
	},
	LEATHER_LEGGINGS: {
		ID:           LEATHER_LEGGINGS,
		Name:         "Leather Leggings",
		IconColor:    color.RGBA{139, 69, 19, 255},
		Description:  "Basic leather leggings",
		StackSize:    1,
		Durability:   76,
		IsArmor:      true,
		ArmorType:    "leggings",
		ArmorDefense: 2.0,
	},
	LEATHER_BOOTS: {
		ID:           LEATHER_BOOTS,
		Name:         "Leather Boots",
		IconColor:    color.RGBA{139, 69, 19, 255},
		Description:  "Basic leather boots",
		StackSize:    1,
		Durability:   66,
		IsArmor:      true,
		ArmorType:    "boots",
		ArmorDefense: 1.0,
	},
	IRON_HELMET: {
		ID:           IRON_HELMET,
		Name:         "Iron Helmet",
		IconColor:    color.RGBA{169, 166, 150, 255},
		Description:  "Iron helmet",
		StackSize:    1,
		Durability:   166,
		IsArmor:      true,
		ArmorType:    "helmet",
		ArmorDefense: 2.0,
	},
	IRON_CHESTPLATE: {
		ID:           IRON_CHESTPLATE,
		Name:         "Iron Chestplate",
		IconColor:    color.RGBA{169, 166, 150, 255},
		Description:  "Iron chestplate",
		StackSize:    1,
		Durability:   241,
		IsArmor:      true,
		ArmorType:    "chestplate",
		ArmorDefense: 6.0,
	},
	IRON_LEGGINGS: {
		ID:           IRON_LEGGINGS,
		Name:         "Iron Leggings",
		IconColor:    color.RGBA{169, 166, 150, 255},
		Description:  "Iron leggings",
		StackSize:    1,
		Durability:   226,
		IsArmor:      true,
		ArmorType:    "leggings",
		ArmorDefense: 5.0,
	},
	IRON_BOOTS: {
		ID:           IRON_BOOTS,
		Name:         "Iron Boots",
		IconColor:    color.RGBA{169, 166, 150, 255},
		Description:  "Iron boots",
		StackSize:    1,
		Durability:   196,
		IsArmor:      true,
		ArmorType:    "boots",
		ArmorDefense: 2.0,
	},
	DIAMOND_HELMET: {
		ID:           DIAMOND_HELMET,
		Name:         "Diamond Helmet",
		IconColor:    color.RGBA{0, 255, 255, 255},
		Description:  "Diamond helmet",
		StackSize:    1,
		Durability:   364,
		IsArmor:      true,
		ArmorType:    "helmet",
		ArmorDefense: 3.0,
	},
	DIAMOND_CHESTPLATE: {
		ID:           DIAMOND_CHESTPLATE,
		Name:         "Diamond Chestplate",
		IconColor:    color.RGBA{0, 255, 255, 255},
		Description:  "Diamond chestplate",
		StackSize:    1,
		Durability:   529,
		IsArmor:      true,
		ArmorType:    "chestplate",
		ArmorDefense: 8.0,
	},
	DIAMOND_LEGGINGS: {
		ID:           DIAMOND_LEGGINGS,
		Name:         "Diamond Leggings",
		IconColor:    color.RGBA{0, 255, 255, 255},
		Description:  "Diamond leggings",
		StackSize:    1,
		Durability:   496,
		IsArmor:      true,
		ArmorType:    "leggings",
		ArmorDefense: 6.0,
	},
	DIAMOND_BOOTS: {
		ID:           DIAMOND_BOOTS,
		Name:         "Diamond Boots",
		IconColor:    color.RGBA{0, 255, 255, 255},
		Description:  "Diamond boots",
		StackSize:    1,
		Durability:   430,
		IsArmor:      true,
		ArmorType:    "boots",
		ArmorDefense: 3.0,
	},
}

// Item represents a stack of items
type Item struct {
	Type       ItemType
	Quantity   int
	Durability int // For tools
}

// Inventory represents a player's inventory
type Inventory struct {
	Slots    []Item
	Selected int
}

// NewInventory creates a new inventory with the specified number of slots
func NewInventory(slotCount int) *Inventory {
	slots := make([]Item, slotCount)
	for i := range slots {
		slots[i] = Item{Type: NONE, Quantity: 0, Durability: -1}
	}
	return &Inventory{
		Slots:    slots,
		Selected: 0,
	}
}

// DefaultHotbarItems returns the default items for the hotbar
func DefaultHotbarItems() []Item {
	return []Item{
		{Type: WOODEN_PICKAXE, Quantity: 1, Durability: 60},
		{Type: STONE_PICKAXE, Quantity: 1, Durability: 132},
		{Type: DIRT_BLOCK, Quantity: 10, Durability: -1},
		{Type: STONE_BLOCK, Quantity: 10, Durability: -1},
		{Type: LOG_BLOCK, Quantity: 5, Durability: -1},
		{Type: NONE, Quantity: 0, Durability: -1},
		{Type: NONE, Quantity: 0, Durability: -1},
		{Type: NONE, Quantity: 0, Durability: -1},
	}
}

// AddItem adds an item to the inventory
func (inv *Inventory) AddItem(itemType ItemType, quantity int) bool {
	props := ItemDefinitions[itemType]
	if props == nil {
		return false
	}

	remaining := quantity

	// First, try to stack with existing items
	if props.StackSize > 1 {
		for i := range inv.Slots {
			if inv.Slots[i].Type == itemType && inv.Slots[i].Quantity < props.StackSize {
				canAdd := props.StackSize - inv.Slots[i].Quantity
				add := min(canAdd, remaining)
				inv.Slots[i].Quantity += add
				remaining -= add
				if remaining == 0 {
					return true
				}
			}
		}
	}

	// Then, try to find empty slots
	for i := range inv.Slots {
		if inv.Slots[i].Type == NONE {
			inv.Slots[i].Type = itemType
			inv.Slots[i].Durability = props.Durability
			inv.Slots[i].Quantity = min(remaining, props.StackSize)
			remaining -= inv.Slots[i].Quantity
			if remaining == 0 {
				return true
			}
		}
	}

	return remaining == 0
}

// RemoveItem removes items from the selected slot
func (inv *Inventory) RemoveItem(quantity int) bool {
	if inv.Selected >= len(inv.Slots) {
		return false
	}

	slot := &inv.Slots[inv.Selected]
	if slot.Quantity < quantity {
		return false
	}

	slot.Quantity -= quantity
	if slot.Quantity <= 0 {
		slot.Type = NONE
		slot.Durability = -1
	}

	return true
}

// GetSelectedItem returns the currently selected item
func (inv *Inventory) GetSelectedItem() *Item {
	if inv.Selected >= len(inv.Slots) {
		return nil
	}
	return &inv.Slots[inv.Selected]
}

// UseItem uses the selected item (returns true if successful)
func (inv *Inventory) UseItem() bool {
	item := inv.GetSelectedItem()
	if item == nil || item.Type == NONE {
		return false
	}

	props := ItemDefinitions[item.Type]
	if props.IsTool && item.Durability > 0 {
		item.Durability--
		if item.Durability <= 0 {
			item.Type = NONE
			item.Quantity = 0
			item.Durability = -1
		}
		return true
	}

	if !props.IsTool && item.Quantity > 0 {
		item.Quantity--
		if item.Quantity <= 0 {
			item.Type = NONE
			item.Durability = -1
		}
		return true
	}

	return false
}

// SelectSlot selects a slot by index
func (inv *Inventory) SelectSlot(index int) bool {
	if index < 0 || index >= len(inv.Slots) {
		return false
	}
	inv.Selected = index
	return true
}

// NextSlot selects the next slot
func (inv *Inventory) NextSlot() {
	inv.Selected = (inv.Selected + 1) % len(inv.Slots)
}

// PrevSlot selects the previous slot
func (inv *Inventory) PrevSlot() {
	inv.Selected = (inv.Selected - 1 + len(inv.Slots)) % len(inv.Slots)
}

// HasItem checks if inventory has at least the specified quantity of an item
func (inv *Inventory) HasItem(itemType ItemType, quantity int) bool {
	total := 0
	for _, slot := range inv.Slots {
		if slot.Type == itemType {
			total += slot.Quantity
			if total >= quantity {
				return true
			}
		}
	}
	return false
}

// RemoveItemType removes items of a specific type from anywhere in inventory
func (inv *Inventory) RemoveItemType(itemType ItemType, quantity int) bool {
	remaining := quantity

	// Remove from slots that have the item
	for i := range inv.Slots {
		if inv.Slots[i].Type == itemType {
			if inv.Slots[i].Quantity <= remaining {
				remaining -= inv.Slots[i].Quantity
				inv.Slots[i].Type = NONE
				inv.Slots[i].Quantity = 0
				inv.Slots[i].Durability = -1
			} else {
				inv.Slots[i].Quantity -= remaining
				remaining = 0
			}

			if remaining == 0 {
				return true
			}
		}
	}

	return remaining == 0
}

// ItemNameByID returns the name of an item type
func ItemNameByID(itemType ItemType) string {
	if props, ok := ItemDefinitions[itemType]; ok {
		return props.Name
	}
	return "Unknown"
}

// ItemColorByID returns the icon color of an item type
func ItemColorByID(itemType ItemType) color.RGBA {
	if props, ok := ItemDefinitions[itemType]; ok {
		return props.IconColor
	}
	return color.RGBA{0, 0, 0, 255}
}

// GetItemProperties returns the properties of an item type
func GetItemProperties(itemType ItemType) *ItemProperties {
	return ItemDefinitions[itemType]
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
